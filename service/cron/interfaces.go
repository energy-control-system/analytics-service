package cron

import (
	"analytics-service/service/analytics"
	"time"

	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/golog"
)

type AnalyticsService interface {
	CreateBasicReport(ctx goctx.Context, log golog.Logger, periodStart, periodEnd time.Time) (analytics.Report, error)
}
