package models

import "time"

// Product represents a product in our system.
// The struct tags (json:"...") control how fields
// are serialised to JSON in API responses.
type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateProductRequest is what the client sends
// when creating a product. Separate from Product
// because the client doesn't set ID or CreatedAt.
type CreateProductRequest struct {
	Name        string  `json:"name"        binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price"       binding:"required,gt=0"`
	Stock       int     `json:"stock"       binding:"required,gte=0"`
}