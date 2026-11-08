package ledger

import (
	"bytes"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/girishsaraf/xpense-trackr/backend/internal/category"
	"gorm.io/gorm"
)

const (
	MaxTransactionCSVSize = 5 << 20
	MaxTransactionCSVRows = 5000
)

var (
	ErrInvalidCSV  = errors.New("CSV must contain account_id, type, amount, currency, and date columns")
	ErrCSVRowLimit = errors.New("CSV exceeds the 5000 transaction row limit")
)

type CSVImportRow struct {
	Line        int          `json:"line"`
	Status      string       `json:"status"`
	Message     string       `json:"message,omitempty"`
	Transaction *Transaction `json:"transaction,omitempty"`
}

type CSVImportReport struct {
	Total      int            `json:"total"`
	Ready      int            `json:"ready"`
	Duplicates int            `json:"duplicates"`
	Invalid    int            `json:"invalid"`
	Imported   int            `json:"imported"`
	Rows       []CSVImportRow `json:"rows"`
}

func ParseAmountMinor(value, currency string) (int64, error) {
	value = strings.TrimSpace(value)
	currency = strings.ToUpper(strings.TrimSpace(currency))
	digits, ok := currencyMinorUnitDigits[currency]
	if value == "" || !ok || strings.HasPrefix(value, "-") || strings.HasPrefix(value, "+") {
		return 0, ErrInvalidTxn
	}
	parts := strings.Split(value, ".")
	if len(parts) > 2 || parts[0] == "" {
		return 0, ErrInvalidTxn
	}
	for _, character := range parts[0] {
		if character < '0' || character > '9' {
			return 0, ErrInvalidTxn
		}
	}
	fraction := ""
	if len(parts) == 2 {
		fraction = parts[1]
		if fraction == "" || len(fraction) > digits {
			return 0, ErrInvalidTxn
		}
		for _, character := range fraction {
			if character < '0' || character > '9' {
				return 0, ErrInvalidTxn
			}
		}
	}
	whole, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, ErrInvalidTxn
	}
	unit := int64(1)
	for range digits {
		unit *= 10
	}
	fraction += strings.Repeat("0", digits-len(fraction))
	var fractional int64
	if fraction != "" {
		fractional, err = strconv.ParseInt(fraction, 10, 64)
		if err != nil {
			return 0, ErrInvalidTxn
		}
	}
	if whole > (math.MaxInt64-fractional)/unit {
		return 0, ErrInvalidTxn
	}
	amount := whole*unit + fractional
	if amount <= 0 {
		return 0, ErrInvalidTxn
	}
	return amount, nil
}

func ImportTransactionsCSV(data []byte, commit bool) (*CSVImportReport, error) {
	if len(data) == 0 || len(data) > MaxTransactionCSVSize {
		return nil, ErrInvalidCSV
	}
	reader := csv.NewReader(bytes.NewReader(data))
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1
	header, err := reader.Read()
	if err != nil {
		return nil, ErrInvalidCSV
	}
	columns := make(map[string]int, len(header))
	for index, value := range header {
		name := strings.ToLower(strings.TrimSpace(value))
		if name == "" {
			return nil, ErrInvalidCSV
		}
		if _, exists := columns[name]; exists {
			return nil, ErrInvalidCSV
		}
		columns[name] = index
	}
	for _, required := range []string{"account_id", "type", "amount", "currency", "date"} {
		if _, exists := columns[required]; !exists {
			return nil, ErrInvalidCSV
		}
	}

	accounts, err := GetAllAccounts()
	if err != nil {
		return nil, err
	}
	accountByID := make(map[uint]Account, len(accounts))
	for _, account := range accounts {
		accountByID[account.ID] = account
	}
	report := &CSVImportReport{Rows: make([]CSVImportRow, 0)}
	seen := make(map[string]struct{})
	for {
		record, readErr := reader.Read()
		if errors.Is(readErr, io.EOF) {
			break
		}
		if readErr != nil {
			line := report.Total + 2
			var parseErr *csv.ParseError
			if errors.As(readErr, &parseErr) {
				line = parseErr.Line
			}
			report.Rows = append(report.Rows, CSVImportRow{Line: line, Status: "invalid", Message: "Malformed CSV row"})
			report.Invalid++
			report.Total++
			continue
		}
		if len(report.Rows) >= MaxTransactionCSVRows {
			return nil, ErrCSVRowLimit
		}
		report.Total++
		line, _ := reader.FieldPos(0)
		row := CSVImportRow{Line: line}
		if len(record) != len(header) {
			row.Status = "invalid"
			row.Message = "Row has a different number of columns than the header"
			report.Invalid++
			report.Rows = append(report.Rows, row)
			continue
		}
		transaction, fingerprint, rowErr := parseCSVTransaction(record, columns, accountByID)
		if rowErr != nil {
			row.Status = "invalid"
			row.Message = rowErr.Error()
			report.Invalid++
			report.Rows = append(report.Rows, row)
			continue
		}
		row.Transaction = transaction
		if _, duplicateInFile := seen[fingerprint]; duplicateInFile {
			row.Status = "duplicate"
			row.Message = "Duplicate row in this file"
			report.Duplicates++
			report.Rows = append(report.Rows, row)
			continue
		}
		seen[fingerprint] = struct{}{}
		var count int64
		if err := db.Model(&Transaction{}).Where("import_fingerprint = ?", fingerprint).Count(&count).Error; err != nil {
			return nil, err
		}
		if count > 0 {
			row.Status = "duplicate"
			row.Message = "Already imported"
			report.Duplicates++
			report.Rows = append(report.Rows, row)
			continue
		}
		row.Status = "ready"
		report.Ready++
		transaction.ImportFingerprint = &fingerprint
		report.Rows = append(report.Rows, row)
	}
	if report.Total == 0 {
		return nil, ErrInvalidCSV
	}
	if !commit || report.Invalid > 0 || report.Ready == 0 {
		return report, nil
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		for index := range report.Rows {
			row := &report.Rows[index]
			if row.Status != "ready" || row.Transaction == nil {
				continue
			}
			transaction := row.Transaction
			if err := normalizeAndValidateTransaction(transaction); err != nil {
				return fmt.Errorf("row %d: %w", row.Line, err)
			}
			if err := validateTransactionAccounts(tx, transaction); err != nil {
				return fmt.Errorf("row %d: %w", row.Line, err)
			}
			if err := applyCategoryRule(tx, transaction); err != nil {
				return fmt.Errorf("row %d: %w", row.Line, err)
			}
			if err := tx.Create(transaction).Error; err != nil {
				return fmt.Errorf("row %d: %w", row.Line, err)
			}
			row.Status = "imported"
			report.Imported++
		}
		return nil
	})
	if err != nil {
		report.Imported = 0
		for index := range report.Rows {
			if report.Rows[index].Status == "imported" {
				report.Rows[index].Status = "ready"
			}
		}
		return nil, err
	}
	report.Ready = 0
	return report, nil
}

