package models

import "time"

type Product struct {
	ID          uint      `gorm:"column:id_prod;primaryKey" json:"id_prod"`
	Name        string    `gorm:"column:name_prod;not null" json:"name_prod"`
	Description string    `gorm:"column:description_prod" json:"description_prod"`
	Price       float64   `gorm:"column:price_prod;not null" json:"price_prod"`
	Stock       int       `gorm:"column:stock_prod;not null" json:"stock_prod"`
	CategoryID  uint      `gorm:"column:category_prod;not null" json:"category_prod"`
	Image       string    `gorm:"column:image_prod" json:"image_prod"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`

	Category Category `gorm:"foreignKey:CategoryID;references:ID"`
}

func (Product) TableName() string {
	return "products"
}
