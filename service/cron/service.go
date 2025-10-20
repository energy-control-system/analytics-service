package cron

import (
	"analytics-service/config"
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/golog"
	"github.com/sunshineOfficial/golib/gotime"
)

type Service struct {
	scheduler        gocron.Scheduler
	settings         config.Cron
	analyticsService AnalyticsService
	running          *atomic.Bool
}

func NewService(settings config.Cron, analyticsService AnalyticsService) *Service {
	return &Service{
		settings:         settings,
		analyticsService: analyticsService,
		running:          &atomic.Bool{},
	}
}

func (s *Service) Start(ctx context.Context, log golog.Logger) error {
	if s.running.Load() {
		return errors.New("already running")
	}

	s.running.Store(true)

	dailyReportTime, err := time.Parse(gotime.TimeOnlyNet, s.settings.DailyReportTime)
	if err != nil {
		return fmt.Errorf("parse daily report time: %w", err)
	}

	s.scheduler, err = gocron.NewScheduler(
		gocron.WithLocation(time.UTC),
		gocron.WithLogger(logger{
			log: log,
		}),
	)
	if err != nil {
		return fmt.Errorf("create cron scheduler: %w", err)
	}

	reportJob, err := s.scheduler.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(uint(dailyReportTime.Hour()), uint(dailyReportTime.Minute()), 0),
			),
		),
		gocron.NewTask(s.dailyReportTask, ctx, log.WithTags("dailyReportTask")),
	)
	if err != nil {
		return fmt.Errorf("create report job: %w", err)
	}

	s.scheduler.Start()

	log.Debugf("started report job %s", reportJob.ID())

	return nil
}

func (s *Service) Stop() error {
	if !s.running.Load() {
		return errors.New("not running")
	}

	s.running.Store(false)

	if err := s.scheduler.Shutdown(); err != nil {
		return fmt.Errorf("shutdown cron scheduler: %w", err)
	}

	return nil
}

func (s *Service) dailyReportTask(ctx context.Context, log golog.Logger) {
	wrappedCtx, cancel := goctx.Wrap(ctx).WithTimeout(time.Duration(s.settings.TaskTimeout))
	defer cancel()

	now := time.Now()

	report, err := s.analyticsService.CreateBasicReport(wrappedCtx, log, now, now.AddDate(0, 0, 1))
	if err != nil {
		log.Errorf("failed to create daily basic report: %v", err)
		return
	}

	log.Debugf("created daily report %q at %v", report.File.FileName, report.CreatedAt)
}
