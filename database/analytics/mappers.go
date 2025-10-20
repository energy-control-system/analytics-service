package analytics

import (
	"analytics-service/cluster/file"
	"analytics-service/cluster/inspection"
	"analytics-service/cluster/subscriber"
	"analytics-service/service/analytics"
)

func MapFinishedTaskToDB(t analytics.FinishedTask) FinishedTask {
	return FinishedTask{
		TaskID:                            t.TaskID,
		Comment:                           t.Comment,
		PlanVisitAt:                       t.PlanVisitAt,
		StartedAt:                         t.StartedAt,
		FinishedAt:                        t.FinishedAt,
		InspectionID:                      t.Inspection.ID,
		InspectionType:                    int(t.Inspection.Type),
		InspectionResolution:              int(t.Inspection.Resolution),
		InspectionLimitReason:             t.Inspection.LimitReason,
		InspectionMethod:                  t.Inspection.Method,
		InspectionMethodBy:                int(t.Inspection.MethodBy),
		InspectionReasonType:              int(t.Inspection.ReasonType),
		InspectionReasonDescription:       t.Inspection.ReasonDescription,
		InspectionIsRestrictionChecked:    t.Inspection.IsRestrictionChecked,
		InspectionIsViolationDetected:     t.Inspection.IsViolationDetected,
		InspectionIsExpenseAvailable:      t.Inspection.IsExpenseAvailable,
		InspectionViolationDescription:    t.Inspection.ViolationDescription,
		InspectionIsUnauthorizedConsumers: t.Inspection.IsUnauthorizedConsumers,
		InspectionUnauthorizedDescription: t.Inspection.UnauthorizedDescription,
		InspectionUnauthorizedExplanation: t.Inspection.UnauthorizedExplanation,
		InspectionInspectAt:               t.Inspection.InspectAt,
		InspectionEnergyActionAt:          t.Inspection.EnergyActionAt,
		BrigadeID:                         t.Brigade.ID,
		BrigadeInspectors:                 MapInspectorSliceToDB(t.Brigade.Inspectors),
		ObjectID:                          t.Object.ID,
		ObjectAddress:                     t.Object.Address,
		ObjectHaveAutomaton:               t.Object.HaveAutomaton,
		SubscriberID:                      t.Object.Subscriber.ID,
		SubscriberAccountNumber:           t.Object.Subscriber.AccountNumber,
		SubscriberSurname:                 t.Object.Subscriber.Surname,
		SubscriberName:                    t.Object.Subscriber.Name,
		SubscriberPatronymic:              t.Object.Subscriber.Patronymic,
		SubscriberPhoneNumber:             t.Object.Subscriber.PhoneNumber,
		SubscriberEmail:                   t.Object.Subscriber.Email,
		SubscriberINN:                     t.Object.Subscriber.INN,
		SubscriberBirthDate:               t.Object.Subscriber.BirthDate,
		SubscriberStatus:                  int(t.Object.Subscriber.Status),
	}
}

func MapInspectorToDB(i analytics.Inspector) Inspector {
	return Inspector{
		ID:          i.ID,
		Surname:     i.Surname,
		Name:        i.Name,
		Patronymic:  i.Patronymic,
		PhoneNumber: i.PhoneNumber,
		Email:       i.Email,
		AssignedAt:  i.AssignedAt,
	}
}

func MapInspectorSliceToDB(inspectors []analytics.Inspector) []Inspector {
	result := make([]Inspector, 0, len(inspectors))
	for _, inspector := range inspectors {
		result = append(result, MapInspectorToDB(inspector))
	}

	return result
}

