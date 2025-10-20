package inspection

import (
	"fmt"
	"net/http"

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

func (c *Client) GetInspectionByTaskID(ctx goctx.Context, taskID int) (Inspection, error) {
	var response Inspection
	status, err := c.client.DoJson(ctx, http.MethodGet, fmt.Sprintf("%s/inspections/task/%d", c.baseURL, taskID), nil, &response)
	if err != nil {
		return Inspection{}, fmt.Errorf("c.client.DoJson: %w", err)
	}

	if status != http.StatusOK {
		return Inspection{}, fmt.Errorf("unexpected status code: %d", status)
	}

	return response, nil
}
