package recurring

import (
	"errors"
	"testing"
	"time"

	"github.com/girishsaraf/xpense-trackr/backend/internal/expense"
	"github.com/girishsaraf/xpense-trackr/backend/internal/investments"
)

func TestRecurringExpenseCRUD(t *testing.T) {
	database := newRecurringTestDB(t)

	expense := &RecurringExpense{
		CategoryID:  1,
		Amount:      49.99,
		Description: "Internet",
		StartDate:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
		Frequency:   "monthly",
	}
	if err := CreateRecurringExpense(expense); err != nil {
		t.Fatalf("CreateRecurringExpense() error = %v", err)
	}

	expenses, err := GetAllRecurringExpenses()
	if err != nil {
		t.Fatalf("GetAllRecurringExpenses() error = %v", err)
	}
	if len(expenses) != 1 {
		t.Fatalf("GetAllRecurringExpenses() returned %d expenses, want 1", len(expenses))
	}

	got, err := GetRecurringExpenseByID(int(expense.ID))
	if err != nil {
		t.Fatalf("GetRecurringExpenseByID() error = %v", err)
	}
	if got.Description != expense.Description {
		t.Errorf("GetRecurringExpenseByID() = %+v, want %+v", got, expense)
	}

	if err := UpdateRecurringExpense(int(expense.ID), &RecurringExpense{Description: "Updated", Amount: 59.99}); err != nil {
		t.Fatalf("UpdateRecurringExpense() error = %v", err)
	}
	got, err = GetRecurringExpenseByID(int(expense.ID))
	if err != nil {
		t.Fatalf("GetRecurringExpenseByID() after update error = %v", err)
	}
	if got.Description != "Updated" || got.Amount != 59.99 {
		t.Errorf("updated recurring expense = %+v", got)
	}

	if err := DeleteRecurringExpense(int(expense.ID)); err != nil {
		t.Fatalf("DeleteRecurringExpense() error = %v", err)
	}
	if _, err := GetRecurringExpenseByID(int(expense.ID)); !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("GetRecurringExpenseByID() after delete error = %v, want ErrRecordNotFound", err)
	}

	var count int64
	if err := database.Model(&RecurringExpense{}).Count(&count).Error; err != nil {
		t.Fatalf("count query failed: %v", err)
	}
	if count != 0 {
		t.Errorf("recurring expense count = %d, want 0", count)
	}
}

func TestRecurringInvestmentCRUD(t *testing.T) {
	newRecurringTestDB(t)

	investment := &RecurringInvestment{
		Amount:         500,
		InvestmentType: "401k",
		Description:    "Contribution",
		StartDate:      time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:        time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
		Frequency:      "monthly",
	}
	if err := CreateRecurringInvestment(investment); err != nil {
		t.Fatalf("CreateRecurringInvestment() error = %v", err)
	}

	investmentsList, err := GetAllRecurringInvestments()
	if err != nil {
		t.Fatalf("GetAllRecurringInvestments() error = %v", err)
	}
	if len(investmentsList) != 1 {
		t.Fatalf("GetAllRecurringInvestments() returned %d investments, want 1", len(investmentsList))
	}

	got, err := GetRecurringInvestmentByID(int(investment.ID))
	if err != nil {
		t.Fatalf("GetRecurringInvestmentByID() error = %v", err)
	}
	if got.InvestmentType != investment.InvestmentType {
		t.Errorf("GetRecurringInvestmentByID() = %+v, want %+v", got, investment)
	}

	if err := UpdateRecurringInvestment(int(investment.ID), &RecurringInvestment{Description: "Updated", Amount: 750}); err != nil {
		t.Fatalf("UpdateRecurringInvestment() error = %v", err)
	}
	got, err = GetRecurringInvestmentByID(int(investment.ID))
	if err != nil {
		t.Fatalf("GetRecurringInvestmentByID() after update error = %v", err)
	}
	if got.Description != "Updated" || got.Amount != 750 {
		t.Errorf("updated recurring investment = %+v", got)
	}

	if err := DeleteRecurringInvestment(int(investment.ID)); err != nil {
		t.Fatalf("DeleteRecurringInvestment() error = %v", err)
	}
	if _, err := GetRecurringInvestmentByID(int(investment.ID)); !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("GetRecurringInvestmentByID() after delete error = %v, want ErrRecordNotFound", err)
	}
}

func TestRecurringNotFoundErrors(t *testing.T) {
	newRecurringTestDB(t)

	if _, err := GetRecurringExpenseByID(404); !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("GetRecurringExpenseByID() error = %v, want ErrRecordNotFound", err)
	}
	if err := UpdateRecurringExpense(404, &RecurringExpense{Description: "x"}); !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("UpdateRecurringExpense() error = %v, want ErrRecordNotFound", err)
	}
	if err := DeleteRecurringExpense(404); !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("DeleteRecurringExpense() error = %v, want ErrRecordNotFound", err)
	}
	if _, err := GetRecurringInvestmentByID(404); !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("GetRecurringInvestmentByID() error = %v, want ErrRecordNotFound", err)
	}
	if err := UpdateRecurringInvestment(404, &RecurringInvestment{Description: "x"}); !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("UpdateRecurringInvestment() error = %v, want ErrRecordNotFound", err)
	}
	if err := DeleteRecurringInvestment(404); !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("DeleteRecurringInvestment() error = %v, want ErrRecordNotFound", err)
	}
}

