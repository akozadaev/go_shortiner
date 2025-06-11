package queue

import (
	"context"
	"errors"
	"fmt"
	"go_shurtiner/internal/app/model"
	"log"
	"runtime"
	"sync"
	"time"
)

type QueueService interface {
	GetQueue(ctx context.Context) (model.JobQueue, error)
	CompleteJob(ctx context.Context, job model.JobQueue) (model.JobQueue, error)
}

type Queue struct {
	ctx          context.Context
	cancel       context.CancelFunc
	queueService QueueService
	mu           sync.RWMutex
	jobs         map[string]QueueJob
}

func NewQueue(
	queueService QueueService,
) *Queue {
	ctx, cancel := context.WithCancel(context.Background())

	return &Queue{
		ctx:          ctx,
		cancel:       cancel,
		queueService: queueService,
		jobs:         make(map[string]QueueJob),
	}
}

func (q *Queue) Run(ctx context.Context) error {
	workerPoolSize := runtime.NumCPU()

	errCh := make(chan error)
	defer close(errCh)

	go func(errCh <-chan error) {
		for err := range errCh {
			log.Printf("[ERR] %v\n", err)
		}
	}(errCh)

	var wg sync.WaitGroup
	for i := 0; i < workerPoolSize; i++ {
		wg.Add(1)
		go q.worker(errCh, &wg)
	}

	cleanTicker := time.NewTicker(time.Hour)
	defer cleanTicker.Stop()

	wg.Wait()

	return nil
}

func (q *Queue) Shutdown() error {
	q.cancel()

	return nil
}

// worker функция обработки задач
func (q *Queue) worker(errCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	stop := false
	for !stop {
		select {
		case <-q.ctx.Done():
			stop = true
		default:
			task, err := q.queueService.GetQueue(q.ctx)
			if err != nil {
				errCh <- fmt.Errorf("cannot load job queue: %v", err)
				break
			}
			fmt.Println("task.Name: " + task.Name)
			if err = q.processTask(&task); err != nil {
				errCh <- fmt.Errorf("cannot process task: %v", err)
			}

			time.Sleep(1 * time.Second)
		}
	}
}

func (q *Queue) processTask(task *model.JobQueue) error {
	job, err := q.resolve(task.Name)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Start job %s...\n", task.ID)

	if err := job.Process(*task); err != nil {
		return fmt.Errorf("cannot process job %v: %v", task.ID, err)
	}

	_, err = q.queueService.CompleteJob(q.ctx, *task)
	if err != nil {
		return fmt.Errorf("cannot complete job %s: %v", task.ID, err)
	}

	log.Printf("[INFO] Finished job %v\n", task.ID)
	return nil
}

func (q *Queue) AddJob(jobName string, jobs ...QueueJob) {
	q.mu.Lock()
	for _, job := range jobs {
		q.jobs[jobName] = job
	}
	q.mu.Unlock()
}

func (q *Queue) resolve(name string) (QueueJob, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if job, ok := q.jobs[name]; ok {
		return job, nil
	}

	return nil, errors.New("Cannot find job with name " + name)
}
