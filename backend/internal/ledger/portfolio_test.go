package ledger

import (
	"errors"
	"testing"
	"time"
)

func TestPortfolioHoldingValuationsFlowIntoHistoricalNetWorth(t *testing.T) {
	database := setupLedgerDB(t)
	now := time.Now().UTC()
	yesterday := now.AddDate(0, 0, -1)
	dateBefore := yesterday.Format("2006-01-02")
	dateToday := now.Format("2006-01-02")

	holding, err := CreatePortfolioHolding(CreatePortfolioHoldingRequest{
		Name: "Index fund", Ticker: "vti", AssetClass: "equity", Currency: "USD",
		Quantity: "2.125", CostBasisMinor: 20000, ValueMinor: 24000,
		Date: dateBefore, Source: "broker statement",
	})
	if err != nil {
		t.Fatalf("CreatePortfolioHolding() error = %v", err)
	}
	if holding.Ticker != "VTI" || !holding.HasValuation {
		t.Fatalf("created holding = %+v", holding)
	}
	if err := database.Model(holding).Update("created_at", yesterday.Add(-time.Hour)).Error; err != nil {
		t.Fatalf("set historical creation date: %v", err)
	}
	if err := CreateHoldingValuation(holding.ID, &HoldingValuation{
		ValueMinor: 27500, Date: dateToday, Source: "broker statement",
	}); err != nil {
		t.Fatalf("CreateHoldingValuation() error = %v", err)
	}

	before, err := GetNetWorthReport(yesterday)
	if err != nil {
		t.Fatalf("GetNetWorthReport(yesterday) error = %v", err)
	}
	if len(before.Currencies) != 1 || before.Currencies[0].PortfolioMinor != 24000 ||
		before.Currencies[0].NetWorthMinor != 24000 {
		t.Fatalf("historical net worth = %+v, want USD 24000", before.Currencies)
	}
	after, err := GetNetWorthReport(now)
	if err != nil {
		t.Fatalf("GetNetWorthReport(now) error = %v", err)
	}
	if len(after.Currencies) != 1 || after.Currencies[0].PortfolioMinor != 27500 ||
		after.Currencies[0].NetWorthMinor != 27500 || len(after.Currencies[0].Holdings) != 1 {
		t.Fatalf("current net worth = %+v, want USD 27500 and holding detail", after.Currencies)
	}

	holdings, err := GetPortfolioHoldings()
	if err != nil || len(holdings) != 1 || holdings[0].LatestValueMinor != 27500 {
		t.Fatalf("GetPortfolioHoldings() = %+v, %v", holdings, err)
	}
	valuations, err := GetHoldingValuations(holding.ID)
	if err != nil || len(valuations) != 2 || valuations[0].Date != dateToday {
		t.Fatalf("GetHoldingValuations() = %+v, %v", valuations, err)
	}
}

func TestPortfolioHoldingValidationAndDeletion(t *testing.T) {
	setupLedgerDB(t)
	for _, invalid := range []CreatePortfolioHoldingRequest{
		{Name: "Bad quantity", AssetClass: "equity", Currency: "USD", Quantity: "0", Date: "2026-10-01"},
		{Name: "Bad precision", AssetClass: "equity", Currency: "USD", Quantity: "1.123456789", Date: "2026-10-01"},
		{Name: "Bad class", AssetClass: "unknown", Currency: "USD", Quantity: "1", Date: "2026-10-01"},
		{Name: "Bad currency", AssetClass: "equity", Currency: "ZZZ", Quantity: "1", Date: "2026-10-01"},
		{Name: "Negative value", AssetClass: "equity", Currency: "USD", Quantity: "1", ValueMinor: -1, Date: "2026-10-01"},
	} {
		invalid.Name = "Invalid holding"
		if _, err := CreatePortfolioHolding(invalid); !errors.Is(err, ErrInvalidPortfolioHolding) {
			t.Errorf("CreatePortfolioHolding(%+v) error = %v, want validation error", invalid, err)
		}
	}

	holding, err := CreatePortfolioHolding(CreatePortfolioHoldingRequest{
		Name: "Bond", AssetClass: "fixed_income", Currency: "USD", Quantity: "1",
		Date: "2026-10-01", ValueMinor: 10000,
	})
	if err != nil {
		t.Fatalf("CreatePortfolioHolding() error = %v", err)
	}
	if err := DeletePortfolioHolding(holding.ID); err != nil {
		t.Fatalf("DeletePortfolioHolding() error = %v", err)
	}
	if _, err := GetHoldingValuations(holding.ID); err == nil {
		t.Fatal("GetHoldingValuations() after delete should fail")
	}
}
