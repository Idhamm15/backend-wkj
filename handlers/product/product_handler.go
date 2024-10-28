package handlers

import (
	"net/http"
	// "os"
	"path/filepath"
	"time"

	"backend-wkj/database"
	"backend-wkj/models"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct{}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var Product models.Product

	// Parsing form data
	if err := c.ShouldBind(&Product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle image upload
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image is required"})
		return
	}

	// Save the file to the server
	filename := filepath.Base(file.Filename)
	filepath := filepath.Join("uploads", time.Now().Format("20060102150405")+"_"+filename)
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	Product.ImageURL = "/" + filepath

	if err := database.DB.Create(&Product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product created successfully", "Product": Product})
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	var Product []models.Product

	if err := database.DB.Preload("ProductCategory").Find(&Product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Product": Product})
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

	// Handle image upload if exists
	file, err := c.FormFile("image")
	if err == nil {
		// Save the new file to the server
		filename := filepath.Base(file.Filename)
		filepath := filepath.Join("uploads", time.Now().Format("20060102150405")+"_"+filename)
		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
		Product.ImageURL = "/" + filepath
	}

	if err := database.DB.Save(&Product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully", "Product": Product})
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
