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

type ServiceHandler struct{}

func NewServiceHandler() *ServiceHandler {
	return &ServiceHandler{}
}

func (h *ServiceHandler) CreateService(c *gin.Context) {
	var Service models.Service

	// Parsing form data
	if err := c.ShouldBind(&Service); err != nil {
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

	Service.ImageURL = "/" + filepath

	if err := database.DB.Create(&Service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service created successfully", "Service": Service})
}

func (h *ServiceHandler) GetService(c *gin.Context) {
	var Service []models.Service

	if err := database.DB.Preload("ServiceCategory").Find(&Service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Services"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Service": Service})
}


func (h *ServiceHandler) GetServiceByID(c *gin.Context) {
	var Service models.Service
	id := c.Param("id")

	if err := database.DB.Preload("ServiceCategory").Find(&Service, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Service": Service})
}

func (h *ServiceHandler) UpdateService(c *gin.Context) {
	var Service models.Service
	id := c.Param("id")

	if err := database.DB.First(&Service, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	if err := c.ShouldBind(&Service); err != nil {
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
		Service.ImageURL = "/" + filepath
	}

	if err := database.DB.Save(&Service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service updated successfully", "Service": Service})
}

func (h *ServiceHandler) DeleteService(c *gin.Context) {
	var Service models.Service
	id := c.Param("id")

	if err := database.DB.First(&Service, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	if err := database.DB.Delete(&Service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service deleted successfully"})
}
