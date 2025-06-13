package job

import (
	"context"
	"encoding/json"
	"fmt"
	"go_shurtiner/internal/app/model"
	"go_shurtiner/internal/app/repository"
	"go_shurtiner/pkg/config"
	"go_shurtiner/pkg/report"
	"time"
)

type QueueService interface {
	NextJob(name string, startAfter time.Duration, payload any) error
}

type DataJobPayload struct {
	Data string `json:"data"`
}

type DataJob struct {
	ctx          context.Context
	repository   repository.ReportRepository
	queueService QueueService
	cfg          config.PrepareDataConfig
}

func NewDataJob(
	ctx context.Context,
	repository repository.ReportRepository,
	queueService QueueService,
	cfg config.PrepareDataConfig,
) *DataJob {
	return &DataJob{
		ctx:          ctx,
		repository:   repository,
		queueService: queueService,
		cfg:          cfg,
	}
}

func (j *DataJob) Process(job model.JobQueue) error {
	var payload DataJobPayload
	var err error
	err = json.Unmarshal(job.Params, &payload)
	if err != nil {
		return fmt.Errorf("cannot load job params: %v", err)
	}

	if "prepare.data" == job.Name {
		links := make([]model.Link, 0)
		links, err = j.repository.PrepareReportData(j.ctx)
		cntUsers := 0
		for _, link := range links {
			reportData := model.PreparedReport{}
			reportData.Timestamp = time.Now()
			reportData.Shortened = link.Shortened
			reportData.Source = link.Source
			for _, u := range link.User {
				reportData.UserEmail = u.Email
				reportData.UserFullName = fmt.Sprintf("%s %s %s", u.LastName, u.Name, u.MiddleName)
				err = j.repository.SaveReportData(j.ctx, &reportData)
				cntUsers++
			}
			if len(link.User) == 0 {
				err = j.repository.SaveReportData(j.ctx, &reportData)
			}
		}

		type Output struct {
			Count int `json:"count_links"`
			Users int `json:"count_users"`
		}

		output := Output{Count: len(links), Users: cntUsers}
		jsonData, _ := json.Marshal(output)
		jsonString := fmt.Sprintf("%s", string(jsonData))
		payload := DataJobPayload{Data: jsonString}
		json.Unmarshal([]byte(payload.Data), &output)

		err = j.queueService.NextJob(job.Name, j.cfg.TimeRange, output)
	}

	if "create.report" == job.Name {
		var prepadredReportData *[]model.PreparedReport
		now := time.Now()
		previousMonth := now.AddDate(0, -1, 0)
		prepadredReportData, err = j.repository.GetReportData(j.ctx, previousMonth)

		err = report.NewStatReport().GenerateReport(prepadredReportData)
		if err != nil {
			fmt.Println("ERROR")
			return err
		}

	}

	return err
}
