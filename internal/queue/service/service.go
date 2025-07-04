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
	CreateJob(ctx context.Context, job *model.JobQueue) error
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

func (s *QueueService) NextJob(name string, startAfter time.Duration, payload json.RawMessage) error {
	queueJob := model.JobQueue{
		Name:               name,
		Params:             payload,
		ScheduledStartedAt: time.Now().Unix(),
		LaunchedAt:         null.IntFrom(time.Now().Add(startAfter).Unix()),
	}

	err := s.repository.CreateJob(context.Background(), &queueJob)
	if err != nil {
		return err
	}

	return nil
}

func (s *QueueService) CreateJob(ctx context.Context, name string, startAfter time.Duration, payload json.RawMessage) (model.JobQueue, error) {
	queueJob := model.JobQueue{
		Name:               name,
		ScheduledStartedAt: time.Now().Add(startAfter).Unix(),
	}

	queueJob.Params = payload

	err := s.repository.CreateJob(ctx, &queueJob)
	if err != nil {
		return queueJob, err
	}

	return queueJob, nil
}

func (s *QueueService) GetJob(ctx context.Context, id string) (model.JobQueue, error) {
	return s.repository.GetJob(ctx, id)
}
