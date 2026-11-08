package legacyimport_test

import (
	"errors"
	"testing"
	"time"

	"github.com/girishsaraf/xpense-trackr/backend/internal/budget"
	"github.com/girishsaraf/xpense-trackr/backend/internal/category"
	"github.com/girishsaraf/xpense-trackr/backend/internal/expense"
	"github.com/girishsaraf/xpense-trackr/backend/internal/goal"
	"github.com/girishsaraf/xpense-trackr/backend/internal/investments"
	"github.com/girishsaraf/xpense-trackr/backend/internal/ledger"
	"github.com/girishsaraf/xpense-trackr/backend/internal/legacyimport"
	"github.com/girishsaraf/xpense-trackr/backend/internal/planning"
	"github.com/girishsaraf/xpense-trackr/backend/internal/recurring"
	"github.com/girishsaraf/xpense-trackr/backend/migrations"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestRunPreviewsThenTransactionallyImportsLegacyFinance(t *testing.T) {
	source, destination := openDatabase(t, "source"), openDatabase(t, "destination")
	if err := source.AutoMigrate(
		&category.Category{}, &expense.Expense{}, &budget.Budget{}, &goal.Goal{},
		&investments.Investment{}, &recurring.RecurringExpense{}, &recurring.RecurringInvestment{},
	); err != nil {
		t.Fatalf("migrate source fixture: %v", err)
	}
	if err := migrations.Migrate(destination); err != nil {
		t.Fatalf("migrate destination: %v", err)
	}
	if err := seedLegacySource(source); err != nil {
		t.Fatalf("seed legacy source: %v", err)
	}

	report, err := legacyimport.Run(source, destination, "legacy-db", "USD", false)
	if err != nil {
		t.Fatalf("preview import: %v", err)
	}
	if report.Applied || report.Categories != 1 || report.Expenses != 1 ||
		report.ExpensesTotalMinor != 1234 || report.BudgetCategories != 1 || !report.UsersNotImported {
		t.Fatalf("unexpected preview report: %+v", report)
	}
	var count int64
	if err := destination.Model(&category.Category{}).Count(&count).Error; err != nil || count != 0 {
		t.Fatalf("preview changed destination: count=%d error=%v", count, err)
	}

	report, err = legacyimport.Run(source, destination, "legacy-db", "USD", true)
	if err != nil {
		t.Fatalf("apply import: %v", err)
	}
	if !report.Applied {
		t.Fatal("applied import report is not marked applied")
	}
	var transactions []ledger.Transaction
	if err := destination.Find(&transactions).Error; err != nil {
		t.Fatalf("read imported ledger: %v", err)
	}
	if len(transactions) != 1 || transactions[0].AmountMinor != 1234 ||
		transactions[0].Type != "expense" || transactions[0].CategoryID == nil || *transactions[0].CategoryID != 1 {
		t.Fatalf("unexpected imported ledger transactions: %+v", transactions)
	}
	var envelopes []planning.Envelope
	if err := destination.Find(&envelopes).Error; err != nil || len(envelopes) != 1 ||
		envelopes[0].MonthlyLimitMinor != 15000 || envelopes[0].Currency != "USD" {
		t.Fatalf("unexpected imported budget envelopes: %+v error=%v", envelopes, err)
	}
	var importedGoals []planning.Goal
	if err := destination.Find(&importedGoals).Error; err != nil || len(importedGoals) != 2 {
		t.Fatalf("unexpected imported goals: %+v error=%v", importedGoals, err)
	}
	var contributions []planning.GoalContribution
	if err := destination.Where("goal_id = ?", importedGoals[0].ID).Find(&contributions).Error; err != nil ||
		len(contributions) != 1 || contributions[0].AmountMinor != 5000 {
		t.Fatalf("unexpected imported goal balance: %+v error=%v", contributions, err)
	}
	if err := destination.Model(&investments.Investment{}).Count(&count).Error; err != nil || count != 1 {
		t.Fatalf("legacy investment history not preserved: count=%d error=%v", count, err)
	}
	if err := destination.Model(&recurring.RecurringExpense{}).Count(&count).Error; err != nil || count != 1 {
		t.Fatalf("legacy recurring rule not preserved: count=%d error=%v", count, err)
	}
	if err := destination.Model(&legacyimport.ImportRecord{}).Count(&count).Error; err != nil || count != 1 {
		t.Fatalf("import marker missing: count=%d error=%v", count, err)
	}
	if _, err := legacyimport.Run(source, destination, "legacy-db", "USD", true); !errors.Is(err, legacyimport.ErrImportAlreadyRun) {
		t.Fatalf("second import error = %v, want ErrImportAlreadyRun", err)
	}
}

func TestRunRejectsInvalidSourceAndNonEmptyDestination(t *testing.T) {
	source, destination := openDatabase(t, "invalid-source"), openDatabase(t, "invalid-destination")
	if err := source.AutoMigrate(&category.Category{}, &expense.Expense{}); err != nil {
		t.Fatalf("migrate source fixture: %v", err)
	}
	if err := migrations.Migrate(destination); err != nil {
		t.Fatalf("migrate destination: %v", err)
	}
	if err := source.Create(&expense.Expense{
		ID: 1, CategoryID: 77, Amount: 1.25, Description: "Invalid category", Date: "2026-01-02",
	}).Error; err != nil {
		t.Fatalf("seed invalid expense: %v", err)
	}
	if _, err := legacyimport.Run(source, destination, "invalid-db", "USD", true); err == nil {
		t.Fatal("expected invalid source reference to fail")
	}
	var count int64
	if err := destination.Model(&ledger.Account{}).Count(&count).Error; err != nil || count != 0 {
		t.Fatalf("invalid source changed destination: count=%d error=%v", count, err)
	}

	if err := destination.Create(&category.Category{Name: "Existing"}).Error; err != nil {
		t.Fatalf("seed non-empty destination: %v", err)
	}
	if err := source.Create(&category.Category{ID: 77, Name: "Other"}).Error; err != nil {
		t.Fatalf("complete valid source fixture: %v", err)
	}
	if _, err := legacyimport.Run(source, destination, "different-db", "USD", false); !errors.Is(err, legacyimport.ErrDestinationNotEmpty) {
		t.Fatalf("non-empty destination error = %v, want ErrDestinationNotEmpty", err)
	}
}

func openDatabase(t *testing.T, name string) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+name+"?mode=memory&cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open %s database: %v", name, err)
	}
	t.Cleanup(func() {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	})
	return db
}

