package prepare_data

import (
	"context"
	"go_shurtiner/internal/app/model"
	"go_shurtiner/internal/app/repository"
)

type PrepareData interface {
	Prepare(ctx context.Context) ([]model.PreparedReport, error)
}

type prepareData struct {
	payload *any
	//db      *gorm.DB
	prepareRepository repository.ReportRepository
}

func NewPrepareData(payload *any, prepareRepository repository.ReportRepository) PrepareData {
	return &prepareData{
		payload: payload,
		//db:      db1,
		prepareRepository: prepareRepository,
	}
}

func (p prepareData) Prepare(ctx context.Context) ([]model.PreparedReport, error) {
	//db := database.FromContext(ctx)
	//INSERT INTO PREPARED_report (id, created_at, updated_at, timestamp, source, shortened, user_email, user_fullname)
	//VALUES (1, '2025-03-15T12:00:00Z', now(), '2025-03-15T12:00:00Z', 'https://example.com', 'exmpl', 'user@example.com', 'John Doe');

	report := make([]model.PreparedReport, 0)

	//err := db.Find(&report).Error
	reportData := model.PreparedReport{}
	_, err := p.prepareRepository.PrepareReportData(ctx, &reportData)
	//err := db.Find(&report, "created_at >= 2025-03-15T12:00:00Z").Error
	if err != nil {
		return report, err
	}

	return report, err
}
