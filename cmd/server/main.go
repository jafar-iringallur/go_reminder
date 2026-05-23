package main

import (
	"log"

	"go-reminder/internal/config"
	"go-reminder/internal/database"
	httpRouter "go-reminder/internal/http/router"
	"go-reminder/internal/repository"
	"go-reminder/internal/scheduler"
	"go-reminder/internal/service"
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

	ruleRepo := repository.NewReminderRuleRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	auditRepo := repository.NewAuditLogRepository(db)

	reminderService := service.NewReminderService(ruleRepo, taskRepo, auditRepo)
	taskService := service.NewTaskService(taskRepo)
	auditService := service.NewAuditService(auditRepo)

	reminderScheduler := scheduler.New(reminderService, cfg.SchedulerInterval)
	reminderScheduler.Start()
	defer reminderScheduler.Stop()

	router := httpRouter.New(reminderService, taskService, auditService)

	log.Printf("server listening on %s", cfg.HTTPAddress)
	if err := router.Run(cfg.HTTPAddress); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
