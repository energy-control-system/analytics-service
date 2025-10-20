package analytics

import (
	"analytics-service/cluster/brigade"
	"analytics-service/cluster/file"
	"analytics-service/cluster/inspection"
	"analytics-service/cluster/subscriber"
	"context"
	"io"
	"time"

	"github.com/sunshineOfficial/golib/goctx"
)

type Repository interface {
	AddFinishedTask(ctx context.Context, t FinishedTask) error
	GetFinishedTasksByPeriod(ctx context.Context, periodStart, periodEnd time.Time) ([]FinishedTask, error)
	AddReport(ctx context.Context, r Report) (Report, error)
	GetAllReports(ctx context.Context) ([]Report, error)
}

type InspectionService interface {
	GetInspectionByTaskID(ctx goctx.Context, taskID int) (inspection.Inspection, error)
}

type BrigadeService interface {
	GetBrigadeByID(ctx goctx.Context, id int) (brigade.Brigade, error)
}

type SubscriberService interface {
	GetObjectExtendedByID(ctx goctx.Context, id int) (subscriber.ObjectExtended, error)
}

type FileService interface {
	Upload(ctx goctx.Context, fileName string, file io.Reader) (file.File, error)
}
