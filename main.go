package main

import (
	"analytics-service/config"
	"context"
	"os"

	"github.com/shopspring/decimal"
	"github.com/sunshineOfficial/golib/golog"
)

func main() {
	configureDecimal()

	log := golog.NewLogger(serviceName)
	log.Debug("service up")

	settings, err := config.Get(log)
	if err != nil {
		log.Errorf("failed to get config: %v", err)
		return
	}

	mainCtx, cancelMainCtx := context.WithCancel(context.Background())
	defer cancelMainCtx()

	app := NewApp(mainCtx, log, settings)

	if err = app.InitDatabases(os.DirFS("./"), "database/migrations/clickhouse"); err != nil {
		log.Errorf("failed to init databases: %v", err)
		return
	}

	log.Debug("service down")
}

func configureDecimal() {
	decimal.DivisionPrecision = 2
	decimal.MarshalJSONWithoutQuotes = true
}
