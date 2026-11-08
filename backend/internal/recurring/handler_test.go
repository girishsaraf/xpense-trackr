package recurring

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func recurringRequest(t *testing.T, method, path string, body []byte) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(method, path, bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	newRecurringTestRouter().ServeHTTP(recorder, request)
	return recorder
}

func TestRecurringExpenseHandlers(t *testing.T) {
	database := newRecurringTestDB(t)
	body := []byte(`{"category_id":1,"amount":49.99,"description":"Internet","start_date":"2026-01-01T00:00:00Z","end_date":"2026-12-31T00:00:00Z","frequency":"monthly"}`)

	created := recurringRequest(t, http.MethodPost, "/api/recurring/expenses", body)
	if created.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d; body = %s", created.Code, http.StatusCreated, created.Body.String())
	}
	var recurringExpense RecurringExpense
	if err := json.NewDecoder(created.Body).Decode(&recurringExpense); err != nil {
		t.Fatalf("failed to decode create response: %v", err)
	}
	if recurringExpense.ID == 0 || recurringExpense.Description != "Internet" {
		t.Fatalf("created recurring expense = %+v", recurringExpense)
	}

	list := recurringRequest(t, http.MethodGet, "/api/recurring/expenses", nil)
	if list.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d", list.Code, http.StatusOK)
	}
	var listResponse []RecurringExpense
	if err := json.NewDecoder(list.Body).Decode(&listResponse); err != nil {
		t.Fatalf("failed to decode list response: %v", err)
	}
	if len(listResponse) != 1 {
		t.Fatalf("list returned %d expenses, want 1", len(listResponse))
	}

	path := fmt.Sprintf("/api/recurring/expenses/%d", recurringExpense.ID)
	got := recurringRequest(t, http.MethodGet, path, nil)
	if got.Code != http.StatusOK {
		t.Fatalf("get status = %d, want %d", got.Code, http.StatusOK)
	}

	updated := recurringRequest(t, http.MethodPut, path, []byte(`{"category_id":1,"description":"Updated","amount":59.99,"start_date":"2026-01-01T00:00:00Z","end_date":null,"frequency":"monthly"}`))
	if updated.Code != http.StatusOK {
		t.Fatalf("update status = %d, want %d", updated.Code, http.StatusOK)
	}
	var updatedResponse RecurringExpense
	if err := json.NewDecoder(updated.Body).Decode(&updatedResponse); err != nil {
		t.Fatalf("failed to decode update response: %v", err)
	}
	if updatedResponse.Description != "Updated" || updatedResponse.Amount != 59.99 || updatedResponse.EndDate != nil {
		t.Errorf("updated recurring expense = %+v", updatedResponse)
	}

	deleted := recurringRequest(t, http.MethodDelete, path, nil)
	if deleted.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want %d", deleted.Code, http.StatusNoContent)
	}
	var count int64
	if err := database.Model(&RecurringExpense{}).Count(&count).Error; err != nil {
		t.Fatalf("count query failed: %v", err)
	}
	if count != 0 {
		t.Errorf("recurring expense count = %d, want 0", count)
	}
}

func TestRecurringExpenseHandlerErrors(t *testing.T) {
	newRecurringTestDB(t)

	invalidBody := recurringRequest(t, http.MethodPost, "/api/recurring/expenses", []byte(`{`))
	if invalidBody.Code != http.StatusBadRequest {
		t.Errorf("invalid body status = %d, want %d", invalidBody.Code, http.StatusBadRequest)
	}

	notFound := recurringRequest(t, http.MethodGet, "/api/recurring/expenses/999", nil)
	if notFound.Code != http.StatusNotFound {
		t.Errorf("not found status = %d, want %d", notFound.Code, http.StatusNotFound)
	}

	updateMissing := recurringRequest(t, http.MethodPut, "/api/recurring/expenses/999", []byte(`{"description":"x"}`))
	if updateMissing.Code != http.StatusNotFound {
		t.Errorf("update missing status = %d, want %d", updateMissing.Code, http.StatusNotFound)
	}

	deleteMissing := recurringRequest(t, http.MethodDelete, "/api/recurring/expenses/999", nil)
	if deleteMissing.Code != http.StatusNotFound {
		t.Errorf("delete missing status = %d, want %d", deleteMissing.Code, http.StatusNotFound)
	}
}

