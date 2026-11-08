package legacyimport

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
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

var (
	ErrDestinationNotEmpty = errors.New("legacy import requires an empty finance database")
	ErrImportAlreadyRun    = errors.New("this legacy source has already been imported")
	ErrNoLegacyData        = errors.New("legacy source contains no supported finance records")
)

type ImportRecord struct {
	ID         uint      `gorm:"primaryKey"`
	SourceName string    `gorm:"size:120;not null;uniqueIndex"`
	Currency   string    `gorm:"size:3;not null"`
	ImportedAt time.Time `gorm:"not null"`
	ReportJSON string    `gorm:"type:text;not null"`
}

func (ImportRecord) TableName() string {
	return "legacy_import_records"
}

type Report struct {
	SourceName           string   `json:"source_name"`
	Currency             string   `json:"currency"`
	Categories           int      `json:"categories"`
	Expenses             int      `json:"expenses"`
	ExpensesTotalMinor   int64    `json:"expenses_total_minor"`
	Budgets              int      `json:"budgets"`
	BudgetCategories     int      `json:"budget_categories"`
	Goals                int      `json:"goals"`
	Investments          int      `json:"investments"`
	RecurringExpenses    int      `json:"recurring_expenses"`
	RecurringInvestments int      `json:"recurring_investments"`
	UsersNotImported     bool     `json:"users_not_imported"`
	Applied              bool     `json:"applied"`
	Warnings             []string `json:"warnings,omitempty"`
}

type sourceData struct {
	categories           []category.Category
	expenses             []expense.Expense
	budgets              []budget.Budget
	goals                []goal.Goal
	investments          []investments.Investment
	recurringExpenses    []recurring.RecurringExpense
	recurringInvestments []recurring.RecurringInvestment
}

// Run validates a legacy database and returns a report without writing unless apply is true.
func Run(source, destination *gorm.DB, sourceName, currency string, apply bool) (Report, error) {
	sourceName = strings.TrimSpace(sourceName)
	currency = strings.ToUpper(strings.TrimSpace(currency))
	if sourceName == "" || len(sourceName) > 120 || !ledger.SupportedCurrency(currency) {
		return Report{}, errors.New("source name and supported currency are required")
	}
	data, err := readSource(source)
	if err != nil {
		return Report{}, err
	}
	report, categoryIDs, err := validateSource(source, data, sourceName, currency)
	if err != nil {
		return Report{}, err
	}
	if report.Categories+report.Expenses+report.Budgets+report.Goals+report.Investments+
		report.RecurringExpenses+report.RecurringInvestments == 0 {
		return Report{}, ErrNoLegacyData
	}
	if err := ensureDestinationReady(destination, sourceName); err != nil {
		return report, err
	}
	if !apply {
		return report, nil
	}
	if err := applySource(destination, sourceName, currency, data, categoryIDs, &report); err != nil {
		return report, err
	}
	report.Applied = true
	return report, nil
}

func readSource(source *gorm.DB) (sourceData, error) {
	var data sourceData
	for _, entry := range []struct {
		model interface{}
		dest  interface{}
	}{
		{&category.Category{}, &data.categories},
		{&expense.Expense{}, &data.expenses},
		{&budget.Budget{}, &data.budgets},
		{&goal.Goal{}, &data.goals},
		{&investments.Investment{}, &data.investments},
	} {
		if !source.Migrator().HasTable(entry.model) {
			continue
		}
		if err := source.Find(entry.dest).Error; err != nil {
			return sourceData{}, fmt.Errorf("read legacy table %T: %w", entry.model, err)
		}
	}
	for _, candidate := range []struct {
		name string
		dest interface{}
	}{
		{"recurring_expenses", &data.recurringExpenses},
		{"recurring_expense", &data.recurringExpenses},
	} {
		if source.Migrator().HasTable(candidate.name) {
			if err := source.Table(candidate.name).Find(candidate.dest).Error; err != nil {
				return sourceData{}, fmt.Errorf("read legacy table %s: %w", candidate.name, err)
			}
			break
		}
	}
	for _, candidate := range []struct {
		name string
		dest interface{}
	}{
		{"recurring_investments", &data.recurringInvestments},
		{"recurring_investment", &data.recurringInvestments},
	} {
		if source.Migrator().HasTable(candidate.name) {
			if err := source.Table(candidate.name).Find(candidate.dest).Error; err != nil {
				return sourceData{}, fmt.Errorf("read legacy table %s: %w", candidate.name, err)
			}
			break
		}
	}
	return data, nil
}

