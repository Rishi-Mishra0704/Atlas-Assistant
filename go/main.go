package main

import (
	"atlas/config"
	"atlas/server"
	"log"
)

func main() {
	// Load environment/configuration variables from the project root (e.g., .env)
	cfg, err := config.LoadConfig("../")
	if err != nil {
		log.Fatal("config load failed:", err)
	}

	// Initialize a PostgreSQL connection pool using pgx with the DB URL from config
	// _, err = pgxpool.New(context.Background(), cfg.POSTGRES_URL)
	// if err != nil {
	// 	log.Fatal("Failed to connect to the database:", err)
	// }

	server, err := server.NewServer(cfg)
	if err != nil {
		log.Fatal("Failed to create server:", err)
	}

	log.Printf("Starting server on port %s...", cfg.Port)
	if err := server.Start(cfg.Port); err != nil {
		log.Fatal("Server failed to start:", err)
	}

}
