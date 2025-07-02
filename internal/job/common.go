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
	NextJob(name string, startAfter time.Duration, payload json.RawMessage) error
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
	var err error

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
				if err != nil {
					fmt.Println("Ошибка:", err)
				}
				cntUsers++
			}
			if len(link.User) == 0 {
				err = j.repository.SaveReportData(j.ctx, &reportData)
				if err != nil {
					fmt.Println("Ошибка:", err)
				}
			}
		}

		type Output struct {
			Count int `json:"count_links"`
			Users int `json:"count_users"`
		}

		output := Output{Count: len(links), Users: cntUsers}
		jsonData, _ := json.Marshal(output)

		err = j.queueService.NextJob(job.Name, j.cfg.TimeRange, jsonData)
	}

	if "create.report" == job.Name {
		var prepadredReportData *[]model.PreparedReport
		now := time.Now()
		previousMonth := now.AddDate(0, -1, 0)
		var result map[string]json.RawMessage

		err = json.Unmarshal([]byte(job.Params), &result)
		if err != nil {
			fmt.Println("Ошибка:", err)
		}

		fmt.Println("resultresultresultresultresult")
		fmt.Println(result)
		fmt.Println("resultresultresultresultresult")
		// Для параметров предусмотрено job.Params, но в рамках данной реализации просто за месяц
		prepadredReportData, _ = j.repository.GetReportData(j.ctx, previousMonth)

		err = report.NewStatReport().GenerateReport(prepadredReportData)
		if err != nil {
			fmt.Println("ERROR")
			return err
		}

	}

	return err
}
