package analytics

import "time"

type Report struct {
	ID          int       `db:"id"`
	Type        int       `db:"type"`
	PeriodStart time.Time `db:"period_start"`
	PeriodEnd   time.Time `db:"period_end"`
	CreatedAt   time.Time `db:"created_at"`
}

type Attachment struct {
	ID        int       `db:"id"`
	ReportID  int       `db:"report_id"`
	FileID    int       `db:"file_id"`
	CreatedAt time.Time `db:"created_at"`
}

type FinishedTask struct {
	TaskID                            int64       `ch:"task_id"`
	Comment                           *string     `ch:"comment"`
	PlanVisitAt                       *time.Time  `ch:"plan_visit_at"`
	StartedAt                         time.Time   `ch:"started_at"`
	FinishedAt                        time.Time   `ch:"finished_at"`
	InspectionID                      int64       `ch:"inspection_id"`
	InspectionType                    int8        `ch:"inspection_type"`
	InspectionResolution              int8        `ch:"inspection_resolution"`
	InspectionLimitReason             *string     `ch:"inspection_limit_reason"`
	InspectionMethod                  string      `ch:"inspection_method"`
	InspectionMethodBy                int8        `ch:"inspection_method_by"`
	InspectionReasonType              int8        `ch:"inspection_reason_type"`
	InspectionReasonDescription       *string     `ch:"inspection_reason_description"`
	InspectionIsRestrictionChecked    bool        `ch:"inspection_is_restriction_checked"`
	InspectionIsViolationDetected     bool        `ch:"inspection_is_violation_detected"`
	InspectionIsExpenseAvailable      bool        `ch:"inspection_is_expense_available"`
	InspectionViolationDescription    *string     `ch:"inspection_violation_description"`
	InspectionIsUnauthorizedConsumers bool        `ch:"inspection_is_unauthorized_consumers"`
	InspectionUnauthorizedDescription *string     `ch:"inspection_unauthorized_description"`
	InspectionUnauthorizedExplanation *string     `ch:"inspection_unauthorized_explanation"`
	InspectionInspectAt               time.Time   `ch:"inspection_inspect_at"`
	InspectionEnergyActionAt          time.Time   `ch:"inspection_energy_action_at"`
	BrigadeID                         int64       `ch:"brigade_id"`
	BrigadeInspectors                 []Inspector `ch:"brigade_inspectors"`
	ObjectID                          int64       `ch:"object_id"`
	ObjectAddress                     string      `ch:"object_address"`
	ObjectHaveAutomaton               bool        `ch:"object_have_automaton"`
	SubscriberID                      int64       `ch:"subscriber_id"`
	SubscriberAccountNumber           string      `ch:"subscriber_account_number"`
	SubscriberSurname                 string      `ch:"subscriber_surname"`
	SubscriberName                    string      `ch:"subscriber_name"`
	SubscriberPatronymic              string      `ch:"subscriber_patronymic"`
	SubscriberPhoneNumber             string      `ch:"subscriber_phone_number"`
	SubscriberEmail                   string      `ch:"subscriber_email"`
	SubscriberINN                     string      `ch:"subscriber_inn"`
	SubscriberBirthDate               time.Time   `ch:"subscriber_birth_date"`
	SubscriberStatus                  int8        `ch:"subscriber_status"`
}

type Inspector struct {
	ID          int64     `ch:"id"`
	Surname     string    `ch:"surname"`
	Name        string    `ch:"name"`
	Patronymic  string    `ch:"patronymic"`
	PhoneNumber string    `ch:"phone_number"`
	Email       string    `ch:"email"`
	AssignedAt  time.Time `ch:"assigned_at"`
}
