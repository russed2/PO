package models

import "time"

type CartItem struct {
	ID        uint      `gorm:"column:id_cart;primaryKey" json:"id_cart"`
	UserID    uint      `gorm:"column:user_cart;not null" json:"user_cart"`
	ProductID uint      `gorm:"column:product_cart;not null" json:"product_cart"`
	Quantity  int       `gorm:"column:quantity_cart;not null" json:"quantity_cart"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`

	User    User    `gorm:"foreignKey:UserID;references:ID"`
	Product Product `gorm:"foreignKey:ProductID;references:ID"`
}

func (CartItem) TableName() string {
	return "cart_items"
}
