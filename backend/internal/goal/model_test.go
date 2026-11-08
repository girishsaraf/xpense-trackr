package goal

import (
	"errors"
	"testing"

	"github.com/girishsaraf/xpense-trackr/backend/internal/testutil"
	"gorm.io/gorm"
)

func setupGoalDB(t *testing.T) *gorm.DB {
	t.Helper()
	database := testutil.NewDB(t, &Goal{})
	SetDB(database)
	return database
}

func TestGoalCRUD(t *testing.T) {
	setupGoalDB(t)
	first := &Goal{Name: "Emergency fund", TargetAmount: 5000, CurrentAmount: 1250, TargetDate: "2027-01-01"}
	second := &Goal{Name: "Vacation", TargetAmount: 1500}

	if err := CreateGoal(first); err != nil {
		t.Fatalf("CreateGoal() error = %v", err)
	}
	if err := CreateGoal(second); err != nil {
		t.Fatalf("CreateGoal() error = %v", err)
	}
	goals, err := GetAllGoals()
	if err != nil || len(goals) != 2 {
		t.Fatalf("GetAllGoals() = %d goals, error = %v; want 2 goals", len(goals), err)
	}

	got, err := GetGoalByID(int(first.ID))
	if err != nil || got.CurrentAmount != 1250 {
		t.Fatalf("GetGoalByID() = %+v, error = %v", got, err)
	}
	if err := UpdateGoal(int(first.ID), &Goal{Name: "Emergency fund", TargetAmount: 5000, CurrentAmount: 0}); err != nil {
		t.Fatalf("UpdateGoal() error = %v", err)
	}
	got, err = GetGoalByID(int(first.ID))
	if err != nil || got.CurrentAmount != 0 || got.TargetDate != "" {
		t.Fatalf("UpdateGoal() failed to persist zero values: goal = %+v, error = %v", got, err)
	}
	if err := DeleteGoal(int(second.ID)); err != nil {
		t.Fatalf("DeleteGoal() error = %v", err)
	}
	if err := DeleteGoal(int(second.ID)); !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("DeleteGoal() error = %v, want ErrRecordNotFound", err)
	}
}

func TestCreateGoalValidation(t *testing.T) {
	setupGoalDB(t)
	tests := []struct {
		name string
		goal Goal
		want error
	}{
		{"missing name", Goal{TargetAmount: 1}, ErrInvalidGoal},
		{"non-positive target", Goal{Name: "Goal", TargetAmount: 0}, ErrInvalidGoal},
		{"negative progress", Goal{Name: "Goal", TargetAmount: 1, CurrentAmount: -1}, ErrInvalidGoal},
		{"invalid date", Goal{Name: "Goal", TargetAmount: 1, TargetDate: "tomorrow"}, ErrInvalidDate},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := CreateGoal(&test.goal); !errors.Is(err, test.want) {
				t.Errorf("CreateGoal() error = %v, want %v", err, test.want)
			}
		})
	}
}
