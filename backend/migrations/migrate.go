package migrations

import (
	"errors"
	"time"

	"github.com/girishsaraf/xpense-trackr/backend/internal/budget"
	"github.com/girishsaraf/xpense-trackr/backend/internal/category"
	"github.com/girishsaraf/xpense-trackr/backend/internal/expense"
	"github.com/girishsaraf/xpense-trackr/backend/internal/goal"
	"github.com/girishsaraf/xpense-trackr/backend/internal/investments"
	"github.com/girishsaraf/xpense-trackr/backend/internal/ledger"
	"github.com/girishsaraf/xpense-trackr/backend/internal/legacyimport"
	"github.com/girishsaraf/xpense-trackr/backend/internal/planning"
	"github.com/girishsaraf/xpense-trackr/backend/internal/receipt"
	"github.com/girishsaraf/xpense-trackr/backend/internal/recurring"
	"gorm.io/gorm"
)

type schemaMigration struct {
	Version   uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:120;not null"`
	AppliedAt time.Time `gorm:"not null"`
}

type migration struct {
	version uint
	name    string
	apply   func(*gorm.DB) error
}

var schemaMigrations = []migration{
	{
		version: 1,
		name:    "legacy-finance-schema",
		apply: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				&category.Category{},
				&expense.Expense{},
				&budget.Budget{},
				&investments.Investment{},
				&goal.Goal{},
				&recurring.RecurringExpense{},
				&recurring.RecurringInvestment{},
			)
		},
	},
	{
		version: 2,
		name:    "unified-account-transaction-ledger",
		apply: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&ledger.Account{}, &ledger.Transaction{})
		},
	},
	{
		version: 3,
		name:    "merchant-categorization-rules",
		apply: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&ledger.CategoryRule{}, &ledger.Transaction{})
		},
	},
	{
		version: 4,
		name:    "scheduled-cashflow",
		apply: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&ledger.ScheduledFlow{})
		},
	},
	{
		version: 5,
		name:    "transaction-receipts",
		apply: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&receipt.Receipt{})
		},
	},
	{
		version: 6,
		name:    "transaction-csv-import-fingerprints",
		apply: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&ledger.Transaction{})
		},
	},
	{
		version: 7,
		name:    "goal-contributions-and-budget-envelopes",
		apply: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&planning.Goal{}, &planning.GoalContribution{}, &planning.Envelope{})
		},
	},
	{
		version: 8,
		name:    "bill-occurrences-and-reminders",
		apply: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&ledger.ScheduledFlow{}, &ledger.ScheduledOccurrence{})
		},
	},
	{
		version: 9,
		name:    "legacy-import-records",
		apply: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&legacyimport.ImportRecord{})
		},
	},
	{
		version: 10,
		name:    "manual-assets-and-valuations",
		apply: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&ledger.ManualAsset{}, &ledger.AssetValuation{})
		},
	},
	{
		version: 11,
		name:    "portfolio-holdings-and-valuations",
		apply: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&ledger.PortfolioHolding{}, &ledger.HoldingValuation{})
		},
	},
}

func (schemaMigration) TableName() string {
	return "schema_migrations"
}

// Migrate applies additive schema changes once and records the baseline version.
// GORM AutoMigrate does not drop existing columns or tables.
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&schemaMigration{}); err != nil {
		return err
	}

	for _, pending := range schemaMigrations {
		if err := db.Transaction(func(tx *gorm.DB) error {
			var applied schemaMigration
			err := tx.First(&applied, pending.version).Error
			if err == nil {
				return nil
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if err := pending.apply(tx); err != nil {
				return err
			}
			return tx.Create(&schemaMigration{
				Version:   pending.version,
				Name:      pending.name,
				AppliedAt: time.Now().UTC(),
			}).Error
		}); err != nil {
			return err
		}
	}
	return nil
}
