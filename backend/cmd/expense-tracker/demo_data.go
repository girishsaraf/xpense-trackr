package main

import (
	"fmt"
	"time"

	"github.com/girishsaraf/xpense-trackr/backend/internal/budget"
	"github.com/girishsaraf/xpense-trackr/backend/internal/category"
	"github.com/girishsaraf/xpense-trackr/backend/internal/expense"
	"github.com/girishsaraf/xpense-trackr/backend/internal/goal"
	"github.com/girishsaraf/xpense-trackr/backend/internal/investments"
	"github.com/girishsaraf/xpense-trackr/backend/internal/ledger"
	"github.com/girishsaraf/xpense-trackr/backend/internal/planning"
	"github.com/girishsaraf/xpense-trackr/backend/internal/recurring"
	"gorm.io/gorm"
)

func seedDemoData(db *gorm.DB) error {
	var existingRecords int64
	if err := db.Model(&category.Category{}).Count(&existingRecords).Error; err != nil {
		return fmt.Errorf("count existing categories: %w", err)
	}
	if existingRecords > 0 {
		return nil
	}

	now := time.Now().UTC()
	date := func(days int) string {
		return now.AddDate(0, 0, days).Format("2006-01-02")
	}
	monthDate := func(months, day int) string {
		value := now.AddDate(0, months, 0)
		if day > 28 {
			day = 28
		}
		return time.Date(value.Year(), value.Month(), day, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	}

	return db.Transaction(func(tx *gorm.DB) error {
		categories := []category.Category{
			{Name: "Groceries"},
			{Name: "Dining"},
			{Name: "Housing"},
			{Name: "Utilities"},
			{Name: "Transport"},
			{Name: "Health"},
			{Name: "Entertainment"},
			{Name: "Shopping"},
		}
		if err := tx.Create(&categories).Error; err != nil {
			return fmt.Errorf("create demo categories: %w", err)
		}

		expenses := []expense.Expense{
			{CategoryID: categories[0].ID, Amount: 82.45, Description: "Weekly grocery run", Date: date(-2)},
			{CategoryID: categories[1].ID, Amount: 36.20, Description: "Dinner with friends", Date: date(-4)},
			{CategoryID: categories[3].ID, Amount: 114.60, Description: "Electricity bill", Date: date(-7)},
			{CategoryID: categories[4].ID, Amount: 48.00, Description: "Fuel", Date: date(-10)},
			{CategoryID: categories[0].ID, Amount: 96.30, Description: "Monthly groceries", Date: monthDate(-1, 12)},
			{CategoryID: categories[6].ID, Amount: 24.99, Description: "Streaming subscription", Date: monthDate(-1, 19)},
			{CategoryID: categories[5].ID, Amount: 65.00, Description: "Pharmacy", Date: monthDate(-2, 8)},
			{CategoryID: categories[7].ID, Amount: 120.00, Description: "Home supplies", Date: monthDate(-3, 16)},
		}
		if err := tx.Create(&expenses).Error; err != nil {
			return fmt.Errorf("create demo expenses: %w", err)
		}

		budgets := []budget.Budget{
			{CategoryID: categories[0].ID, Amount: 500},
			{CategoryID: categories[1].ID, Amount: 250},
			{CategoryID: categories[3].ID, Amount: 300},
			{CategoryID: categories[4].ID, Amount: 180},
			{CategoryID: categories[6].ID, Amount: 150},
		}
		if err := tx.Create(&budgets).Error; err != nil {
			return fmt.Errorf("create demo budgets: %w", err)
		}

		investmentRecords := []investments.Investment{
			{Amount: 400, InvestmentType: "stocks", Date: monthDate(-4, 1)},
			{Amount: 400, InvestmentType: "stocks", Date: monthDate(-3, 1)},
			{Amount: 400, InvestmentType: "stocks", Date: monthDate(-2, 1)},
			{Amount: 400, InvestmentType: "stocks", Date: monthDate(-1, 1)},
			{Amount: 400, InvestmentType: "stocks", Date: monthDate(0, 1)},
			{Amount: 250, InvestmentType: "401k", Date: monthDate(-2, 15)},
			{Amount: 250, InvestmentType: "401k", Date: monthDate(-1, 15)},
			{Amount: 250, InvestmentType: "401k", Date: monthDate(0, 15)},
		}
		if err := tx.Create(&investmentRecords).Error; err != nil {
			return fmt.Errorf("create demo investments: %w", err)
		}

		legacyGoals := []goal.Goal{
			{Name: "Emergency fund", TargetAmount: 10000, CurrentAmount: 4200, TargetDate: date(365)},
			{Name: "Summer trip", TargetAmount: 2500, CurrentAmount: 900, TargetDate: date(180)},
		}
		if err := tx.Create(&legacyGoals).Error; err != nil {
			return fmt.Errorf("create demo goals: %w", err)
		}

		recurringExpenses := []recurring.RecurringExpense{
			{CategoryID: categories[3].ID, Amount: 65, Description: "Internet", StartDate: now.AddDate(0, -3, 0), Frequency: "monthly"},
			{CategoryID: categories[6].ID, Amount: 14.99, Description: "Music subscription", StartDate: now.AddDate(0, -2, 0), Frequency: "monthly"},
		}
		if err := tx.Create(&recurringExpenses).Error; err != nil {
			return fmt.Errorf("create demo recurring expenses: %w", err)
		}
		recurringInvestments := []recurring.RecurringInvestment{
			{Amount: 400, InvestmentType: "stocks", Description: "Index fund contribution", StartDate: now.AddDate(0, -4, 0), Frequency: "monthly"},
		}
		if err := tx.Create(&recurringInvestments).Error; err != nil {
			return fmt.Errorf("create demo recurring investments: %w", err)
		}

		accounts := []ledger.Account{
			{Name: "Everyday Checking", Type: "checking", Currency: "USD", OpeningBalanceMinor: 285000},
			{Name: "Rainy Day Savings", Type: "savings", Currency: "USD", OpeningBalanceMinor: 640000},
			{Name: "Travel Card", Type: "credit_card", Currency: "USD", OpeningBalanceMinor: -32500},
		}
		if err := tx.Create(&accounts).Error; err != nil {
			return fmt.Errorf("create demo accounts: %w", err)
		}

		rules := []ledger.CategoryRule{
			{Name: "Groceries", MatchValue: "Market", CategoryID: categories[0].ID, Priority: 20, Enabled: true},
			{Name: "Dining", MatchValue: "Cafe", CategoryID: categories[1].ID, Priority: 10, Enabled: true},
		}
		if err := tx.Create(&rules).Error; err != nil {
			return fmt.Errorf("create demo category rules: %w", err)
		}

		transactions := []ledger.Transaction{
			{AccountID: accounts[0].ID, Type: "income", AmountMinor: 420000, Currency: "USD", Payee: "Acme Payroll", Description: "Paycheck", Date: monthDate(-3, 1)},
			{AccountID: accounts[0].ID, Type: "expense", AmountMinor: 142500, Currency: "USD", Payee: "Apartment Management", Description: "Rent", Date: monthDate(-3, 2), CategoryID: &categories[2].ID},
			{AccountID: accounts[0].ID, Type: "expense", AmountMinor: 8250, Currency: "USD", Payee: "Neighborhood Market", Description: "Groceries", Date: monthDate(-2, 8), CategoryID: &categories[0].ID},
			{AccountID: accounts[0].ID, Type: "income", AmountMinor: 420000, Currency: "USD", Payee: "Acme Payroll", Description: "Paycheck", Date: monthDate(-2, 1)},
			{AccountID: accounts[0].ID, Type: "expense", AmountMinor: 142500, Currency: "USD", Payee: "Apartment Management", Description: "Rent", Date: monthDate(-2, 2), CategoryID: &categories[2].ID},
			{AccountID: accounts[0].ID, Type: "expense", AmountMinor: 10350, Currency: "USD", Payee: "Neighborhood Market", Description: "Groceries", Date: monthDate(-1, 8), CategoryID: &categories[0].ID},
			{AccountID: accounts[0].ID, Type: "income", AmountMinor: 420000, Currency: "USD", Payee: "Acme Payroll", Description: "Paycheck", Date: monthDate(-1, 1)},
			{AccountID: accounts[0].ID, Type: "expense", AmountMinor: 142500, Currency: "USD", Payee: "Apartment Management", Description: "Rent", Date: monthDate(-1, 2), CategoryID: &categories[2].ID},
			{AccountID: accounts[0].ID, Type: "expense", AmountMinor: 11975, Currency: "USD", Payee: "Neighborhood Market", Description: "Groceries", Date: date(-3), CategoryID: &categories[0].ID},
			{AccountID: accounts[0].ID, Type: "expense", AmountMinor: 4780, Currency: "USD", Payee: "Corner Cafe", Description: "Lunch", Date: date(-5), CategoryID: &categories[1].ID},
			{AccountID: accounts[0].ID, Type: "expense", AmountMinor: 7125, Currency: "USD", Payee: "Power Company", Description: "Utilities", Date: date(-8), CategoryID: &categories[3].ID},
			{AccountID: accounts[0].ID, TransferAccountID: &accounts[1].ID, Type: "transfer", AmountMinor: 60000, Currency: "USD", Payee: "Savings transfer", Description: "Emergency fund", Date: date(-12)},
			{AccountID: accounts[2].ID, Type: "expense", AmountMinor: 8950, Currency: "USD", Payee: "Airline", Description: "Flight deposit", Date: date(-18), CategoryID: &categories[7].ID},
		}
		if err := tx.Create(&transactions).Error; err != nil {
			return fmt.Errorf("create demo ledger transactions: %w", err)
		}

		flows := []ledger.ScheduledFlow{
			{AccountID: accounts[0].ID, CategoryID: &categories[2].ID, Label: "Rent", Type: "expense", AmountMinor: 142500, Currency: "USD", Frequency: "monthly", NextDate: monthDate(1, 2), ReminderDays: 5, Enabled: true},
			{AccountID: accounts[0].ID, CategoryID: &categories[3].ID, Label: "Internet", Type: "expense", AmountMinor: 6500, Currency: "USD", Frequency: "monthly", NextDate: date(6), ReminderDays: 3, Enabled: true},
			{AccountID: accounts[0].ID, Label: "Paycheck", Type: "income", AmountMinor: 420000, Currency: "USD", Frequency: "biweekly", NextDate: date(8), ReminderDays: 0, Enabled: true},
		}
		if err := tx.Create(&flows).Error; err != nil {
			return fmt.Errorf("create demo scheduled flows: %w", err)
		}

		planningGoals := []planning.Goal{
			{Name: "House down payment", Currency: "USD", TargetAmountMinor: 6000000, TargetDate: date(730)},
			{Name: "New laptop", Currency: "USD", TargetAmountMinor: 220000, TargetDate: date(150)},
		}
		if err := tx.Create(&planningGoals).Error; err != nil {
			return fmt.Errorf("create demo planning goals: %w", err)
		}
		contributions := []planning.GoalContribution{
			{GoalID: planningGoals[0].ID, AmountMinor: 1450000, Direction: "contribution", Date: monthDate(-3, 15), Note: "Initial savings"},
			{GoalID: planningGoals[0].ID, AmountMinor: 60000, Direction: "contribution", Date: date(-12), Note: "Monthly contribution"},
			{GoalID: planningGoals[1].ID, AmountMinor: 85000, Direction: "contribution", Date: monthDate(-1, 20), Note: "Set aside bonus"},
		}
		if err := tx.Create(&contributions).Error; err != nil {
			return fmt.Errorf("create demo goal contributions: %w", err)
		}
		envelopes := []planning.Envelope{
			{CategoryID: categories[0].ID, Currency: "USD", MonthlyLimitMinor: 50000, Rollover: true},
			{CategoryID: categories[1].ID, Currency: "USD", MonthlyLimitMinor: 25000, Rollover: false},
			{CategoryID: categories[6].ID, Currency: "USD", MonthlyLimitMinor: 15000, Rollover: true},
		}
		if err := tx.Create(&envelopes).Error; err != nil {
			return fmt.Errorf("create demo budget envelopes: %w", err)
		}

		assets := []ledger.ManualAsset{
			{Name: "Car", Kind: "asset", Currency: "USD"},
			{Name: "Student loan", Kind: "liability", Currency: "USD"},
		}
		if err := tx.Create(&assets).Error; err != nil {
			return fmt.Errorf("create demo manual assets: %w", err)
		}
		assetValuations := []ledger.AssetValuation{
			{AssetID: assets[0].ID, ValueMinor: 1850000, Date: monthDate(-3, 1), Source: "manual"},
			{AssetID: assets[0].ID, ValueMinor: 1780000, Date: date(-1), Source: "manual"},
			{AssetID: assets[1].ID, ValueMinor: 920000, Date: monthDate(-3, 1), Source: "manual"},
			{AssetID: assets[1].ID, ValueMinor: 845000, Date: date(-1), Source: "manual"},
		}
		if err := tx.Create(&assetValuations).Error; err != nil {
			return fmt.Errorf("create demo asset valuations: %w", err)
		}

		holdings := []ledger.PortfolioHolding{
			{Name: "Vanguard Total Stock Market ETF", Ticker: "VTI", AssetClass: "equity", Currency: "USD", Quantity: "18.5", CostBasisMinor: 420000},
			{Name: "US Treasury Bond Fund", Ticker: "BND", AssetClass: "fixed_income", Currency: "USD", Quantity: "40", CostBasisMinor: 290000},
		}
		if err := tx.Create(&holdings).Error; err != nil {
			return fmt.Errorf("create demo portfolio holdings: %w", err)
		}
		holdingValuations := []ledger.HoldingValuation{
			{HoldingID: holdings[0].ID, ValueMinor: 435000, Date: monthDate(-2, 1), Source: "manual"},
			{HoldingID: holdings[0].ID, ValueMinor: 468000, Date: date(-1), Source: "manual"},
			{HoldingID: holdings[1].ID, ValueMinor: 294000, Date: monthDate(-2, 1), Source: "manual"},
			{HoldingID: holdings[1].ID, ValueMinor: 301000, Date: date(-1), Source: "manual"},
		}
		if err := tx.Create(&holdingValuations).Error; err != nil {
			return fmt.Errorf("create demo holding valuations: %w", err)
		}
		return nil
	})
}
