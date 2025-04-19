package model

import (
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type JobQueue struct {
	gorm.Model
	Name               string      `gorm:"column:name;type:VARCHAR;size:256;" json:"name"`
	Params             []byte      `gorm:"column:params;type:JSON;" json:"params"`
	Output             null.String `gorm:"column:output;type:JSON;" json:"output"`
	ScheduledStartedAt int64       `gorm:"column:scheduled_started_at;type:INT8;" json:"scheduled_started_at"`
	LaunchedAt         null.Int    `gorm:"column:launched_at;type:INT8;" json:"launched_at"`
	ExpireAt           null.Int    `gorm:"column:expire_at;type:INT8;" json:"expire_at"`
	CompletedAt        null.Int    `gorm:"column:completed_at;type:INT8;" json:"completed_at"`
}

func (q *JobQueue) TableName() string {
	return "job_queue"
}
