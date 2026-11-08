package expense

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/girishsaraf/xpense-trackr/backend/internal/category"
	"github.com/girishsaraf/xpense-trackr/backend/internal/testutil"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func newExpenseHandlerDB(t *testing.T) *gorm.DB {
	t.Helper()
	database := testutil.NewDB(t, &Expense{}, &category.Category{})
	SetDB(database)
	return database
}

func newExpenseTestRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/expenses", GetAllExpensesHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/expenses", CreateExpenseHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/expenses/cumulative", CumulativeExpensesByMonthHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/expenses/{id}", GetExpenseByIDHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/expenses/{id}", UpdateExpenseHandler).Methods(http.MethodPut)
	router.HandleFunc("/api/expenses/{id}", DeleteExpenseHandler).Methods(http.MethodDelete)
	return router
}

func doRequest(t *testing.T, router http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(method, path, bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder
}

func seedExpense(t *testing.T, database *gorm.DB, expense *Expense) *Expense {
	t.Helper()
	if err := database.Create(expense).Error; err != nil {
		t.Fatalf("failed to seed expense: %v", err)
	}
	return expense
}

func TestGetAllExpensesHandler(t *testing.T) {
	database := newExpenseHandlerDB(t)
	seedExpense(t, database, &Expense{CategoryID: 1, Amount: 10, Description: "One", Date: "2026-01-01"})
	seedExpense(t, database, &Expense{CategoryID: 1, Amount: 20, Description: "Two", Date: "2026-01-02"})

	recorder := doRequest(t, newExpenseTestRouter(), http.MethodGet, "/api/expenses", nil)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var got []Expense
	if err := json.NewDecoder(recorder.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("returned %d expenses, want 2", len(got))
	}
}

func TestGetExpenseByIDHandler(t *testing.T) {
	database := newExpenseHandlerDB(t)
	seeded := seedExpense(t, database, &Expense{CategoryID: 1, Amount: 42.5, Description: "Lunch", Date: "2026-01-03"})

	recorder := doRequest(t, newExpenseTestRouter(), http.MethodGet, fmt.Sprintf("/api/expenses/%d", seeded.ID), nil)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var got Expense
	if err := json.NewDecoder(recorder.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got.ID != seeded.ID || got.Description != seeded.Description {
		t.Errorf("response = %+v, want %+v", got, seeded)
	}
}

func TestGetExpenseByIDHandlerInvalidID(t *testing.T) {
	newExpenseHandlerDB(t)
	request := mux.SetURLVars(httptest.NewRequest(http.MethodGet, "/api/expenses/abc", nil), map[string]string{"id": "abc"})
	recorder := httptest.NewRecorder()

	GetExpenseByIDHandler(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestGetExpenseByIDHandlerNotFound(t *testing.T) {
	newExpenseHandlerDB(t)

	recorder := doRequest(t, newExpenseTestRouter(), http.MethodGet, "/api/expenses/999", nil)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusNotFound)
	}
}

func TestCreateExpenseHandler(t *testing.T) {
	database := newExpenseHandlerDB(t)
	body := []byte(`{"category_id":"1","amount":"15.75","description":"Snack","date":"2026-01-04"}`)

	recorder := doRequest(t, newExpenseTestRouter(), http.MethodPost, "/api/expenses", body)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body = %s", recorder.Code, http.StatusCreated, recorder.Body.String())
	}
	var created Expense
	if err := json.NewDecoder(recorder.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if created.ID == 0 || created.Description != "Snack" || created.Amount != 15.75 {
		t.Errorf("created expense = %+v", created)
	}
	var count int64
	if err := database.Model(&Expense{}).Count(&count).Error; err != nil {
		t.Fatalf("count query failed: %v", err)
	}
	if count != 1 {
		t.Errorf("expense count = %d, want 1", count)
	}
}

func TestCreateExpenseHandlerInvalidBody(t *testing.T) {
	newExpenseHandlerDB(t)

	recorder := doRequest(t, newExpenseTestRouter(), http.MethodPost, "/api/expenses", []byte(`{`))

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestUpdateExpenseHandler(t *testing.T) {
	database := newExpenseHandlerDB(t)
	seeded := seedExpense(t, database, &Expense{CategoryID: 1, Amount: 10, Description: "Old", Date: "2026-01-01"})
	body := []byte(`{"category_id":"2","amount":"25.5","description":"New","date":"2026-01-05"}`)

	recorder := doRequest(t, newExpenseTestRouter(), http.MethodPut, fmt.Sprintf("/api/expenses/%d", seeded.ID), body)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	var updated Expense
	if err := json.NewDecoder(recorder.Body).Decode(&updated); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if updated.Description != "New" || updated.Amount != 25.5 {
		t.Errorf("updated expense = %+v", updated)
	}
}

func TestUpdateExpenseHandlerNotFound(t *testing.T) {
	newExpenseHandlerDB(t)
	body := []byte(`{"description":"Missing"}`)

	recorder := doRequest(t, newExpenseTestRouter(), http.MethodPut, "/api/expenses/999", body)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusNotFound)
	}
}

func TestDeleteExpenseHandler(t *testing.T) {
	database := newExpenseHandlerDB(t)
	seeded := seedExpense(t, database, &Expense{CategoryID: 1, Amount: 10, Description: "Delete me", Date: "2026-01-01"})

	recorder := doRequest(t, newExpenseTestRouter(), http.MethodDelete, fmt.Sprintf("/api/expenses/%d", seeded.ID), nil)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNoContent)
	}
	var count int64
	if err := database.Model(&Expense{}).Count(&count).Error; err != nil {
		t.Fatalf("count query failed: %v", err)
	}
	if count != 0 {
		t.Errorf("expense count = %d, want 0", count)
	}
}

