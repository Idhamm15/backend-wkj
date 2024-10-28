package handlers

import (
	"net/http"

	"backend-wkj/database"
	"backend-wkj/models"

	"github.com/gin-gonic/gin"
)

type ArticleCategoryHandler struct{}

func NewArticleCategoryHandler() *ArticleCategoryHandler {
	return &ArticleCategoryHandler{}
}

func (h *ArticleCategoryHandler) CreateArticleCategory(c *gin.Context) {
	var ArticleCategory models.ArticleCategory

	// Parsing form data
	if err := c.ShouldBind(&ArticleCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&ArticleCategory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ArticleCategory"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ArticleCategory created successfully", "ArticleCategory": ArticleCategory})
}

func (h *ArticleCategoryHandler) GetArticleCategory(c *gin.Context) {
	var ArticleCategory []models.ArticleCategory

	if err := database.DB.Find(&ArticleCategory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ArticleCategorys"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ArticleCategory": ArticleCategory})
}


func (h *ArticleCategoryHandler) GetArticleCategoryByID(c *gin.Context) {
	var ArticleCategory models.ArticleCategory
	id := c.Param("id")

	if err := database.DB.First(&ArticleCategory, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ArticleCategory not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ArticleCategory": ArticleCategory})
}

func (h *ArticleCategoryHandler) UpdateArticleCategory(c *gin.Context) {
	var ArticleCategory models.ArticleCategory
	id := c.Param("id")

	if err := database.DB.First(&ArticleCategory, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ArticleCategory not found"})
		return
	}

	if err := c.ShouldBind(&ArticleCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&ArticleCategory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ArticleCategory"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ArticleCategory updated successfully", "ArticleCategory": ArticleCategory})
}

func (h *ArticleCategoryHandler) DeleteArticleCategory(c *gin.Context) {
	var ArticleCategory models.ArticleCategory
	id := c.Param("id")

	if err := database.DB.First(&ArticleCategory, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ArticleCategory not found"})
		return
	}

	if err := database.DB.Delete(&ArticleCategory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete ArticleCategory"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article Category deleted successfully"})
}
