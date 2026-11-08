package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDatabaseConfig(t *testing.T) {
	t.Setenv("DB_DRIVER", "mysql")
	t.Setenv("DB_PATH", "/tmp/expenses.db")
	t.Setenv("DB_HOST", "db.example.com")
	t.Setenv("DB_PORT", "3307")
	t.Setenv("DB_USER", "app")
	t.Setenv("DB_PASS", "secret")
	t.Setenv("DB_NAME", "expenses")

	config := LoadDatabaseConfig()

	if config.Driver != "mysql" {
		t.Errorf("Driver = %q, want %q", config.Driver, "mysql")
	}
	if config.Path != "/tmp/expenses.db" {
		t.Errorf("Path = %q, want %q", config.Path, "/tmp/expenses.db")
	}
	if config.Host != "db.example.com" {
		t.Errorf("Host = %q, want %q", config.Host, "db.example.com")
	}
	if config.Port != "3307" {
		t.Errorf("Port = %q, want %q", config.Port, "3307")
	}
	if config.User != "app" {
		t.Errorf("User = %q, want %q", config.User, "app")
	}
	if config.Password != "secret" {
		t.Errorf("Password = %q, want %q", config.Password, "secret")
	}
	if config.Name != "expenses" {
		t.Errorf("Name = %q, want %q", config.Name, "expenses")
	}
}

func TestLoadDatabaseConfigDefaultsToSQLite(t *testing.T) {
	t.Setenv("DB_DRIVER", "")
	t.Setenv("DB_PATH", "")
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASS", "")
	t.Setenv("DB_NAME", "")

	config := LoadDatabaseConfig()

	if config.Driver != "sqlite" {
		t.Errorf("Driver = %q, want %q", config.Driver, "sqlite")
	}
	if config.Path != defaultDatabasePath {
		t.Errorf("Path = %q, want %q", config.Path, defaultDatabasePath)
	}
	if config.Host != "" || config.Port != "" || config.User != "" || config.Password != "" || config.Name != "" {
		t.Errorf("MySQL settings should be empty when unset: %+v", config)
	}
}

func TestOpenDatabaseCreatesSQLiteDatabase(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data", "xpense.db")
	database, err := OpenDatabase(DatabaseConfig{Driver: "sqlite", Path: path})
	if err != nil {
		t.Fatalf("OpenDatabase() error = %v", err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		t.Fatalf("database.DB() error = %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("SQLite ping error = %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("SQLite database file was not created: %v", err)
	}
}

func TestOpenDatabaseRejectsUnknownDriver(t *testing.T) {
	if _, err := OpenDatabase(DatabaseConfig{Driver: "postgres"}); err == nil {
		t.Fatal("OpenDatabase() error = nil, want unsupported driver error")
	}
}
