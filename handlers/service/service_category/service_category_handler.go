package handlers

import (
	"net/http"

	"backend-wkj/database"
	"backend-wkj/models"

	"github.com/gin-gonic/gin"
)

type ServiceCategoryHandler struct{}

func NewServiceCategoryHandler() *ServiceCategoryHandler {
	return &ServiceCategoryHandler{}
}

func (h *ServiceCategoryHandler) CreateServiceCategory(c *gin.Context) {
	var ServiceCategory models.ServiceCategory

	// Parsing form data
	if err := c.ShouldBind(&ServiceCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&ServiceCategory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ServiceCategory"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ServiceCategory created successfully", "ServiceCategory": ServiceCategory})
}

func (h *ServiceCategoryHandler) GetServiceCategory(c *gin.Context) {
	var ServiceCategory []models.ServiceCategory

	if err := database.DB.Find(&ServiceCategory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ServiceCategorys"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ServiceCategory": ServiceCategory})
}


func (h *ServiceCategoryHandler) GetServiceCategoryByID(c *gin.Context) {
	var ServiceCategory models.ServiceCategory
	id := c.Param("id")

	if err := database.DB.First(&ServiceCategory, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ServiceCategory not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ServiceCategory": ServiceCategory})
}

func (h *ServiceCategoryHandler) UpdateServiceCategory(c *gin.Context) {
	var ServiceCategory models.ServiceCategory
	id := c.Param("id")

	if err := database.DB.First(&ServiceCategory, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ServiceCategory not found"})
		return
	}

	if err := c.ShouldBind(&ServiceCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&ServiceCategory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ServiceCategory"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ServiceCategory updated successfully", "ServiceCategory": ServiceCategory})
}

func (h *ServiceCategoryHandler) DeleteServiceCategory(c *gin.Context) {
	var ServiceCategory models.ServiceCategory
	id := c.Param("id")

	if err := database.DB.First(&ServiceCategory, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ServiceCategory not found"})
		return
	}

	if err := database.DB.Delete(&ServiceCategory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete ServiceCategory"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ServiceCategory deleted successfully"})
}