func TestDeleteExpenseHandlerNotFound(t *testing.T) {
	newExpenseHandlerDB(t)

	recorder := doRequest(t, newExpenseTestRouter(), http.MethodDelete, "/api/expenses/999", nil)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusNotFound)
	}
}

func TestCumulativeExpensesByMonthHandler(t *testing.T) {
	database := newExpenseHandlerDB(t)
	if err := database.Create(&category.Category{ID: 1, Name: "Food"}).Error; err != nil {
		t.Fatalf("failed to seed category: %v", err)
	}
	if err := database.Create(&category.Category{ID: 2, Name: "Travel"}).Error; err != nil {
		t.Fatalf("failed to seed category: %v", err)
	}

	now := time.Now()
	currentMonth := now.Format("2006-01-02")
	previousMonth := now.AddDate(0, -1, 0).Format("2006-01-02")
	seedExpense(t, database, &Expense{CategoryID: 1, Amount: 10, Description: "Food now", Date: currentMonth})
	seedExpense(t, database, &Expense{CategoryID: 1, Amount: 20, Description: "Food prior", Date: previousMonth})
	seedExpense(t, database, &Expense{CategoryID: 2, Amount: 5, Description: "Travel now", Date: currentMonth})

	recorder := doRequest(t, newExpenseTestRouter(), http.MethodGet, "/api/expenses/cumulative?lastMonths=6", nil)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	var response []struct {
		Month    string  `json:"month"`
		Category string  `json:"category"`
		Amount   float64 `json:"amount"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	totals := make(map[string]float64)
	for _, item := range response {
		totals[item.Category] = item.Amount
	}
	if totals["Food"] != 30 {
		t.Errorf("Food cumulative amount = %v, want 30", totals["Food"])
	}
	if totals["Travel"] != 5 {
		t.Errorf("Travel cumulative amount = %v, want 5", totals["Travel"])
	}
}

func TestCumulativeExpensesByMonthHandlerInvalidLastMonths(t *testing.T) {
	newExpenseHandlerDB(t)

	for _, path := range []string{"/api/expenses/cumulative", "/api/expenses/cumulative?lastMonths=0", "/api/expenses/cumulative?lastMonths=abc"} {
		recorder := doRequest(t, newExpenseTestRouter(), http.MethodGet, path, nil)
		if recorder.Code != http.StatusBadRequest {
			t.Errorf("%s status = %d, want %d", path, recorder.Code, http.StatusBadRequest)
		}
	}
}
