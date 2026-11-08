package budget

import (
	"errors"
	"testing"

	"github.com/girishsaraf/xpense-trackr/backend/internal/testutil"
	"gorm.io/gorm"
)

func newBudgetModelDB(t *testing.T) *gorm.DB {
	t.Helper()
	database := testutil.NewDB(t, &Budget{})
	SetDB(database)
	return database
}

func TestBudgetCRUD(t *testing.T) {
	newBudgetModelDB(t)

	first := &Budget{CategoryID: 1, Amount: 500}
	second := &Budget{CategoryID: 2, Amount: 250}
	if err := CreateBudget(first); err != nil {
		t.Fatalf("CreateBudget() error = %v", err)
	}
	if err := CreateBudget(second); err != nil {
		t.Fatalf("CreateBudget() error = %v", err)
	}

	budgets, err := GetAllBudgets()
	if err != nil {
		t.Fatalf("GetAllBudgets() error = %v", err)
	}
	if len(budgets) != 2 {
		t.Fatalf("GetAllBudgets() returned %d budgets, want 2", len(budgets))
	}

	got, err := GetBudgetByID(int(first.ID))
	if err != nil {
		t.Fatalf("GetBudgetByID() error = %v", err)
	}
	if got.Amount != first.Amount {
		t.Errorf("GetBudgetByID() = %+v, want %+v", got, first)
	}

	if err := UpdateBudget(int(first.ID), &Budget{Amount: 750}); err != nil {
		t.Fatalf("UpdateBudget() error = %v", err)
	}
	got, err = GetBudgetByID(int(first.ID))
	if err != nil {
		t.Fatalf("GetBudgetByID() after update error = %v", err)
	}
	if got.Amount != 750 {
		t.Errorf("updated budget amount = %v, want 750", got.Amount)
	}

	if err := DeleteBudget(int(second.ID)); err != nil {
		t.Fatalf("DeleteBudget() error = %v", err)
	}
	if _, err := GetBudgetByID(int(second.ID)); !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("GetBudgetByID() after delete error = %v, want ErrRecordNotFound", err)
	}
}

func TestGetBudgetByIDNotFound(t *testing.T) {
	newBudgetModelDB(t)

	if _, err := GetBudgetByID(404); !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("GetBudgetByID() error = %v, want ErrRecordNotFound", err)
	}
}

func TestUpdateBudgetNotFound(t *testing.T) {
	newBudgetModelDB(t)

	if err := UpdateBudget(404, &Budget{Amount: 1}); !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("UpdateBudget() error = %v, want ErrRecordNotFound", err)
	}
}

func TestDeleteBudgetNotFound(t *testing.T) {
	newBudgetModelDB(t)

	if err := DeleteBudget(404); !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("DeleteBudget() error = %v, want ErrRecordNotFound", err)
	}
}
