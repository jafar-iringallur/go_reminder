package main

import (
	"log"

	"go-reminder/internal/config"
	"go-reminder/internal/database"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	if err := database.Seed(db); err != nil {
		log.Fatalf("database seeding failed: %v", err)
	}

	log.Println("database seeded successfully")
}
