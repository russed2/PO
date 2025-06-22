package models

import "time"

type Review struct {
	ID        uint      `gorm:"column:id_rev;primaryKey" json:"id_rev"`
	UserID    uint      `gorm:"column:user_rev;not null" json:"user_rev"`
	ProductID uint      `gorm:"column:product_rev;not null" json:"product_rev"`
	Rating    int       `gorm:"column:rating_rev;not null" json:"rating_rev"`
	Comment   string    `gorm:"column:comment_rev" json:"comment_rev"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`

	User    User    `gorm:"foreignKey:UserID;references:ID"`
	Product Product `gorm:"foreignKey:ProductID;references:ID"`
}

func (Review) TableName() string {
	return "reviews"
}
