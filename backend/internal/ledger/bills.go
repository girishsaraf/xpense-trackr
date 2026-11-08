package ledger

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ScheduledOccurrence struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	ScheduledFlowID   uint      `gorm:"not null;uniqueIndex:idx_flow_occurrence,priority:1" json:"scheduled_flow_id"`
	ScheduledDate     string    `gorm:"size:10;not null;uniqueIndex:idx_flow_occurrence,priority:2" json:"scheduled_date"`
	Status            string    `gorm:"size:12;not null" json:"status"`
	TransactionID     *uint     `json:"transaction_id,omitempty"`
	ActualAmountMinor int64     `gorm:"not null;default:0" json:"actual_amount_minor"`
	CreatedAt         time.Time `json:"created_at"`
}

type BillInsight struct {
	ScheduledFlow
	AccountName         string  `json:"account_name"`
	Status              string  `json:"status"`
	DaysUntilDue        int     `json:"days_until_due"`
	HasPriceChange      bool    `json:"has_price_change"`
	PreviousAmountMinor int64   `json:"previous_amount_minor,omitempty"`
	LatestAmountMinor   int64   `json:"latest_amount_minor,omitempty"`
	ChangePercent       float64 `json:"change_percent,omitempty"`
}

type ProcessOccurrenceRequest struct {
	ScheduledDate string `json:"scheduled_date"`
	Action        string `json:"action"`
	AmountMinor   int64  `json:"amount_minor,omitempty"`
}

func (ScheduledOccurrence) TableName() string {
	return "scheduled_occurrences"
}

func GetBillInsights(now time.Time) ([]BillInsight, error) {
	today := time.Date(now.UTC().Year(), now.UTC().Month(), now.UTC().Day(), 0, 0, 0, 0, time.UTC)
	flows := make([]ScheduledFlow, 0)
	if err := db.Where("enabled = ? AND type = ?", true, "expense").
		Order("next_date ASC, id ASC").Find(&flows).Error; err != nil {
		return nil, err
	}
	insights := make([]BillInsight, 0, len(flows))
	for _, flow := range flows {
		dueDate, err := time.Parse("2006-01-02", flow.NextDate)
		if err != nil {
			return nil, fmt.Errorf("scheduled flow %d has invalid date: %w", flow.ID, err)
		}
		daysUntilDue := int(dueDate.Sub(today).Hours() / 24)
		status := "upcoming"
		switch {
		case daysUntilDue < 0:
			status = "overdue"
		case daysUntilDue <= flow.ReminderDays:
			status = "due_soon"
		}
		var account Account
		if err := db.First(&account, flow.AccountID).Error; err != nil {
			return nil, fmt.Errorf("scheduled flow %d account: %w", flow.ID, err)
		}
		insight := BillInsight{
			ScheduledFlow: flow, AccountName: account.Name, Status: status, DaysUntilDue: daysUntilDue,
		}
		if err := addBillPriceChange(&insight); err != nil {
			return nil, err
		}
		insights = append(insights, insight)
	}
	return insights, nil
}

func ProcessScheduledOccurrence(flowID uint, request ProcessOccurrenceRequest) (*ScheduledOccurrence, error) {
	request.Action = strings.ToLower(strings.TrimSpace(request.Action))
	if !validDate(request.ScheduledDate) ||
		(request.Action != "record" && request.Action != "skip") ||
		(request.Action == "skip" && request.AmountMinor != 0) ||
		(request.Action == "record" && request.AmountMinor < 0) {
		return nil, ErrInvalidOccurrence
	}
	var result ScheduledOccurrence
	err := db.Transaction(func(tx *gorm.DB) error {
		var flow ScheduledFlow
		if err := tx.First(&flow, flowID).Error; err != nil {
			return err
		}
		var existing ScheduledOccurrence
		err := tx.Where("scheduled_flow_id = ? AND scheduled_date = ?", flowID, request.ScheduledDate).
			First(&existing).Error
		if err == nil {
			if existing.Status != request.Action ||
				(request.Action == "record" && request.AmountMinor > 0 && request.AmountMinor != existing.ActualAmountMinor) {
				return ErrOccurrenceExists
			}
			result = existing
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if !flow.Enabled || flow.NextDate != request.ScheduledDate {
			return ErrInvalidOccurrence
		}

		occurrence := ScheduledOccurrence{
			ScheduledFlowID: flow.ID,
			ScheduledDate:   request.ScheduledDate,
			Status:          request.Action,
		}
		if request.Action == "record" {
			amount := request.AmountMinor
			if amount == 0 {
				amount = flow.AmountMinor
			}
			transaction := Transaction{
				AccountID: flow.AccountID, CategoryID: flow.CategoryID, Type: flow.Type,
				AmountMinor: amount, Currency: flow.Currency, Payee: flow.Label,
				Description: flow.Label, Date: request.ScheduledDate,
			}
			if err := normalizeAndValidateTransaction(&transaction); err != nil {
				return err
			}
			if err := validateTransactionAccounts(tx, &transaction); err != nil {
				return err
			}
			if err := applyCategoryRule(tx, &transaction); err != nil {
				return err
			}
			if err := tx.Create(&transaction).Error; err != nil {
				return err
			}
			occurrence.TransactionID = &transaction.ID
			occurrence.ActualAmountMinor = amount
		}
		if err := tx.Create(&occurrence).Error; err != nil {
			return err
		}

		if flow.Frequency == "once" {
			flow.Enabled = false
		} else {
			currentDate, err := time.Parse("2006-01-02", flow.NextDate)
			if err != nil {
				return ErrInvalidOccurrence
			}
			nextDate, ok := nextFlowOccurrence(currentDate, flow.Frequency)
			if !ok {
				return ErrInvalidSchedule
			}
			if flow.EndDate != nil && nextDate.Format("2006-01-02") > *flow.EndDate {
				flow.Enabled = false
			} else {
				flow.NextDate = nextDate.Format("2006-01-02")
			}
		}
		if err := tx.Model(&flow).Select("NextDate", "Enabled").Updates(&flow).Error; err != nil {
			return err
		}
		result = occurrence
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func addBillPriceChange(insight *BillInsight) error {
	normalized := strings.ToLower(strings.TrimSpace(insight.Label))
	transactions := make([]Transaction, 0, 2)
	if err := db.Where(
		"transaction_type = ? AND currency = ? AND (LOWER(payee) = ? OR LOWER(description) = ?)",
		"expense", insight.Currency, normalized, normalized,
	).Order("date DESC, id DESC").Limit(2).Find(&transactions).Error; err != nil {
		return err
	}
	if len(transactions) < 2 {
		return nil
	}
	latest, previous := transactions[0].AmountMinor, transactions[1].AmountMinor
	if previous <= 0 {
		return nil
	}
	change := latest - previous
	insight.PreviousAmountMinor = previous
	insight.LatestAmountMinor = latest
	insight.ChangePercent = math.Round(float64(change)/float64(previous)*1000) / 10
	insight.HasPriceChange = math.Abs(insight.ChangePercent) >= 5
	return nil
}
