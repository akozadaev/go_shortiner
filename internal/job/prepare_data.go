package job

import (
	"context"
	"encoding/json"
	"fmt"
	"go_shurtiner/internal/app/model"
	"go_shurtiner/internal/app/repository"
	"go_shurtiner/pkg/config"
	"time"
)

type QueueService interface {
	NextJob(name string, startAfter time.Duration, payload any) error
}

type PrepareDataJobPayload struct {
	Data string `json:"data"`
}

type PrepareDataJob struct {
	ctx          context.Context
	repository   repository.PrepareReportRepository
	queueService QueueService
	cfg          config.PrepareDataConfig
}

func NewPrepareDataJob(
	ctx context.Context,
	repository repository.PrepareReportRepository,
	queueService QueueService,
	cfg config.PrepareDataConfig,
) *PrepareDataJob {
	return &PrepareDataJob{
		ctx:          ctx,
		repository:   repository,
		queueService: queueService,
		cfg:          cfg,
	}
}

func (j *PrepareDataJob) Process(job model.JobQueue) error {
	var payload PrepareDataJobPayload
	err := json.Unmarshal(job.Params, &payload)
	if err != nil {
		return fmt.Errorf("cannot load job params: %v", err)
	}

	links, err := j.repository.PrepareReportData(j.ctx)
	for _, link := range links {
		reportData := model.PreparedReport{}
		reportData.Timestamp = time.Now()
		reportData.Shortened = link.Shortened
		reportData.Source = link.Source
		for _, u := range link.User {
			reportData.UserEmail = u.Email
			reportData.UserFullName = fmt.Sprintf("%s %s %s", u.LastName, u.Name, u.MiddleName)
			err = j.repository.SaveReportData(j.ctx, &reportData)
		}
		if len(link.User) == 0 {
			err = j.repository.SaveReportData(j.ctx, &reportData)
		}
	}

	err = j.queueService.NextJob(j.Name(), j.cfg.TimeRange, nil)

	return nil
}

func (j *PrepareDataJob) Name() string {
	return "prepare.data"
}
