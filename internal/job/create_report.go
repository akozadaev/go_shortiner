package job

//
//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"go_shurtiner/internal/app/model"
//	"go_shurtiner/internal/app/repository"
//	"go_shurtiner/pkg/config"
//	"time"
//)
//
////type QueueService interface {
////	NextJob(name string, startAfter time.Duration, payload any) error
////}
//
//type CreateReportJobJobPayload struct {
//	Data []any `json:"data"`
//}
//
//type CreateReportJob struct {
//	ctx          context.Context
//	repository   repository.PrepareReportRepository
//	queueService QueueService
//	cfg          config.PrepareDataConfig
//}
//
//func NewCreateReportJob(
//	ctx context.Context,
//	repository repository.PrepareReportRepository,
//	queueService QueueService,
//	cfg config.PrepareDataConfig,
//) *CreateReportJob {
//	return &CreateReportJob{
//		ctx:          ctx,
//		repository:   repository,
//		queueService: queueService,
//		cfg:          cfg,
//	}
//}
//
//func (j *PrepareDataJob) Process(job model.JobQueue) error {
//	var payload PrepareDataJobPayload
//	fmt.Println("AAAAAAAAAAAA")
//	err := json.Unmarshal(job.Params, &payload)
//	if err != nil {
//		return fmt.Errorf("cannot load job params: %v", err)
//	}
//	reportData, err := j.repository.GetReportData(j.ctx, time.Now())
//	if err != nil {
//		return fmt.Errorf("cannot get report data: %v", err)
//	}
//	fmt.Println("reportData")
//	fmt.Println(reportData)
//	fmt.Println(reportData)
//
//	//j.prepareData.Prepare(j.ctx)
//	//INSERT INTO PREPARED_report (id, created_at, updated_at, timestamp, source, shortened, user_email, user_fullname)
//	//VALUES (1, '2025-03-15T12:00:00Z', now(), '2025-03-15T12:00:00Z', 'https://example.com', 'exmpl', 'user@example.com', 'John Doe');
//
//	report := make([]model.PreparedReport, 0)
//	//db := database2.FromContext(j.ctx, j.db)
//	//err = db.Find(&report, "created_at >= 2025-03-15T12:00:00Z").Error
//
//	for _, link := range report {
//		//result = append(result, model.UserApi{}
//		fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++=")
//		fmt.Println(link)
//	}
//
//	/*	if err := db.WithContext(j.ctx).Create(slot).Error; err != nil {
//		logger.Errorw("dont save prepared data", "err", err)
//		if database.IsKeyConflictErr(err) {
//			return database.ErrKeyConflict
//		}
//		return err
//	}*/
//	return nil
//}
//
//func (j *CreateReportJob) Name() string {
//	return "create.report"
//}