func MapFinishedTaskFromDB(t FinishedTask) analytics.FinishedTask {
	return analytics.FinishedTask{
		TaskID:      t.TaskID,
		Comment:     t.Comment,
		PlanVisitAt: t.PlanVisitAt,
		StartedAt:   t.StartedAt,
		FinishedAt:  t.FinishedAt,
		Inspection: analytics.Inspection{
			ID:                      t.InspectionID,
			Type:                    inspection.Type(t.InspectionType),
			Resolution:              inspection.Resolution(t.InspectionResolution),
			LimitReason:             t.InspectionLimitReason,
			Method:                  t.InspectionMethod,
			MethodBy:                inspection.MethodBy(t.InspectionMethodBy),
			ReasonType:              inspection.ReasonType(t.InspectionReasonType),
			ReasonDescription:       t.InspectionReasonDescription,
			IsRestrictionChecked:    t.InspectionIsRestrictionChecked,
			IsViolationDetected:     t.InspectionIsViolationDetected,
			IsExpenseAvailable:      t.InspectionIsExpenseAvailable,
			ViolationDescription:    t.InspectionViolationDescription,
			IsUnauthorizedConsumers: t.InspectionIsUnauthorizedConsumers,
			UnauthorizedDescription: t.InspectionUnauthorizedDescription,
			UnauthorizedExplanation: t.InspectionUnauthorizedExplanation,
			InspectAt:               t.InspectionInspectAt,
			EnergyActionAt:          t.InspectionEnergyActionAt,
		},
		Brigade: analytics.Brigade{
			ID:         t.BrigadeID,
			Inspectors: MapInspectorSliceFromDB(t.BrigadeInspectors),
		},
		Object: analytics.Object{
			ID:            t.ObjectID,
			Address:       t.ObjectAddress,
			HaveAutomaton: t.ObjectHaveAutomaton,
			Subscriber: analytics.Subscriber{
				ID:            t.SubscriberID,
				AccountNumber: t.SubscriberAccountNumber,
				Surname:       t.SubscriberSurname,
				Name:          t.SubscriberName,
				Patronymic:    t.SubscriberPatronymic,
				PhoneNumber:   t.SubscriberPhoneNumber,
				Email:         t.SubscriberEmail,
				INN:           t.SubscriberINN,
				BirthDate:     t.SubscriberBirthDate,
				Status:        subscriber.Status(t.SubscriberStatus),
			},
		},
	}
}

func MapFinishedTaskSliceFromDB(tasks []FinishedTask) []analytics.FinishedTask {
	result := make([]analytics.FinishedTask, 0, len(tasks))
	for _, t := range tasks {
		result = append(result, MapFinishedTaskFromDB(t))
	}

	return result
}

func MapInspectorFromDB(i Inspector) analytics.Inspector {
	return analytics.Inspector{
		ID:          i.ID,
		Surname:     i.Surname,
		Name:        i.Name,
		Patronymic:  i.Patronymic,
		PhoneNumber: i.PhoneNumber,
		Email:       i.Email,
		AssignedAt:  i.AssignedAt,
	}
}

func MapInspectorSliceFromDB(inspectors []Inspector) []analytics.Inspector {
	result := make([]analytics.Inspector, 0, len(inspectors))
	for _, i := range inspectors {
		result = append(result, MapInspectorFromDB(i))
	}

	return result
}

func MapReportToDB(r analytics.Report) Report {
	return Report{
		ID:          r.ID,
		Type:        int(r.Type),
		PeriodStart: r.PeriodStart,
		PeriodEnd:   r.PeriodEnd,
		CreatedAt:   r.CreatedAt,
	}
}

func MapReportFromDB(r Report) analytics.Report {
	return analytics.Report{
		ID:          r.ID,
		Type:        analytics.ReportType(r.Type),
		PeriodStart: r.PeriodStart,
		PeriodEnd:   r.PeriodEnd,
		CreatedAt:   r.CreatedAt,
	}
}

func MapAttachmentToDB(f file.File, reportID int) Attachment {
	return Attachment{
		ReportID: reportID,
		FileID:   f.ID,
	}
}

func MapAttachmentSliceToDB(files []file.File, reportID int) []Attachment {
	result := make([]Attachment, 0, len(files))
	for _, f := range files {
		result = append(result, MapAttachmentToDB(f, reportID))
	}

	return result
}

func MapAttachmentFromDB(a Attachment) file.File {
	return file.File{
		ID: a.FileID,
	}
}

func MapAttachmentSliceFromDB(attachments []Attachment) []file.File {
	result := make([]file.File, 0, len(attachments))
	for _, a := range attachments {
		result = append(result, MapAttachmentFromDB(a))
	}

	return result
}
