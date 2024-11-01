	package handlers

	import (
		"net/http"
		"strings"

		"backend-wkj/database"
		"backend-wkj/models"

		"github.com/gin-gonic/gin"
	)

	type ProductCategoryHandler struct{}

	func NewProductCategoryHandler() *ProductCategoryHandler {
		return &ProductCategoryHandler{}
	}

	func (h *ProductCategoryHandler) CreateProductCategory(c *gin.Context) {
		var ProductCategory models.ProductCategory

		// Parsing form data
		if err := c.ShouldBind(&ProductCategory); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := database.DB.Create(&ProductCategory).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ProductCategory"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "ProductCategory created successfully", "ProductCategory": ProductCategory})
	}

	func (h *ProductCategoryHandler) GetProductCategory(c *gin.Context) {
		var ProductCategory []models.ProductCategory

		if err := database.DB.Find(&ProductCategory).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ProductCategorys"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"ProductCategory": ProductCategory})
	}


	func (h *ProductCategoryHandler) GetProductCategoryByID(c *gin.Context) {
		var ProductCategory models.ProductCategory
		id := c.Param("id")

		if err := database.DB.First(&ProductCategory, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ProductCategory not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"ProductCategory": ProductCategory})
	}

	func (h *ProductCategoryHandler) UpdateProductCategory(c *gin.Context) {
		var ProductCategory models.ProductCategory
		id := c.Param("id")

		if err := database.DB.First(&ProductCategory, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ProductCategory not found"})
			return
		}

		if err := c.ShouldBind(&ProductCategory); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := database.DB.Save(&ProductCategory).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ProductCategory"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "ProductCategory updated successfully", "ProductCategory": ProductCategory})
	}

	func (h *ProductCategoryHandler) DeleteProductCategory(c *gin.Context) {
		var productCategory models.ProductCategory
		id := c.Param("id")

		// Cari kategori produk berdasarkan ID
		if err := database.DB.First(&productCategory, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product category not found"})
			return
		}

		// Coba hapus kategori produk
		if err := database.DB.Delete(&productCategory).Error; err != nil {
			// Periksa apakah error disebabkan oleh constraint foreign key
			if strings.Contains(err.Error(), "1451") {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Cannot delete product category",
					"message": "Tidak bisa menghapus kategori produk karena ada produk yang terkait dengan kategori ini.",
				})
				return
			}

			// Jika error lain
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product category", "details": err.Error()})
			return
		}

		// Jika berhasil dihapus
		c.JSON(http.StatusOK, gin.H{"message": "Product category deleted successfully"})
	}

