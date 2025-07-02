package repository

import (
	"context"
	"go_shurtiner/internal/app/model"
	"go_shurtiner/internal/database"
	_ "go_shurtiner/internal/queue"
	"go_shurtiner/pkg/config"
	"time"

	"gorm.io/gorm"
)

// ERR Unable to find 'QueueRepository' in any go files under this path dry-run=false version=v2.53.4
//
//go:generate mockery --name=QueueRepository --output=../../../test/mocks --outpkg=mocks --filename=mock_queue_repository.go
func NewQueueRepository(db *gorm.DB) *queueRepository {
	return &queueRepository{db: db}
}

type queueRepository struct {
	db  *gorm.DB
	cfg config.PrepareDataConfig
}

func (a *queueRepository) GetQueue(ctx context.Context) (model.JobQueue, error) {
	db := database.FromContext(ctx, a.db)
	// CTE для того, чтобы не использовать UPDATE напрямую
	sql := `
WITH cte_next_job AS
	(SELECT id
	FROM job_queue
	WHERE
			(
				launched_at IS NULL
				OR
				launched_at < @timeout
			)
			AND 
				completed_at IS NULL
			AND
			(
				scheduled_started_at IS NULL
				OR
				scheduled_started_at < @now
			)
		    AND deleted_at IS NULL
		ORDER BY 
			created_at ASC
		FOR UPDATE SKIP LOCKED
		LIMIT 1
	  )
	  UPDATE job_queue jq
	  SET launched_at = @now
	  FROM cte_next_job
	  WHERE cte_next_job.id = jq.id
	  RETURNING jq.id, name, params, scheduled_started_at, launched_at, completed_at, created_at
`

	params := map[string]any{
		"now":     time.Now().Unix(),
		"timeout": time.Now().Add(a.cfg.TimeRange).Unix(),
	}

	var item model.JobQueue
	tx := db.WithContext(ctx).Raw(sql, params).Scan(&item)
	if item.ID == 0 {
		return item, nil // не будем отдавать ошибку, если ничего не получили
	}

	return item, tx.Error
}

func (r *queueRepository) CompleteJob(ctx context.Context, job model.JobQueue) (model.JobQueue, error) {
	db := database.FromContext(ctx, r.db)
	if err := db.WithContext(ctx).Model(&job).Updates(map[string]interface{}{"completed_at": job.CompletedAt, "output": job.Output}).Error; err != nil {
		return job, err
	}

	return job, nil
}

func (r *queueRepository) CreateJob(ctx context.Context, job *model.JobQueue) error {
	db := database.FromContext(ctx, r.db)

	if err := db.WithContext(ctx).Create(job).Error; err != nil {
		return err
	}

	return nil
}

func (r *queueRepository) GetJob(ctx context.Context, id string) (model.JobQueue, error) {
	db := database.FromContext(ctx, r.db)

	var item model.JobQueue
	if err := db.WithContext(ctx).Take(&item, "id = ?", id).Error; err != nil {
		return item, err
	}

	return item, nil
}
