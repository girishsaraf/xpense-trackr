package recurring

import (
	"time"
)

type RecurringExpense struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CategoryID  uint      `json:"category_id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Frequency   string    `json:"frequency"`
}

type RecurringInvestment struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Amount         float64   `json:"amount"`
	InvestmentType string    `json:"investment_type"`
	Description    string    `json:"description"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Frequency      string    `json:"frequency"`
}
