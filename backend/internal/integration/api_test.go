package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/girishsaraf/xpense-trackr/backend/internal"
	"github.com/girishsaraf/xpense-trackr/backend/internal/budget"
	"github.com/girishsaraf/xpense-trackr/backend/internal/category"
	"github.com/girishsaraf/xpense-trackr/backend/internal/expense"
	"github.com/girishsaraf/xpense-trackr/backend/internal/goal"
	"github.com/girishsaraf/xpense-trackr/backend/internal/investments"
	"github.com/girishsaraf/xpense-trackr/backend/internal/ledger"
	"github.com/girishsaraf/xpense-trackr/backend/internal/planning"
	"github.com/girishsaraf/xpense-trackr/backend/internal/receipt"
	"github.com/girishsaraf/xpense-trackr/backend/internal/recurring"
	"github.com/girishsaraf/xpense-trackr/backend/internal/trends"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var integrationDBCounter int64

// setupIntegrationDB uses an isolated SQLite database by default and MySQL
// when TEST_MYSQL_DSN is set.
func setupIntegrationDB(t *testing.T) (*gorm.DB, http.Handler) {
	t.Helper()

	dsn := os.Getenv("TEST_MYSQL_DSN")
	dialector := sqlite.Open(fmt.Sprintf("file:integration-%d?mode=memory&cache=shared", atomic.AddInt64(&integrationDBCounter, 1)))
	if dsn != "" {
		dialector = mysql.Open(dsn)
	}
	database, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to connect to MySQL: %v", err)
	}

	models := []interface{}{
		&category.Category{},
		&expense.Expense{},
		&budget.Budget{},
		&goal.Goal{},
		&investments.Investment{},
		&ledger.Account{},
		&ledger.Transaction{},
		&ledger.CategoryRule{},
		&ledger.ScheduledFlow{},
		&ledger.ScheduledOccurrence{},
		&planning.Goal{},
		&planning.GoalContribution{},
		&planning.Envelope{},
		&receipt.Receipt{},
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
	ledger.SetDB(database)
	planning.SetDB(database)
	recurring.SetDB(database)
	receipt.SetDB(database)
	goal.SetDB(database)
	trends.SetDB(database)

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

func apiMultipartRequest(t *testing.T, handler http.Handler, path, transactionID, fileName string, content []byte) *httptest.ResponseRecorder {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("transaction_id", transactionID); err != nil {
		t.Fatalf("write transaction id: %v", err)
	}
	file, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		t.Fatalf("create multipart file: %v", err)
	}
	if _, err := file.Write(content); err != nil {
		t.Fatalf("write multipart file: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, path, &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
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

func TestAccountAndUnifiedTransactionAPI(t *testing.T) {
	database, handler := setupIntegrationDB(t)
	categoryValue := category.Category{Name: "Groceries"}
	if err := database.Create(&categoryValue).Error; err != nil {
		t.Fatalf("create category: %v", err)
	}

	checkingResponse := apiRequest(t, handler, http.MethodPost, "/api/accounts", map[string]interface{}{
		"name": "Checking", "type": "checking", "currency": "USD", "opening_balance_minor": 50000,
	})
	if checkingResponse.Code != http.StatusCreated {
		t.Fatalf("create checking status = %d, want %d; body = %s", checkingResponse.Code, http.StatusCreated, checkingResponse.Body.String())
	}
	checking := decodeResponse[ledger.Account](t, checkingResponse)

	savingsResponse := apiRequest(t, handler, http.MethodPost, "/api/accounts", map[string]interface{}{
		"name": "Savings", "type": "savings", "currency": "USD", "opening_balance_minor": 1000,
	})
	if savingsResponse.Code != http.StatusCreated {
		t.Fatalf("create savings status = %d, want %d; body = %s", savingsResponse.Code, http.StatusCreated, savingsResponse.Body.String())
	}
	savings := decodeResponse[ledger.Account](t, savingsResponse)

	expenseResponse := apiRequest(t, handler, http.MethodPost, "/api/transactions", map[string]interface{}{
		"account_id": checking.ID, "category_id": categoryValue.ID, "type": "expense",
		"amount_minor": 1250, "currency": "USD", "payee": "Market", "date": "2026-10-12",
	})
	if expenseResponse.Code != http.StatusCreated {
		t.Fatalf("create expense transaction status = %d, want %d; body = %s", expenseResponse.Code, http.StatusCreated, expenseResponse.Body.String())
	}

	transferResponse := apiRequest(t, handler, http.MethodPost, "/api/transactions", map[string]interface{}{
		"account_id": checking.ID, "transfer_account_id": savings.ID, "type": "transfer",
		"amount_minor": 1000, "currency": "USD", "date": "2026-10-13",
	})
	if transferResponse.Code != http.StatusCreated {
		t.Fatalf("create transfer status = %d, want %d; body = %s", transferResponse.Code, http.StatusCreated, transferResponse.Body.String())
	}

	accountsResponse := apiRequest(t, handler, http.MethodGet, "/api/accounts", nil)
	if accountsResponse.Code != http.StatusOK {
		t.Fatalf("list accounts status = %d, want %d", accountsResponse.Code, http.StatusOK)
	}
	accounts := decodeResponse[[]ledger.Account](t, accountsResponse)
	balances := map[uint]int64{}
	for _, account := range accounts {
		balances[account.ID] = account.BalanceMinor
	}
	if balances[checking.ID] != 47750 || balances[savings.ID] != 2000 {
		t.Errorf("balances = %+v, want checking 47750 and savings 2000", balances)
	}

	transactionsResponse := apiRequest(t, handler, http.MethodGet, "/api/transactions?type=transfer", nil)
	if transactionsResponse.Code != http.StatusOK {
		t.Fatalf("list transactions status = %d, want %d", transactionsResponse.Code, http.StatusOK)
	}
	transactions := decodeResponse[[]ledger.Transaction](t, transactionsResponse)
	if len(transactions) != 1 || transactions[0].Type != "transfer" {
		t.Errorf("transfer transactions = %+v, want exactly one transfer", transactions)
	}
}

func TestMerchantCategoryRuleAPI(t *testing.T) {
	database, handler := setupIntegrationDB(t)
	categoryValue := category.Category{Name: "Groceries"}
	if err := database.Create(&categoryValue).Error; err != nil {
		t.Fatalf("create category: %v", err)
	}
	accountResponse := apiRequest(t, handler, http.MethodPost, "/api/accounts", map[string]interface{}{
		"name": "Checking", "type": "checking", "currency": "USD",
	})
	if accountResponse.Code != http.StatusCreated {
		t.Fatalf("create account status = %d: %s", accountResponse.Code, accountResponse.Body)
	}
	account := decodeResponse[ledger.Account](t, accountResponse)
	ruleResponse := apiRequest(t, handler, http.MethodPost, "/api/category-rules", map[string]interface{}{
		"name": "Groceries", "match_value": "market", "category_id": categoryValue.ID,
		"priority": 5, "enabled": true,
	})
	if ruleResponse.Code != http.StatusCreated {
		t.Fatalf("create category rule status = %d, want %d: %s", ruleResponse.Code, http.StatusCreated, ruleResponse.Body)
	}
	rule := decodeResponse[ledger.CategoryRule](t, ruleResponse)

	transactionResponse := apiRequest(t, handler, http.MethodPost, "/api/transactions", map[string]interface{}{
		"account_id": account.ID, "type": "expense", "amount_minor": 1250,
		"currency": "USD", "payee": "Neighborhood Market", "date": "2026-10-12",
	})
	if transactionResponse.Code != http.StatusCreated {
		t.Fatalf("create transaction status = %d, want %d: %s", transactionResponse.Code, http.StatusCreated, transactionResponse.Body)
	}
	transaction := decodeResponse[ledger.Transaction](t, transactionResponse)
	if transaction.CategoryID == nil || *transaction.CategoryID != categoryValue.ID || transaction.CategoryRuleID == nil || *transaction.CategoryRuleID != rule.ID {
		t.Errorf("categorized transaction = %+v, want category/rule %d/%d", transaction, categoryValue.ID, rule.ID)
	}

	listResponse := apiRequest(t, handler, http.MethodGet, "/api/category-rules", nil)
	if listResponse.Code != http.StatusOK {
		t.Fatalf("list category rules status = %d, want %d", listResponse.Code, http.StatusOK)
	}
	rules := decodeResponse[[]ledger.CategoryRule](t, listResponse)
	if len(rules) != 1 || rules[0].ID != rule.ID {
		t.Fatalf("category rules = %+v, want the created rule", rules)
	}
}

func TestScheduledCashflowAPI(t *testing.T) {
	_, handler := setupIntegrationDB(t)
	accountResponse := apiRequest(t, handler, http.MethodPost, "/api/accounts", map[string]interface{}{
		"name": "Checking", "type": "checking", "currency": "USD", "opening_balance_minor": 5000,
	})
	if accountResponse.Code != http.StatusCreated {
		t.Fatalf("create account status = %d: %s", accountResponse.Code, accountResponse.Body)
	}
	account := decodeResponse[ledger.Account](t, accountResponse)
	nextDate := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	flowResponse := apiRequest(t, handler, http.MethodPost, "/api/scheduled-flows", map[string]interface{}{
		"account_id": account.ID, "label": "Phone bill", "type": "expense",
		"amount_minor": 1250, "currency": "USD", "frequency": "once",
		"next_date": nextDate, "enabled": true,
	})
	if flowResponse.Code != http.StatusCreated {
		t.Fatalf("create scheduled flow status = %d, want %d: %s", flowResponse.Code, http.StatusCreated, flowResponse.Body)
	}

	forecastResponse := apiRequest(t, handler, http.MethodGet, "/api/forecast/cashflow?days=30", nil)
	if forecastResponse.Code != http.StatusOK {
		t.Fatalf("cashflow forecast status = %d, want %d: %s", forecastResponse.Code, http.StatusOK, forecastResponse.Body)
	}
	forecast := decodeResponse[ledger.CashflowForecast](t, forecastResponse)
	if len(forecast.Events) != 1 || forecast.Accounts[0].ProjectedBalanceMinor != 3750 {
		t.Fatalf("cashflow forecast = %+v, want one bill and projected balance 3750", forecast)
	}
	billsResponse := apiRequest(t, handler, http.MethodGet, "/api/bills", nil)
	if billsResponse.Code != http.StatusOK {
		t.Fatalf("bill reminders status = %d, want %d: %s", billsResponse.Code, http.StatusOK, billsResponse.Body)
	}
	bills := decodeResponse[[]ledger.BillInsight](t, billsResponse)
	if len(bills) != 1 || bills[0].Status != "due_soon" {
		t.Fatalf("bill reminders = %+v, want one due-soon bill", bills)
	}
	occurrenceResponse := apiRequest(t, handler, http.MethodPost,
		fmt.Sprintf("/api/scheduled-flows/%d/occurrences", bills[0].ID),
		map[string]interface{}{"scheduled_date": nextDate, "action": "record", "amount_minor": 1300},
	)
	if occurrenceResponse.Code != http.StatusOK {
		t.Fatalf("record bill occurrence status = %d, want %d: %s", occurrenceResponse.Code, http.StatusOK, occurrenceResponse.Body)
	}
	occurrence := decodeResponse[ledger.ScheduledOccurrence](t, occurrenceResponse)
	if occurrence.TransactionID == nil || occurrence.ActualAmountMinor != 1300 {
		t.Fatalf("recorded occurrence = %+v", occurrence)
	}
	retryResponse := apiRequest(t, handler, http.MethodPost,
		fmt.Sprintf("/api/scheduled-flows/%d/occurrences", bills[0].ID),
		map[string]interface{}{"scheduled_date": nextDate, "action": "record", "amount_minor": 1300},
	)
	retry := decodeResponse[ledger.ScheduledOccurrence](t, retryResponse)
	if retry.ID != occurrence.ID || retry.TransactionID == nil || *retry.TransactionID != *occurrence.TransactionID {
		t.Fatalf("idempotent occurrence retry = %+v, initial = %+v", retry, occurrence)
	}
	accountsResponse := apiRequest(t, handler, http.MethodGet, "/api/accounts", nil)
	accounts := decodeResponse[[]ledger.Account](t, accountsResponse)
	if len(accounts) != 1 || accounts[0].BalanceMinor != 3700 {
		t.Fatalf("balance after recording bill = %+v, want 3700", accounts)
	}
	if invalidDays := apiRequest(t, handler, http.MethodGet, "/api/forecast/cashflow?days=500", nil); invalidDays.Code != http.StatusBadRequest {
		t.Errorf("invalid forecast range status = %d, want %d", invalidDays.Code, http.StatusBadRequest)
	}
}

func TestReceiptUploadDownloadAndDelete(t *testing.T) {
	_, handler := setupIntegrationDB(t)
	receipt.SetStorageDir(t.TempDir())
	accountResponse := apiRequest(t, handler, http.MethodPost, "/api/accounts", map[string]interface{}{
		"name": "Checking", "type": "checking", "currency": "USD",
	})
	account := decodeResponse[ledger.Account](t, accountResponse)
	transactionResponse := apiRequest(t, handler, http.MethodPost, "/api/transactions", map[string]interface{}{
		"account_id": account.ID, "type": "expense", "amount_minor": 1250,
		"currency": "USD", "payee": "Market", "date": "2026-10-12",
	})
	if transactionResponse.Code != http.StatusCreated {
		t.Fatalf("create transaction status = %d: %s", transactionResponse.Code, transactionResponse.Body)
	}
	transaction := decodeResponse[ledger.Transaction](t, transactionResponse)
	pdfContent := []byte("%PDF-1.4\nreceipt fixture")
	upload := apiMultipartRequest(t, handler, "/api/receipts", fmt.Sprint(transaction.ID), "receipt.pdf", pdfContent)
	if upload.Code != http.StatusCreated {
		t.Fatalf("upload receipt status = %d, want %d: %s", upload.Code, http.StatusCreated, upload.Body)
	}
	item := decodeResponse[receipt.Receipt](t, upload)
	if item.TransactionID != transaction.ID || item.FileName != "receipt.pdf" || item.SizeBytes != int64(len(pdfContent)) {
		t.Fatalf("receipt metadata = %+v", item)
	}

	fileResponse := apiRequest(t, handler, http.MethodGet, fmt.Sprintf("/api/receipts/%d/file", item.ID), nil)
	if fileResponse.Code != http.StatusOK || !bytes.Equal(fileResponse.Body.Bytes(), pdfContent) {
		t.Fatalf("download status/body = %d/%q", fileResponse.Code, fileResponse.Body.Bytes())
	}
	if fileResponse.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Errorf("receipt download missing nosniff header")
	}
	listResponse := apiRequest(t, handler, http.MethodGet, fmt.Sprintf("/api/receipts?transaction_id=%d", transaction.ID), nil)
	listed := decodeResponse[[]receipt.Receipt](t, listResponse)
	if len(listed) != 1 || listed[0].ID != item.ID {
		t.Fatalf("listed receipts = %+v", listed)
	}

	deleteResponse := apiRequest(t, handler, http.MethodDelete, fmt.Sprintf("/api/receipts/%d", item.ID), nil)
	if deleteResponse.Code != http.StatusNoContent {
		t.Fatalf("delete receipt status = %d, want %d", deleteResponse.Code, http.StatusNoContent)
	}
	if missing := apiRequest(t, handler, http.MethodGet, fmt.Sprintf("/api/receipts/%d/file", item.ID), nil); missing.Code != http.StatusNotFound {
		t.Errorf("download deleted receipt status = %d, want %d", missing.Code, http.StatusNotFound)
	}
	if invalid := apiMultipartRequest(t, handler, "/api/receipts", fmt.Sprint(transaction.ID), "bad.txt", []byte("not a receipt")); invalid.Code != http.StatusBadRequest {
		t.Errorf("upload invalid receipt type status = %d, want %d", invalid.Code, http.StatusBadRequest)
	}
}

func TestCashflowAndNetWorthReportsAPI(t *testing.T) {
	_, handler := setupIntegrationDB(t)
	accountResponse := apiRequest(t, handler, http.MethodPost, "/api/accounts", map[string]interface{}{
		"name": "Checking", "type": "checking", "currency": "USD", "opening_balance_minor": 10000,
	})
	account := decodeResponse[ledger.Account](t, accountResponse)
	transactionResponse := apiRequest(t, handler, http.MethodPost, "/api/transactions", map[string]interface{}{
		"account_id": account.ID, "type": "expense", "amount_minor": 1250,
		"currency": "USD", "date": "2026-10-12", "payee": "Market",
	})
	if transactionResponse.Code != http.StatusCreated {
		t.Fatalf("create transaction status = %d: %s", transactionResponse.Code, transactionResponse.Body)
	}

	reportResponse := apiRequest(t, handler, http.MethodGet, "/api/reports/cashflow?from=2026-10-01&to=2026-10-31", nil)
	if reportResponse.Code != http.StatusOK {
		t.Fatalf("cashflow report status = %d: %s", reportResponse.Code, reportResponse.Body)
	}
	report := decodeResponse[ledger.CashflowReport](t, reportResponse)
	if len(report.Currencies) != 1 || report.Currencies[0].ExpenseMinor != 1250 || report.Currencies[0].NetMinor != -1250 {
		t.Fatalf("cashflow report = %+v", report)
	}

	csvResponse := apiRequest(t, handler, http.MethodGet, "/api/reports/cashflow.csv?from=2026-10-01&to=2026-10-31", nil)
	if csvResponse.Code != http.StatusOK || !strings.Contains(csvResponse.Body.String(), "USD,expense,Uncategorized,1250") {
		t.Fatalf("cashflow CSV status/body = %d/%s", csvResponse.Code, csvResponse.Body.String())
	}
	netWorthResponse := apiRequest(t, handler, http.MethodGet, "/api/net-worth", nil)
	netWorth := decodeResponse[ledger.NetWorthReport](t, netWorthResponse)
	if len(netWorth.Currencies) != 1 || netWorth.Currencies[0].NetWorthMinor != 8750 {
		t.Fatalf("net worth report = %+v", netWorth)
	}
}

func TestTransactionCSVImportPreviewCommitAndExport(t *testing.T) {
	_, handler := setupIntegrationDB(t)
	accountResponse := apiRequest(t, handler, http.MethodPost, "/api/accounts", map[string]interface{}{
		"name": "Checking", "type": "checking", "currency": "USD",
	})
	account := decodeResponse[ledger.Account](t, accountResponse)
	csvContent := []byte(fmt.Sprintf(
		"account_id,type,amount,currency,date,payee,description,category_id,transfer_account_id\n"+
			"%d,expense,12.34,USD,2026-10-12,Market,Groceries,,\n",
		account.ID,
	))

	previewResponse := apiMultipartRequest(t, handler, "/api/transactions/import.csv", "", "transactions.csv", csvContent)
	if previewResponse.Code != http.StatusOK {
		t.Fatalf("preview status = %d, want %d: %s", previewResponse.Code, http.StatusOK, previewResponse.Body)
	}
	preview := decodeResponse[ledger.CSVImportReport](t, previewResponse)
	if preview.Ready != 1 || preview.Imported != 0 {
		t.Fatalf("preview report = %+v", preview)
	}

	commitResponse := apiMultipartRequest(t, handler, "/api/transactions/import.csv?commit=true", "", "transactions.csv", csvContent)
	if commitResponse.Code != http.StatusOK {
		t.Fatalf("commit status = %d, want %d: %s", commitResponse.Code, http.StatusOK, commitResponse.Body)
	}
	committed := decodeResponse[ledger.CSVImportReport](t, commitResponse)
	if committed.Imported != 1 {
		t.Fatalf("commit report = %+v", committed)
	}

	exportResponse := apiRequest(t, handler, http.MethodGet, "/api/transactions/export.csv", nil)
	if exportResponse.Code != http.StatusOK || !strings.Contains(exportResponse.Body.String(), "12.34") {
		t.Fatalf("export status/body = %d/%s", exportResponse.Code, exportResponse.Body.String())
	}
	repeatResponse := apiMultipartRequest(t, handler, "/api/transactions/import.csv?commit=true", "", "transactions.csv", csvContent)
	repeated := decodeResponse[ledger.CSVImportReport](t, repeatResponse)
	if repeated.Imported != 0 || repeated.Duplicates != 1 {
		t.Fatalf("repeated import report = %+v", repeated)
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

func TestSavingsGoalsAndSpendingTrendsAPI(t *testing.T) {
	_, handler := setupIntegrationDB(t)

	categoryResponse := apiRequest(t, handler, http.MethodPost, "/api/categories", map[string]string{"name": "Food"})
	if categoryResponse.Code != http.StatusCreated {
		t.Fatalf("create category status = %d, want %d", categoryResponse.Code, http.StatusCreated)
	}
	createdCategory := decodeResponse[category.Category](t, categoryResponse)
	expenseResponse := apiRequest(t, handler, http.MethodPost, "/api/expenses", map[string]string{
		"category_id": fmt.Sprint(createdCategory.ID),
		"amount":      "24.75",
		"description": "Lunch",
		"date":        time.Now().Format("2006-01-02"),
	})
	if expenseResponse.Code != http.StatusCreated {
		t.Fatalf("create expense status = %d, want %d", expenseResponse.Code, http.StatusCreated)
	}

	trendsResponse := apiRequest(t, handler, http.MethodGet, "/api/trends/spending?months=3", nil)
	if trendsResponse.Code != http.StatusOK {
		t.Fatalf("trends status = %d, want %d: %s", trendsResponse.Code, http.StatusOK, trendsResponse.Body)
	}
	trendPoints := decodeResponse[[]struct {
		Month    string  `json:"month"`
		Category string  `json:"category"`
		Amount   float64 `json:"amount"`
	}](t, trendsResponse)
	if len(trendPoints) != 1 || trendPoints[0].Category != "Food" || trendPoints[0].Amount != 24.75 {
		t.Fatalf("trend points = %+v", trendPoints)
	}

	goalResponse := apiRequest(t, handler, http.MethodPost, "/api/goals", map[string]interface{}{
		"name": "Emergency fund", "target_amount": 5000, "current_amount": 250, "target_date": "2027-01-01",
	})
	if goalResponse.Code != http.StatusCreated {
		t.Fatalf("create goal status = %d, want %d: %s", goalResponse.Code, http.StatusCreated, goalResponse.Body)
	}
	createdGoal := decodeResponse[goal.Goal](t, goalResponse)
	updateResponse := apiRequest(t, handler, http.MethodPut, fmt.Sprintf("/api/goals/%d", createdGoal.ID), map[string]interface{}{
		"name": "Emergency fund", "target_amount": 5000, "current_amount": 500, "target_date": "",
	})
	if updateResponse.Code != http.StatusOK {
		t.Fatalf("update goal status = %d, want %d: %s", updateResponse.Code, http.StatusOK, updateResponse.Body)
	}
	listResponse := apiRequest(t, handler, http.MethodGet, "/api/goals", nil)
	goals := decodeResponse[[]goal.Goal](t, listResponse)
	if len(goals) != 1 || goals[0].CurrentAmount != 500 || goals[0].TargetDate != "" {
		t.Fatalf("listed goals = %+v", goals)
	}
	deleteResponse := apiRequest(t, handler, http.MethodDelete, fmt.Sprintf("/api/goals/%d", createdGoal.ID), nil)
	if deleteResponse.Code != http.StatusNoContent {
		t.Fatalf("delete goal status = %d, want %d", deleteResponse.Code, http.StatusNoContent)
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

func TestGoalContributionsAndBudgetEnvelopesAPI(t *testing.T) {
	_, handler := setupIntegrationDB(t)
	targetDate := time.Now().UTC().AddDate(1, 0, 0).Format("2006-01-02")
	goalResponse := apiRequest(t, handler, http.MethodPost, "/api/planning/goals", map[string]interface{}{
		"name": "Emergency fund", "currency": "USD", "target_amount_minor": 10000, "target_date": targetDate,
	})
	if goalResponse.Code != http.StatusCreated {
		t.Fatalf("create financial goal status = %d, want %d; body = %s", goalResponse.Code, http.StatusCreated, goalResponse.Body.String())
	}
	createdGoal := decodeResponse[planning.Goal](t, goalResponse)

	contributionResponse := apiRequest(t, handler, http.MethodPost,
		fmt.Sprintf("/api/planning/goals/%d/contributions", createdGoal.ID),
		map[string]interface{}{"amount_minor": 2500, "direction": "contribution", "date": time.Now().UTC().Format("2006-01-02"), "note": "Monthly saving"},
	)
	if contributionResponse.Code != http.StatusCreated {
		t.Fatalf("create contribution status = %d, want %d; body = %s", contributionResponse.Code, http.StatusCreated, contributionResponse.Body.String())
	}
	overdrawResponse := apiRequest(t, handler, http.MethodPost,
		fmt.Sprintf("/api/planning/goals/%d/contributions", createdGoal.ID),
		map[string]interface{}{"amount_minor": 3000, "direction": "withdrawal", "date": time.Now().UTC().Format("2006-01-02")},
	)
	if overdrawResponse.Code != http.StatusConflict {
		t.Fatalf("overdraw status = %d, want %d", overdrawResponse.Code, http.StatusConflict)
	}
	goalsResponse := apiRequest(t, handler, http.MethodGet, "/api/planning/goals", nil)
	goals := decodeResponse[[]planning.Goal](t, goalsResponse)
	if len(goals) != 1 || goals[0].CurrentAmountMinor != 2500 {
		t.Fatalf("goal progress = %+v", goals)
	}

	categoryResponse := apiRequest(t, handler, http.MethodPost, "/api/categories", map[string]string{"name": "Food"})
	createdCategory := decodeResponse[category.Category](t, categoryResponse)
	accountResponse := apiRequest(t, handler, http.MethodPost, "/api/accounts", map[string]interface{}{
		"name": "Checking", "type": "checking", "currency": "USD", "opening_balance_minor": 10000,
	})
	if accountResponse.Code != http.StatusCreated {
		t.Fatalf("create account status = %d, want %d", accountResponse.Code, http.StatusCreated)
	}
	today := time.Now().UTC().Format("2006-01-02")
	transactionResponse := apiRequest(t, handler, http.MethodPost, "/api/transactions", map[string]interface{}{
		"account_id": 1, "type": "expense", "amount_minor": 300, "currency": "USD", "date": today,
		"category_id": createdCategory.ID,
	})
	if transactionResponse.Code != http.StatusCreated {
		t.Fatalf("create transaction status = %d, want %d; body = %s", transactionResponse.Code, http.StatusCreated, transactionResponse.Body.String())
	}
	envelopeResponse := apiRequest(t, handler, http.MethodPost, "/api/planning/envelopes", map[string]interface{}{
		"category_id": createdCategory.ID, "currency": "USD", "monthly_limit_minor": 1000, "rollover": false,
	})
	if envelopeResponse.Code != http.StatusCreated {
		t.Fatalf("create envelope status = %d, want %d; body = %s", envelopeResponse.Code, http.StatusCreated, envelopeResponse.Body.String())
	}
	summaryResponse := apiRequest(t, handler, http.MethodGet, "/api/planning/envelopes?period="+time.Now().UTC().Format("2006-01"), nil)
	summaries := decodeResponse[[]planning.EnvelopeSummary](t, summaryResponse)
	if len(summaries) != 1 || summaries[0].SpentMinor != 300 || summaries[0].AvailableMinor != 700 {
		t.Fatalf("envelope summary = %+v", summaries)
	}
}
