package main

import (
	"atlas/config"
	"atlas/server"
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

	// Create a new SQLC store instance for interacting with the database
	// store := db.NewStore(conn)

	// Initialize a new HTTP server with all required dependencies (config, store, auth, etc.)
	server, err := server.NewServer(cfg /*store*/)
	if err != nil {
		log.Fatal("Failed to create new server:", err)
	}

	// Start the HTTP server on the specified port
	if err := server.Start(cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
