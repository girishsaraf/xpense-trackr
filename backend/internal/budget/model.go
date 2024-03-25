// budget.go

package models

type Budget struct {
	ID         int     `json:"id"`
	UserID     int     `json:"user_id"`
	CategoryID int     `json:"category_id"`
	Amount     float64 `json:"amount"`
}
