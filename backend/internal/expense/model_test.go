package expense

import (
	"errors"
	"testing"

	"github.com/girishsaraf/xpense-trackr/backend/internal/testutil"
	"gorm.io/gorm"
)

func newExpenseModelDB(t *testing.T) *gorm.DB {
	t.Helper()
	database := testutil.NewDB(t, &Expense{})
	SetDB(database)
	return database
}

func TestExpenseCRUD(t *testing.T) {
	newExpenseModelDB(t)

	first := &Expense{CategoryID: 1, Amount: 12.5, Description: "Coffee", Date: "2026-01-10"}
	second := &Expense{CategoryID: 2, Amount: 30, Description: "Train", Date: "2026-01-11"}
	if err := CreateExpense(first); err != nil {
		t.Fatalf("CreateExpense() error = %v", err)
	}
	if err := CreateExpense(second); err != nil {
		t.Fatalf("CreateExpense() error = %v", err)
	}

	expenses, err := GetAllExpenses()
	if err != nil {
		t.Fatalf("GetAllExpenses() error = %v", err)
	}
	if len(expenses) != 2 {
		t.Fatalf("GetAllExpenses() returned %d expenses, want 2", len(expenses))
	}

	got, err := GetExpenseByID(int(first.ID))
	if err != nil {
		t.Fatalf("GetExpenseByID() error = %v", err)
	}
	if got.Description != first.Description || got.Amount != first.Amount {
		t.Errorf("GetExpenseByID() = %+v, want %+v", got, first)
	}

	updated := &Expense{CategoryID: 3, Amount: 99.99, Description: "Updated", Date: "2026-02-01"}
	if err := UpdateExpense(int(first.ID), updated); err != nil {
		t.Fatalf("UpdateExpense() error = %v", err)
	}
	got, err = GetExpenseByID(int(first.ID))
	if err != nil {
		t.Fatalf("GetExpenseByID() after update error = %v", err)
	}
	if got.CategoryID != updated.CategoryID || got.Amount != updated.Amount || got.Description != updated.Description {
		t.Errorf("updated expense = %+v, want %+v", got, updated)
	}

	if err := DeleteExpense(int(second.ID)); err != nil {
		t.Fatalf("DeleteExpense() error = %v", err)
	}
	if _, err := GetExpenseByID(int(second.ID)); !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("GetExpenseByID() after delete error = %v, want ErrRecordNotFound", err)
	}
}

func TestGetExpenseByIDNotFound(t *testing.T) {
	newExpenseModelDB(t)

	if _, err := GetExpenseByID(404); !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("GetExpenseByID() error = %v, want ErrRecordNotFound", err)
	}
}

func TestUpdateExpenseNotFound(t *testing.T) {
	newExpenseModelDB(t)

	err := UpdateExpense(404, &Expense{Description: "Missing"})
	if !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("UpdateExpense() error = %v, want ErrRecordNotFound", err)
	}
}

func TestDeleteExpenseNotFound(t *testing.T) {
	newExpenseModelDB(t)

	if err := DeleteExpense(404); !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("DeleteExpense() error = %v, want ErrRecordNotFound", err)
	}
}
