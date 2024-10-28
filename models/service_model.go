package models

import (
	"time"
)

type Service struct {
	ID          		uint      		`gorm:"primaryKey" json:"id"`
	Name        		string    		`form:"name" json:"name"`
	Description 		string    		`form:"description" json:"description"`
	Price 				string    		`form:"price" json:"price"`
	ImageURL    		string    		`json:"image_url"`
	PhoneNumber 		string    		`json:"phone_number"`
	ServiceCategoryID 	uint           	`json:"service_category_id"`
	ServiceCategory  	ServiceCategory `gorm:"foreignKey:ServiceCategoryID" json:"service_category"`
	CreatedAt   		time.Time 		`json:"created_at"`
	UpdatedAt   		time.Time 		`json:"updated_at"`
}

type ServiceCategory struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `form:"name" json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
