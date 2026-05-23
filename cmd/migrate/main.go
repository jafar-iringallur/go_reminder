package main

import (
	"log"
	"os"

	"go-reminder/internal/config"
	"go-reminder/internal/database"
)

func main() {
	action := "up"
	if len(os.Args) > 1 {
		action = os.Args[1]
	}

	cfg := config.Load()
	db, err := database.Connect(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	switch action {
	case "up":
		err = database.RunMigrations(db)
	case "down":
		err = database.RollbackLastMigration(db)
	default:
		log.Fatalf("unknown migration action %q; use up or down", action)
	}

	if err != nil {
		log.Fatalf("migration %s failed: %v", action, err)
	}

	log.Printf("migration %s completed", action)
}