func TestCheckIfExpenseOrInvestmentIsDue(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name      string
		frequency string
		lastRun   time.Time
		wantIsDue bool
	}{
		{name: "daily overdue", frequency: "daily", lastRun: now.AddDate(0, 0, -2), wantIsDue: true},
		{name: "daily current", frequency: "daily", lastRun: now.Add(time.Hour), wantIsDue: false},
		{name: "weekly overdue", frequency: "weekly", lastRun: now.AddDate(0, 0, -8), wantIsDue: true},
		{name: "weekly current", frequency: "weekly", lastRun: now.Add(time.Hour), wantIsDue: false},
		{name: "biweekly overdue", frequency: "biweekly", lastRun: now.AddDate(0, 0, -15), wantIsDue: true},
		{name: "biweekly current", frequency: "biweekly", lastRun: now.Add(time.Hour), wantIsDue: false},
		{name: "monthly overdue", frequency: "monthly", lastRun: now.AddDate(0, -2, 0), wantIsDue: true},
		{name: "monthly current", frequency: "monthly", lastRun: now.Add(time.Hour), wantIsDue: false},
		{name: "unsupported frequency", frequency: "yearly", lastRun: now.AddDate(-1, 0, 0), wantIsDue: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := checkIfExpenseOrInvestmentIsDue(test.frequency, test.lastRun); got != test.wantIsDue {
				t.Errorf("checkIfExpenseOrInvestmentIsDue(%q, %v) = %v, want %v", test.frequency, test.lastRun, got, test.wantIsDue)
			}
		})
	}
}

func TestFetchRecurringExpenseInsertsDueExpenses(t *testing.T) {
	database := newRecurringTestDB(t)
	now := time.Now()
	due := &RecurringExpense{
		CategoryID:  1,
		Amount:      12,
		Description: "Due expense",
		StartDate:   now.AddDate(0, 0, -2),
		EndDate:     now.AddDate(1, 0, 0),
		Frequency:   "daily",
	}
	future := &RecurringExpense{
		CategoryID:  1,
		Amount:      20,
		Description: "Future expense",
		StartDate:   now.AddDate(0, 0, 2),
		EndDate:     now.AddDate(1, 0, 0),
		Frequency:   "daily",
	}
	if err := database.Create(due).Error; err != nil {
		t.Fatalf("failed to seed recurring expense: %v", err)
	}
	if err := database.Create(future).Error; err != nil {
		t.Fatalf("failed to seed recurring expense: %v", err)
	}

	if err := FetchRecurringExpense(); err != nil {
		t.Fatalf("FetchRecurringExpense() error = %v", err)
	}

	var inserted []expense.Expense
	if err := database.Find(&inserted).Error; err != nil {
		t.Fatalf("failed to query expenses: %v", err)
	}
	if len(inserted) != 1 {
		t.Fatalf("inserted %d expenses, want 1", len(inserted))
	}
	if inserted[0].Amount != due.Amount || inserted[0].CategoryID != due.CategoryID {
		t.Errorf("inserted expense = %+v, want amount %v category %v", inserted[0], due.Amount, due.CategoryID)
	}

	var updated RecurringExpense
	if err := database.First(&updated, due.ID).Error; err != nil {
		t.Fatalf("failed to query recurring expense: %v", err)
	}
	if !updated.StartDate.After(due.StartDate) {
		t.Errorf("recurring expense start date was not advanced: %v", updated.StartDate)
	}
}

func TestFetchRecurringInvestmentInsertsDueInvestments(t *testing.T) {
	database := newRecurringTestDB(t)
	now := time.Now()
	due := &RecurringInvestment{
		Amount:         500,
		InvestmentType: "401k",
		Description:    "Due investment",
		StartDate:      now.AddDate(0, 0, -2),
		EndDate:        now.AddDate(1, 0, 0),
		Frequency:      "daily",
	}
	future := &RecurringInvestment{
		Amount:         100,
		InvestmentType: "stocks",
		Description:    "Future investment",
		StartDate:      now.AddDate(0, 0, 2),
		EndDate:        now.AddDate(1, 0, 0),
		Frequency:      "daily",
	}
	if err := database.Create(due).Error; err != nil {
		t.Fatalf("failed to seed recurring investment: %v", err)
	}
	if err := database.Create(future).Error; err != nil {
		t.Fatalf("failed to seed recurring investment: %v", err)
	}

	if err := FetchRecurringInvestment(); err != nil {
		t.Fatalf("FetchRecurringInvestment() error = %v", err)
	}

	var inserted []investments.Investment
	if err := database.Find(&inserted).Error; err != nil {
		t.Fatalf("failed to query investments: %v", err)
	}
	if len(inserted) != 1 {
		t.Fatalf("inserted %d investments, want 1", len(inserted))
	}
	if inserted[0].Amount != due.Amount || inserted[0].InvestmentType != due.InvestmentType {
		t.Errorf("inserted investment = %+v", inserted[0])
	}
}
