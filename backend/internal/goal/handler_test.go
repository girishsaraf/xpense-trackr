package goal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/girishsaraf/xpense-trackr/backend/internal/testutil"
	"github.com/gorilla/mux"
)

func goalRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/goals", GetAllGoalsHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/goals", CreateGoalHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/goals/{id}", GetGoalByIDHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/goals/{id}", UpdateGoalHandler).Methods(http.MethodPut)
	router.HandleFunc("/api/goals/{id}", DeleteGoalHandler).Methods(http.MethodDelete)
	return router
}

func requestGoal(t *testing.T, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	goalRouter().ServeHTTP(recorder, request)
	return recorder
}

func TestGoalHandlers(t *testing.T) {
	database := testutil.NewDB(t, &Goal{})
	SetDB(database)

	created := requestGoal(t, http.MethodPost, "/api/goals", `{"name":"Home","target_amount":20000,"current_amount":500,"target_date":"2028-01-01"}`)
	if created.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d: %s", created.Code, http.StatusCreated, created.Body)
	}
	var goal Goal
	if err := json.NewDecoder(created.Body).Decode(&goal); err != nil || goal.ID == 0 {
		t.Fatalf("create response = %+v, error = %v", goal, err)
	}

	updated := requestGoal(t, http.MethodPut, "/api/goals/1", `{"name":"Home","target_amount":20000,"current_amount":1500}`)
	if updated.Code != http.StatusOK {
		t.Fatalf("update status = %d, want %d: %s", updated.Code, http.StatusOK, updated.Body)
	}
	got, err := GetGoalByID(1)
	if err != nil || got.CurrentAmount != 1500 || got.TargetDate != "" {
		t.Fatalf("updated goal = %+v, error = %v", got, err)
	}

	list := requestGoal(t, http.MethodGet, "/api/goals", "")
	if list.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d", list.Code, http.StatusOK)
	}
	deleted := requestGoal(t, http.MethodDelete, "/api/goals/1", "")
	if deleted.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want %d", deleted.Code, http.StatusNoContent)
	}
}

func TestCreateGoalRejectsInvalidInput(t *testing.T) {
	database := testutil.NewDB(t, &Goal{})
	SetDB(database)

	for _, body := range []string{`{`, `{"name":"Goal","target_amount":0}`} {
		response := requestGoal(t, http.MethodPost, "/api/goals", body)
		if response.Code != http.StatusBadRequest {
			t.Errorf("POST body %q status = %d, want %d", body, response.Code, http.StatusBadRequest)
		}
	}
}
