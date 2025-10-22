package analytics

import (
	"analytics-service/service/analytics"
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/jmoiron/sqlx"
)

var (
	//go:embed sql/add_attachment.sql
	addAttachmentSQL string

	//go:embed sql/add_finished_task.sql
	addFinishedTaskSQL string

	//go:embed sql/add_report.sql
	addReportSQL string

	//go:embed sql/get_all_reports.sql
	getAllReportsSQL string

	//go:embed sql/get_attachments_by_reports.sql
	getAttachmentsByReportSQL string

	//go:embed sql/get_finished_tasks_by_period.sql
	getFinishedTasksByPeriodSQL string
)

type Repository struct {
	postgres   *sqlx.DB
	clickhouse driver.Conn
}

func NewRepository(postgres *sqlx.DB, clickhouse driver.Conn) *Repository {
	return &Repository{
		postgres:   postgres,
		clickhouse: clickhouse,
	}
}

func (r *Repository) AddFinishedTask(ctx context.Context, t analytics.FinishedTask) error {
	dbTask := MapFinishedTaskToDB(t)

	batch, err := r.clickhouse.PrepareBatch(ctx, addFinishedTaskSQL)
	if err != nil {
		return fmt.Errorf("r.clickhouse.PrepareBatch: %w", err)
	}
	defer func() {
		err = errors.Join(err, batch.Close())
	}()

	err = batch.AppendStruct(&dbTask)
	if err != nil {
		err = fmt.Errorf("batch.AppendStruct: %w", err)
		return err
	}

	err = batch.Send()
	if err != nil {
		err = fmt.Errorf("batch.Send: %w", err)
		return err
	}

	return err
}

func (r *Repository) GetFinishedTasksByPeriod(ctx context.Context, periodStart, periodEnd time.Time) ([]analytics.FinishedTask, error) {
	var tasks []FinishedTask
	err := r.clickhouse.Select(ctx, &tasks, getFinishedTasksByPeriodSQL, periodStart, periodEnd)
	if err != nil {
		return nil, fmt.Errorf("r.clickhouse.Select: %w", err)
	}

	return MapFinishedTaskSliceFromDB(tasks), nil
}

func (r *Repository) AddReport(ctx context.Context, report analytics.Report) (analytics.Report, error) {
	tx, err := r.postgres.BeginTxx(ctx, nil)
	if err != nil {
		return analytics.Report{}, fmt.Errorf("r.postgres.BeginTxx: %w", err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		}
	}()

	rows, err := tx.NamedQuery(addReportSQL, MapReportToDB(report))
	if err != nil {
		err = fmt.Errorf("tx.NamedQuery: %w", err)
		return analytics.Report{}, err
	}

	if !rows.Next() {
		err = errors.New("rows.Next == false")
		return analytics.Report{}, err
	}

	var dbReport Report
	if err = rows.StructScan(&dbReport); err != nil {
		err = fmt.Errorf("rows.StructScan: %w", err)
		return analytics.Report{}, err
	}

	if err = rows.Close(); err != nil {
		err = fmt.Errorf("rows.Close: %w", err)
		return analytics.Report{}, err
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("rows.Err: %w", err)
		return analytics.Report{}, err
	}

	newReport := MapReportFromDB(dbReport)
	newReport.Files = report.Files

	_, err = tx.NamedExecContext(ctx, addAttachmentSQL, MapAttachmentSliceToDB(report.Files, newReport.ID))
	if err != nil {
		err = fmt.Errorf("tx.NamedExecContext: %w", err)
		return analytics.Report{}, err
	}

	if err = tx.Commit(); err != nil {
		err = fmt.Errorf("tx.Commit: %w", err)
		return analytics.Report{}, err
	}

	return newReport, err
}

func (r *Repository) GetAllReports(ctx context.Context) ([]analytics.Report, error) {
	tx, err := r.postgres.BeginTxx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, fmt.Errorf("r.postgres.BeginTxx: %w", err)
	}

	var dbReports []Report
	if err = tx.SelectContext(ctx, &dbReports, getAllReportsSQL); err != nil {
		return nil, fmt.Errorf("tx.SelectContext: %w", err)
	}

	ids := make([]int, 0, len(dbReports))
	for _, dbReport := range dbReports {
		ids = append(ids, dbReport.ID)
	}

	var dbAttachments []Attachment
	if err = tx.SelectContext(ctx, &dbAttachments, getAttachmentsByReportSQL, ids); err != nil {
		return nil, fmt.Errorf("tx.SelectContext: %w", err)
	}

	attachmentsMap := make(map[int][]Attachment, len(dbReports))
	for _, dbAttachment := range dbAttachments {
		attachmentsMap[dbAttachment.ReportID] = append(attachmentsMap[dbAttachment.ReportID], dbAttachment)
	}

	reports := make([]analytics.Report, 0, len(dbReports))
	for _, dbReport := range dbReports {
		attachments, ok := attachmentsMap[dbReport.ID]
		if !ok {
			return nil, fmt.Errorf("attachments for report %d not found", dbReport.ID)
		}

		report := MapReportFromDB(dbReport)
		report.Files = MapAttachmentSliceFromDB(attachments)

		reports = append(reports, report)
	}

	return reports, nil
}
