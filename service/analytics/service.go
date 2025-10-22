package analytics

import (
	"analytics-service/cluster/file"
	"analytics-service/cluster/inspection"
	"analytics-service/cluster/task"
	"analytics-service/config"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/gokafka"
	"github.com/sunshineOfficial/golib/golog"
	"github.com/sunshineOfficial/golib/gotime"
	"github.com/xuri/excelize/v2"
)

const kafkaSubscribeTimeout = 2 * time.Minute

type Service struct {
	repository        Repository
	inspectionService InspectionService
	brigadeService    BrigadeService
	subscriberService SubscriberService
	fileService       FileService
	templates         config.Templates
}

func NewService(repository Repository, inspectionService InspectionService, brigadeService BrigadeService,
	subscriberService SubscriberService, fileService FileService, templates config.Templates) *Service {
	return &Service{
		repository:        repository,
		inspectionService: inspectionService,
		brigadeService:    brigadeService,
		subscriberService: subscriberService,
		fileService:       fileService,
		templates:         templates,
	}
}

func (s *Service) CreateBasicReport(ctx goctx.Context, log golog.Logger, periodStart, periodEnd time.Time) (Report, error) {
	periodStart = time.Date(periodStart.Year(), periodStart.Month(), periodStart.Day(), 0, 0, 0, 0, gotime.Moscow)
	periodEnd = time.Date(periodEnd.Year(), periodEnd.Month(), periodEnd.Day(), 0, 0, 0, 0, gotime.Moscow)

	fmt.Printf("period start: %v, end: %v\n", periodStart, periodEnd)

	if days := gotime.Days(periodEnd, periodStart); days < 1 {
		return Report{}, fmt.Errorf("period days must be positive, got: %f", days)
	}

	tasks, err := s.repository.GetFinishedTasksByPeriod(ctx, periodStart, periodEnd)
	if err != nil {
		return Report{}, fmt.Errorf("get finished tasks: %w", err)
	}
	if len(tasks) == 0 {
		return Report{}, fmt.Errorf("no finished tasks found from %s to %s", periodStart, periodEnd)
	}

	f, err := excelize.OpenFile(s.templates.BasicReport)
	if err != nil {
		return Report{}, fmt.Errorf("open template file: %w", err)
	}

	defer func() {
		if fErr := f.Close(); fErr != nil {
			log.Errorf("close template file: %v", fErr)
		}
	}()

	sheet := f.GetSheetName(0)
	for i, t := range tasks {
		cell, cellErr := excelize.CoordinatesToCellName(1, i+2)
		if cellErr != nil {
			return Report{}, fmt.Errorf("coordinates to cell name: %w", cellErr)
		}

		workType := ""
		workResult := ""
		switch t.Inspection.Type {
		case inspection.TypeResumption:
			workType = "Возобновление"

			if t.Inspection.Resolution == inspection.ResolutionResumed {
				workResult = "Возобновление"
			} else {
				workResult = "Недопуск"
			}

		case inspection.TypeLimitation:
			workType = "Отключение"

			if t.Inspection.Resolution != inspection.ResolutionLimited {
				workResult = "Отключение"
			} else {
				workResult = "Недопуск"
			}

		default:
			workType = "Контроль ранее введенного ограничения"

			if t.Inspection.IsViolationDetected {
				workResult = "Нарушено"
			} else {
				workResult = "Не нарушено"
			}
		}

		inspectors := make([]string, 0, len(t.Brigade.Inspectors))
		for _, inspector := range t.Brigade.Inspectors {
			inspectors = append(inspectors, fullFIO(inspector.Surname, inspector.Name, inspector.Patronymic))
		}

		err = f.SetSheetRow(sheet, cell, &[]any{
			i + 1,
			t.Object.Address,
			fullFIO(t.Object.Subscriber.Surname, t.Object.Subscriber.Name, t.Object.Subscriber.Patronymic),
			t.Object.Subscriber.AccountNumber,
			t.StartedAt.In(gotime.Moscow).Format(gotime.DateTimeNet),
			t.FinishedAt.In(gotime.Moscow).Format(gotime.DateTimeNet),
			workType,
			workResult,
			strings.Join(inspectors, ", "),
		})
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return Report{}, fmt.Errorf("write file to buffer: %w", err)
	}

	fileName := fmt.Sprintf("Отчет за %s-%s.xlsx", periodStart.Format(gotime.DateOnlyNet), periodEnd.Format(gotime.DateOnlyNet))

	uploadedFile, err := s.fileService.Upload(ctx, fileName, buf)
	if err != nil {
		return Report{}, fmt.Errorf("upload file: %w", err)
	}

	report := Report{
		Type:        ReportTypeBasic,
		Files:       []file.File{uploadedFile},
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
	}

	report, err = s.repository.AddReport(ctx, report)
	if err != nil {
		return Report{}, fmt.Errorf("add report: %w", err)
	}

	return report, nil
}

