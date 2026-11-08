package ledger

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

var ErrInvalidSchedule = errors.New("scheduled flow requires a valid account, label, type, positive amount, currency, date, and frequency")

type ScheduledFlow struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	AccountID   uint      `gorm:"not null;index" json:"account_id"`
	Label       string    `gorm:"size:180;not null" json:"label"`
	Type        string    `gorm:"size:16;not null" json:"type"`
	AmountMinor int64     `gorm:"not null" json:"amount_minor"`
	Currency    string    `gorm:"size:3;not null" json:"currency"`
	Frequency   string    `gorm:"size:16;not null" json:"frequency"`
	NextDate    string    `gorm:"size:10;not null;index" json:"next_date"`
	EndDate     *string   `gorm:"size:10" json:"end_date,omitempty"`
	Enabled     bool      `gorm:"not null;default:true" json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ForecastPoint struct {
	Date         string `json:"date"`
	BalanceMinor int64  `json:"balance_minor"`
}

type AccountForecast struct {
	AccountID             uint            `json:"account_id"`
	Name                  string          `json:"name"`
	Currency              string          `json:"currency"`
	CurrentBalanceMinor   int64           `json:"current_balance_minor"`
	ProjectedBalanceMinor int64           `json:"projected_balance_minor"`
	Points                []ForecastPoint `json:"points"`
}

type ForecastEvent struct {
	ScheduledFlowID uint   `json:"scheduled_flow_id"`
	AccountID       uint   `json:"account_id"`
	Label           string `json:"label"`
	Type            string `json:"type"`
	AmountMinor     int64  `json:"amount_minor"`
	Currency        string `json:"currency"`
	Date            string `json:"date"`
}

type CashflowForecast struct {
	StartDate string            `json:"start_date"`
	EndDate   string            `json:"end_date"`
	Days      int               `json:"days"`
	Accounts  []AccountForecast `json:"accounts"`
	Events    []ForecastEvent   `json:"events"`
}

func GetScheduledFlows() ([]ScheduledFlow, error) {
	flows := make([]ScheduledFlow, 0)
	if err := db.Order("next_date ASC, id ASC").Find(&flows).Error; err != nil {
		return nil, err
	}
	return flows, nil
}

func CreateScheduledFlow(flow *ScheduledFlow) error {
	if err := validateScheduledFlow(flow); err != nil {
		return err
	}
	return db.Create(flow).Error
}

func UpdateScheduledFlow(id int, updated *ScheduledFlow) error {
	if err := validateScheduledFlow(updated); err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		var existing ScheduledFlow
		if err := tx.First(&existing, id).Error; err != nil {
			return err
		}
		return tx.Model(&existing).Select(
			"AccountID", "Label", "Type", "AmountMinor", "Currency", "Frequency", "NextDate", "EndDate", "Enabled",
		).Updates(updated).Error
	})
}

func DeleteScheduledFlow(id int) error {
	result := db.Delete(&ScheduledFlow{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func BuildCashflowForecast(start time.Time, days int) (*CashflowForecast, error) {
	if days < 1 || days > 365 {
		return nil, ErrInvalidTxn
	}
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 0, days)

	accounts, err := GetAllAccounts()
	if err != nil {
		return nil, err
	}
	flows := make([]ScheduledFlow, 0)
	if err := db.Where("enabled = ?", true).Find(&flows).Error; err != nil {
		return nil, err
	}

	forecast := &CashflowForecast{
		StartDate: start.Format("2006-01-02"),
		EndDate:   end.Format("2006-01-02"),
		Days:      days,
		Accounts:  make([]AccountForecast, 0, len(accounts)),
		Events:    make([]ForecastEvent, 0),
	}
	byAccount := make(map[uint]int, len(accounts))
	for _, account := range accounts {
		byAccount[account.ID] = len(forecast.Accounts)
		forecast.Accounts = append(forecast.Accounts, AccountForecast{
			AccountID: account.ID, Name: account.Name, Currency: account.Currency,
			CurrentBalanceMinor: account.BalanceMinor, ProjectedBalanceMinor: account.BalanceMinor,
			Points: []ForecastPoint{{Date: start.Format("2006-01-02"), BalanceMinor: account.BalanceMinor}},
		})
	}

	deltas := make(map[string]map[uint]int64)
	for _, flow := range flows {
		accountIndex, found := byAccount[flow.AccountID]
		if !found {
			return nil, fmt.Errorf("scheduled flow %d references missing account %d", flow.ID, flow.AccountID)
		}
		account := accounts[accountIndex]
		if account.Currency != flow.Currency {
			return nil, fmt.Errorf("scheduled flow %d currency no longer matches account", flow.ID)
		}
		occurrence, err := time.Parse("2006-01-02", flow.NextDate)
		if err != nil {
			return nil, fmt.Errorf("scheduled flow %d has invalid date: %w", flow.ID, err)
		}
		var lastOccurrence time.Time
		if flow.EndDate != nil {
			lastOccurrence, err = time.Parse("2006-01-02", *flow.EndDate)
			if err != nil {
				return nil, fmt.Errorf("scheduled flow %d has invalid end date: %w", flow.ID, err)
			}
		}
		for occurrence.Before(start) && flow.Frequency != "once" {
			occurrence, _ = nextFlowOccurrence(occurrence, flow.Frequency)
		}
		for !occurrence.Before(start) && !occurrence.After(end) && (flow.EndDate == nil || !occurrence.After(lastOccurrence)) {
			date := occurrence.Format("2006-01-02")
			if deltas[date] == nil {
				deltas[date] = make(map[uint]int64)
			}
			delta := flow.AmountMinor
			if flow.Type == "expense" {
				delta = -delta
			}
			deltas[date][flow.AccountID] += delta
			forecast.Events = append(forecast.Events, ForecastEvent{
				ScheduledFlowID: flow.ID, AccountID: flow.AccountID, Label: flow.Label,
				Type: flow.Type, AmountMinor: flow.AmountMinor, Currency: flow.Currency, Date: date,
			})
			next, ok := nextFlowOccurrence(occurrence, flow.Frequency)
			if !ok {
				break
			}
			occurrence = next
		}
	}

	for day := 1; day <= days; day++ {
		date := start.AddDate(0, 0, day)
		dateText := date.Format("2006-01-02")
		for accountID, delta := range deltas[dateText] {
			index := byAccount[accountID]
			forecast.Accounts[index].ProjectedBalanceMinor += delta
		}
		for i := range forecast.Accounts {
			account := &forecast.Accounts[i]
			account.Points = append(account.Points, ForecastPoint{
				Date: dateText, BalanceMinor: account.ProjectedBalanceMinor,
			})
		}
	}
	sort.SliceStable(forecast.Events, func(i, j int) bool {
		return forecast.Events[i].Date < forecast.Events[j].Date
	})
	return forecast, nil
}

func validateScheduledFlow(flow *ScheduledFlow) error {
	flow.Label = strings.TrimSpace(flow.Label)
	flow.Type = strings.ToLower(strings.TrimSpace(flow.Type))
	flow.Currency = strings.ToUpper(strings.TrimSpace(flow.Currency))
	flow.Frequency = strings.ToLower(strings.TrimSpace(flow.Frequency))
	if flow.AccountID == 0 || flow.Label == "" || len(flow.Label) > 180 ||
		(flow.Type != "income" && flow.Type != "expense") ||
		flow.AmountMinor <= 0 || !currencyCodePattern.MatchString(flow.Currency) ||
		(flow.Frequency != "once" && flow.Frequency != "weekly" && flow.Frequency != "biweekly" && flow.Frequency != "monthly" && flow.Frequency != "yearly") {
		return ErrInvalidSchedule
	}
	if _, err := time.Parse("2006-01-02", flow.NextDate); err != nil {
		return ErrInvalidSchedule
	}
	if flow.EndDate != nil {
		end, err := time.Parse("2006-01-02", *flow.EndDate)
		next, nextErr := time.Parse("2006-01-02", flow.NextDate)
		if err != nil || nextErr != nil || end.Before(next) {
			return ErrInvalidSchedule
		}
	}
	var account Account
	if err := db.First(&account, flow.AccountID).Error; err != nil {
		return fmt.Errorf("scheduled flow account: %w", err)
	}
	if account.Currency != flow.Currency {
		return ErrInvalidSchedule
	}
	return nil
}

func nextFlowOccurrence(date time.Time, frequency string) (time.Time, bool) {
	switch frequency {
	case "once":
		return date, false
	case "weekly":
		return date.AddDate(0, 0, 7), true
	case "biweekly":
		return date.AddDate(0, 0, 14), true
	case "monthly":
		return addMonthsClamped(date, 1), true
	case "yearly":
		return addMonthsClamped(date, 12), true
	default:
		return date, false
	}
}

func addMonthsClamped(date time.Time, months int) time.Time {
	firstOfTargetMonth := time.Date(date.Year(), date.Month()+time.Month(months), 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstOfTargetMonth.AddDate(0, 1, -1).Day()
	day := date.Day()
	if day > lastDay {
		day = lastDay
	}
	return time.Date(firstOfTargetMonth.Year(), firstOfTargetMonth.Month(), day, 0, 0, 0, 0, time.UTC)
}
