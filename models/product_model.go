package models

import (
	"time"
)

type Product struct {
	ID               	uint            `gorm:"primaryKey" json:"id"`
	Name             	string          `form:"name" json:"name"`
	Description      	string          `form:"description" json:"description"`
	Price            	string          `form:"price" json:"price"`
	UnitType         	string          `form:"unit_type" json:"unit_type"`
	ImageURL         	string          `json:"image_url"`
	ProductCategoryID 	uint           	`json:"product_category_id"`
	ProductCategory  	ProductCategory `gorm:"foreignKey:ProductCategoryID" json:"product_category"`
	CreatedAt        	time.Time       `json:"created_at"`
	UpdatedAt        	time.Time       `json:"updated_at"`
}

type ProductCategory struct {
	ID        	uint      `gorm:"primaryKey" json:"id"`
	Name      	string    `form:"name" json:"name"`
	CreatedAt 	time.Time `json:"created_at"`
	UpdatedAt 	time.Time `json:"updated_at"`
}
