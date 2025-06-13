package repository

import (
	"context"
	"go_shurtiner/internal/app/model"
	database2 "go_shurtiner/internal/database"
	"go_shurtiner/pkg/logging"
	"gorm.io/gorm"
	"time"
)

type ReportRepository interface {
	PrepareReportData(ctx context.Context) ([]model.Link, error)
	SaveReportData(ctx context.Context, reportData *model.PreparedReport) error
	GetReportData(ctx context.Context, startDate time.Time) (*[]model.PreparedReport, error)
	CreateReport(ctx context.Context, startDate time.Time) (*[]model.PreparedReport, error)
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &prepareReportRepository{db: db}
}

type prepareReportRepository struct {
	db *gorm.DB
}

func (p prepareReportRepository) PrepareReportData(ctx context.Context) ([]model.Link, error) {
	logger := logging.FromContext(ctx)
	db := database2.FromContext(ctx, p.db)
	links := make([]model.Link, 0)
	if err := db.Preload("User").Find(&links).Error; err != nil {
		logger.Errorw("failed to prepare report data", "err", err)
		return nil, err
	}

	return links, nil
}

func (p prepareReportRepository) SaveReportData(ctx context.Context, reportData *model.PreparedReport) error {
	logger := logging.FromContext(ctx)
	db := database2.FromContext(ctx, p.db)

	if err := db.Create(reportData).Error; err != nil {
		logger.Errorw("failed to save report data", "err", err)
		if database2.IsKeyConflictErr(err) {
			return database2.ErrKeyConflict
		}
		return err
	}
	return nil
}

func (p prepareReportRepository) GetReportData(ctx context.Context, startDate time.Time) (*[]model.PreparedReport, error) {
	logger := logging.FromContext(ctx)
	db := database2.FromContext(ctx, p.db)
	var err error
	var report []model.PreparedReport
	if err = db.WithContext(ctx).Find(&report /*, fmt.Sprintf("created_at <= timestamptz(%d)", startDate)*/).Error; err != nil {
		//if err = db.WithContext(ctx).Find(&report,"created_at >= 2025-03-15T12:00:00Z").Error; err != nil {
		logger.Errorw("failed to get report data", "err", err)
	}

	return &report, err
}
func (p prepareReportRepository) CreateReport(ctx context.Context, startDate time.Time) (*[]model.PreparedReport, error) {
	logger := logging.FromContext(ctx)
	db := database2.FromContext(ctx, p.db)
	var err error
	var report []model.PreparedReport
	if err = db.WithContext(ctx).Find(&report /*, fmt.Sprintf("created_at <= timestamptz(%d)", startDate)*/).Error; err != nil {
		//if err = db.WithContext(ctx).Find(&report,"created_at >= 2025-03-15T12:00:00Z").Error; err != nil {
		logger.Errorw("failed to get report data", "err", err)
	}

	return &report, err
}
