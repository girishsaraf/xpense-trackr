// internal/expense/handler.go

package expense

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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
