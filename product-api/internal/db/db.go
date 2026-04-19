package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool is a connection pool — reused across all requests.
// Never create a new connection per request (expensive).
var Pool *pgxpool.Pool

// Init connects to PostgreSQL and creates the products
// table if it doesn't exist yet.
func Init() {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	Pool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	// Verify the connection is actually alive
	if err = Pool.Ping(context.Background()); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	log.Println("Connected to PostgreSQL")
	createTable()
}

// createTable runs on startup — idempotent (safe to run multiple times)
func createTable() {
	query := `
		CREATE TABLE IF NOT EXISTS products (
			id          SERIAL PRIMARY KEY,
			name        VARCHAR(255) NOT NULL,
			description TEXT,
			price       DECIMAL(10,2) NOT NULL,
			stock       INTEGER NOT NULL DEFAULT 0,
			created_at  TIMESTAMP DEFAULT NOW()
		)
	`
	_, err := Pool.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to create products table: %v", err)
	}
	log.Println("Products table ready")
}