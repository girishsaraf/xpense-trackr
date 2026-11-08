package investments

import (
	"errors"
	"testing"

	"github.com/girishsaraf/xpense-trackr/backend/internal/testutil"
	"gorm.io/gorm"
)

func newInvestmentModelDB(t *testing.T) *gorm.DB {
	t.Helper()
	database := testutil.NewDB(t, &Investment{})
	SetDB(database)
	return database
}

func TestInvestmentCRUD(t *testing.T) {
	newInvestmentModelDB(t)

	first := &Investment{Amount: 100, InvestmentType: "stocks", Date: "2026-01-10"}
	second := &Investment{Amount: 250, InvestmentType: "bonds", Date: "2026-01-11"}
	if err := CreateInvestment(first); err != nil {
		t.Fatalf("CreateInvestment() error = %v", err)
	}
	if err := CreateInvestment(second); err != nil {
		t.Fatalf("CreateInvestment() error = %v", err)
	}

	investments, err := GetAllInvestments()
	if err != nil {
		t.Fatalf("GetAllInvestments() error = %v", err)
	}
	if len(investments) != 2 {
		t.Fatalf("GetAllInvestments() returned %d investments, want 2", len(investments))
	}

	got, err := GetInvestmentByID(int(first.ID))
	if err != nil {
		t.Fatalf("GetInvestmentByID() error = %v", err)
	}
	if got.InvestmentType != first.InvestmentType || got.Amount != first.Amount {
		t.Errorf("GetInvestmentByID() = %+v, want %+v", got, first)
	}

	if err := UpdateInvestment(int(first.ID), &Investment{Amount: 175, InvestmentType: "savings"}); err != nil {
		t.Fatalf("UpdateInvestment() error = %v", err)
	}
	got, err = GetInvestmentByID(int(first.ID))
	if err != nil {
		t.Fatalf("GetInvestmentByID() after update error = %v", err)
	}
	if got.Amount != 175 || got.InvestmentType != "savings" {
		t.Errorf("updated investment = %+v", got)
	}

	if err := DeleteInvestment(int(second.ID)); err != nil {
		t.Fatalf("DeleteInvestment() error = %v", err)
	}
	if _, err := GetInvestmentByID(int(second.ID)); !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("GetInvestmentByID() after delete error = %v, want ErrRecordNotFound", err)
	}
}

func TestGetInvestmentByIDNotFound(t *testing.T) {
	newInvestmentModelDB(t)

	if _, err := GetInvestmentByID(404); !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("GetInvestmentByID() error = %v, want ErrRecordNotFound", err)
	}
}

func TestUpdateInvestmentNotFound(t *testing.T) {
	newInvestmentModelDB(t)

	if err := UpdateInvestment(404, &Investment{Amount: 1}); !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("UpdateInvestment() error = %v, want ErrRecordNotFound", err)
	}
}

func TestDeleteInvestmentNotFound(t *testing.T) {
	newInvestmentModelDB(t)

	if err := DeleteInvestment(404); !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("DeleteInvestment() error = %v, want ErrRecordNotFound", err)
	}
}
