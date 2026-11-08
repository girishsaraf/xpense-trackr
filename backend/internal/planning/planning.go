package planning

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/girishsaraf/xpense-trackr/backend/internal/category"
	"github.com/girishsaraf/xpense-trackr/backend/internal/ledger"
	"gorm.io/gorm"
)

var (
	db                       *gorm.DB
	ErrInvalidGoal           = errors.New("goal name, supported currency, positive target, and valid target date are required")
	ErrInvalidEntry          = errors.New("contribution requires a positive amount, supported direction, and valid date")
	ErrInsufficientGoalFunds = errors.New("withdrawal exceeds the goal's available contributions")
	ErrInvalidEnvelope       = errors.New("envelope requires a category, supported currency, and positive monthly amount")
	ErrInvalidPeriod         = errors.New("period must use YYYY-MM format")
)

type Goal struct {
	ID                   uint      `gorm:"primaryKey" json:"id"`
	Name                 string    `gorm:"size:120;not null" json:"name"`
	Currency             string    `gorm:"size:3;not null;index" json:"currency"`
	TargetAmountMinor    int64     `gorm:"not null" json:"target_amount_minor"`
	TargetDate           string    `gorm:"size:10" json:"target_date,omitempty"`
	CurrentAmountMinor   int64     `gorm:"->" json:"current_amount_minor"`
	RemainingMinor       int64     `gorm:"->" json:"remaining_minor"`
	RequiredMonthlyMinor int64     `gorm:"->" json:"required_monthly_minor,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type GoalContribution struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	GoalID      uint      `gorm:"not null;index" json:"goal_id"`
	AmountMinor int64     `gorm:"not null" json:"amount_minor"`
	Direction   string    `gorm:"size:12;not null" json:"direction"`
	Date        string    `gorm:"size:10;not null;index" json:"date"`
	Note        string    `gorm:"size:300" json:"note,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type Envelope struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	CategoryID        uint      `gorm:"not null;uniqueIndex:idx_envelope_category_currency" json:"category_id"`
	Currency          string    `gorm:"size:3;not null;uniqueIndex:idx_envelope_category_currency" json:"currency"`
	MonthlyLimitMinor int64     `gorm:"not null" json:"monthly_limit_minor"`
	Rollover          bool      `gorm:"not null;default:false" json:"rollover"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type EnvelopeSummary struct {
	Envelope
	Period         string `json:"period"`
	SpentMinor     int64  `json:"spent_minor"`
	RolloverMinor  int64  `json:"rollover_minor"`
	AvailableMinor int64  `json:"available_minor"`
}

func SetDB(database *gorm.DB) {
	db = database
}

func (Goal) TableName() string {
	return "financial_goals"
}

func (GoalContribution) TableName() string {
	return "goal_contributions"
}

func (Envelope) TableName() string {
	return "budget_envelopes"
}

func GetGoals() ([]Goal, error) {
	goals := make([]Goal, 0)
	if err := db.Order("target_date = '', target_date ASC, id ASC").Find(&goals).Error; err != nil {
		return nil, err
	}
	for i := range goals {
		if err := populateGoalProgress(&goals[i]); err != nil {
			return nil, err
		}
	}
	return goals, nil
}

func CreateGoal(goal *Goal) error {
	goal.Name = strings.TrimSpace(goal.Name)
	goal.Currency = strings.ToUpper(strings.TrimSpace(goal.Currency))
	if err := validateGoal(goal); err != nil {
		return err
	}
	goal.ID = 0
	goal.CreatedAt = time.Time{}
	goal.UpdatedAt = time.Time{}
	if err := db.Create(goal).Error; err != nil {
		return err
	}
	return populateGoalProgress(goal)
}

func AddContribution(goalID uint, entry *GoalContribution) error {
	entry.Note = strings.TrimSpace(entry.Note)
	entry.Direction = strings.ToLower(strings.TrimSpace(entry.Direction))
	if entry.AmountMinor <= 0 ||
		(entry.Direction != "contribution" && entry.Direction != "withdrawal") ||
		!validDate(entry.Date) || len(entry.Note) > 300 {
		return ErrInvalidEntry
	}
	return db.Transaction(func(tx *gorm.DB) error {
		var goal Goal
		if err := tx.First(&goal, goalID).Error; err != nil {
			return err
		}
		if entry.Direction == "withdrawal" {
			var progress int64
			if err := tx.Model(&GoalContribution{}).
				Select("COALESCE(SUM(CASE WHEN direction = 'withdrawal' THEN -amount_minor ELSE amount_minor END), 0)").
				Where("goal_id = ?", goalID).Scan(&progress).Error; err != nil {
				return err
			}
			if entry.AmountMinor > progress {
				return ErrInsufficientGoalFunds
			}
		}
		entry.GoalID = goalID
		return tx.Create(entry).Error
	})
}

func GetContributions(goalID uint) ([]GoalContribution, error) {
	var goal Goal
	if err := db.First(&goal, goalID).Error; err != nil {
		return nil, err
	}
	entries := make([]GoalContribution, 0)
	if err := db.Where("goal_id = ?", goalID).Order("date DESC, id DESC").Find(&entries).Error; err != nil {
		return nil, err
	}
	return entries, nil
}

func DeleteGoal(goalID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&GoalContribution{}, "goal_id = ?", goalID).Error; err != nil {
			return err
		}
		result := tx.Delete(&Goal{}, goalID)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

func GetEnvelopes(period string) ([]EnvelopeSummary, error) {
	if period == "" {
		period = time.Now().UTC().Format("2006-01")
	}
	start, err := parsePeriod(period)
	if err != nil {
		return nil, err
	}
	end := start.AddDate(0, 1, 0)
	envelopes := make([]Envelope, 0)
	if err := db.Order("id ASC").Find(&envelopes).Error; err != nil {
		return nil, err
	}
	summaries := make([]EnvelopeSummary, 0, len(envelopes))
	for _, envelope := range envelopes {
		summary := EnvelopeSummary{Envelope: envelope, Period: period}
		if err := db.Table("ledger_transactions").
			Select("COALESCE(SUM(amount_minor), 0)").
			Where("transaction_type = ? AND category_id = ? AND currency = ? AND date >= ? AND date < ?",
				"expense", envelope.CategoryID, envelope.Currency, start.Format("2006-01-02"), end.Format("2006-01-02")).
			Scan(&summary.SpentMinor).Error; err != nil {
			return nil, err
		}
		if envelope.Rollover {
			if err := calculateRollover(&summary, start); err != nil {
				return nil, err
			}
		}
		summary.AvailableMinor = envelope.MonthlyLimitMinor + summary.RolloverMinor - summary.SpentMinor
		summaries = append(summaries, summary)
	}
	return summaries, nil
}

func CreateEnvelope(envelope *Envelope) error {
	envelope.Currency = strings.ToUpper(strings.TrimSpace(envelope.Currency))
	if envelope.CategoryID == 0 || !ledger.SupportedCurrency(envelope.Currency) || envelope.MonthlyLimitMinor <= 0 {
		return ErrInvalidEnvelope
	}
	envelope.ID = 0
	envelope.CreatedAt = time.Time{}
	envelope.UpdatedAt = time.Time{}
	var cat category.Category
	if err := db.First(&cat, envelope.CategoryID).Error; err != nil {
		return fmt.Errorf("envelope category: %w", err)
	}
	return db.Create(envelope).Error
}

func UpdateEnvelope(id uint, updated *Envelope) error {
	updated.Currency = strings.ToUpper(strings.TrimSpace(updated.Currency))
	if updated.CategoryID == 0 || !ledger.SupportedCurrency(updated.Currency) || updated.MonthlyLimitMinor <= 0 {
		return ErrInvalidEnvelope
	}
	var cat category.Category
	if err := db.First(&cat, updated.CategoryID).Error; err != nil {
		return fmt.Errorf("envelope category: %w", err)
	}
	return db.Transaction(func(tx *gorm.DB) error {
		var existing Envelope
		if err := tx.First(&existing, id).Error; err != nil {
			return err
		}
		return tx.Model(&existing).Select(
			"CategoryID", "Currency", "MonthlyLimitMinor", "Rollover",
		).Updates(updated).Error
	})
}

func DeleteEnvelope(id uint) error {
	result := db.Delete(&Envelope{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func validateGoal(goal *Goal) error {
	if goal.Name == "" || len(goal.Name) > 120 || !ledger.SupportedCurrency(goal.Currency) || goal.TargetAmountMinor <= 0 {
		return ErrInvalidGoal
	}
	if goal.TargetDate != "" && !validDate(goal.TargetDate) {
		return ErrInvalidGoal
	}
	return nil
}

func populateGoalProgress(goal *Goal) error {
	var progress int64
	err := db.Table("goal_contributions").
		Select("COALESCE(SUM(CASE WHEN direction = 'withdrawal' THEN -amount_minor ELSE amount_minor END), 0)").
		Where("goal_id = ?", goal.ID).
		Scan(&progress).Error
	if err != nil {
		return err
	}
	goal.CurrentAmountMinor = progress
	goal.RemainingMinor = goal.TargetAmountMinor - progress
	if goal.RemainingMinor < 0 {
		goal.RemainingMinor = 0
	}
	if goal.TargetDate != "" && goal.RemainingMinor > 0 {
		target, _ := time.Parse("2006-01-02", goal.TargetDate)
		now := time.Now().UTC()
		months := (target.Year()-now.Year())*12 + int(target.Month()-now.Month())
		if target.Day() > now.Day() {
			months++
		}
		if months < 1 {
			months = 1
		}
		goal.RequiredMonthlyMinor = (goal.RemainingMinor + int64(months) - 1) / int64(months)
	}
	return nil
}

func calculateRollover(summary *EnvelopeSummary, currentMonth time.Time) error {
	startMonth := time.Date(summary.CreatedAt.Year(), summary.CreatedAt.Month(), 1, 0, 0, 0, 0, time.UTC)
	periods := monthDifference(startMonth, currentMonth)
	if periods <= 0 {
		return nil
	}
	end := currentMonth.Format("2006-01-02")
	var spent int64
	if err := db.Table("ledger_transactions").
		Select("COALESCE(SUM(amount_minor), 0)").
		Where("transaction_type = ? AND category_id = ? AND currency = ? AND date >= ? AND date < ?",
			"expense", summary.CategoryID, summary.Currency, startMonth.Format("2006-01-02"), end).
		Scan(&spent).Error; err != nil {
		return err
	}
	summary.RolloverMinor = int64(periods)*summary.MonthlyLimitMinor - spent
	return nil
}

func parsePeriod(period string) (time.Time, error) {
	if len(period) != 7 {
		return time.Time{}, ErrInvalidPeriod
	}
	value, err := time.Parse("2006-01", period)
	if err != nil {
		return time.Time{}, ErrInvalidPeriod
	}
	return value, nil
}

func monthDifference(start, end time.Time) int {
	return (end.Year()-start.Year())*12 + int(end.Month()-start.Month())
}

func validDate(value string) bool {
	date, err := time.Parse("2006-01-02", value)
	return err == nil && date.Format("2006-01-02") == value
}
