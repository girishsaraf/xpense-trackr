package migrations

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestMigrateCreatesLegacySchemaAndIsRepeatable(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:migration-test?mode=memory&cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}

	if err := Migrate(db); err != nil {
		t.Fatalf("first Migrate() error = %v", err)
	}
	if err := Migrate(db); err != nil {
		t.Fatalf("second Migrate() error = %v", err)
	}

	var count int64
	if err := db.Model(&schemaMigration{}).Count(&count).Error; err != nil {
		t.Fatalf("count migration records: %v", err)
	}
	if count != 6 {
		t.Errorf("migration record count = %d, want 6", count)
	}

	for _, table := range []string{
		"categories",
		"expenses",
		"budgets",
		"investments",
		"accounts",
		"ledger_transactions",
		"category_rules",
		"scheduled_flows",
		"receipts",
		"goals",
		"recurring_expenses",
		"recurring_investments",
	} {
		if !db.Migrator().HasTable(table) {
			t.Errorf("migration did not create %q table", table)
		}
	}
}
