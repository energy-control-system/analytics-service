package main

import (
	"analytics-service/api"
	"analytics-service/cluster/brigade"
	"analytics-service/cluster/file"
	"analytics-service/cluster/inspection"
	"analytics-service/cluster/subscriber"
	"analytics-service/config"
	dbanalytics "analytics-service/database/analytics"
	"analytics-service/service/analytics"
	"analytics-service/service/cron"
	"context"
	"fmt"
	"io/fs"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/jmoiron/sqlx"
	"github.com/sunshineOfficial/golib/db"
	"github.com/sunshineOfficial/golib/gohttp"
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
	postgres         *sqlx.DB
	clickhouse       *sqlx.DB
	clickhouseNative driver.Conn
	kafka            gokafka.Kafka
	taskConsumer     gokafka.Consumer

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
	postgresCtx, cancelPostgresCtx := context.WithTimeout(a.mainCtx, dbTimeout)
	defer cancelPostgresCtx()

	a.postgres, err = db.NewPgx(postgresCtx, a.settings.Databases.Postgres)
	if err != nil {
		return fmt.Errorf("init postgres: %w", err)
	}

	err = db.Migrate(fs, a.log, a.postgres, path+"/postgres", "postgres")
	if err != nil {
		return fmt.Errorf("migrate postgres: %w", err)
	}

	clickCtx, cancelClickCtx := context.WithTimeout(a.mainCtx, dbTimeout)
	defer cancelClickCtx()

	a.clickhouse, err = db.NewClickhouse(clickCtx, a.settings.Databases.Clickhouse.ConnectionString)
	if err != nil {
		return fmt.Errorf("init clickhouse: %w", err)
	}

	err = db.Migrate(fs, a.log, a.clickhouse, path+"/clickhouse", "clickhouse")
	if err != nil {
		return fmt.Errorf("migrate clickhouse: %w", err)
	}

	nativeClickCtx, cancelNativeClickCtx := context.WithTimeout(a.mainCtx, dbTimeout)
	defer cancelNativeClickCtx()

	a.clickhouseNative, err = db.NewNativeClickhouse(nativeClickCtx, db.NewClickhouseOptions(
		a.settings.Databases.Clickhouse.Host,
		a.settings.Databases.Clickhouse.Port,
		a.settings.Databases.Clickhouse.Database,
		a.settings.Databases.Clickhouse.Username,
		a.settings.Databases.Clickhouse.Password,
	))
	if err != nil {
		return fmt.Errorf("init native clickhouse: %w", err)
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

func (a *App) InitServices() error {
	analyticsRepository := dbanalytics.NewRepository(a.postgres, a.clickhouseNative)

	httpClient := gohttp.NewClient(gohttp.WithTimeout(1 * time.Minute))

	brigadeClient := brigade.NewClient(httpClient, a.settings.Cluster.BrigadeService)
	fileClient := file.NewClient(httpClient, a.settings.Cluster.FileService)
	inspectionClient := inspection.NewClient(httpClient, a.settings.Cluster.InspectionService)
	subscriberClient := subscriber.NewClient(httpClient, a.settings.Cluster.SubscriberService)

	a.analyticsService = analytics.NewService(
		analyticsRepository,
		inspectionClient,
		brigadeClient,
		subscriberClient,
		fileClient,
		a.settings.Templates,
	)

	a.cronService = cron.NewService(a.settings.Cron, a.analyticsService)

	return nil
}

func (a *App) InitServer() {
	sb := api.NewServerBuilder(a.mainCtx, a.log, a.settings)
	sb.AddDebug()
	sb.AddReports(a.analyticsService)

	a.server = sb.Build()
}

func (a *App) Start() error {
	a.server.Start()
	a.taskConsumer.Subscribe(a.analyticsService.SubscriberOnTaskEvent(a.mainCtx, a.log.WithTags("taskSubscriber")))

	if err := a.cronService.Start(a.mainCtx, a.log.WithTags("cronService")); err != nil {
		return fmt.Errorf("start cron: %w", err)
	}

	return nil
}

func (a *App) Stop(ctx context.Context) {
	err := a.cronService.Stop()
	if err != nil {
		a.log.Errorf("failed to stop cron: %v", err)
	}

	consumerCtx, cancelConsumerCtx := context.WithTimeout(ctx, dbTimeout)
	defer cancelConsumerCtx()

	if err = a.taskConsumer.Close(consumerCtx); err != nil {
		a.log.Errorf("failed to close task consumer: %v", err)
	}

	a.server.Stop()

	if err = a.clickhouseNative.Close(); err != nil {
		a.log.Errorf("failed to close clickhouse native connection: %v", err)
	}

	if err = a.clickhouse.Close(); err != nil {
		a.log.Errorf("failed to close clickhouse connection: %v", err)
	}

	if err = a.postgres.Close(); err != nil {
		a.log.Errorf("failed to close postgres connection: %v", err)
	}
}
