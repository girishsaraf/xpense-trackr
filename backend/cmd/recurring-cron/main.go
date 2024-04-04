package main

import (
	"github.com/girishsaraf/xpense-trackr/backend/config"
	"github.com/girishsaraf/xpense-trackr/backend/recurring"
	"github.com/robfig/cron/v3"
	"log"
	"os"
)

func main() {

	// Setting Environment Variables - Will be replaced with Dockerfile.mysql changes in the future
	if err := os.Setenv("DB_HOST", "mysql-container"); err != nil {
		// Handle error
		log.Fatalf("Error setting DB_HOST environment variable: %v", err)
	}
	if err := os.Setenv("DB_USER", "root"); err != nil {
		// Handle error
		log.Fatalf("Error setting DB_USER environment variable: %v", err)
	}
	if err := os.Setenv("DB_PASS", "root_password"); err != nil {
		// Handle error
		log.Fatalf("Error setting DB_PASS environment variable: %v", err)
	}
	if err := os.Setenv("DB_NAME", "xpense_trackr_db"); err != nil {
		// Handle error
		log.Fatalf("Error setting DB_NAME environment variable: %v", err)
	}
	if err := os.Setenv("DB_PORT", "3306"); err != nil {
		// Handle error
		log.Fatalf("Error setting DB_PORT environment variable: %v", err)
	}

	// Setting Database Connections for all models
	db := config.GetDatabaseConnectionObject(config.LoadDatabaseConfig())
	recurring.SetDB(db)

	// Initialize the recurring scheduler
	cronInstance := cron.New()

	// Add a scheduled task to handle recurring expenses
	_, err := cronInstance.AddFunc("@every 24h", recurring.HandleRecurringExpense)
	if err != nil {
		log.Fatal("Error adding recurring expenses task:", err)
	}

	// Add a scheduled task to handle recurring investments
	_, err = cronInstance.AddFunc("@every 24h", recurring.HandleRecurringInvestment)
	if err != nil {
		log.Fatal("Error adding recurring investments task:", err)
	}

	// Start the recurring scheduler
	cronInstance.Start()

	// Keep the main function running to allow the recurring scheduler to execute the tasks
	select {}
}
