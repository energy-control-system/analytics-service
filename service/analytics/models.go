package analytics

import (
	"analytics-service/cluster/file"
	"analytics-service/cluster/inspection"
	"analytics-service/cluster/subscriber"
	"time"
)

type ReportType int

const (
	ReportTypeUnknown ReportType = iota
	ReportTypeBasic
)

type Report struct {
	ID          int         `json:"ID"`
	Type        ReportType  `json:"Type"`
	Files       []file.File `json:"Files"`
	PeriodStart time.Time   `json:"PeriodStart"`
	PeriodEnd   time.Time   `json:"PeriodEnd"`
	CreatedAt   time.Time   `json:"CreatedAt"`
}

type FinishedTask struct {
	TaskID      int        `json:"TaskID"`
	Comment     *string    `json:"Comment"`
	PlanVisitAt *time.Time `json:"PlanVisitAt"`
	StartedAt   time.Time  `json:"StartedAt"`
	FinishedAt  time.Time  `json:"FinishedAt"`
	Inspection  Inspection `json:"Inspection"`
	Brigade     Brigade    `json:"Brigade"`
	Object      Object     `json:"Object"`
}

type Inspection struct {
	ID                      int                   `json:"ID"`
	Type                    inspection.Type       `json:"Type"`
	Resolution              inspection.Resolution `json:"Resolution"`
	LimitReason             *string               `json:"LimitReason"`
	Method                  string                `json:"Method"`
	MethodBy                inspection.MethodBy   `json:"MethodBy"`
	ReasonType              inspection.ReasonType `json:"ReasonType"`
	ReasonDescription       *string               `json:"ReasonDescription"`
	IsRestrictionChecked    bool                  `json:"IsRestrictionChecked"`
	IsViolationDetected     bool                  `json:"IsViolationDetected"`
	IsExpenseAvailable      bool                  `json:"IsExpenseAvailable"`
	ViolationDescription    *string               `json:"ViolationDescription"`
	IsUnauthorizedConsumers bool                  `json:"IsUnauthorizedConsumers"`
	UnauthorizedDescription *string               `json:"UnauthorizedDescription"`
	UnauthorizedExplanation *string               `json:"UnauthorizedExplanation"`
	InspectAt               time.Time             `json:"InspectAt"`
	EnergyActionAt          time.Time             `json:"EnergyActionAt"`
}

type Brigade struct {
	ID         int         `json:"ID"`
	Inspectors []Inspector `json:"Inspectors"`
}

type Inspector struct {
	ID          int       `json:"ID"`
	Surname     string    `json:"Surname"`
	Name        string    `json:"Name"`
	Patronymic  string    `json:"Patronymic"`
	PhoneNumber string    `json:"PhoneNumber"`
	Email       string    `json:"Email"`
	AssignedAt  time.Time `json:"AssignedAt"`
}

type Object struct {
	ID            int        `json:"ID"`
	Address       string     `json:"Address"`
	HaveAutomaton bool       `json:"HaveAutomaton"`
	Subscriber    Subscriber `json:"Subscriber"`
}

type Subscriber struct {
	ID            int               `json:"ID"`
	AccountNumber string            `json:"AccountNumber"`
	Surname       string            `json:"Surname"`
	Name          string            `json:"Name"`
	Patronymic    string            `json:"Patronymic"`
	PhoneNumber   string            `json:"PhoneNumber"`
	Email         string            `json:"Email"`
	INN           string            `json:"INN"`
	BirthDate     string            `json:"BirthDate"`
	Status        subscriber.Status `json:"Status"`
}
