package analytics

import (
	_ "embed"

	"github.com/jmoiron/sqlx"
)

var (
	//go:embed sql/add_finished_task.sql
	addFinishedTaskSQL string

	//go:embed sql/add_report.sql
	addReportSQL string

	//go:embed sql/get_all_reports.sql
	getAllReportsSQL string

	//go:embed sql/get_finished_tasks_by_period.sql
	getFinishedTasksByPeriodSQL string
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}
