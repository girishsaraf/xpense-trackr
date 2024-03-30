// internal/budget/handler.go

package budget

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// GetAllBudgetsHandler handles the request to retrieve all budgets.
func GetAllBudgetsHandler(w http.ResponseWriter, r *http.Request) {
	budgets, err := GetAllBudgets()
	if err != nil {
		http.Error(w, "Failed to fetch budgets", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, budgets)
}

// CreateBudgetHandler handles the request to create a new budget.
func CreateBudgetHandler(w http.ResponseWriter, r *http.Request) {
	var budget Budget
	if err := json.NewDecoder(r.Body).Decode(&budget); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := CreateBudget(&budget); err != nil {
		http.Error(w, "Failed to create budget", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, budget)
}

// GetBudgetByIDHandler handles the request to retrieve a specific budget by ID.
func GetBudgetByIDHandler(w http.ResponseWriter, r *http.Request) {
	budgetID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid budget ID", http.StatusBadRequest)
		return
	}
	budget, err := GetBudgetByID(budgetID)
	if err != nil {
		http.Error(w, "Budget not found", http.StatusNotFound)
		return
	}
	respondWithJSON(w, budget)
}

// UpdateBudgetHandler handles the request to update an existing budget.
func UpdateBudgetHandler(w http.ResponseWriter, r *http.Request) {
	budgetID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid budget ID", http.StatusBadRequest)
		return
	}
	var updatedBudget Budget
	if err := json.NewDecoder(r.Body).Decode(&updatedBudget); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := UpdateBudget(budgetID, &updatedBudget); err != nil {
		http.Error(w, "Failed to update budget", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, updatedBudget)
}

// DeleteBudgetHandler handles the request to delete an existing budget.
func DeleteBudgetHandler(w http.ResponseWriter, r *http.Request) {
	budgetID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid budget ID", http.StatusBadRequest)
		return
	}
	if err := DeleteBudget(budgetID); err != nil {
		http.Error(w, "Failed to delete budget", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// respondWithJSON writes a JSON response with the given data.
func respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