func fullFIO(surname, name, patronymic string) string {
	result := fmt.Sprintf("%s %s", surname, name)

	if len(patronymic) > 0 {
		result = fmt.Sprintf("%s %s", result, patronymic)
	}

	return result
}

func (s *Service) GetAllReports(ctx goctx.Context) ([]Report, error) {
	reports, err := s.repository.GetAllReports(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all reports from db: %w", err)
	}

	fileIDs := make([]int, 0, len(reports))
	for _, report := range reports {
		for _, f := range report.Files {
			fileIDs = append(fileIDs, f.ID)
		}
	}

	files, err := s.fileService.GetFilesByIDs(ctx, fileIDs)
	if err != nil {
		return nil, fmt.Errorf("get files by ids: %w", err)
	}

	filesMap := make(map[int]file.File, len(files))
	for _, f := range files {
		filesMap[f.ID] = f
	}

	for i, report := range reports {
		for j, reportFile := range report.Files {
			f, ok := filesMap[reportFile.ID]
			if !ok {
				return nil, fmt.Errorf("report %d file %d not found", report.ID, reportFile.ID)
			}

			reports[i].Files[j] = f
		}
	}

	return reports, nil
}

func (s *Service) SubscriberOnTaskEvent(mainCtx context.Context, log golog.Logger) gokafka.Subscriber {
	return func(message gokafka.Message, err error) {
		ctx, cancel := context.WithTimeout(mainCtx, kafkaSubscribeTimeout)
		defer cancel()

		if err != nil {
			log.Errorf("got error on task event: %v", err)
			return
		}

		var event task.Event
		err = json.Unmarshal(message.Value, &event)
		if err != nil {
			log.Errorf("failed to unmarshal task event: %v", err)
			return
		}

		switch event.Type {
		case task.EventTypeAdd:
			err = s.handleAddedTask(ctx, event.Task)
		case task.EventTypeStart:
			err = s.handleStartedTask(ctx, event.Task)
		case task.EventTypeFinish:
			err = s.handleFinishedTask(ctx, event.Task)
		default:
			err = fmt.Errorf("unknown event type: %v", event.Type)
		}

		if err != nil {
			log.Errorf("failed to handle task event (type = %d): %v", event.Type, err)
			return
		}
	}
}

func (s *Service) handleAddedTask(ctx context.Context, t task.Task) error {
	return nil
}

func (s *Service) handleStartedTask(ctx context.Context, t task.Task) error {
	return nil
}

func (s *Service) handleFinishedTask(ctx context.Context, t task.Task) error {
	if t.Status != task.StatusDone {
		return fmt.Errorf("invalid task status: %v", t.Status)
	}

	goCtx := goctx.Wrap(ctx)

	ins, err := s.inspectionService.GetInspectionByTaskID(goCtx, t.ID)
	if err != nil {
		return fmt.Errorf("get inspection by task id: %w", err)
	}

	brig, err := s.brigadeService.GetBrigadeByID(goCtx, *t.BrigadeID)
	if err != nil {
		return fmt.Errorf("get brigade by id: %w", err)
	}

	obj, err := s.subscriberService.GetObjectExtendedByID(goCtx, t.ObjectID)
	if err != nil {
		return fmt.Errorf("get object by id: %w", err)
	}

	finishedTask := MapToFinishedTask(t, ins, brig, obj)

	err = s.repository.AddFinishedTask(goCtx, finishedTask)
	if err != nil {
		return fmt.Errorf("add finished task: %w", err)
	}

	return nil
}
