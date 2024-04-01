package main

import (
	"github.com/girishsaraf/xpense-trackr/backend/config"
	"github.com/girishsaraf/xpense-trackr/backend/internal"
	"github.com/girishsaraf/xpense-trackr/backend/internal/budget"
	"github.com/girishsaraf/xpense-trackr/backend/internal/category"
	"github.com/girishsaraf/xpense-trackr/backend/internal/expense"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
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
	budget.SetDB(db)
	expense.SetDB(db)
	category.SetDB(db)

	// Creating API routes and starting server on 8080
	router := mux.NewRouter()
	c := cors.AllowAll()
	handler := c.Handler(router)
	internal.SetRoutes(router)
	log.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