func validateSource(source *gorm.DB, data sourceData, sourceName, currency string) (Report, map[uint]uint, error) {
	report := Report{
		SourceName:           sourceName,
		Currency:             currency,
		Categories:           len(data.categories),
		Expenses:             len(data.expenses),
		Budgets:              len(data.budgets),
		Goals:                len(data.goals),
		Investments:          len(data.investments),
		RecurringExpenses:    len(data.recurringExpenses),
		RecurringInvestments: len(data.recurringInvestments),
		Warnings: []string{
			"Legacy amounts have no stored currency; all imported finance records are assigned " + currency + ".",
			"Legacy budgets are treated as monthly envelopes; duplicate category budgets are combined.",
			"Legacy investment entries remain contribution history and are not treated as current market valuations.",
			"Legacy recurring rules are preserved in their original tables; review them before relying on scheduled generation.",
		},
	}
	categoryIDs := make(map[uint]uint, len(data.categories))
	for _, item := range data.categories {
		if item.ID == 0 || strings.TrimSpace(item.Name) == "" {
			return Report{}, nil, fmt.Errorf("legacy category %d has an invalid ID or name", item.ID)
		}
		categoryIDs[item.ID] = item.ID
	}
	for _, item := range data.expenses {
		if item.ID == 0 || item.Amount <= 0 || !validDate(item.Date) {
			return Report{}, nil, fmt.Errorf("legacy expense %d has an invalid ID, amount, or date", item.ID)
		}
		if _, ok := categoryIDs[item.CategoryID]; !ok {
			return Report{}, nil, fmt.Errorf("legacy expense %d references missing category %d", item.ID, item.CategoryID)
		}
		amount, err := amountMinor(item.Amount, currency)
		if err != nil {
			return Report{}, nil, fmt.Errorf("legacy expense %d amount: %w", item.ID, err)
		}
		report.ExpensesTotalMinor += amount
	}
	budgetTotals := make(map[uint]int64)
	for _, item := range data.budgets {
		if item.ID == 0 || item.Amount <= 0 {
			return Report{}, nil, fmt.Errorf("legacy budget %d has an invalid ID or amount", item.ID)
		}
		if _, ok := categoryIDs[item.CategoryID]; !ok {
			return Report{}, nil, fmt.Errorf("legacy budget %d references missing category %d", item.ID, item.CategoryID)
		}
		amount, err := amountMinor(item.Amount, currency)
		if err != nil {
			return Report{}, nil, fmt.Errorf("legacy budget %d amount: %w", item.ID, err)
		}
		budgetTotals[item.CategoryID] += amount
	}
	report.BudgetCategories = len(budgetTotals)
	for _, item := range data.goals {
		if item.ID == 0 || strings.TrimSpace(item.Name) == "" || item.TargetAmount <= 0 ||
			item.CurrentAmount < 0 || !validOptionalDate(item.TargetDate) {
			return Report{}, nil, fmt.Errorf("legacy goal %d has invalid values", item.ID)
		}
		if _, err := amountMinor(item.TargetAmount, currency); err != nil {
			return Report{}, nil, fmt.Errorf("legacy goal %d target: %w", item.ID, err)
		}
		if item.CurrentAmount > 0 {
			if _, err := amountMinor(item.CurrentAmount, currency); err != nil {
				return Report{}, nil, fmt.Errorf("legacy goal %d balance: %w", item.ID, err)
			}
		}
	}
	for _, item := range data.investments {
		if item.ID == 0 || item.Amount <= 0 || !validDate(item.Date) {
			return Report{}, nil, fmt.Errorf("legacy investment %d has an invalid ID, amount, or date", item.ID)
		}
		if _, err := amountMinor(item.Amount, currency); err != nil {
			return Report{}, nil, fmt.Errorf("legacy investment %d amount: %w", item.ID, err)
		}
	}
	report.UsersNotImported = source.Migrator().HasTable("users")
	if report.UsersNotImported {
		report.Warnings = append(report.Warnings, "Legacy user accounts and password data are intentionally not imported; this app uses local single-user access.")
	}
	return report, categoryIDs, nil
}

func ensureDestinationReady(destination *gorm.DB, sourceName string) error {
	if !destination.Migrator().HasTable(&ImportRecord{}) {
		return errors.New("destination is not initialized; start the app once to create its schema, then stop it before importing")
	}
	var existing ImportRecord
	err := destination.Where("source_name = ?", sourceName).First(&existing).Error
	if err == nil {
		return ErrImportAlreadyRun
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	for _, model := range []interface{}{
		&category.Category{}, &expense.Expense{}, &budget.Budget{}, &goal.Goal{},
		&investments.Investment{}, &recurring.RecurringExpense{}, &recurring.RecurringInvestment{},
		&ledger.Account{}, &ledger.Transaction{}, &ledger.ManualAsset{}, &ledger.AssetValuation{},
		&planning.Goal{}, &planning.Envelope{},
	} {
		var count int64
		if !destination.Migrator().HasTable(model) {
			return fmt.Errorf("destination schema is incomplete; missing table for %T", model)
		}
		if err := destination.Model(model).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return ErrDestinationNotEmpty
		}
	}
	return nil
}

