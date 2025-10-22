package config

import (
	"fmt"
	"os"

	"github.com/sunshineOfficial/golib/config"
	"github.com/sunshineOfficial/golib/golog"
)

func Get(log golog.Logger) (Settings, error) {
	var settings Settings

	env := config.GetEnvironmentName()
	if len(env) == 0 {
		env = "local"
	}

	log.Debugf("Environment: %v", env)

	if err := config.Parse(&settings); err != nil {
		return Settings{}, err
	}

	settings.Databases.Postgres = fmt.Sprintf(settings.Databases.Postgres, os.Getenv("POSTGRES_PASSWORD"))

	settings.Databases.Clickhouse.Username = os.Getenv("CLICKHOUSE_USER")
	settings.Databases.Clickhouse.Password = os.Getenv("CLICKHOUSE_PASSWORD")
	settings.Databases.Clickhouse.ConnectionString = fmt.Sprintf(
		settings.Databases.Clickhouse.ConnectionString,
		settings.Databases.Clickhouse.Username,
		settings.Databases.Clickhouse.Password,
		settings.Databases.Clickhouse.Host,
		settings.Databases.Clickhouse.Port,
		settings.Databases.Clickhouse.Database,
	)

	return settings, nil
}
