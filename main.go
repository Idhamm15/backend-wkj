package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/gin-contrib/cors"

	"backend-wkj/handlers"
	articleHandler "backend-wkj/handlers/article"
	articleCategoryHandler "backend-wkj/handlers/article/article_category"
	productHandler "backend-wkj/handlers/product"
	productCategoryHandler "backend-wkj/handlers/product/product_category"
	serviceHandler "backend-wkj/handlers/service"
	serviceCategoryHandler "backend-wkj/handlers/service/service_category"
	"backend-wkj/middleware"
	"backend-wkj/database"
)

func envPortOr(port string) string {
	if envPort := os.Getenv("PORT"); envPort != "" {
		return ":" + envPort
	}
	return ":" + port
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	database.Init()

	r := gin.Default()

	// Mendapatkan absolute path untuk folder uploads
	uploadsDir := "./uploads"
	absUploadsDir, err := filepath.Abs(uploadsDir)
	if err != nil {
		log.Printf("Warning: Could not get absolute path for uploads directory: %v", err)
	} else {
		log.Printf("Uploads directory absolute path: %s", absUploadsDir)
	}

	// Cek apakah folder uploads ada
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		log.Printf("Warning: Uploads directory does not exist, creating it...")
		if err := os.MkdirAll(uploadsDir, 0755); err != nil {
			log.Printf("Error creating uploads directory: %v", err)
		}
	}

	// Tambahkan middleware untuk logging static file requests
	r.Use(func(c *gin.Context) {
		if filepath.HasPrefix(c.Request.URL.Path, "/uploads/") {
			log.Printf("Static file request: %s", c.Request.URL.Path)
			requestedFile := filepath.Join(uploadsDir, filepath.Clean(c.Request.URL.Path[8:]))
			if _, err := os.Stat(requestedFile); os.IsNotExist(err) {
				log.Printf("File not found: %s", requestedFile)
			}
		}
		c.Next()
	})

	r.Static("/uploads", uploadsDir)

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	
	authHandler := handlers.NewAuthHandler()
	articleHandler := articleHandler.NewArticleHandler()
	articleCategoryHandler := articleCategoryHandler.NewArticleCategoryHandler()
	productHandler := productHandler.NewProductHandler()
	productCategoryHandler := productCategoryHandler.NewProductCategoryHandler()
	serviceHandler := serviceHandler.NewServiceHandler()
	serviceCategoryHandler := serviceCategoryHandler.NewServiceCategoryHandler()

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.POST("/logout", authHandler.Logout)

	r.GET("/check", middleware.AuthMiddleware(), func(c *gin.Context) {
		role, _ := c.Get("role")
		if role == "admin" {
			c.JSON(http.StatusOK, gin.H{"message": "Kamu admin"})
		} else if role == "user" {
			c.JSON(http.StatusOK, gin.H{"message": "Kamu user"})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Kamu belum login"})
		}
	})

	// Public routes

	// Article
	r.GET("/article", articleHandler.GetArticle)
	r.GET("/article/:id", articleHandler.GetArticleByID)

	// Product
	r.GET("/product", productHandler.GetProduct)
	r.GET("/product/:id", productHandler.GetProductByID)

	// Service
	r.GET("/service", serviceHandler.GetService)
	r.GET("/service/:id", serviceHandler.GetServiceByID)
	r.GET("/article-category", articleCategoryHandler.GetArticleCategory)
	r.GET("/product-category", productCategoryHandler.GetProductCategory)

	// Admin routes
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		// Article
		admin.POST("/article-category", articleCategoryHandler.CreateArticleCategory)
		admin.PUT("/article-category/:id", articleCategoryHandler.UpdateArticleCategory)
		admin.DELETE("/article-category/:id", articleCategoryHandler.DeleteArticleCategory)

		admin.POST("/article", articleHandler.CreateArticle)
		admin.PUT("/article/:id", articleHandler.UpdateArticle)
		admin.DELETE("/article/:id", articleHandler.DeleteArticle)

		// Product
		admin.POST("/product-category", productCategoryHandler.CreateProductCategory)
		admin.PUT("/product-category/:id", productCategoryHandler.UpdateProductCategory)
		admin.DELETE("/product-category/:id", productCategoryHandler.DeleteProductCategory)

		admin.POST("/product", productHandler.CreateProduct)
		admin.PUT("/product/:id", productHandler.UpdateProduct)
		admin.DELETE("/product/:id", productHandler.DeleteProduct)

		// Service
		admin.GET("/service-category", serviceCategoryHandler.GetServiceCategory)
		admin.POST("/service-category", serviceCategoryHandler.CreateServiceCategory)
		admin.PUT("/service-category/:id", serviceCategoryHandler.UpdateServiceCategory)
		admin.DELETE("/service-category/:id", serviceCategoryHandler.DeleteServiceCategory)

		admin.POST("/service", serviceHandler.CreateService)
		admin.PUT("/service/:id", serviceHandler.UpdateService)
		admin.DELETE("/service/:id", serviceHandler.DeleteService)
	}
	

	// Set trusted proxies
	if err := r.SetTrustedProxies(nil); err != nil {
		panic(err)
	}

	// Use `PORT` provided in environment or default to 3000
	var port = envPortOr("5000")

	// Mulai server
	log.Fatal(http.ListenAndServe(port, r))
}