func parseCSVTransaction(record []string, columns map[string]int, accounts map[uint]Account) (*Transaction, string, error) {
	value := func(name string) string {
		index, exists := columns[name]
		if !exists {
			return ""
		}
		return strings.TrimSpace(record[index])
	}
	accountValue, err := strconv.ParseUint(value("account_id"), 10, 64)
	if err != nil || accountValue == 0 {
		return nil, "", ErrInvalidAccount
	}
	accountID := uint(accountValue)
	account, exists := accounts[accountID]
	if !exists {
		return nil, "", fmt.Errorf("account_id %d does not exist", accountID)
	}
	currency := strings.ToUpper(value("currency"))
	amountMinor, err := ParseAmountMinor(value("amount"), currency)
	if err != nil {
		return nil, "", fmt.Errorf("amount is invalid for currency %s", currency)
	}
	transaction := &Transaction{
		AccountID:   accountID,
		Type:        value("type"),
		AmountMinor: amountMinor,
		Currency:    currency,
		Date:        value("date"),
		Payee:       value("payee"),
		Description: value("description"),
	}
	if raw := value("category_id"); raw != "" {
		id, err := strconv.ParseUint(raw, 10, 64)
		if err != nil || id == 0 {
			return nil, "", fmt.Errorf("category_id %q is invalid", raw)
		}
		categoryID := uint(id)
		transaction.CategoryID = &categoryID
	}
	if raw := value("transfer_account_id"); raw != "" {
		id, err := strconv.ParseUint(raw, 10, 64)
		if err != nil || id == 0 {
			return nil, "", fmt.Errorf("transfer_account_id %q is invalid", raw)
		}
		transferAccountID := uint(id)
		transaction.TransferAccountID = &transferAccountID
	}
	if err := normalizeAndValidateTransaction(transaction); err != nil {
		return nil, "", err
	}
	if transaction.Currency != account.Currency {
		return nil, "", ErrInvalidTxn
	}
	if err := validateTransactionAccounts(db, transaction); err != nil {
		return nil, "", err
	}
	if transaction.CategoryID != nil {
		var existing category.Category
		if err := db.First(&existing, *transaction.CategoryID).Error; err != nil {
			return nil, "", fmt.Errorf("transaction category: %w", err)
		}
	}
	fingerprint := transactionFingerprint(transaction)
	return transaction, fingerprint, nil
}

func transactionFingerprint(transaction *Transaction) string {
	value := fmt.Sprintf(
		"%d|%v|%s|%d|%s|%s|%s|%s|%v",
		transaction.AccountID, transaction.TransferAccountID, transaction.Type, transaction.AmountMinor,
		transaction.Currency, transaction.Date, strings.ToLower(transaction.Payee),
		strings.ToLower(transaction.Description), transaction.CategoryID,
	)
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func formatAmountForCSV(amount int64, currency string) string {
	digits := currencyMinorUnitDigits[currency]
	unit := int64(1)
	for range digits {
		unit *= 10
	}
	whole := amount / unit
	fraction := amount % unit
	if digits == 0 {
		return strconv.FormatInt(whole, 10)
	}
	return fmt.Sprintf("%d.%0*d", whole, digits, fraction)
}
