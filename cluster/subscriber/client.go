package subscriber

import (
	"time"

	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/gohttp"
)

type Client struct {
	client  gohttp.Client
	baseURL string
}

func NewClient(client gohttp.Client, baseURL string) *Client {
	return &Client{
		client:  client,
		baseURL: baseURL,
	}
}

func (c *Client) GetObjectExtendedByID(ctx goctx.Context, id int) (ObjectExtended, error) {
	// TODO: убрать мок
	return ObjectExtended{
		ID:            id,
		Address:       "г. Пенза, ул. Ворошилова, д. 9",
		HaveAutomaton: true,
		CreatedAt:     time.Date(2025, 10, 10, 9, 13, 14, 0, time.UTC),
		UpdatedAt:     time.Date(2025, 10, 10, 9, 16, 31, 0, time.UTC),
		Subscriber: Subscriber{
			ID:            1,
			AccountNumber: "asf123",
			Surname:       "Шаипов",
			Name:          "Камиль",
			Patronymic:    "Гяряевич",
			PhoneNumber:   "89234567856",
			Email:         "test@gmail.com",
			INN:           "1234567890",
			BirthDate:     time.Date(2004, 9, 1, 0, 0, 0, 0, time.UTC),
			Status:        StatusActive,
			CreatedAt:     time.Date(2025, 10, 10, 9, 13, 14, 0, time.UTC),
			UpdatedAt:     time.Date(2025, 10, 10, 9, 16, 31, 0, time.UTC),
		},
		Devices: []DeviceExtended{{
			ID:               43,
			ObjectID:         123,
			Type:             "Счетчик",
			Number:           "34525",
			PlaceType:        DevicePlaceFlat,
			PlaceDescription: "",
			CreatedAt:        time.Date(2025, 10, 10, 9, 13, 14, 0, time.UTC),
			UpdatedAt:        time.Date(2025, 10, 10, 9, 16, 31, 0, time.UTC),
			Seals: []Seal{{
				ID:        69,
				DeviceID:  43,
				Number:    "39303",
				Place:     "Кухня",
				CreatedAt: time.Date(2025, 10, 10, 9, 13, 14, 0, time.UTC),
				UpdatedAt: time.Date(2025, 10, 10, 9, 16, 31, 0, time.UTC),
			}},
		}},
	}, nil
}
