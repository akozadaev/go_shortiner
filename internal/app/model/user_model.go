package model

import (
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name       string  `gorm:"column:name;type:varchar(128);not null;" json:"name"`
	LastName   string  `gorm:"column:lastname;type:varchar(128);not null;" json:"lastname"`
	MiddleName string  `gorm:"column:middlename;type:varchar(128);default null;" json:"middlename"`
	Password   string  `gorm:"column:password;type:varchar(72);not null;" json:"password"`
	Email      string  `gorm:"column:email;type:varchar(72);not null;unique;" json:"email"`
	Links      []*Link `gorm:"many2many:user_link;"`
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

type AuthUserData struct {
	Id       string      `db:"id"`
	Username string      `db:"username"`
	Password null.String `db:"password"`
}

type AuthUserRequest struct {
	Email    string
	Password string
}

func (u *User) TableName() string {
	return "user"
}
