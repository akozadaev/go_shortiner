package service

import (
	"context"
	"go_shurtiner/internal/app/model"
	"time"
)

type ReportRepository interface {
	PrepareReportData(ctx context.Context) ([]model.Link, error)
	SaveReportData(ctx context.Context, reportData *model.PreparedReport) error
	GetReportData(ctx context.Context, startDate time.Time) (*[]model.PreparedReport, error)
	CreateReport(ctx context.Context, startDate time.Time) (*[]model.PreparedReport, error)
}

type ReportService struct {
	repository ReportRepository
}

// NewReportService фабрика сервиса отчетов
func NewReportService(repo ReportRepository) *ReportService {
	return &ReportService{
		repository: repo,
	}
}
