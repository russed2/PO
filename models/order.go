package models

import "time"

type Order struct {
	ID        uint      `gorm:"column:id_order;primaryKey" json:"id_order"`
	UserID    uint      `gorm:"column:user_order;not null" json:"user_order"`
	Status    string    `gorm:"column:status_order;not null" json:"status_order"`
	Sum       float64   `gorm:"column:sum_order;not null" json:"sum_order"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`

	User User `gorm:"foreignKey:UserID;references:ID"`
}

func (Order) TableName() string {
	return "orders"
}
