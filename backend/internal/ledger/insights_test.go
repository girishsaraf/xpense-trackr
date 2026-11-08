package ledger

import (
	"testing"

	"github.com/girishsaraf/xpense-trackr/backend/internal/category"
)

func TestFinancialInsightsCompareEqualPeriodsWithoutMixingCurrencies(t *testing.T) {
	database := setupLedgerDB(t)
	account := &Account{Name: "Checking", Type: "checking", Currency: "USD"}
	euroAccount := &Account{Name: "Euro Checking", Type: "checking", Currency: "EUR"}
	food := &category.Category{Name: "Food"}
	utilities := &category.Category{Name: "Utilities"}
	for _, record := range []interface{}{account, euroAccount, food, utilities} {
		if err := database.Create(record).Error; err != nil {
			t.Fatalf("create fixture: %v", err)
		}
	}
	foodID, utilitiesID := food.ID, utilities.ID
	for _, transaction := range []Transaction{
		{AccountID: account.ID, CategoryID: &foodID, Type: "expense", AmountMinor: 10000, Currency: "USD", Date: "2026-09-10"},
		{AccountID: account.ID, CategoryID: &utilitiesID, Type: "expense", AmountMinor: 30000, Currency: "USD", Date: "2026-09-11"},
		{AccountID: account.ID, CategoryID: &foodID, Type: "expense", AmountMinor: 15000, Currency: "USD", Date: "2026-10-10"},
		{AccountID: euroAccount.ID, CategoryID: &foodID, Type: "expense", AmountMinor: 8000, Currency: "EUR", Date: "2026-10-10"},
	} {
		if err := CreateTransaction(&transaction); err != nil {
			t.Fatalf("create transaction: %v", err)
		}
	}

	report, err := GetFinancialInsights("2026-10-01", "2026-10-31")
	if err != nil {
		t.Fatalf("GetFinancialInsights() error = %v", err)
	}
	if report.PreviousFrom != "2026-08-31" || report.PreviousTo != "2026-09-30" {
		t.Fatalf("comparison range = %s to %s", report.PreviousFrom, report.PreviousTo)
	}
	if len(report.Insights) != 5 {
		t.Fatalf("insights = %+v, want USD comparisons and separate EUR spending insights", report.Insights)
	}
	foundOverall, foundFood, foundUtilities := false, false, false
	foundEuroSpending := false
	for _, insight := range report.Insights {
		if insight.Currency == "EUR" {
			foundEuroSpending = insight.CurrentAmountMinor == 8000
			continue
		}
		if insight.Currency != "USD" {
			t.Errorf("unexpected insight currency = %q", insight.Currency)
		}
		switch insight.Kind {
		case "spending_change":
			foundOverall = insight.ChangePercent == -62.5 && insight.CurrentAmountMinor == 15000
		case "category_spending_change":
			if insight.CategoryName == "Food" {
				foundFood = insight.ChangePercent == 50
			}
			if insight.CategoryName == "Utilities" {
				foundUtilities = insight.CurrentAmountMinor == 0 && insight.ChangePercent == -100
			}
		}
	}
	if !foundOverall || !foundFood || !foundUtilities || !foundEuroSpending {
		t.Errorf("insights missing expected comparisons: %+v", report.Insights)
	}
}

func TestFinancialInsightsRejectInvalidRange(t *testing.T) {
	setupLedgerDB(t)
	if _, err := GetFinancialInsights("2026-11-02", "2026-11-01"); err == nil {
		t.Fatal("expected reversed comparison range to fail")
	}
}
