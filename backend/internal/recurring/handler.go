package recurring

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// GetAllRecurringExpensesHandler retrieves all expenses.
func GetAllRecurringExpensesHandler(w http.ResponseWriter, r *http.Request) {
	expenses, err := GetAllRecurringExpenses()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(expenses)
	if err != nil {
		return
	}
}

// GetRecurringExpenseByIDHandler retrieves a specific expense by ID.
func GetRecurringExpenseByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid expense ID", http.StatusBadRequest)
		return
	}

	expense, err := GetRecurringExpenseByID(id)
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

// CreateRecurringExpenseHandler creates a new expense.
func CreateRecurringExpenseHandler(w http.ResponseWriter, r *http.Request) {
	var expense RecurringExpense
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := CreateRecurringExpense(&expense); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(expense)
	if err != nil {
		return
	}
}

// UpdateRecurringExpenseHandler updates an existing expense.
func UpdateRecurringExpenseHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid expense ID", http.StatusBadRequest)
		return
	}

	var updatedExpense RecurringExpense
	if err := json.NewDecoder(r.Body).Decode(&updatedExpense); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := UpdateRecurringExpense(id, &updatedExpense); err != nil {
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

// DeleteRecurringExpenseHandler deletes an existing expense.
func DeleteRecurringExpenseHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid expense ID", http.StatusBadRequest)
		return
	}

	if err := DeleteRecurringExpense(id); err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			http.Error(w, "Expense not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetAllRecurringInvestmentsHandler retrieves all expenses.
func GetAllRecurringInvestmentsHandler(w http.ResponseWriter, r *http.Request) {
	expenses, err := GetAllRecurringInvestments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(expenses)
	if err != nil {
		return
	}
}

// GetRecurringInvestmentByIDHandler retrieves a specific expense by ID.
func GetRecurringInvestmentByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid expense ID", http.StatusBadRequest)
		return
	}

	expense, err := GetRecurringInvestmentByID(id)
	if err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			http.Error(w, "Investment not found", http.StatusNotFound)
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

// CreateRecurringInvestmentHandler creates a new expense.
func CreateRecurringInvestmentHandler(w http.ResponseWriter, r *http.Request) {
	var expense RecurringInvestment
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := CreateRecurringInvestment(&expense); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(expense)
	if err != nil {
		return
	}
}

// UpdateRecurringInvestmentHandler updates an existing expense.
func UpdateRecurringInvestmentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid expense ID", http.StatusBadRequest)
		return
	}

	var updatedInvestment RecurringInvestment
	if err := json.NewDecoder(r.Body).Decode(&updatedInvestment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := UpdateRecurringInvestment(id, &updatedInvestment); err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			http.Error(w, "Investment not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(updatedInvestment)
	if err != nil {
		return
	}
}

// DeleteRecurringInvestmentHandler deletes an existing expense.
func DeleteRecurringInvestmentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid expense ID", http.StatusBadRequest)
		return
	}

	if err := DeleteRecurringInvestment(id); err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			http.Error(w, "Investment not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
