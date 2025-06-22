package models

import "time"

type Category struct {
	ID        uint      `gorm:"column:id_cat;primaryKey" json:"id_cat"`
	Name      string    `gorm:"column:name_cat;not null" json:"name_cat"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (Category) TableName() string {
	return "categories"
}
