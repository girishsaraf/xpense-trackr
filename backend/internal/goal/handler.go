package goal

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetAllGoalsHandler(w http.ResponseWriter, _ *http.Request) {
	goals, err := GetAllGoals()
	if err != nil {
		http.Error(w, "Failed to fetch goals", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, goals)
}

func GetGoalByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	goal, err := GetGoalByID(id)
	if err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			http.Error(w, "Goal not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch goal", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, goal)
}

func CreateGoalHandler(w http.ResponseWriter, r *http.Request) {
	var goal Goal
	if err := json.NewDecoder(r.Body).Decode(&goal); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := CreateGoal(&goal); err != nil {
		if errors.Is(err, ErrInvalidGoal) || errors.Is(err, ErrInvalidDate) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to create goal", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, goal)
}

func UpdateGoalHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var updated Goal
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := UpdateGoal(id, &updated); err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			http.Error(w, "Goal not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, ErrInvalidGoal) || errors.Is(err, ErrInvalidDate) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to update goal", http.StatusInternalServerError)
		return
	}
	updated.ID = uint(id)
	writeJSON(w, http.StatusOK, updated)
}

func DeleteGoalHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	if err := DeleteGoal(id); err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			http.Error(w, "Goal not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete goal", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func parseID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil || id < 1 {
		http.Error(w, "Invalid goal ID", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
