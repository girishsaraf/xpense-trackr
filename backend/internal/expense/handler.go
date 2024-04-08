// internal/expense/handler.go

package expense

import (
	"encoding/json"
	"errors"
	"github.com/girishsaraf/xpense-trackr/backend/internal/category"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

// GetAllExpensesHandler retrieves all expenses.
func GetAllExpensesHandler(w http.ResponseWriter, r *http.Request) {
	expenses, err := GetAllExpenses()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(expenses)
	if err != nil {
		return
	}
}

// GetExpenseByIDHandler retrieves a specific expense by ID.
func GetExpenseByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid expense ID", http.StatusBadRequest)
		return
	}

	expense, err := GetExpenseByID(id)
	if err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			http.Error(w, "Expense not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(expense)
	if err != nil {
		return
	}
}

// CreateExpenseHandler creates a new expense.
func CreateExpenseHandler(w http.ResponseWriter, r *http.Request) {
	var expense Expense
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := CreateExpense(&expense); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(expense)
	if err != nil {
		return
	}
}

// UpdateExpenseHandler updates an existing expense.
func UpdateExpenseHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid expense ID", http.StatusBadRequest)
		return
	}

	var updatedExpense Expense
	if err := json.NewDecoder(r.Body).Decode(&updatedExpense); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := UpdateExpense(id, &updatedExpense); err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			http.Error(w, "Expense not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(updatedExpense)
	if err != nil {
		return
	}
}

// DeleteExpenseHandler deletes an existing expense.
func DeleteExpenseHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid expense ID", http.StatusBadRequest)
		return
	}

	if err := DeleteExpense(id); err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			http.Error(w, "Expense not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CumulativeExpensesByMonthHandler handles the request to get cumulative investments by month.
func CumulativeExpensesByMonthHandler(w http.ResponseWriter, r *http.Request) {
	lastMonthsStr := r.URL.Query().Get("lastMonths")
	lastMonths, err := strconv.Atoi(lastMonthsStr)
	if err != nil || lastMonths < 1 {
		http.Error(w, "Invalid value for lastMonths", http.StatusBadRequest)
		return
	}

	// Fetch categories
	var categories []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	err = db.Model(&category.Category{}).Select("id, name").Find(&categories).Error
	if err != nil {
		http.Error(w, "Failed to retrieve categories", http.StatusInternalServerError)
		return
	}

	// Create a map for category ID to name mapping
	categoryMap := make(map[int]string)
	for _, categoryVal := range categories {
		categoryMap[categoryVal.ID] = categoryVal.Name
	}

	var cumulativeExpenses []struct {
		Month    string  `json:"month"`
		Category string  `json:"category"`
		Amount   float64 `json:"amount"`
	}

	// Calculate start date based on lastMonths
	startDate := time.Now().AddDate(0, -lastMonths, 0)

	// Query to get cumulative investments grouped by month
	query := db.Model(&Expense{}).
		Select("DATE_FORMAT(date, '%Y-%m') AS month, category_id AS category, SUM(amount) AS amount").
		Where("date >= ?", startDate).
		Group("month, category").
		Order("month ASC").
		Scan(&cumulativeExpenses)

	err = query.Error
	if err != nil {
		http.Error(w, "Failed to retrieve cumulative investments", http.StatusInternalServerError)
		return
	}

	// Calculate cumulative expenses
	cumulativeAmount := make(map[string]float64)

	for i := range cumulativeExpenses {
		categoryId, _ := strconv.Atoi(cumulativeExpenses[i].Category)
		categoryName := categoryMap[categoryId]
		cumulativeExpenses[i].Category = categoryName
	}

	for i := range cumulativeExpenses {
		categoryName := cumulativeExpenses[i].Category
		amount := cumulativeExpenses[i].Amount
		cumulativeAmount[categoryName] += amount
		cumulativeExpenses[i].Amount = cumulativeAmount[categoryName]
	}
	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cumulativeExpenses)
}
