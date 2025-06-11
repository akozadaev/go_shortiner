package queue

import "go_shurtiner/internal/app/model"

type QueueJob interface {
	Process(job model.JobQueue) error
}
