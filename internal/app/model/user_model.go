package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name       string `gorm:"column:name;type:varchar(128);" json:"name"`
	LastName   string `gorm:"column:lastname;type:varchar(128);" json:"lastname"`
	MiddleName string `gorm:"column:middlename;type:varchar(128);" json:"middlename"`
	Password   string `gorm:"column:password;type:varchar(72);" json:"password"`
	Email      string `gorm:"column:email;type:varchar(72);" json:"email"`
}

type UserApi struct {
	gorm.Model
	Name       string `gorm:"column:name;type:varchar(128);" json:"name"`
	LastName   string `gorm:"column:lastname;type:varchar(128);" json:"lastname"`
	MiddleName string `gorm:"column:middlename;type:varchar(128);" json:"middlename"`
	Email      string `gorm:"column:email;type:varchar(72);" json:"email"`
}

type CreateUser struct {
	Source     string `json:"url"`
	Name       string `json:"name"`
	LastName   string `json:"lastname"`
	MiddleName string `json:"middlename"`
	Password   string `json:"password"`
	Email      string `json:"email"`
}

func (u *User) TableName() string {
	return "user"
}
