// Package testutil provides shared helpers for backend tests.
package testutil

import (
	"fmt"
	"sync/atomic"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbCounter int64

// NewDB creates an isolated in-memory SQLite database for a test and migrates
// the provided models into it. The database is closed automatically when the
// test finishes.
func NewDB(t *testing.T, models ...interface{}) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:testdb%d?mode=memory&cache=shared", atomic.AddInt64(&dbCounter, 1))
	database, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		t.Fatalf("failed to get underlying database handle: %v", err)
	}
	// An in-memory SQLite database lives for as long as a connection to it
	// exists, so keep a single connection open for the duration of the test.
	sqlDB.SetMaxOpenConns(1)
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	if len(models) > 0 {
		if err := database.AutoMigrate(models...); err != nil {
			t.Fatalf("failed to migrate test database: %v", err)
		}
	}

	return database
}
