package model

import (
	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	Source    string  `gorm:"column:source;type:varchar(2048);" json:"source"`
	Shortened string  `gorm:"column:shortened;type:varchar(256);" json:"shortened"`
	User      []*User `gorm:"many2many:user_link;"`
}

type CreateLink struct {
	Source string `json:"url"`
}

func (l *Link) TableName() string {
	return "link"
}
