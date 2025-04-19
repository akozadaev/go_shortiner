package model

import (
	"time"
)

type PreparedReport struct {
	ID           uint `gorm:"primarykey;autoIncrement"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Timestamp    time.Time //*int64 `gorm:"column:timestamp;type:bigint" json:"timestamp,omitempty"`
	Source       string    `gorm:"column:source;type:varchar(2048);" json:"source"`
	Shortened    string    `gorm:"column:shortened;type:varchar(256);" json:"shortened"`
	UserEmail    string    `gorm:"column:user_email;type:varchar(72)" json:"user_email,omitempty"`
	UserFullName string    `gorm:"column:user_fullname;type:varchar(384)" json:"user_fullname,omitempty"`
}

func (r *PreparedReport) TableName() string {
	return "prepared_report"
}
