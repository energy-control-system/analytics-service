package api

import (
	"analytics-service/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sunshineOfficial/golib/golog"
)

func TestReportRoutesAllowUnauthenticatedRequests(t *testing.T) {
	builder := NewServerBuilder(t.Context(), golog.NewLogger("test"), config.Settings{
		Port: 80,
	})
	builder.AddReports(nil)

	routes := []struct {
		method string
		path   string
	}{
		{method: http.MethodPost, path: "/reports/basic/2026-01-01/2026-01-02"},
		{method: http.MethodGet, path: "/reports"},
	}

	for _, route := range routes {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			response := httptest.NewRecorder()
			request := httptest.NewRequest(route.method, route.path, nil)

			builder.router.ServeHTTP(response, request)

			if response.Code == http.StatusUnauthorized {
				t.Fatalf("status = %d, route must be open without authorization", response.Code)
			}
		})
	}
}
