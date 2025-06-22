package models

import "time"

type User struct {
	ID        uint      `gorm:"column:id_user;primaryKey" json:"id_user"`
	Email     string    `gorm:"column:email_user;unique;not null" json:"email_user"`
	Password  string    `gorm:"column:passwordhs_user;not null" json:"passwordhs_user"`
	Name      string    `gorm:"column:name_user;not null" json:"name_user"`
	Admin     bool      `gorm:"column:admin_user;not null" json:"admin_user"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (User) TableName() string {
	return "users"
}
