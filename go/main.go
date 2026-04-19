package main

import (
	"atlas/config"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Load environment/configuration variables from the project root (e.g., .env)
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("config load failed:", err)
	}

	// Initialize a PostgreSQL connection pool using pgx with the DB URL from config
	_, err = pgxpool.New(context.Background(), cfg.POSTGRES_URL)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
}
