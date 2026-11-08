package planning

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetGoalsHandler(w http.ResponseWriter, _ *http.Request) {
	goals, err := GetGoals()
	if err != nil {
		http.Error(w, "Failed to fetch financial goals", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, goals)
}

func CreateGoalHandler(w http.ResponseWriter, r *http.Request) {
	var goal Goal
	if err := json.NewDecoder(r.Body).Decode(&goal); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := CreateGoal(&goal); err != nil {
		if errors.Is(err, ErrInvalidGoal) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to create financial goal", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, goal)
}

func DeleteGoalHandler(w http.ResponseWriter, r *http.Request) {
	goalID, ok := parseID(w, r)
	if !ok {
		return
	}
	if err := DeleteGoal(goalID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Goal not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete financial goal", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func GetContributionsHandler(w http.ResponseWriter, r *http.Request) {
	goalID, ok := parseID(w, r)
	if !ok {
		return
	}
	entries, err := GetContributions(goalID)
	if err != nil {
		http.Error(w, "Failed to fetch goal contributions", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, entries)
}

func CreateContributionHandler(w http.ResponseWriter, r *http.Request) {
	goalID, ok := parseID(w, r)
	if !ok {
		return
	}
	var entry GoalContribution
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := AddContribution(goalID, &entry); err != nil {
		switch {
		case errors.Is(err, ErrInvalidEntry):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, ErrInsufficientGoalFunds):
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, gorm.ErrRecordNotFound):
			http.Error(w, "Goal not found", http.StatusNotFound)
		default:
			http.Error(w, "Failed to add goal contribution", http.StatusInternalServerError)
		}
		return
	}
	writeJSON(w, http.StatusCreated, entry)
}

func GetEnvelopesHandler(w http.ResponseWriter, r *http.Request) {
	envelopes, err := GetEnvelopes(r.URL.Query().Get("period"))
	if err != nil {
		if errors.Is(err, ErrInvalidPeriod) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to fetch budget envelopes", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, envelopes)
}

func CreateEnvelopeHandler(w http.ResponseWriter, r *http.Request) {
	var envelope Envelope
	if err := json.NewDecoder(r.Body).Decode(&envelope); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := CreateEnvelope(&envelope); err != nil {
		if errors.Is(err, ErrInvalidEnvelope) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to create budget envelope", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, envelope)
}

func UpdateEnvelopeHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var envelope Envelope
	if err := json.NewDecoder(r.Body).Decode(&envelope); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := UpdateEnvelope(id, &envelope); err != nil {
		switch {
		case errors.Is(err, ErrInvalidEnvelope):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, gorm.ErrRecordNotFound):
			http.Error(w, "Budget envelope not found", http.StatusNotFound)
		default:
			http.Error(w, "Failed to update budget envelope", http.StatusInternalServerError)
		}
		return
	}
	envelope.ID = id
	writeJSON(w, http.StatusOK, envelope)
}

func DeleteEnvelopeHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	if err := DeleteEnvelope(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Budget envelope not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete budget envelope", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func parseID(w http.ResponseWriter, r *http.Request) (uint, bool) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil || id == 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return 0, false
	}
	return uint(id), true
}

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
