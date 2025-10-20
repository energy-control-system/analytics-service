package analytics

import "time"

type Report struct {
	ID          int       `db:"id"`
	Type        int       `db:"type"`
	FileID      int       `db:"file_id"`
	FileName    string    `db:"file_name"`
	FileSize    int64     `db:"file_size"`
	FileBucket  string    `db:"file_bucket"`
	FileURL     string    `db:"file_url"`
	PeriodStart time.Time `db:"period_start"`
	PeriodEnd   time.Time `db:"period_end"`
	CreatedAt   time.Time `db:"created_at"`
}

type FinishedTask struct {
	TaskID                            int         `db:"task_id"`
	Comment                           *string     `db:"comment"`
	PlanVisitAt                       *time.Time  `db:"plan_visit_at"`
	StartedAt                         time.Time   `db:"started_at"`
	FinishedAt                        time.Time   `db:"finished_at"`
	InspectionID                      int         `db:"inspection_id"`
	InspectionType                    int         `db:"inspection_type"`
	InspectionResolution              int         `db:"inspection_resolution"`
	InspectionLimitReason             *string     `db:"inspection_limit_reason"`
	InspectionMethod                  string      `db:"inspection_method"`
	InspectionMethodBy                int         `db:"inspection_method_by"`
	InspectionReasonType              int         `db:"inspection_reason_type"`
	InspectionReasonDescription       *string     `db:"inspection_reason_description"`
	InspectionIsRestrictionChecked    bool        `db:"inspection_is_restriction_checked"`
	InspectionIsViolationDetected     bool        `db:"inspection_is_violation_detected"`
	InspectionIsExpenseAvailable      bool        `db:"inspection_is_expense_available"`
	InspectionViolationDescription    *string     `db:"inspection_violation_description"`
	InspectionIsUnauthorizedConsumers bool        `db:"inspection_is_unauthorized_consumers"`
	InspectionUnauthorizedDescription *string     `db:"inspection_unauthorized_description"`
	InspectionUnauthorizedExplanation *string     `db:"inspection_unauthorized_explanation"`
	InspectionInspectAt               time.Time   `db:"inspection_inspect_at"`
	InspectionEnergyActionAt          time.Time   `db:"inspection_energy_action_at"`
	BrigadeID                         int         `db:"brigade_id"`
	BrigadeInspectors                 []Inspector `db:"brigade_inspectors"`
	ObjectID                          int         `db:"object_id"`
	ObjectAddress                     string      `db:"object_address"`
	ObjectHaveAutomaton               bool        `db:"object_have_automaton"`
	SubscriberID                      int         `db:"subscriber_id"`
	SubscriberAccountNumber           string      `db:"subscriber_account_number"`
	SubscriberSurname                 string      `db:"subscriber_surname"`
	SubscriberName                    string      `db:"subscriber_name"`
	SubscriberPatronymic              string      `db:"subscriber_patronymic"`
	SubscriberPhoneNumber             string      `db:"subscriber_phone_number"`
	SubscriberEmail                   string      `db:"subscriber_email"`
	SubscriberINN                     string      `db:"subscriber_inn"`
	SubscriberBirthDate               string      `db:"subscriber_birth_date"`
	SubscriberStatus                  int         `db:"subscriber_status"`
}

type Inspector struct {
	ID          int       `db:"id"`
	Surname     string    `db:"surname"`
	Name        string    `db:"name"`
	Patronymic  string    `db:"patronymic"`
	PhoneNumber string    `db:"phone_number"`
	Email       string    `db:"email"`
	AssignedAt  time.Time `db:"assigned_at"`
}
