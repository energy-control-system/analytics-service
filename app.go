package main

import (
	"analytics-service/config"
	"analytics-service/service/analytics"
	"analytics-service/service/cron"
	"context"
	"fmt"
	"io/fs"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sunshineOfficial/golib/db"
	"github.com/sunshineOfficial/golib/gohttp/goserver"
	"github.com/sunshineOfficial/golib/gokafka"
	"github.com/sunshineOfficial/golib/golog"
)

const (
	serviceName = "analytics-service"
	dbTimeout   = 15 * time.Second
)

type App struct {
	/* main */
	mainCtx  context.Context
	log      golog.Logger
	settings config.Settings

	/* http */
	server goserver.Server

	/* db */
	clickhouse   *sqlx.DB
	kafka        gokafka.Kafka
	taskConsumer gokafka.Consumer

	/* services */
	analyticsService *analytics.Service
	cronService      *cron.Service
}

func NewApp(mainCtx context.Context, log golog.Logger, settings config.Settings) *App {
	return &App{
		mainCtx:  mainCtx,
		log:      log,
		settings: settings,
	}
}

func (a *App) InitDatabases(fs fs.FS, path string) (err error) {
	clickCtx, cancelClickCtx := context.WithTimeout(a.mainCtx, dbTimeout)
	defer cancelClickCtx()

	a.clickhouse, err = db.NewClickhouse(clickCtx, a.settings.Databases.Clickhouse)
	if err != nil {
		return fmt.Errorf("init clickhouse: %w", err)
	}

	err = db.Migrate(fs, a.log, a.clickhouse, path, "clickhouse")
	if err != nil {
		return fmt.Errorf("migrate clickhouse: %w", err)
	}

	a.kafka = gokafka.NewKafka(a.settings.Databases.Kafka.Brokers)

	a.taskConsumer, err = a.kafka.Consumer(a.log.WithTags("taskConsumer"), func() (context.Context, context.CancelFunc) {
		return context.WithCancel(a.mainCtx)
	}, gokafka.WithTopic(a.settings.Databases.Kafka.Topics.Tasks), gokafka.WithConsumerGroup(serviceName))
	if err != nil {
		return fmt.Errorf("init task consumer: %w", err)
	}

	return nil
}
