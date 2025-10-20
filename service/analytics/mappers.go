package analytics

import (
	"analytics-service/cluster/brigade"
	"analytics-service/cluster/inspection"
	"analytics-service/cluster/subscriber"
	"analytics-service/cluster/task"
)

func MapToFinishedTask(t task.Task, ins inspection.Inspection, brig brigade.Brigade, obj subscriber.ObjectExtended) FinishedTask {
	return FinishedTask{
		TaskID:      t.ID,
		Comment:     t.Comment,
		PlanVisitAt: t.PlanVisitAt,
		StartedAt:   *t.StartedAt,
		FinishedAt:  *t.FinishedAt,
		Inspection:  MapInspectionToDomain(ins),
		Brigade:     MapBrigadeToDomain(brig),
		Object:      MapObjectToDomain(obj),
	}
}

func MapInspectionToDomain(ins inspection.Inspection) Inspection {
	return Inspection{
		ID:                      ins.ID,
		Type:                    *ins.Type,
		Resolution:              *ins.Resolution,
		LimitReason:             ins.LimitReason,
		Method:                  *ins.Method,
		MethodBy:                *ins.MethodBy,
		ReasonType:              *ins.ReasonType,
		ReasonDescription:       ins.ReasonDescription,
		IsRestrictionChecked:    *ins.IsRestrictionChecked,
		IsViolationDetected:     *ins.IsViolationDetected,
		IsExpenseAvailable:      *ins.IsExpenseAvailable,
		ViolationDescription:    ins.ViolationDescription,
		IsUnauthorizedConsumers: *ins.IsUnauthorizedConsumers,
		UnauthorizedDescription: ins.UnauthorizedDescription,
		UnauthorizedExplanation: ins.UnauthorizedExplanation,
		InspectAt:               *ins.InspectAt,
		EnergyActionAt:          *ins.EnergyActionAt,
	}
}

func MapBrigadeToDomain(brig brigade.Brigade) Brigade {
	return Brigade{
		ID:         brig.ID,
		Inspectors: MapInspectorSliceToDomain(brig.Inspectors),
	}
}

func MapInspectorToDomain(i brigade.Inspector) Inspector {
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

func MapInspectorSliceToDomain(inspectors []brigade.Inspector) []Inspector {
	result := make([]Inspector, 0, len(inspectors))
	for _, i := range inspectors {
		result = append(result, MapInspectorToDomain(i))
	}

	return result
}

func MapObjectToDomain(obj subscriber.ObjectExtended) Object {
	return Object{
		ID:            obj.ID,
		Address:       obj.Address,
		HaveAutomaton: obj.HaveAutomaton,
		Subscriber:    MapSubscriberToDomain(obj.Subscriber),
	}
}

func MapSubscriberToDomain(s subscriber.Subscriber) Subscriber {
	return Subscriber{
		ID:            s.ID,
		AccountNumber: s.AccountNumber,
		Surname:       s.Surname,
		Name:          s.Name,
		Patronymic:    s.Patronymic,
		PhoneNumber:   s.PhoneNumber,
		Email:         s.Email,
		INN:           s.INN,
		BirthDate:     s.BirthDate,
		Status:        s.Status,
	}
}
