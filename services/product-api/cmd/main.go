package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sanaanidhal/nexuscloud/product-api/internal/db"
	"github.com/sanaanidhal/nexuscloud/product-api/internal/handlers"
)

func main() {
	// Load .env file (ignored in production where env vars are injected)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found — using environment variables")
	}

	// Initialise database connection + create tables
	db.Init()

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Health check — same pattern as auth service
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "product-api",
			"version": "1.0.0",
		})
	})

	// Product routes
	v1 := r.Group("/api/v1")
	{
		v1.GET("/products", handlers.GetProducts)
		v1.GET("/products/:id", handlers.GetProduct)
		v1.POST("/products", handlers.CreateProduct)
		v1.DELETE("/products/:id", handlers.DeleteProduct)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("[NexusCloud] Product API running on port %s", port)
	r.Run(":" + port)
}