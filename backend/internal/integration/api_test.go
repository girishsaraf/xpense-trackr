package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/girishsaraf/xpense-trackr/backend/internal"
	"github.com/girishsaraf/xpense-trackr/backend/internal/budget"
	"github.com/girishsaraf/xpense-trackr/backend/internal/category"
	"github.com/girishsaraf/xpense-trackr/backend/internal/expense"
	"github.com/girishsaraf/xpense-trackr/backend/internal/investments"
	"github.com/girishsaraf/xpense-trackr/backend/internal/recurring"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupIntegrationDB connects to the MySQL instance configured by
// TEST_MYSQL_DSN, recreates the schema, and returns the full API handler.
// The test is skipped when TEST_MYSQL_DSN is not set, so `go test ./...`
// works without a database.
func setupIntegrationDB(t *testing.T) (*gorm.DB, http.Handler) {
	t.Helper()

	dsn := os.Getenv("TEST_MYSQL_DSN")
	if dsn == "" {
		t.Skip("set TEST_MYSQL_DSN to run MySQL integration tests")
	}

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to connect to MySQL: %v", err)
	}

	models := []interface{}{
		&category.Category{},
		&expense.Expense{},
		&budget.Budget{},
		&investments.Investment{},
		&recurring.RecurringExpense{},
		&recurring.RecurringInvestment{},
	}
	if err := database.Migrator().DropTable(models...); err != nil {
		t.Fatalf("failed to reset schema: %v", err)
	}
	if err := database.AutoMigrate(models...); err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	budget.SetDB(database)
	category.SetDB(database)
	expense.SetDB(database)
	investments.SetDB(database)
	recurring.SetDB(database)

	router := mux.NewRouter()
	internal.SetRoutes(router)
	return database, router
}

func apiRequest(t *testing.T, handler http.Handler, method, path string, payload interface{}) *httptest.ResponseRecorder {
	t.Helper()

	var body []byte
	if payload != nil {
		var err error
		body, err = json.Marshal(payload)
		if err != nil {
			t.Fatalf("failed to encode request payload: %v", err)
		}
	}
	request := httptest.NewRequest(method, path, bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	return recorder
}

func decodeResponse[T any](t *testing.T, recorder *httptest.ResponseRecorder) T {
	t.Helper()
	var response T
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response %q: %v", recorder.Body.String(), err)
	}
	return response
}

func TestExpenseCategoryAPI(t *testing.T) {
	_, handler := setupIntegrationDB(t)

	categoryResponse := apiRequest(t, handler, http.MethodPost, "/api/categories", map[string]string{"name": "Food"})
	if categoryResponse.Code != http.StatusCreated {
		t.Fatalf("create category status = %d, want %d", categoryResponse.Code, http.StatusCreated)
	}
	createdCategory := decodeResponse[category.Category](t, categoryResponse)

	expenseResponse := apiRequest(t, handler, http.MethodPost, "/api/expenses", map[string]string{
		"category_id": fmt.Sprint(createdCategory.ID),
		"amount":      "42.50",
		"description": "Groceries",
		"date":        time.Now().Format("2006-01-02"),
	})
	if expenseResponse.Code != http.StatusCreated {
		t.Fatalf("create expense status = %d, want %d; body = %s", expenseResponse.Code, http.StatusCreated, expenseResponse.Body.String())
	}
	createdExpense := decodeResponse[expense.Expense](t, expenseResponse)

	listResponse := apiRequest(t, handler, http.MethodGet, "/api/expenses", nil)
	if listResponse.Code != http.StatusOK {
		t.Fatalf("list expenses status = %d, want %d", listResponse.Code, http.StatusOK)
	}
	expenses := decodeResponse[[]expense.Expense](t, listResponse)
	if len(expenses) != 1 || expenses[0].Description != "Groceries" {
		t.Fatalf("expenses = %+v", expenses)
	}

	cumulativeResponse := apiRequest(t, handler, http.MethodGet, "/api/expenses/cumulative?lastMonths=6", nil)
	if cumulativeResponse.Code != http.StatusOK {
		t.Fatalf("cumulative expenses status = %d, want %d", cumulativeResponse.Code, http.StatusOK)
	}
	cumulative := decodeResponse[[]struct {
		Month    string  `json:"month"`
		Category string  `json:"category"`
		Amount   float64 `json:"amount"`
	}](t, cumulativeResponse)
	if len(cumulative) != 1 || cumulative[0].Category != "Food" || cumulative[0].Amount != 42.5 {
		t.Fatalf("cumulative expenses = %+v", cumulative)
	}

	updateResponse := apiRequest(t, handler, http.MethodPut, fmt.Sprintf("/api/expenses/%d", createdExpense.ID), map[string]string{
		"description": "Updated groceries",
		"amount":      "50",
	})
	if updateResponse.Code != http.StatusOK {
		t.Fatalf("update expense status = %d, want %d", updateResponse.Code, http.StatusOK)
	}

	deleteResponse := apiRequest(t, handler, http.MethodDelete, fmt.Sprintf("/api/expenses/%d", createdExpense.ID), nil)
	if deleteResponse.Code != http.StatusNoContent {
		t.Fatalf("delete expense status = %d, want %d", deleteResponse.Code, http.StatusNoContent)
	}

	deleteCategoryResponse := apiRequest(t, handler, http.MethodDelete, fmt.Sprintf("/api/categories/%d", createdCategory.ID), nil)
	if deleteCategoryResponse.Code != http.StatusNoContent {
		t.Fatalf("delete category status = %d, want %d", deleteCategoryResponse.Code, http.StatusNoContent)
	}
}