func seedLegacySource(db *gorm.DB) error {
	if err := db.Create(&category.Category{ID: 1, Name: "Utilities"}).Error; err != nil {
		return err
	}
	if err := db.Create(&expense.Expense{
		ID: 1, CategoryID: 1, Amount: 12.34, Description: "Internet", Date: "2026-01-02",
	}).Error; err != nil {
		return err
	}
	if err := db.Create(&budget.Budget{ID: 1, CategoryID: 1, Amount: 150}).Error; err != nil {
		return err
	}
	if err := db.Create(&goal.Goal{
		ID: 1, Name: "Emergency fund", TargetAmount: 1000, CurrentAmount: 50, TargetDate: "2027-01-01",
		CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}).Error; err != nil {
		return err
	}
	if err := db.Create(&goal.Goal{
		ID: 2, Name: "Travel", TargetAmount: 500, CurrentAmount: 0,
	}).Error; err != nil {
		return err
	}
	if err := db.Create(&investments.Investment{
		ID: 1, Amount: 25, InvestmentType: "stocks", Date: "2026-01-03",
	}).Error; err != nil {
		return err
	}
	if err := db.Create(&recurring.RecurringExpense{
		ID: 1, CategoryID: 1, Amount: 5, Description: "Coffee", Frequency: "weekly",
		StartDate: time.Date(2026, 1, 4, 0, 0, 0, 0, time.UTC),
	}).Error; err != nil {
		return err
	}
	return db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY, password_hash TEXT)").Error
}
