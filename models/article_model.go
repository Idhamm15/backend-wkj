package models

import (
	"time"
)

type Article struct {
	ID               	uint            `gorm:"primaryKey" json:"id"`
	Name             	string          `form:"name" json:"name"`
	Description      	string          `form:"description" json:"description"`
	ImageURL         	string          `json:"image_url"`
	ArticleCategoryID 	uint           	`json:"article_category_id"`
	ArticleCategory  	ArticleCategory `gorm:"foreignKey:ArticleCategoryID" json:"article_category"`
	CreatedAt        	time.Time       `json:"created_at"`
	UpdatedAt        	time.Time       `json:"updated_at"`
}

type ArticleCategory struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `form:"name" json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
