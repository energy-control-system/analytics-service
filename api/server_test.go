package api

import (
	"analytics-service/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sunshineOfficial/golib/golog"
)

func TestReportCreationRequiresAuthorization(t *testing.T) {
	builder := NewServerBuilder(t.Context(), golog.NewLogger("test"), config.Settings{
		Port: 80,
	})
	builder.AddReports(nil)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/reports/basic/2026-01-01/2026-01-02", nil)

	builder.router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusUnauthorized)
	}
}
