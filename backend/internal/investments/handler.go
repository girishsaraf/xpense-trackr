// internal/Investment/handler.go

package investments

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// GetAllInvestmentsHandler handles the request to retrieve all Investments.
func GetAllInvestmentsHandler(w http.ResponseWriter, r *http.Request) {
	Investments, err := GetAllInvestments()
	if err != nil {
		http.Error(w, "Failed to fetch Investments", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, Investments)
}

// CreateInvestmentHandler handles the request to create a new Investment.
func CreateInvestmentHandler(w http.ResponseWriter, r *http.Request) {
	var Investment Investment
	if err := json.NewDecoder(r.Body).Decode(&Investment); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := CreateInvestment(&Investment); err != nil {
		http.Error(w, "Failed to create Investment", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, Investment)
}

// GetInvestmentByIDHandler handles the request to retrieve a specific Investment by ID.
func GetInvestmentByIDHandler(w http.ResponseWriter, r *http.Request) {
	InvestmentID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid Investment ID", http.StatusBadRequest)
		return
	}
	Investment, err := GetInvestmentByID(InvestmentID)
	if err != nil {
		http.Error(w, "Investment not found", http.StatusNotFound)
		return
	}
	respondWithJSON(w, Investment)
}

// UpdateInvestmentHandler handles the request to update an existing Investment.
func UpdateInvestmentHandler(w http.ResponseWriter, r *http.Request) {
	InvestmentID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid Investment ID", http.StatusBadRequest)
		return
	}
	var updatedInvestment Investment
	if err := json.NewDecoder(r.Body).Decode(&updatedInvestment); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := UpdateInvestment(InvestmentID, &updatedInvestment); err != nil {
		http.Error(w, "Failed to update Investment", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, updatedInvestment)
}

// DeleteInvestmentHandler handles the request to delete an existing Investment.
func DeleteInvestmentHandler(w http.ResponseWriter, r *http.Request) {
	InvestmentID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid Investment ID", http.StatusBadRequest)
		return
	}
	if err := DeleteInvestment(InvestmentID); err != nil {
		http.Error(w, "Failed to delete Investment", http.StatusInternalServerError)
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
