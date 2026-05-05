package analytics

import (
	service "analytics-service/service/analytics"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestMapFinishedTaskToDBIncludesInspectionDeviceReadings(t *testing.T) {
	createdAt := time.Date(2026, time.January, 10, 9, 30, 0, 0, time.UTC)
	task := service.FinishedTask{
		Inspection: service.Inspection{
			Devices: []service.InspectedDevice{
				{
					ID:          55,
					DeviceID:    77,
					Value:       decimal.RequireFromString("1542.40"),
					Consumption: decimal.RequireFromString("318.50"),
					CreatedAt:   createdAt,
				},
			},
		},
	}

	dbTask := MapFinishedTaskToDB(task)

	if len(dbTask.InspectedDevices) != 1 {
		t.Fatalf("expected one inspected device, got %d", len(dbTask.InspectedDevices))
	}

	device := dbTask.InspectedDevices[0]
	if device.ID != 55 {
		t.Fatalf("expected inspected device id 55, got %d", device.ID)
	}
	if device.DeviceID != 77 {
		t.Fatalf("expected device id 77, got %d", device.DeviceID)
	}
	if !device.Value.Equal(decimal.RequireFromString("1542.40")) {
		t.Fatalf("unexpected value: %s", device.Value)
	}
	if !device.Consumption.Equal(decimal.RequireFromString("318.50")) {
		t.Fatalf("unexpected consumption: %s", device.Consumption)
	}
	if !device.CreatedAt.Equal(createdAt) {
		t.Fatalf("unexpected created_at: %s", device.CreatedAt)
	}
}
