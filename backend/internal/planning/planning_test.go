package planning

import (
	"errors"
	"testing"
	"time"

	"github.com/girishsaraf/xpense-trackr/backend/internal/category"
	"github.com/girishsaraf/xpense-trackr/backend/internal/ledger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupPlanningDB(t *testing.T) (*gorm.DB, uint) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&category.Category{}, &ledger.Transaction{}, &Goal{}, &GoalContribution{}, &Envelope{}); err != nil {
		t.Fatal(err)
	}
	SetDB(db)
	cat := category.Category{Name: "Food"}
	if err := db.Create(&cat).Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	})
	return db, cat.ID
}

func TestGoalContributionsTrackBalanceAndPreventOverdraft(t *testing.T) {
	_, _ = setupPlanningDB(t)
	goal := Goal{Name: "Emergency fund", Currency: "USD", TargetAmountMinor: 10000, TargetDate: "2027-01-01"}
	if err := CreateGoal(&goal); err != nil {
		t.Fatal(err)
	}
	entry := GoalContribution{AmountMinor: 4000, Direction: "contribution", Date: "2026-01-01"}
	if err := AddContribution(goal.ID, &entry); err != nil {
		t.Fatal(err)
	}
	entry = GoalContribution{AmountMinor: 1000, Direction: "withdrawal", Date: "2026-01-02"}
	if err := AddContribution(goal.ID, &entry); err != nil {
		t.Fatal(err)
	}
	entry = GoalContribution{AmountMinor: 4000, Direction: "withdrawal", Date: "2026-01-03"}
	if err := AddContribution(goal.ID, &entry); !errors.Is(err, ErrInsufficientGoalFunds) {
		t.Fatalf("oversized withdrawal error = %v, want insufficient funds", err)
	}

	goals, err := GetGoals()
	if err != nil {
		t.Fatal(err)
	}
	if len(goals) != 1 || goals[0].CurrentAmountMinor != 3000 || goals[0].RemainingMinor != 7000 {
		t.Fatalf("goal progress = %+v", goals)
	}
}

func TestEnvelopeRolloverUsesPriorMonthsAndCurrentSpend(t *testing.T) {
	db, categoryID := setupPlanningDB(t)
	envelope := Envelope{
		CategoryID: categoryID, Currency: "USD", MonthlyLimitMinor: 10000, Rollover: true,
		CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	if err := CreateEnvelope(&envelope); err != nil {
		t.Fatal(err)
	}
	if err := db.Model(&envelope).Update("created_at", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)).Error; err != nil {
		t.Fatal(err)
	}
	for _, entry := range []struct {
		amount int64
		date   string
	}{
		{amount: 7000, date: "2026-01-15"},
		{amount: 4000, date: "2026-02-15"},
		{amount: 3000, date: "2026-03-15"},
	} {
		transaction := ledger.Transaction{
			CategoryID: &categoryID, Type: "expense", AmountMinor: entry.amount,
			Currency: "USD", Date: entry.date,
		}
		if err := db.Create(&transaction).Error; err != nil {
			t.Fatal(err)
		}
	}

	summaries, err := GetEnvelopes("2026-03")
	if err != nil {
		t.Fatal(err)
	}
	if len(summaries) != 1 {
		t.Fatalf("got %d envelope summaries, want 1", len(summaries))
	}
	summary := summaries[0]
	if summary.SpentMinor != 3000 || summary.RolloverMinor != 9000 || summary.AvailableMinor != 16000 {
		t.Fatalf("envelope summary = %+v", summary)
	}
	if _, err := GetEnvelopes("2026-13"); !errors.Is(err, ErrInvalidPeriod) {
		t.Fatalf("invalid period error = %v, want ErrInvalidPeriod", err)
	}
}
