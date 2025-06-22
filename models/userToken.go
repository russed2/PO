package models

import "time"

type UserToken struct {
	ID           uint      `gorm:"column:id_token;primaryKey" json:"id_token"`
	UserID       uint      `gorm:"column:user_token;not null" json:"user_token"`
	RefreshToken string    `gorm:"column:refresh_token;not null" json:"refresh_token"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	ExpiresAt    time.Time `gorm:"column:expires_at" json:"expires_at"`

	User User `gorm:"foreignKey:UserID;references:ID"`
}

func (UserToken) TableName() string {
	return "user_tokens"
}