func applySource(destination *gorm.DB, sourceName, currency string, data sourceData, categoryIDs map[uint]uint, report *Report) error {
	return destination.Transaction(func(tx *gorm.DB) error {
		if err := ensureDestinationReady(tx, sourceName); err != nil {
			return err
		}
		if len(data.categories) > 0 {
			if err := tx.Create(&data.categories).Error; err != nil {
				return fmt.Errorf("import categories: %w", err)
			}
		}
		if len(data.expenses) > 0 {
			if err := tx.Create(&data.expenses).Error; err != nil {
				return fmt.Errorf("import legacy expenses: %w", err)
			}
			account := ledger.Account{Name: "Imported checking", Type: "checking", Currency: currency}
			if err := tx.Create(&account).Error; err != nil {
				return fmt.Errorf("create imported account: %w", err)
			}
			transactions := make([]ledger.Transaction, 0, len(data.expenses))
			for _, old := range data.expenses {
				categoryID := categoryIDs[old.CategoryID]
				amount, err := amountMinor(old.Amount, currency)
				if err != nil {
					return err
				}
				description := strings.TrimSpace(old.Description)
				transactions = append(transactions, ledger.Transaction{
					AccountID: account.ID, CategoryID: &categoryID, Type: "expense", AmountMinor: amount,
					Currency: currency, Payee: description, Description: description, Date: old.Date,
				})
			}
			if err := tx.Create(&transactions).Error; err != nil {
				return fmt.Errorf("import ledger expenses: %w", err)
			}
		}
		if len(data.budgets) > 0 {
			if err := tx.Create(&data.budgets).Error; err != nil {
				return fmt.Errorf("import legacy budgets: %w", err)
			}
			envelopes := make(map[uint]*planning.Envelope, report.BudgetCategories)
			for _, old := range data.budgets {
				envelope := envelopes[old.CategoryID]
				if envelope == nil {
					envelope = &planning.Envelope{CategoryID: old.CategoryID, Currency: currency}
					envelopes[old.CategoryID] = envelope
				}
				amount, err := amountMinor(old.Amount, currency)
				if err != nil {
					return err
				}
				envelope.MonthlyLimitMinor += amount
			}
			for _, envelope := range envelopes {
				if err := tx.Create(envelope).Error; err != nil {
					return fmt.Errorf("import monthly budget envelope: %w", err)
				}
			}
		}
		if len(data.goals) > 0 {
			if err := tx.Create(&data.goals).Error; err != nil {
				return fmt.Errorf("import legacy goals: %w", err)
			}
			for _, old := range data.goals {
				target, err := amountMinor(old.TargetAmount, currency)
				if err != nil {
					return err
				}
				var current int64
				if old.CurrentAmount > 0 {
					current, err = amountMinor(old.CurrentAmount, currency)
					if err != nil {
						return err
					}
				}
				newGoal := planning.Goal{Name: old.Name, Currency: currency, TargetAmountMinor: target, TargetDate: old.TargetDate}
				if err := tx.Create(&newGoal).Error; err != nil {
					return fmt.Errorf("import planning goal %q: %w", old.Name, err)
				}
				if current > 0 {
					contributionDate := time.Now().UTC().Format("2006-01-02")
					if !old.CreatedAt.IsZero() {
						contributionDate = old.CreatedAt.UTC().Format("2006-01-02")
					}
					contribution := planning.GoalContribution{
						GoalID: newGoal.ID, AmountMinor: current, Direction: "contribution",
						Date: contributionDate, Note: "Imported legacy goal balance",
					}
					if err := tx.Create(&contribution).Error; err != nil {
						return fmt.Errorf("import goal balance for %q: %w", old.Name, err)
					}
				}
			}
		}
		if len(data.investments) > 0 {
			if err := tx.Create(&data.investments).Error; err != nil {
				return fmt.Errorf("import legacy investments: %w", err)
			}
		}
		if len(data.recurringExpenses) > 0 {
			if err := tx.Create(&data.recurringExpenses).Error; err != nil {
				return fmt.Errorf("import recurring expenses: %w", err)
			}
		}
		if len(data.recurringInvestments) > 0 {
			if err := tx.Create(&data.recurringInvestments).Error; err != nil {
				return fmt.Errorf("import recurring investments: %w", err)
			}
		}
		recordReport := *report
		recordReport.Applied = true
		reportJSON, err := json.Marshal(recordReport)
		if err != nil {
			return err
		}
		record := ImportRecord{
			SourceName: sourceName, Currency: currency, ImportedAt: time.Now().UTC(), ReportJSON: string(reportJSON),
		}
		if err := tx.Create(&record).Error; err != nil {
			return fmt.Errorf("record legacy import: %w", err)
		}
		return nil
	})
}

func amountMinor(amount float64, currency string) (int64, error) {
	if math.IsNaN(amount) || math.IsInf(amount, 0) || amount <= 0 {
		return 0, errors.New("amount must be finite and positive")
	}
	minor, err := ledger.ParseAmountMinor(strconv.FormatFloat(amount, 'f', -1, 64), currency)
	if err != nil {
		return 0, errors.New("amount has unsupported precision for the selected currency")
	}
	return minor, nil
}

func validDate(value string) bool {
	parsed, err := time.Parse("2006-01-02", value)
	return err == nil && parsed.Format("2006-01-02") == value
}

func validOptionalDate(value string) bool {
	return value == "" || validDate(value)
}
