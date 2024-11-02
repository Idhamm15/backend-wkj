package models

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

type Product struct {
    ID                  uint              `gorm:"primaryKey" json:"id"`
    Name                string            `form:"name" json:"name"`	
    LatinName           string            `form:"latin_name" json:"latin_name"`	
    Synonym             string            `form:"synonym" json:"synonym"`	
    Familia             string            `form:"familia" json:"familia"`	
    PartUsed            string            `form:"part_used" json:"part_used"`	
    MethodOfReproduction string           `form:"method_of_reproduction" json:"method_of_reproduction"`	
    HarvestAge          string            `form:"harvest_age" json:"harvest_age"`	
    Morphology          string            `form:"morphology" json:"morphology"`	
    AreaName            string            `form:"area_name" json:"area_name"`	
    Efficacy            string            `form:"efficacy" json:"efficacy"`	
    Utilization         json.RawMessage   `gorm:"type:json" json:"utilization"`
    Composition         json.RawMessage   `gorm:"type:json" json:"composition"`
    ImageURL            string            `json:"image_url"`
    ResearchResults     string            `form:"research_results" json:"research_results"`	
    Description         string            `form:"description" json:"description"`
    Price               float32           `form:"price" json:"price"`
    UnitType            string            `form:"unit_type" json:"unit_type"`
    ProductCategoryID   uint              `json:"product_category_id"`
    ProductCategory     ProductCategory    `gorm:"foreignKey:ProductCategoryID" json:"product_category"`
    CreatedAt           time.Time         `json:"created_at"`
    UpdatedAt           time.Time         `json:"updated_at"`
}

type ProductResponse struct {
	ID                  uint        `json:"id"`
	Name                string      `json:"name"`
	LatinName           string      `json:"latin_name"`
	Synonym             string      `json:"synonym"`
	Familia             string      `json:"familia"`
	PartUsed            string      `json:"part_used"`
	MethodOfReproduction string     `json:"method_of_reproduction"`
	HarvestAge          string      `json:"harvest_age"`
	Morphology          string      `json:"morphology"`
	AreaName            string      `json:"area_name"`
	Efficacy            string      `json:"efficacy"`
	Utilization         gin.H       `json:"utilization"`
	Composition         gin.H       `json:"composition"`
	ImageURL            string      `json:"image_url"`
	ResearchResults     string      `json:"research_results"`
	Description         string      `json:"description"`
	Price               float32     `json:"price"`
	UnitType            string      `json:"unit_type"`
	ProductCategoryID   uint        `json:"product_category_id"`
	ProductCategory     interface{} `json:"product_category"`
	CreatedAt           time.Time   `json:"created_at"`
	UpdatedAt           time.Time   `json:"updated_at"`
}

type ProductCategory struct {
	ID        	uint      `gorm:"primaryKey" json:"id"`
	Name      	string    `form:"name" json:"name"`
	CreatedAt 	time.Time `json:"created_at"`
	UpdatedAt 	time.Time `json:"updated_at"`
}
