package trends

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/girishsaraf/xpense-trackr/backend/internal/expense"
	"gorm.io/gorm"
)

var db *gorm.DB

type SpendingPoint struct {
	Month    string  `json:"month"`
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
}

func SetDB(database *gorm.DB) {
	db = database
}

func SpendingHandler(w http.ResponseWriter, r *http.Request) {
	months := 6
	if raw := r.URL.Query().Get("months"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 || parsed > 24 {
			http.Error(w, "months must be between 1 and 24", http.StatusBadRequest)
			return
		}
		months = parsed
	}

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, -(months - 1), 0)
	monthExpression := "DATE_FORMAT(expenses.date, '%Y-%m')"
	if db.Dialector.Name() == "sqlite" {
		monthExpression = "strftime('%Y-%m', expenses.date)"
	}

	points := make([]SpendingPoint, 0)
	err := db.Model(&expense.Expense{}).
		Select(monthExpression+" AS month, categories.name AS category, SUM(expenses.amount) AS amount").
		Joins("JOIN categories ON categories.id = expenses.category_id").
		Where("expenses.date >= ?", start.Format("2006-01-02")).
		Group("month, categories.name").
		Order("month ASC, categories.name ASC").
		Scan(&points).Error
	if err != nil {
		http.Error(w, "Failed to fetch spending trends", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(points); err != nil {
		http.Error(w, "Failed to encode spending trends", http.StatusInternalServerError)
	}
}
