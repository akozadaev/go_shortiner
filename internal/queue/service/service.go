package service

import (
	"context"
	"encoding/json"
	"go_shurtiner/internal/app/model"
	"gopkg.in/guregu/null.v4"

	"time"
)

type QueueRepository interface {
	GetQueue(ctx context.Context) (model.JobQueue, error)
	GetJob(ctx context.Context, id string) (model.JobQueue, error)
	CreateJob(ctx context.Context, job model.JobQueue) error
	CompleteJob(ctx context.Context, job model.JobQueue) (model.JobQueue, error)
}

type QueueService struct {
	repository QueueRepository
}

func NewQueueService(repo QueueRepository) *QueueService {
	return &QueueService{
		repository: repo,
	}
}

func (s *QueueService) GetQueue(ctx context.Context) (model.JobQueue, error) {
	return s.repository.GetQueue(ctx)
}

func (s *QueueService) CompleteJob(ctx context.Context, job model.JobQueue) (model.JobQueue, error) {
	job.CompletedAt = null.IntFrom(time.Now().Unix())
	return s.repository.CompleteJob(ctx, job)
}

func (s *QueueService) NextJob(name string, startAfter time.Duration, payload any) error {
	params, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	queueJob := model.JobQueue{
		Name:   name,
		Params: params,
	}

	err = s.repository.CreateJob(context.Background(), queueJob)
	if err != nil {
		return err
	}

	return nil
}

func (s *QueueService) CreateJob(ctx context.Context, name string, startAfter time.Duration, payload any) (model.JobQueue, error) {
	queueJob := model.JobQueue{
		Name:               name,
		ScheduledStartedAt: null.IntFrom(time.Now().Add(startAfter).Unix()),
	}

	params, err := json.Marshal(payload)
	if err != nil {
		return queueJob, err
	}

	queueJob.Params = params

	err = s.repository.CreateJob(ctx, queueJob)
	if err != nil {
		return queueJob, err
	}

	return queueJob, nil
}

func (s *QueueService) GetJob(ctx context.Context, id string) (model.JobQueue, error) {
	return s.repository.GetJob(ctx, id)
}
