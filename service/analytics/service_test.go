package analytics

import (
	"context"
	"io"
	"testing"
	"time"

	"analytics-service/cluster/file"

	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/pagination"
)

type repositoryMock struct {
	reports []Report
}

func (m repositoryMock) AddFinishedTask(context.Context, FinishedTask) error {
	return nil
}

func (m repositoryMock) GetFinishedTasksByPeriod(context.Context, time.Time, time.Time) ([]FinishedTask, error) {
	return nil, nil
}

func (m repositoryMock) AddReport(context.Context, Report) (Report, error) {
	return Report{}, nil
}

func (m repositoryMock) GetAllReports(context.Context, pagination.Pagination) ([]Report, error) {
	return m.reports, nil
}

type fileServiceMock struct {
	filesByID  map[int]file.File
	gotIDs     []int
	gotHeaders file.ForwardedHeaders
}

func (m *fileServiceMock) Upload(goctx.Context, string, io.Reader, file.ForwardedHeaders) (file.File, error) {
	return file.File{}, nil
}

func (m *fileServiceMock) GetFilesByIDs(_ goctx.Context, ids []int, headers file.ForwardedHeaders) ([]file.File, error) {
	m.gotIDs = append([]int(nil), ids...)
	m.gotHeaders = headers

	files := make([]file.File, 0, len(ids))
	for _, id := range ids {
		if f, ok := m.filesByID[id]; ok {
			files = append(files, f)
		}
	}

	return files, nil
}

func TestGetAllReportsPassesForwardedHeadersToFileService(t *testing.T) {
	fileService := &fileServiceMock{
		filesByID: map[int]file.File{
			70: {ID: 70, URL: "https://api.example.test/storage/report.xlsx"},
		},
	}

	service := &Service{
		repository: repositoryMock{
			reports: []Report{
				{
					ID:    42,
					Files: []file.File{{ID: 70}},
				},
			},
		},
		fileService: fileService,
	}

	headers := file.ForwardedHeaders{Host: "api.example.test", Proto: "https"}
	got, err := service.GetAllReports(goctx.Wrap(context.Background()), pagination.Pagination{}, headers)
	if err != nil {
		t.Fatalf("GetAllReports returned error: %v", err)
	}

	if got[0].Files[0].URL != "https://api.example.test/storage/report.xlsx" {
		t.Fatalf("got[0].Files[0].URL = %q, want %q", got[0].Files[0].URL, "https://api.example.test/storage/report.xlsx")
	}
	if len(fileService.gotIDs) != 1 || fileService.gotIDs[0] != 70 {
		t.Fatalf("fileService.gotIDs = %+v, want [70]", fileService.gotIDs)
	}
	if fileService.gotHeaders != headers {
		t.Fatalf("fileService.gotHeaders = %+v, want %+v", fileService.gotHeaders, headers)
	}
}
