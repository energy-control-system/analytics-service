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

// CreateBasicReport godoc
// @Summary Create basic report
// @Description Generates a basic analytics report for the inclusive date period.
// @Tags reports
// @Produce json
// @Param periodStart path string true "Period start date in YYYY-MM-DD format"
// @Param periodEnd path string true "Period end date in YYYY-MM-DD format"
// @Success 200 {object} analytics.Report
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Router /reports/basic/{periodStart}/{periodEnd} [post]
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

// GetAllReports godoc
// @Summary List reports
// @Description Returns all generated analytics reports.
// @Tags reports
// @Produce json
// @Success 200 {array} analytics.Report
// @Failure 500 {object} gorouter.ErrorResponse
// @Router /reports [get]
func GetAllReports(s *analytics.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		response, err := s.GetAllReports(c.Ctx())
		if err != nil {
			return fmt.Errorf("failed to get reports: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}
