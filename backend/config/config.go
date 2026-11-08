package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const defaultDatabasePath = "./data/xpense-trackr.db"

// LoadDatabaseConfig loads database settings from the environment. SQLite is
// the default so a local installation needs no separately provisioned service.
func LoadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Driver:   envOrDefault("DB_DRIVER", "sqlite"),
		Path:     envOrDefault("DB_PATH", defaultDatabasePath),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Name:     os.Getenv("DB_NAME"),
	}
}

// OpenDatabase opens the configured database and returns connection errors to
// the caller so startup can fail with context instead of terminating here.
func OpenDatabase(dbConfig DatabaseConfig) (*gorm.DB, error) {
	switch strings.ToLower(dbConfig.Driver) {
	case "sqlite":
		path := dbConfig.Path
		if path == "" {
			path = defaultDatabasePath
		}
		if path != ":memory:" && !strings.HasPrefix(path, "file:") {
			if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
				return nil, fmt.Errorf("create database directory: %w", err)
			}
		}
		db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("open SQLite database: %w", err)
		}
		return db, nil
	case "mysql":
		if dbConfig.Host == "" || dbConfig.Port == "" || dbConfig.User == "" || dbConfig.Name == "" {
			return nil, fmt.Errorf("MySQL requires DB_HOST, DB_PORT, DB_USER, and DB_NAME")
		}
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbConfig.User,
			dbConfig.Password,
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.Name,
		)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("open MySQL database: %w", err)
		}
		return db, nil
	default:
		return nil, fmt.Errorf("unsupported DB_DRIVER %q: expected sqlite or mysql", dbConfig.Driver)
	}
}

// DatabaseConfig represents database connection settings.
type DatabaseConfig struct {
	Driver   string
	Path     string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
