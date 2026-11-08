package budget

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/girishsaraf/xpense-trackr/backend/internal/testutil"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func newBudgetHandlerDB(t *testing.T) *gorm.DB {
	t.Helper()
	database := testutil.NewDB(t, &Budget{})
	SetDB(database)
	return database
}

func newBudgetTestRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/budgets", GetAllBudgetsHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/budgets", CreateBudgetHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/budgets/{id}", GetBudgetByIDHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/budgets/{id}", UpdateBudgetHandler).Methods(http.MethodPut)
	router.HandleFunc("/api/budgets/{id}", DeleteBudgetHandler).Methods(http.MethodDelete)
	return router
}

func budgetRequest(t *testing.T, method, path string, body []byte) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(method, path, bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	newBudgetTestRouter().ServeHTTP(recorder, request)
	return recorder
}

func TestGetAllBudgetsHandler(t *testing.T) {
	database := newBudgetHandlerDB(t)
	if err := database.Create(&Budget{CategoryID: 1, Amount: 100}).Error; err != nil {
		t.Fatalf("failed to seed budget: %v", err)
	}

	recorder := budgetRequest(t, http.MethodGet, "/api/budgets", nil)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var budgets []Budget
	if err := json.NewDecoder(recorder.Body).Decode(&budgets); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(budgets) != 1 || budgets[0].Amount != 100 {
		t.Errorf("budgets = %+v", budgets)
	}
}

func TestCreateBudgetHandler(t *testing.T) {
	newBudgetHandlerDB(t)

	recorder := budgetRequest(t, http.MethodPost, "/api/budgets", []byte(`{"category_id":1,"amount":300}`))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	var created Budget
	if err := json.NewDecoder(recorder.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if created.ID == 0 || created.Amount != 300 {
		t.Errorf("created budget = %+v", created)
	}
}

func TestCreateBudgetHandlerInvalidBody(t *testing.T) {
	newBudgetHandlerDB(t)

	recorder := budgetRequest(t, http.MethodPost, "/api/budgets", []byte(`{`))

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestGetBudgetByIDHandler(t *testing.T) {
	database := newBudgetHandlerDB(t)
	seeded := &Budget{CategoryID: 1, Amount: 125}
	if err := database.Create(seeded).Error; err != nil {
		t.Fatalf("failed to seed budget: %v", err)
	}

	recorder := budgetRequest(t, http.MethodGet, fmt.Sprintf("/api/budgets/%d", seeded.ID), nil)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var got Budget
	if err := json.NewDecoder(recorder.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got.Amount != seeded.Amount {
		t.Errorf("response = %+v, want %+v", got, seeded)
	}
}

func TestGetBudgetByIDHandlerInvalidID(t *testing.T) {
	newBudgetHandlerDB(t)
	request := mux.SetURLVars(httptest.NewRequest(http.MethodGet, "/api/budgets/abc", nil), map[string]string{"id": "abc"})
	recorder := httptest.NewRecorder()

	GetBudgetByIDHandler(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestGetBudgetByIDHandlerNotFound(t *testing.T) {
	newBudgetHandlerDB(t)

	recorder := budgetRequest(t, http.MethodGet, "/api/budgets/999", nil)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusNotFound)
	}
}

func TestUpdateBudgetHandler(t *testing.T) {
	database := newBudgetHandlerDB(t)
	seeded := &Budget{CategoryID: 1, Amount: 100}
	if err := database.Create(seeded).Error; err != nil {
		t.Fatalf("failed to seed budget: %v", err)
	}

	recorder := budgetRequest(t, http.MethodPut, fmt.Sprintf("/api/budgets/%d", seeded.ID), []byte(`{"amount":450}`))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var updated Budget
	if err := json.NewDecoder(recorder.Body).Decode(&updated); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if updated.Amount != 450 {
		t.Errorf("updated budget = %+v", updated)
	}
}

func TestUpdateBudgetHandlerNotFound(t *testing.T) {
	newBudgetHandlerDB(t)

	recorder := budgetRequest(t, http.MethodPut, "/api/budgets/999", []byte(`{"amount":1}`))

	if recorder.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusNotFound)
	}
}

func TestDeleteBudgetHandler(t *testing.T) {
	database := newBudgetHandlerDB(t)
	seeded := &Budget{CategoryID: 1, Amount: 100}
	if err := database.Create(seeded).Error; err != nil {
		t.Fatalf("failed to seed budget: %v", err)
	}

	recorder := budgetRequest(t, http.MethodDelete, fmt.Sprintf("/api/budgets/%d", seeded.ID), nil)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
}

func TestDeleteBudgetHandlerNotFound(t *testing.T) {
	newBudgetHandlerDB(t)

	recorder := budgetRequest(t, http.MethodDelete, "/api/budgets/999", nil)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusNotFound)
	}
}
