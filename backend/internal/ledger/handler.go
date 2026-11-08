package ledger

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetAccountsHandler(w http.ResponseWriter, _ *http.Request) {
	accounts, err := GetAllAccounts()
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, accounts)
}

func CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	var account Account
	if err := decodeRequest(r, &account); err != nil {
		http.Error(w, "Invalid account payload", http.StatusBadRequest)
		return
	}
	if err := CreateAccount(&account); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, account)
}

func GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	var accountID uint
	if raw := query.Get("account_id"); raw != "" {
		value, err := strconv.ParseUint(raw, 10, 64)
		if err != nil || value == 0 {
			http.Error(w, "Invalid account_id", http.StatusBadRequest)
			return
		}
		accountID = uint(value)
	}
	limit, err := queryInt(query.Get("limit"), 50)
	if err != nil {
		http.Error(w, "Invalid limit", http.StatusBadRequest)
		return
	}
	offset, err := queryInt(query.Get("offset"), 0)
	if err != nil || offset < 0 {
		http.Error(w, "Invalid offset", http.StatusBadRequest)
		return
	}

	transactions, err := GetTransactions(accountID, query.Get("type"), query.Get("from"), query.Get("to"), limit, offset)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, transactions)
}

func GetTransactionHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := transactionID(w, r)
	if !ok {
		return
	}
	transaction, err := GetTransactionByID(id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, transaction)
}

func CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var transaction Transaction
	if err := decodeRequest(r, &transaction); err != nil {
		http.Error(w, "Invalid transaction payload", http.StatusBadRequest)
		return
	}
	if err := CreateTransaction(&transaction); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, transaction)
}

func UpdateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := transactionID(w, r)
	if !ok {
		return
	}
	var transaction Transaction
	if err := decodeRequest(r, &transaction); err != nil {
		http.Error(w, "Invalid transaction payload", http.StatusBadRequest)
		return
	}
	if err := UpdateTransaction(id, &transaction); err != nil {
		writeError(w, err)
		return
	}
	transaction.ID = uint(id)
	writeJSON(w, http.StatusOK, transaction)
}

func DeleteTransactionHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := transactionID(w, r)
	if !ok {
		return
	}
	if err := DeleteTransaction(id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func GetCategoryRulesHandler(w http.ResponseWriter, _ *http.Request) {
	rules, err := GetCategoryRules()
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, rules)
}

func CreateCategoryRuleHandler(w http.ResponseWriter, r *http.Request) {
	var rule CategoryRule
	if err := decodeRequest(r, &rule); err != nil {
		http.Error(w, "Invalid category rule payload", http.StatusBadRequest)
		return
	}
	if err := CreateCategoryRule(&rule); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, rule)
}

func UpdateCategoryRuleHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := transactionID(w, r)
	if !ok {
		return
	}
	var rule CategoryRule
	if err := decodeRequest(r, &rule); err != nil {
		http.Error(w, "Invalid category rule payload", http.StatusBadRequest)
		return
	}
	if err := UpdateCategoryRule(id, &rule); err != nil {
		writeError(w, err)
		return
	}
	rule.ID = uint(id)
	writeJSON(w, http.StatusOK, rule)
}

func DeleteCategoryRuleHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := transactionID(w, r)
	if !ok {
		return
	}
	if err := DeleteCategoryRule(id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func GetScheduledFlowsHandler(w http.ResponseWriter, _ *http.Request) {
	flows, err := GetScheduledFlows()
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, flows)
}

func CreateScheduledFlowHandler(w http.ResponseWriter, r *http.Request) {
	var flow ScheduledFlow
	if err := decodeRequest(r, &flow); err != nil {
		http.Error(w, "Invalid scheduled flow payload", http.StatusBadRequest)
		return
	}
	if err := CreateScheduledFlow(&flow); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, flow)
}

func UpdateScheduledFlowHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := transactionID(w, r)
	if !ok {
		return
	}
	var flow ScheduledFlow
	if err := decodeRequest(r, &flow); err != nil {
		http.Error(w, "Invalid scheduled flow payload", http.StatusBadRequest)
		return
	}
	if err := UpdateScheduledFlow(id, &flow); err != nil {
		writeError(w, err)
		return
	}
	flow.ID = uint(id)
	writeJSON(w, http.StatusOK, flow)
}

func DeleteScheduledFlowHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := transactionID(w, r)
	if !ok {
		return
	}
	if err := DeleteScheduledFlow(id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func GetBillsHandler(w http.ResponseWriter, _ *http.Request) {
	bills, err := GetBillInsights(time.Now())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, bills)
}

func ProcessScheduledOccurrenceHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := transactionID(w, r)
	if !ok {
		return
	}
	var request ProcessOccurrenceRequest
	if err := decodeRequest(r, &request); err != nil {
		http.Error(w, "Invalid scheduled occurrence payload", http.StatusBadRequest)
		return
	}
	occurrence, err := ProcessScheduledOccurrence(uint(id), request)
	switch {
	case err == nil:
		writeJSON(w, http.StatusOK, occurrence)
	case errors.Is(err, ErrInvalidOccurrence), errors.Is(err, ErrInvalidTxn):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrOccurrenceExists):
		http.Error(w, err.Error(), http.StatusConflict)
	case errors.Is(err, gorm.ErrRecordNotFound):
		http.Error(w, "Scheduled item not found", http.StatusNotFound)
	default:
		writeError(w, err)
	}
}

func CashflowForecastHandler(w http.ResponseWriter, r *http.Request) {
	days, err := queryInt(r.URL.Query().Get("days"), 90)
	if err != nil || days < 1 || days > 365 {
		http.Error(w, "days must be between 1 and 365", http.StatusBadRequest)
		return
	}
	forecast, err := BuildCashflowForecast(time.Now(), days)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, forecast)
}

