package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/girishsaraf/xpense-trackr/backend/config"
	"github.com/girishsaraf/xpense-trackr/backend/internal"
	"github.com/girishsaraf/xpense-trackr/backend/internal/backup"
	"github.com/girishsaraf/xpense-trackr/backend/internal/budget"
	"github.com/girishsaraf/xpense-trackr/backend/internal/category"
	"github.com/girishsaraf/xpense-trackr/backend/internal/expense"
	"github.com/girishsaraf/xpense-trackr/backend/internal/goal"
	"github.com/girishsaraf/xpense-trackr/backend/internal/investments"
	"github.com/girishsaraf/xpense-trackr/backend/internal/ledger"
	"github.com/girishsaraf/xpense-trackr/backend/internal/planning"
	"github.com/girishsaraf/xpense-trackr/backend/internal/receipt"
	"github.com/girishsaraf/xpense-trackr/backend/internal/recurring"
	"github.com/girishsaraf/xpense-trackr/backend/internal/trends"
	"github.com/girishsaraf/xpense-trackr/backend/migrations"
	"github.com/gorilla/mux"
	"github.com/robfig/cron/v3"
	"github.com/rs/cors"
	"gorm.io/gorm"
)

func main() {
	var err error
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "backup":
			err = runBackupCommand(os.Args[2:])
		case "restore":
			err = runRestoreCommand(os.Args[2:])
		default:
			err = fmt.Errorf("unknown command %q; use backup, restore, or start the server without a command", os.Args[1])
		}
	} else {
		err = run()
	}
	if err != nil {
		log.Fatal(err)
	}
}

func runBackupCommand(args []string) error {
	flags := flag.NewFlagSet("backup", flag.ContinueOnError)
	output := flags.String("output", "", "destination ZIP archive")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *output == "" || flags.NArg() != 0 {
		return errors.New("usage: xpense-trackr backup --output <archive.zip>")
	}
	database := config.LoadDatabaseConfig()
	if database.Driver != "sqlite" {
		return errors.New("backup currently supports SQLite databases only")
	}
	dataDir := envOrDefault("APP_DATA_DIR", "./data")
	if err := backup.Create(database.Path, dataDir, *output); err != nil {
		return fmt.Errorf("create backup: %w", err)
	}
	fmt.Printf("Backup saved to %s\n", *output)
	return nil
}

func runRestoreCommand(args []string) error {
	flags := flag.NewFlagSet("restore", flag.ContinueOnError)
	input := flags.String("input", "", "source ZIP archive")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *input == "" || flags.NArg() != 0 {
		return errors.New("usage: xpense-trackr restore --input <archive.zip>")
	}
	database := config.LoadDatabaseConfig()
	if database.Driver != "sqlite" {
		return errors.New("restore currently supports SQLite databases only")
	}
	dataDir := envOrDefault("APP_DATA_DIR", "./data")
	if err := backup.Restore(*input, database.Path, dataDir); err != nil {
		return fmt.Errorf("restore backup: %w", err)
	}
	fmt.Printf("Restored data from %s. The previous database and receipt files were preserved.\n", *input)
	return nil
}

func run() error {
	db, err := config.OpenDatabase(config.LoadDatabaseConfig())
	if err != nil {
		return fmt.Errorf("configure database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get database connection: %w", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}
	if err := migrations.Migrate(db); err != nil {
		return fmt.Errorf("migrate database schema: %w", err)
	}

	budget.SetDB(db)
	expense.SetDB(db)
	category.SetDB(db)
	investments.SetDB(db)
	ledger.SetDB(db)
	planning.SetDB(db)
	goal.SetDB(db)
	recurring.SetDB(db)
	receipt.SetDB(db)
	receipt.SetStorageDir(filepath.Join(envOrDefault("APP_DATA_DIR", "./data"), "receipts"))
	trends.SetDB(db)

	scheduler := cron.New()
	if _, err := scheduler.AddFunc("@every 12h", recurring.HandleRecurringExpense); err != nil {
		return fmt.Errorf("schedule recurring expenses: %w", err)
	}
	if _, err := scheduler.AddFunc("@every 12h", recurring.HandleRecurringInvestment); err != nil {
		return fmt.Errorf("schedule recurring investments: %w", err)
	}
	scheduler.Start()
	defer func() {
		<-scheduler.Stop().Done()
	}()

	apiRouter := mux.NewRouter()
	internal.SetRoutes(apiRouter)
	apiRouter.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}).Methods(http.MethodGet)
	apiRouter.HandleFunc("/readyz", databaseReadiness(db)).Methods(http.MethodGet)

	webDir := envOrDefault("WEB_DIR", "./web")
	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}).Handler(withFrontend(apiRouter, webDir))
	server := &http.Server{
		Addr:              envOrDefault("APP_ADDR", "127.0.0.1:8080"),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("server listening on %s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve HTTP: %w", err)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shut down HTTP server: %w", err)
		}
	}
	return nil
}

func databaseReadiness(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sqlDB, err := db.DB()
		if err != nil {
			http.Error(w, "database unavailable", http.StatusServiceUnavailable)
			return
		}
		if err := sqlDB.PingContext(r.Context()); err != nil {
			http.Error(w, "database unavailable", http.StatusServiceUnavailable)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
	}
}

func withFrontend(api http.Handler, webDir string) http.Handler {
	static := http.FileServer(http.Dir(webDir))
	indexPath := filepath.Join(webDir, "index.html")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/v1/") {
			versionedRequest := r.Clone(r.Context())
			versionedURL := *r.URL
			versionedURL.Path = "/api" + strings.TrimPrefix(r.URL.Path, "/api/v1")
			versionedRequest.URL = &versionedURL
			api.ServeHTTP(w, versionedRequest)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api") || r.URL.Path == "/healthz" || r.URL.Path == "/readyz" {
			api.ServeHTTP(w, r)
			return
		}

		if _, err := os.Stat(indexPath); err != nil {
			api.ServeHTTP(w, r)
			return
		}
		if r.URL.Path == "/" || path.Ext(r.URL.Path) == "" {
			http.ServeFile(w, r, indexPath)
			return
		}
		static.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("write HTTP response: %v", err)
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