func TestBudgetAndInvestmentAPI(t *testing.T) {
	_, handler := setupIntegrationDB(t)

	budgetResponse := apiRequest(t, handler, http.MethodPost, "/api/budgets", map[string]interface{}{
		"category_id": 1,
		"amount":      500,
	})
	if budgetResponse.Code != http.StatusOK {
		t.Fatalf("create budget status = %d, want %d", budgetResponse.Code, http.StatusOK)
	}
	createdBudget := decodeResponse[budget.Budget](t, budgetResponse)

	updateBudgetResponse := apiRequest(t, handler, http.MethodPut, fmt.Sprintf("/api/budgets/%d", createdBudget.ID), map[string]interface{}{
		"amount": 750,
	})
	if updateBudgetResponse.Code != http.StatusOK {
		t.Fatalf("update budget status = %d, want %d", updateBudgetResponse.Code, http.StatusOK)
	}

	deleteBudgetResponse := apiRequest(t, handler, http.MethodDelete, fmt.Sprintf("/api/budgets/%d", createdBudget.ID), nil)
	if deleteBudgetResponse.Code != http.StatusOK {
		t.Fatalf("delete budget status = %d, want %d", deleteBudgetResponse.Code, http.StatusOK)
	}

	investmentResponse := apiRequest(t, handler, http.MethodPost, "/api/investments", map[string]string{
		"amount":          "250",
		"investment_type": "stocks",
		"date":            time.Now().Format("2006-01-02"),
	})
	if investmentResponse.Code != http.StatusOK {
		t.Fatalf("create investment status = %d, want %d", investmentResponse.Code, http.StatusOK)
	}
	createdInvestment := decodeResponse[investments.Investment](t, investmentResponse)

	cumulativeResponse := apiRequest(t, handler, http.MethodGet, "/api/investments/cumulative?lastMonths=6", nil)
	if cumulativeResponse.Code != http.StatusOK {
		t.Fatalf("cumulative investments status = %d, want %d", cumulativeResponse.Code, http.StatusOK)
	}
	cumulative := decodeResponse[[]struct {
		Month  string  `json:"month"`
		Type   string  `json:"type"`
		Amount float64 `json:"amount"`
	}](t, cumulativeResponse)
	if len(cumulative) != 1 || cumulative[0].Type != "stocks" || cumulative[0].Amount != 250 {
		t.Fatalf("cumulative investments = %+v", cumulative)
	}

	deleteInvestmentResponse := apiRequest(t, handler, http.MethodDelete, fmt.Sprintf("/api/investments/%d", createdInvestment.ID), nil)
	if deleteInvestmentResponse.Code != http.StatusOK {
		t.Fatalf("delete investment status = %d, want %d", deleteInvestmentResponse.Code, http.StatusOK)
	}
}

func TestRecurringAPI(t *testing.T) {
	_, handler := setupIntegrationDB(t)
	now := time.Now()

	recurringExpenseResponse := apiRequest(t, handler, http.MethodPost, "/api/recurring/expenses", map[string]interface{}{
		"category_id": 1,
		"amount":      49.99,
		"description": "Internet",
		"start_date":  now.Format(time.RFC3339),
		"end_date":    now.AddDate(1, 0, 0).Format(time.RFC3339),
		"frequency":   "monthly",
	})
	if recurringExpenseResponse.Code != http.StatusCreated {
		t.Fatalf("create recurring expense status = %d, want %d; body = %s", recurringExpenseResponse.Code, http.StatusCreated, recurringExpenseResponse.Body.String())
	}
	createdExpense := decodeResponse[recurring.RecurringExpense](t, recurringExpenseResponse)

	listResponse := apiRequest(t, handler, http.MethodGet, "/api/recurring/expenses", nil)
	if listResponse.Code != http.StatusOK {
		t.Fatalf("list recurring expenses status = %d, want %d", listResponse.Code, http.StatusOK)
	}

	deleteResponse := apiRequest(t, handler, http.MethodDelete, fmt.Sprintf("/api/recurring/expenses/%d", createdExpense.ID), nil)
	if deleteResponse.Code != http.StatusNoContent {
		t.Fatalf("delete recurring expense status = %d, want %d", deleteResponse.Code, http.StatusNoContent)
	}

	recurringInvestmentResponse := apiRequest(t, handler, http.MethodPost, "/api/recurring/investments", map[string]interface{}{
		"amount":          500,
		"investment_type": "401k",
		"description":     "Contribution",
		"start_date":      now.Format(time.RFC3339),
		"end_date":        now.AddDate(1, 0, 0).Format(time.RFC3339),
		"frequency":       "monthly",
	})
	if recurringInvestmentResponse.Code != http.StatusCreated {
		t.Fatalf("create recurring investment status = %d, want %d; body = %s", recurringInvestmentResponse.Code, http.StatusCreated, recurringInvestmentResponse.Body.String())
	}
	createdInvestment := decodeResponse[recurring.RecurringInvestment](t, recurringInvestmentResponse)

	deleteInvestmentResponse := apiRequest(t, handler, http.MethodDelete, fmt.Sprintf("/api/recurring/investments/%d", createdInvestment.ID), nil)
	if deleteInvestmentResponse.Code != http.StatusNoContent {
		t.Fatalf("delete recurring investment status = %d, want %d", deleteInvestmentResponse.Code, http.StatusNoContent)
	}
}
