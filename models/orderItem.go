package models

import "time"

type OrderItem struct {
	ID        uint      `gorm:"column:id_orditem;primaryKey" json:"id_orditem"`
	OrderID   uint      `gorm:"column:order_orditem;not null" json:"order_orditem"`
	ProductID uint      `gorm:"column:product_orditem;not null" json:"product_orditem"`
	Quantity  int       `gorm:"column:quantity_orditem;not null" json:"quantity_orditem"`
	Price     float64   `gorm:"column:price_orditem;not null" json:"price_orditem"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`

	Order   Order   `gorm:"foreignKey:OrderID;references:ID"`
	Product Product `gorm:"foreignKey:ProductID;references:ID"`
}

func (OrderItem) TableName() string {
	return "order_items"
}
