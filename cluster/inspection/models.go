package inspection

import "time"

type Status int

const (
	StatusUnknown Status = iota
	StatusInWork
	StatusDone
)

type Type int

const (
	TypeUnknown Type = iota
	TypeLimitation
	TypeResumption
	TypeVerification
	TypeUnauthorizedConnection
)

type Resolution int

const (
	ResolutionUnknown Resolution = iota
	ResolutionLimited
	ResolutionStopped
	ResolutionResumed
)

type MethodBy int

const (
	MethodByUnknown MethodBy = iota
	MethodByConsumer
	MethodByInspector
)

type ReasonType int

const (
	ReasonTypeUnknown ReasonType = iota
	ReasonTypeNotIntroduced
	ReasonTypeConsumerLimited
	ReasonTypeInspectorLimited
	ReasonTypeResumed
)

type Inspection struct {
	ID                      int         `json:"ID"`
	TaskID                  int         `json:"TaskID"`
	Status                  Status      `json:"Status"`
	Type                    *Type       `json:"Type"`
	Resolution              *Resolution `json:"Resolution"`
	LimitReason             *string     `json:"LimitReason"`
	Method                  *string     `json:"Method"`
	MethodBy                *MethodBy   `json:"MethodBy"`
	ReasonType              *ReasonType `json:"ReasonType"`
	ReasonDescription       *string     `json:"ReasonDescription"`
	IsRestrictionChecked    *bool       `json:"IsRestrictionChecked"`
	IsViolationDetected     *bool       `json:"IsViolationDetected"`
	IsExpenseAvailable      *bool       `json:"IsExpenseAvailable"`
	ViolationDescription    *string     `json:"ViolationDescription"`
	IsUnauthorizedConsumers *bool       `json:"IsUnauthorizedConsumers"`
	UnauthorizedDescription *string     `json:"UnauthorizedDescription"`
	UnauthorizedExplanation *string     `json:"UnauthorizedExplanation"`
	InspectAt               *time.Time  `json:"InspectAt"`
	EnergyActionAt          *time.Time  `json:"EnergyActionAt"`
	CreatedAt               time.Time   `json:"CreatedAt"`
	UpdatedAt               time.Time   `json:"UpdatedAt"`
}
