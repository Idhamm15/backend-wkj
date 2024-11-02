package handlers

import (
	// "encoding/json"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
    _ "image/png" 
    _ "image/gif" 
	"io/ioutil"
	"net/http"
	"os"
	"strconv"	

	// "os"
	"path/filepath"
	// "time"

	"backend-wkj/database"
	"backend-wkj/models"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct{}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

func getNextImageNumber(folderPath string) (int, error) {
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return 0, err
	}

	maxNumber := 0
	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if ext == ".jpg" {
			name := file.Name()
			number, err := strconv.Atoi(name[:len(name)-len(ext)])
			if err == nil && number > maxNumber {
				maxNumber = number
			}
		}
	}
	return maxNumber + 1, nil
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var Product models.Product

	if err := c.ShouldBind(&Product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	utilizationString := c.PostForm("utilization")
	if utilizationString != "" {
		if err := json.Unmarshal([]byte(utilizationString), &Product.Utilization); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid utilization format"})
			return
		}
	} else {
		Product.Utilization = nil
	}

	compositionString := c.PostForm("composition")
	if compositionString != "" {
		if err := json.Unmarshal([]byte(compositionString), &Product.Composition); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid composition format"})
			return
		}
	} else {
	
		Product.Composition = nil
	}


	file, err := c.FormFile("image_url")
	if err == nil {
	
		srcFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded image"})
			return
		}
		defer srcFile.Close()

		
		img, format, err := image.Decode(srcFile)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image format: " + err.Error()})
			return
		}
		
		fmt.Println("Image format:", format)

		
		srcFile, _ = file.Open()
		defer srcFile.Close()


		nextNumber, err := getNextImageNumber("uploads")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate image filename"})
			return
		}
		imagePath := filepath.Join("uploads", strconv.Itoa(nextNumber)+".jpg")

	
		outFile, err := os.Create(imagePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save converted image"})
			return
		}
		defer outFile.Close()

		var options jpeg.Options
		options.Quality = 80
		if err := jpeg.Encode(outFile, img, &options); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode image as JPG"})
			return
		}

		Product.ImageURL = "/" + imagePath
	} else if err == http.ErrMissingFile {
		Product.ImageURL = ""
	}

	// Save product to the database
	if err := database.DB.Create(&Product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Product", "details": err.Error()})
		return
	}

	response := gin.H{
		"message": "Product created successfully",
		"Product": gin.H{
			"id":                  Product.ID,
			"name":                Product.Name,
			"latin_name":          Product.LatinName,
			"synonym":             Product.Synonym,
			"familia":             Product.Familia,
			"part_used":           Product.PartUsed,
			"method_of_reproduction": Product.MethodOfReproduction,
			"harvest_age":         Product.HarvestAge,
			"morphology":          Product.Morphology,
			"area_name":           Product.AreaName,
			"efficacy":            Product.Efficacy,
			"utilization": gin.H{
				"values": Product.Utilization,
			},
			"composition": gin.H{
				"values": Product.Composition,
			},
			"image_url":          Product.ImageURL,
			"research_results":   Product.ResearchResults,
			"description":        Product.Description,
			"price":              Product.Price,
			"unit_type":          Product.UnitType,
			"product_category_id": Product.ProductCategoryID,
			"created_at":         Product.CreatedAt,
			"updated_at":         Product.UpdatedAt,
		},
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product created successfully", "Product": response})
}


func (h *ProductHandler) GetProduct(c *gin.Context) {
	var products []models.Product

	if err := database.DB.Preload("ProductCategory").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Products"})
		return
	}

	response := []models.ProductResponse{}

	for _, product := range products {
		response = append(response, models.ProductResponse{
			ID:                  product.ID,
			Name:                product.Name,
			LatinName:           product.LatinName,
			Synonym:             product.Synonym,
			Familia:             product.Familia,
			PartUsed:            product.PartUsed,
			MethodOfReproduction: product.MethodOfReproduction,
			HarvestAge:          product.HarvestAge,
			Morphology:          product.Morphology,
			AreaName:            product.AreaName,
			Efficacy:            product.Efficacy,
			Utilization: gin.H{
				"values": product.Utilization,
			},
			Composition: gin.H{
				"values": product.Composition,
			},
			ImageURL:         product.ImageURL,
			ResearchResults:  product.ResearchResults,
			Description:      product.Description,
			Price:            product.Price,
			UnitType:         product.UnitType,
			ProductCategoryID: product.ProductCategoryID,
			ProductCategory:  product.ProductCategory,
			CreatedAt:        product.CreatedAt,
			UpdatedAt:        product.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"Products": response})
}



func (h *ProductHandler) GetProductByID(c *gin.Context) {
	var Product models.Product
	id := c.Param("id")

	if err := database.DB.Preload("ProductCategory").Find(&Product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Product": Product})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	var Product models.Product
	id := c.Param("id")

	if err := database.DB.First(&Product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if err := c.ShouldBind(&Product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	utilizationString := c.PostForm("utilization")
	if utilizationString != "" {
		if err := json.Unmarshal([]byte(utilizationString), &Product.Utilization); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid utilization format"})
			return
		}
	} else {
		Product.Utilization = nil
	}

	compositionString := c.PostForm("composition")
	if compositionString != "" {
		if err := json.Unmarshal([]byte(compositionString), &Product.Composition); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid composition format"})
			return
		}
	} else {
		Product.Composition = nil
	}

	file, err := c.FormFile("image_url")
	if err == nil {
		srcFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded image"})
			return
		}
		defer srcFile.Close()

		img, format, err := image.Decode(srcFile)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image format: " + err.Error()})
			return
		}
		fmt.Println("Image format:", format)

		srcFile, _ = file.Open()
		defer srcFile.Close()

		nextNumber, err := getNextImageNumber("uploads")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate image filename"})
			return
		}
		imagePath := filepath.Join("uploads", strconv.Itoa(nextNumber)+".jpg")

		outFile, err := os.Create(imagePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save converted image"})
			return
		}
		defer outFile.Close()

		var options jpeg.Options
		options.Quality = 80
		if err := jpeg.Encode(outFile, img, &options); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode image as JPG"})
			return
		}

		Product.ImageURL = "/" + imagePath
	} else if err == http.ErrMissingFile {
		Product.ImageURL = ""
	}

	// Save product to the database
	if err := database.DB.Create(&Product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Product", "details": err.Error()})
		return
	}

	response := gin.H{
		"message": "Product created successfully",
		"Product": gin.H{
			"id":                  Product.ID,
			"name":                Product.Name,
			"latin_name":          Product.LatinName,
			"synonym":             Product.Synonym,
			"familia":             Product.Familia,
			"part_used":           Product.PartUsed,
			"method_of_reproduction": Product.MethodOfReproduction,
			"harvest_age":         Product.HarvestAge,
			"morphology":          Product.Morphology,
			"area_name":           Product.AreaName,
			"efficacy":            Product.Efficacy,
			"utilization": gin.H{
				"values": Product.Utilization,
			},
			"composition": gin.H{ // Menampilkan composition dalam objek
				"values": Product.Composition,
			},
			"image_url":          Product.ImageURL,
			"research_results":   Product.ResearchResults,
			"description":        Product.Description,
			"price":              Product.Price,
			"unit_type":          Product.UnitType,
			"product_category_id": Product.ProductCategoryID,
			"created_at":         Product.CreatedAt,
			"updated_at":         Product.UpdatedAt,
		},
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully", "Product": response})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	var Product models.Product
	id := c.Param("id")

	if err := database.DB.First(&Product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if err := database.DB.Delete(&Product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
