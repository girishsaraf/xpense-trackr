package investments

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/girishsaraf/xpense-trackr/backend/internal/testutil"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func newInvestmentHandlerDB(t *testing.T) *gorm.DB {
	t.Helper()
	database := testutil.NewDB(t, &Investment{})
	SetDB(database)
	return database
}

func newInvestmentTestRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/investments", GetAllInvestmentsHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/investments", CreateInvestmentHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/investments/cumulative", CumulativeInvestmentsByMonthHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/investments/{id}", GetInvestmentByIDHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/investments/{id}", UpdateInvestmentHandler).Methods(http.MethodPut)
	router.HandleFunc("/api/investments/{id}", DeleteInvestmentHandler).Methods(http.MethodDelete)
	return router
}

func investmentRequest(t *testing.T, method, path string, body []byte) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(method, path, bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	newInvestmentTestRouter().ServeHTTP(recorder, request)
	return recorder
}

func TestGetAllInvestmentsHandler(t *testing.T) {
	database := newInvestmentHandlerDB(t)
	if err := database.Create(&Investment{Amount: 50, InvestmentType: "stocks", Date: "2026-01-01"}).Error; err != nil {
		t.Fatalf("failed to seed investment: %v", err)
	}

	recorder := investmentRequest(t, http.MethodGet, "/api/investments", nil)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var investments []Investment
	if err := json.NewDecoder(recorder.Body).Decode(&investments); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(investments) != 1 || investments[0].InvestmentType != "stocks" {
		t.Errorf("investments = %+v", investments)
	}
}

func TestCreateInvestmentHandler(t *testing.T) {
	newInvestmentHandlerDB(t)

	recorder := investmentRequest(t, http.MethodPost, "/api/investments", []byte(`{"amount":"200","investment_type":"401k","date":"2026-01-02"}`))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	var created Investment
	if err := json.NewDecoder(recorder.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if created.ID == 0 || created.Amount != 200 || created.InvestmentType != "401k" {
		t.Errorf("created investment = %+v", created)
	}
}

func TestCreateInvestmentHandlerInvalidBody(t *testing.T) {
	newInvestmentHandlerDB(t)

	recorder := investmentRequest(t, http.MethodPost, "/api/investments", []byte(`{`))

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestGetInvestmentByIDHandler(t *testing.T) {
	database := newInvestmentHandlerDB(t)
	seeded := &Investment{Amount: 75, InvestmentType: "bonds", Date: "2026-01-03"}
	if err := database.Create(seeded).Error; err != nil {
		t.Fatalf("failed to seed investment: %v", err)
	}

	recorder := investmentRequest(t, http.MethodGet, fmt.Sprintf("/api/investments/%d", seeded.ID), nil)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var got Investment
	if err := json.NewDecoder(recorder.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got.InvestmentType != seeded.InvestmentType {
		t.Errorf("response = %+v, want %+v", got, seeded)
	}
}

func TestGetInvestmentByIDHandlerInvalidID(t *testing.T) {
	newInvestmentHandlerDB(t)
	request := mux.SetURLVars(httptest.NewRequest(http.MethodGet, "/api/investments/abc", nil), map[string]string{"id": "abc"})
	recorder := httptest.NewRecorder()

	GetInvestmentByIDHandler(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestGetInvestmentByIDHandlerNotFound(t *testing.T) {
	newInvestmentHandlerDB(t)

	recorder := investmentRequest(t, http.MethodGet, "/api/investments/999", nil)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusNotFound)
	}
}

func TestUpdateInvestmentHandler(t *testing.T) {
	database := newInvestmentHandlerDB(t)
	seeded := &Investment{Amount: 100, InvestmentType: "stocks", Date: "2026-01-01"}
	if err := database.Create(seeded).Error; err != nil {
		t.Fatalf("failed to seed investment: %v", err)
	}

	recorder := investmentRequest(t, http.MethodPut, fmt.Sprintf("/api/investments/%d", seeded.ID), []byte(`{"amount":"300","investment_type":"real_estate"}`))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var updated Investment
	if err := json.NewDecoder(recorder.Body).Decode(&updated); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if updated.Amount != 300 || updated.InvestmentType != "real_estate" {
		t.Errorf("updated investment = %+v", updated)
	}
}

func TestUpdateInvestmentHandlerNotFound(t *testing.T) {
	newInvestmentHandlerDB(t)

	recorder := investmentRequest(t, http.MethodPut, "/api/investments/999", []byte(`{"amount":"1"}`))

	if recorder.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusNotFound)
	}
}

func TestDeleteInvestmentHandler(t *testing.T) {
	database := newInvestmentHandlerDB(t)
	seeded := &Investment{Amount: 100, InvestmentType: "savings", Date: "2026-01-01"}
	if err := database.Create(seeded).Error; err != nil {
		t.Fatalf("failed to seed investment: %v", err)
	}

	recorder := investmentRequest(t, http.MethodDelete, fmt.Sprintf("/api/investments/%d", seeded.ID), nil)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
}

func TestDeleteInvestmentHandlerNotFound(t *testing.T) {
	newInvestmentHandlerDB(t)

	recorder := investmentRequest(t, http.MethodDelete, "/api/investments/999", nil)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusNotFound)
	}
}

func TestCumulativeInvestmentsByMonthHandler(t *testing.T) {
	database := newInvestmentHandlerDB(t)
	now := time.Now()
	currentMonth := now.Format("2006-01-02")
	previousMonth := now.AddDate(0, -1, 0).Format("2006-01-02")
	if err := database.Create(&Investment{Amount: 100, InvestmentType: "stocks", Date: currentMonth}).Error; err != nil {
		t.Fatalf("failed to seed investment: %v", err)
	}
	if err := database.Create(&Investment{Amount: 50, InvestmentType: "stocks", Date: previousMonth}).Error; err != nil {
		t.Fatalf("failed to seed investment: %v", err)
	}
	if err := database.Create(&Investment{Amount: 25, InvestmentType: "bonds", Date: currentMonth}).Error; err != nil {
		t.Fatalf("failed to seed investment: %v", err)
	}

	recorder := investmentRequest(t, http.MethodGet, "/api/investments/cumulative?lastMonths=6", nil)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	var response []struct {
		Month  string  `json:"month"`
		Type   string  `json:"type"`
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	totals := make(map[string]float64)
	for _, item := range response {
		totals[item.Type] = item.Amount
	}
	if totals["stocks"] != 150 {
		t.Errorf("stocks cumulative amount = %v, want 150", totals["stocks"])
	}
	if totals["bonds"] != 25 {
		t.Errorf("bonds cumulative amount = %v, want 25", totals["bonds"])
	}
}

func TestCumulativeInvestmentsByMonthHandlerInvalidLastMonths(t *testing.T) {
	newInvestmentHandlerDB(t)

	for _, path := range []string{"/api/investments/cumulative", "/api/investments/cumulative?lastMonths=0", "/api/investments/cumulative?lastMonths=abc"} {
		recorder := investmentRequest(t, http.MethodGet, path, nil)
		if recorder.Code != http.StatusBadRequest {
			t.Errorf("%s status = %d, want %d", path, recorder.Code, http.StatusBadRequest)
		}
	}
}
