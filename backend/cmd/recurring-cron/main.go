package main

import (
	"github.com/girishsaraf/xpense-trackr/backend/config"
	"github.com/girishsaraf/xpense-trackr/backend/internal/recurring"
	"github.com/robfig/cron/v3"
	"log"
)

func main() {
	db, err := config.OpenDatabase(config.LoadDatabaseConfig())
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Error getting database connection:", err)
	}
	defer sqlDB.Close()
	recurring.SetDB(db)

	// Initialize the recurring scheduler
	cronInstance := cron.New()

	// Add a scheduled task to handle recurring expenses
	_, err = cronInstance.AddFunc("@every 12h", recurring.HandleRecurringExpense)
	if err != nil {
		log.Fatal("Error adding recurring expenses task:", err)
	}

	// Add a scheduled task to handle recurring investments
	_, err = cronInstance.AddFunc("@every 12h", recurring.HandleRecurringInvestment)
	if err != nil {
		log.Fatal("Error adding recurring investments task:", err)
	}

	// Start the recurring scheduler
	cronInstance.Start()

	// Keep the main function running to allow the recurring scheduler to execute the tasks
	select {}
}
