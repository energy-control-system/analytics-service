package config

import "github.com/sunshineOfficial/golib/gotime"

type Settings struct {
	Port      int       `json:"port"`
	Databases Databases `json:"databases"`
	Cluster   Cluster   `json:"cluster"`
	Templates Templates `json:"templates"`
	Cron      Cron      `json:"cron"`
}

type Databases struct {
	Postgres   string     `json:"postgres"`
	Clickhouse Clickhouse `json:"clickhouse"`
	Kafka      Kafka      `json:"kafka"`
}

type Clickhouse struct {
	ConnectionString string `json:"connectionString"`
	Host             string `json:"host"`
	Port             int    `json:"port"`
	Database         string `json:"database"`
	Username         string `json:"username"`
	Password         string `json:"password"`
}

type Kafka struct {
	Brokers []string `json:"brokers"`
	Topics  Topics   `json:"topics"`
}

type Topics struct {
	Tasks string `json:"tasks"`
}

type Cluster struct {
	BrigadeService    string `json:"brigadeService"`
	FileService       string `json:"fileService"`
	InspectionService string `json:"inspectionService"`
	SubscriberService string `json:"subscriberService"`
	TaskService       string `json:"taskService"`
}

type Templates struct {
	BasicReport string `json:"basicReport"`
}

type Cron struct {
	DailyReportTime string          `json:"dailyReportTime"`
	TaskTimeout     gotime.Duration `json:"taskTimeout"`
}