func CashflowReportHandler(w http.ResponseWriter, r *http.Request) {
	from, to := reportDateRange(r)
	report, err := GetCashflowReport(from, to)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func CashflowCSVHandler(w http.ResponseWriter, r *http.Request) {
	from, to := reportDateRange(r)
	report, err := GetCashflowReport(from, to)
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=\"cashflow-report.csv\"")
	writer := csv.NewWriter(w)
	if err := writer.Write([]string{"from", "to", "currency", "type", "category", "amount_minor"}); err != nil {
		writeError(w, err)
		return
	}
	for _, currency := range report.Currencies {
		if err := writer.Write([]string{from, to, currency.Currency, "income", "All income", strconv.FormatInt(currency.IncomeMinor, 10)}); err != nil {
			writeError(w, err)
			return
		}
		for _, category := range currency.Categories {
			categoryName := category.CategoryName
			if categoryName == "" {
				categoryName = "Uncategorized"
			}
			if err := writer.Write([]string{
				from, to, currency.Currency, "expense", categoryName, strconv.FormatInt(category.AmountMinor, 10),
			}); err != nil {
				writeError(w, err)
				return
			}
		}
		if err := writer.Write([]string{from, to, currency.Currency, "net", "", strconv.FormatInt(currency.NetMinor, 10)}); err != nil {
			writeError(w, err)
			return
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		log.Printf("write cashflow CSV: %v", err)
	}
}

func NetWorthHandler(w http.ResponseWriter, _ *http.Request) {
	report, err := GetNetWorthReport(time.Now())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func NetWorthHistoryHandler(w http.ResponseWriter, r *http.Request) {
	months, err := queryInt(r.URL.Query().Get("months"), 12)
	if err != nil || months < 1 || months > 36 {
		http.Error(w, "months must be between 1 and 36", http.StatusBadRequest)
		return
	}
	report, err := GetNetWorthHistory(months, time.Now())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func ExportTransactionsCSVHandler(w http.ResponseWriter, _ *http.Request) {
	transactions := make([]Transaction, 0)
	if err := db.Order("date ASC, id ASC").Find(&transactions).Error; err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=\"transactions.csv\"")
	writer := csv.NewWriter(w)
	if err := writer.Write([]string{
		"date", "account_id", "type", "amount", "currency", "payee", "description", "category_id", "transfer_account_id",
	}); err != nil {
		log.Printf("write transaction CSV header: %v", err)
		return
	}
	for _, transaction := range transactions {
		categoryID := ""
		if transaction.CategoryID != nil {
			categoryID = strconv.FormatUint(uint64(*transaction.CategoryID), 10)
		}
		transferAccountID := ""
		if transaction.TransferAccountID != nil {
			transferAccountID = strconv.FormatUint(uint64(*transaction.TransferAccountID), 10)
		}
		record := []string{
			transaction.Date,
			strconv.FormatUint(uint64(transaction.AccountID), 10),
			transaction.Type,
			formatAmountForCSV(transaction.AmountMinor, transaction.Currency),
			transaction.Currency,
			transaction.Payee,
			transaction.Description,
			categoryID,
			transferAccountID,
		}
		if err := writer.Write(record); err != nil {
			log.Printf("write transaction CSV row: %v", err)
			return
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		log.Printf("flush transaction CSV: %v", err)
	}
}

func ImportTransactionsCSVHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxTransactionCSVSize+(1<<20))
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		http.Error(w, "Upload a CSV file no larger than 5 MB", http.StatusBadRequest)
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "A CSV file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()
	content, err := io.ReadAll(io.LimitReader(file, MaxTransactionCSVSize+1))
	if err != nil {
		http.Error(w, "Unable to read CSV file", http.StatusBadRequest)
		return
	}
	if len(content) > MaxTransactionCSVSize {
		http.Error(w, "CSV file must be no larger than 5 MB", http.StatusBadRequest)
		return
	}
	report, err := ImportTransactionsCSV(content, r.URL.Query().Get("commit") == "true")
	if err != nil {
		if errors.Is(err, ErrInvalidCSV) || errors.Is(err, ErrCSVRowLimit) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeError(w, err)
		return
	}
	if r.URL.Query().Get("commit") == "true" && report.Invalid > 0 {
		writeJSON(w, http.StatusBadRequest, report)
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func reportDateRange(r *http.Request) (string, string) {
	now := time.Now()
	from := strings.TrimSpace(r.URL.Query().Get("from"))
	to := strings.TrimSpace(r.URL.Query().Get("to"))
	if from == "" {
		from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
	}
	if to == "" {
		to = now.Format("2006-01-02")
	}
	return from, to
}

func transactionID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil || id < 1 {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

func decodeRequest(r *http.Request, target interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}

func queryInt(raw string, fallback int) (int, error) {
	if raw == "" {
		return fallback, nil
	}
	return strconv.Atoi(raw)
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidAccount), errors.Is(err, ErrInvalidTxn),
		errors.Is(err, ErrInvalidTransfer), errors.Is(err, ErrInvalidTxnTarget),
		errors.Is(err, ErrInvalidRule), errors.Is(err, ErrInvalidSchedule), errors.Is(err, ErrInvalidCSV):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, gorm.ErrRecordNotFound):
		http.Error(w, "Account, category, or transaction not found", http.StatusNotFound)
	default:
		log.Printf("ledger request failed: %v", err)
		http.Error(w, "Unable to complete ledger request", http.StatusInternalServerError)
	}
}

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("encode ledger response: %v", err)
	}
}
