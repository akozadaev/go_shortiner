package repository

import (
	"context"
	"go_shurtiner/internal/app/model"
	"go_shurtiner/internal/database"
	_ "go_shurtiner/internal/queue"
	"time"

	"gorm.io/gorm"
)

const JobTimeout = time.Minute * 5

func NewQueueRepository(db *gorm.DB) *queueRepository {
	return &queueRepository{db: db}
}

type queueRepository struct {
	db *gorm.DB
}

func (a *queueRepository) GetQueue(ctx context.Context) (model.JobQueue, error) {
	db := database.FromContext(ctx, a.db)

	sql := `
	UPDATE job_queue
	SET processed_at = @now
	WHERE id IN (
		SELECT id
		FROM job_queue
		WHERE
			(
				processed_at IS NULL
				OR
				processed_at < @timeout
			)
			AND 
				completed_at IS NULL
			AND
			(
				scheduled_started_at IS NULL
				OR
				scheduled_started_at < @now
			)
		ORDER BY 
			created_at ASC
		FOR UPDATE SKIP LOCKED
		LIMIT 1
	  )
	  RETURNING id, name, params, scheduled_started_at, processed_at, completed_at, created_at
	`

	params := map[string]any{
		"now":     time.Now().Unix(),
		"timeout": time.Now().Add(-JobTimeout).Unix(),
	}

	var item model.JobQueue
	tx := db.WithContext(ctx).Raw(sql, params).Scan(&item)
	if tx.Error == gorm.ErrRecordNotFound {
		return item, nil
	}

	return item, tx.Error
}

func (r *queueRepository) CompleteJob(ctx context.Context, job model.JobQueue) (model.JobQueue, error) {
	db := database.FromContext(ctx, r.db)

	if err := db.WithContext(ctx).Model(&job).Update("completed_at", job.CompletedAt).Error; err != nil {
		return job, err
	}

	return job, nil
}

func (r *queueRepository) CreateJob(ctx context.Context, job model.JobQueue) error {
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
