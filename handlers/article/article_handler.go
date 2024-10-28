package handlers

import (
	"net/http"
	"path/filepath"
	"time"

	"backend-wkj/database"
	"backend-wkj/models"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct{}

func NewArticleHandler() *ArticleHandler {
	return &ArticleHandler{}
}

func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var article models.Article

	// Parsing form data
	if err := c.ShouldBind(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle image upload
	file, err := c.FormFile("image_url")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image Url is required"})
		return
	}

	// Save the file to the server
	filename := filepath.Base(file.Filename)
	filepath := filepath.Join("uploads", time.Now().Format("20060102150405")+"_"+filename)
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	article.ImageURL = "/" + filepath

	if err := database.DB.Create(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article created successfully", "article": article})
}

func (h *ArticleHandler) GetArticle(c *gin.Context) {
	var article []models.Article

	if err := database.DB.Preload("ArticleCategory").Find(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch articles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"article": article})
}


func (h *ArticleHandler) GetArticleByID(c *gin.Context) {
	var article models.Article
	id := c.Param("id")

	if err := database.DB.Preload("ArticleCategory").Find(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"article": article})
}

func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	var article models.Article
	id := c.Param("id")

	if err := database.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	if err := c.ShouldBind(&article); err != nil {
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
		article.ImageURL = "/" + filepath
	}

	if err := database.DB.Save(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article updated successfully", "article": article})
}

func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	var article models.Article
	id := c.Param("id")

	if err := database.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	if err := database.DB.Delete(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}
