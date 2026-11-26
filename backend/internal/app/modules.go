package app

import (
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
	"gorm.io/gorm"
)

// ConfigureModules wires the shared application dependencies for the
// feature modules used by the HTTP server and integration tests.
func ConfigureModules(database *gorm.DB, receiptStorageDir string) {
	budget.SetDB(database)
	category.SetDB(database)
	expense.SetDB(database)
	investments.SetDB(database)
	ledger.SetDB(database)
	planning.SetDB(database)
	goal.SetDB(database)
	recurring.SetDB(database)
	receipt.SetDB(database)
	receipt.SetStorageDir(receiptStorageDir)
	trends.SetDB(database)
}
