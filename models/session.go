package models

import "time"

type Session struct {
	ID        uint      `gorm:"column:id_session;primaryKey" json:"id_session"`
	UserID    uint      `gorm:"column:user_session;not null" json:"user_session"`
	Token     string    `gorm:"column:session_token;not null" json:"session_token"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	ExpiresAt time.Time `gorm:"column:expires_at" json:"expires_at"`

	User User `gorm:"foreignKey:UserID;references:ID"`
}

func (Session) TableName() string {
	return "sessions"
}