func TestRecurringInvestmentHandlers(t *testing.T) {
	database := newRecurringTestDB(t)
	body := []byte(`{"amount":500,"investment_type":"401k","description":"Contribution","start_date":"2026-01-01T00:00:00Z","end_date":"2026-12-31T00:00:00Z","frequency":"monthly"}`)

	created := recurringRequest(t, http.MethodPost, "/api/recurring/investments", body)
	if created.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d; body = %s", created.Code, http.StatusCreated, created.Body.String())
	}
	var recurringInvestment RecurringInvestment
	if err := json.NewDecoder(created.Body).Decode(&recurringInvestment); err != nil {
		t.Fatalf("failed to decode create response: %v", err)
	}
	if recurringInvestment.ID == 0 || recurringInvestment.InvestmentType != "401k" {
		t.Fatalf("created recurring investment = %+v", recurringInvestment)
	}

	list := recurringRequest(t, http.MethodGet, "/api/recurring/investments", nil)
	if list.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d", list.Code, http.StatusOK)
	}

	path := fmt.Sprintf("/api/recurring/investments/%d", recurringInvestment.ID)
	got := recurringRequest(t, http.MethodGet, path, nil)
	if got.Code != http.StatusOK {
		t.Fatalf("get status = %d, want %d", got.Code, http.StatusOK)
	}

	updated := recurringRequest(t, http.MethodPut, path, []byte(`{"description":"Updated","amount":750,"investment_type":"401k","start_date":"2026-01-01T00:00:00Z","end_date":null,"frequency":"monthly"}`))
	if updated.Code != http.StatusOK {
		t.Fatalf("update status = %d, want %d", updated.Code, http.StatusOK)
	}
	var updatedResponse RecurringInvestment
	if err := json.NewDecoder(updated.Body).Decode(&updatedResponse); err != nil {
		t.Fatal(err)
	}
	if updatedResponse.EndDate != nil {
		t.Errorf("update end_date = %v, want nil", updatedResponse.EndDate)
	}

	deleted := recurringRequest(t, http.MethodDelete, path, nil)
	if deleted.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want %d", deleted.Code, http.StatusNoContent)
	}
	var count int64
	if err := database.Model(&RecurringInvestment{}).Count(&count).Error; err != nil {
		t.Fatalf("count query failed: %v", err)
	}
	if count != 0 {
		t.Errorf("recurring investment count = %d, want 0", count)
	}
}

func TestRecurringInvestmentHandlerErrors(t *testing.T) {
	newRecurringTestDB(t)

	invalidBody := recurringRequest(t, http.MethodPost, "/api/recurring/investments", []byte(`{`))
	if invalidBody.Code != http.StatusBadRequest {
		t.Errorf("invalid body status = %d, want %d", invalidBody.Code, http.StatusBadRequest)
	}

	notFound := recurringRequest(t, http.MethodGet, "/api/recurring/investments/999", nil)
	if notFound.Code != http.StatusNotFound {
		t.Errorf("not found status = %d, want %d", notFound.Code, http.StatusNotFound)
	}

	updateMissing := recurringRequest(t, http.MethodPut, "/api/recurring/investments/999", []byte(`{"description":"x"}`))
	if updateMissing.Code != http.StatusNotFound {
		t.Errorf("update missing status = %d, want %d", updateMissing.Code, http.StatusNotFound)
	}

	deleteMissing := recurringRequest(t, http.MethodDelete, "/api/recurring/investments/999", nil)
	if deleteMissing.Code != http.StatusNotFound {
		t.Errorf("delete missing status = %d, want %d", deleteMissing.Code, http.StatusNotFound)
	}
}

func TestRecurringTimeFieldsRoundTrip(t *testing.T) {
	database := newRecurringTestDB(t)
	start := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	recurringExpense := &RecurringExpense{Description: "Round trip", StartDate: start, EndDate: &end, Frequency: "monthly"}
	if err := database.Create(recurringExpense).Error; err != nil {
		t.Fatalf("failed to create recurring expense: %v", err)
	}

	got, err := GetRecurringExpenseByID(int(recurringExpense.ID))
	if err != nil {
		t.Fatalf("GetRecurringExpenseByID() error = %v", err)
	}
	if got.EndDate == nil || !got.StartDate.Equal(start) || !got.EndDate.Equal(end) {
		t.Errorf("time fields = %v - %v, want %v - %v", got.StartDate, got.EndDate, start, end)
	}
}
