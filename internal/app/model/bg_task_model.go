package model

import (
	"gorm.io/datatypes"
)

type Task struct {
	ID          int            `gorm:"column:id;primaryKey;autoIncrement:false;not null" json:"id"`
	Name        string         `gorm:"column:name;type:varchar(256);not null" json:"name"`
	UserID      int            `gorm:"column:user_id;type:int" json:"user_id,omitempty"`
	Params      datatypes.JSON `gorm:"column:params;type:jsonb" json:"params,omitempty"`
	Output      datatypes.JSON `gorm:"column:output;type:json" json:"output,omitempty"`
	CreatedAt   int64          `gorm:"column:created_at;type:bigint;not null" json:"created_at"`
	ProcessAt   *int64         `gorm:"column:process_at;type:bigint" json:"process_at,omitempty"`
	ExpireAt    *int64         `gorm:"column:expire_at;type:bigint" json:"expire_at,omitempty"`
	CompletedAt *int64         `gorm:"column:completed_at;type:bigint" json:"completed_at,omitempty"`
}

func (u *Task) TableName() string {
	return "task"
}
