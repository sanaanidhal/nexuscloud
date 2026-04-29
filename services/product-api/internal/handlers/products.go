package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sanaanidhal/nexuscloud/product-api/internal/db"
	"github.com/sanaanidhal/nexuscloud/product-api/internal/models"
)

// GetProducts handles GET /products
// Returns all products from the database
func GetProducts(c *gin.Context) {
	rows, err := db.Pool.Query(
		context.Background(),
		"SELECT id, name, description, price, stock, created_at FROM products ORDER BY id",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	defer rows.Close()

	// Build the products slice from rows
	products := []models.Product{}
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan product"})
			return
		}
		products = append(products, p)
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"count":    len(products),
	})
}

// GetProduct handles GET /products/:id
// Returns a single product by ID
func GetProduct(c *gin.Context) {
	// Extract the :id URL parameter and convert to int
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var p models.Product
	err = db.Pool.QueryRow(
		context.Background(),
		"SELECT id, name, description, price, stock, created_at FROM products WHERE id = $1",
		id,
	).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, p)
}

// CreateProduct handles POST /products
// Creates a new product in the database
func CreateProduct(c *gin.Context) {
	var req models.CreateProductRequest

	// ShouldBindJSON parses the body AND validates struct tags
	// Returns 400 automatically if required fields are missing
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var p models.Product
	err := db.Pool.QueryRow(
		context.Background(),
		`INSERT INTO products (name, description, price, stock)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, name, description, price, stock, created_at`,
		req.Name, req.Description, req.Price, req.Stock,
	).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, p)
}

// DeleteProduct handles DELETE /products/:id
func DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	result, err := db.Pool.Exec(
		context.Background(),
		"DELETE FROM products WHERE id = $1",
		id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	// RowsAffected tells us if the product actually existed
	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}