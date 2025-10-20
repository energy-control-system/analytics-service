package handler

import (
	"analytics-service/service/analytics"
	"fmt"
	"net/http"
	"time"

	"github.com/sunshineOfficial/golib/gohttp/gorouter"
)

type periodVars struct {
	PeriodStart string `path:"periodStart"`
	PeriodEnd   string `path:"periodEnd"`
}

func CreateBasicReport(s *analytics.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var vars periodVars
		if err := c.Vars(&vars); err != nil {
			return fmt.Errorf("failed to read period: %w", err)
		}

		periodStart, err := time.Parse(time.DateOnly, vars.PeriodStart)
		if err != nil {
			return fmt.Errorf("failed to parse periodStart: %w", err)
		}

		periodEnd, err := time.Parse(time.DateOnly, vars.PeriodEnd)
		if err != nil {
			return fmt.Errorf("failed to parse periodEnd: %w", err)
		}

		response, err := s.CreateBasicReport(c.Ctx(), c.Log().WithTags("basicReport"), periodStart, periodEnd)
		if err != nil {
			return fmt.Errorf("failed to create report: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}

func GetAllReports(s *analytics.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		response, err := s.GetAllReports(c.Ctx())
		if err != nil {
			return fmt.Errorf("failed to get reports: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}
