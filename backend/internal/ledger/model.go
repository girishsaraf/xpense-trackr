package ledger

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/girishsaraf/xpense-trackr/backend/internal/category"
	"gorm.io/gorm"
)

var (
	db                  *gorm.DB
	ErrRecordNotFound   = gorm.ErrRecordNotFound
	ErrInvalidAccount   = errors.New("account name, supported type, and supported currency are required")
	ErrInvalidTxn       = errors.New("transaction type, account, positive amount, currency, and YYYY-MM-DD date are required")
	ErrInvalidTransfer  = errors.New("transfers require a different account with the same currency and no category")
	ErrInvalidTxnTarget = errors.New("income and expense transactions cannot have a transfer account")
	ErrInvalidRule      = errors.New("rule name, payee match, and category are required")
	currencyCodePattern = regexp.MustCompile(`^[A-Z]{3}$`)
)

var accountTypes = map[string]struct{}{
	"cash": {}, "checking": {}, "savings": {}, "investment": {}, "credit_card": {}, "loan": {},
}

var currencyMinorUnitDigits = map[string]int{
	"AED": 2, "AUD": 2, "BHD": 3, "BRL": 2, "CAD": 2, "CHF": 2, "CLP": 0,
	"CNY": 2, "DKK": 2, "EUR": 2, "GBP": 2, "HKD": 2, "IDR": 2, "INR": 2,
	"IQD": 3, "ISK": 0, "JPY": 0, "JOD": 3, "KRW": 0, "KWD": 3, "MXN": 2,
	"MYR": 2, "NOK": 2, "NZD": 2, "OMR": 3, "PHP": 2, "PLN": 2, "SAR": 2,
	"SEK": 2, "SGD": 2, "THB": 2, "TND": 3, "TRY": 2, "USD": 2, "VND": 0,
	"ZAR": 2,
}

// SetDB assigns the database used by ledger operations.
func SetDB(database *gorm.DB) {
	db = database
}

type Account struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	Name                string    `gorm:"size:120;not null" json:"name"`
	Type                string    `gorm:"column:account_type;size:32;not null" json:"type"`
	Currency            string    `gorm:"size:3;not null" json:"currency"`
	OpeningBalanceMinor int64     `gorm:"not null;default:0" json:"opening_balance_minor"`
	Archived            bool      `gorm:"not null;default:false" json:"archived"`
	BalanceMinor        int64     `gorm:"->;column:balance_minor" json:"balance_minor"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type Transaction struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	AccountID         uint      `gorm:"not null;index:idx_ledger_tx_account_date,priority:1" json:"account_id"`
	TransferAccountID *uint     `gorm:"index" json:"transfer_account_id,omitempty"`
	CategoryID        *uint     `json:"category_id,omitempty"`
	CategoryRuleID    *uint     `gorm:"index" json:"category_rule_id,omitempty"`
	ImportFingerprint *string   `gorm:"size:64;uniqueIndex" json:"-"`
	Type              string    `gorm:"column:transaction_type;size:16;not null" json:"type"`
	AmountMinor       int64     `gorm:"not null" json:"amount_minor"`
	Currency          string    `gorm:"size:3;not null" json:"currency"`
	Payee             string    `gorm:"size:180" json:"payee,omitempty"`
	Description       string    `gorm:"size:500" json:"description,omitempty"`
	Date              string    `gorm:"size:10;not null;index:idx_ledger_tx_account_date,priority:2" json:"date"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CategoryRule struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"size:120;not null" json:"name"`
	MatchValue string    `gorm:"size:180;not null" json:"match_value"`
	CategoryID uint      `gorm:"not null;index" json:"category_id"`
	Priority   int       `gorm:"not null;default:0" json:"priority"`
	Enabled    bool      `gorm:"not null;default:true" json:"enabled"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Transaction) TableName() string {
	return "ledger_transactions"
}

func CreateAccount(account *Account) error {
	account.Name = strings.TrimSpace(account.Name)
	account.Type = strings.ToLower(strings.TrimSpace(account.Type))
	account.Currency = strings.ToUpper(strings.TrimSpace(account.Currency))
	if account.Name == "" || len(account.Name) > 120 ||
		!validAccountType(account.Type) || !validCurrency(account.Currency) {
		return ErrInvalidAccount
	}
	return db.Create(account).Error
}

func GetAllAccounts() ([]Account, error) {
	accounts := make([]Account, 0)
	err := db.Model(&Account{}).
		Select(`accounts.id, accounts.name, accounts.account_type, accounts.currency,
			accounts.opening_balance_minor, accounts.archived, accounts.created_at, accounts.updated_at,
			accounts.opening_balance_minor + COALESCE(SUM(CASE
				WHEN ledger_transactions.transaction_type = 'income' AND ledger_transactions.account_id = accounts.id THEN ledger_transactions.amount_minor
				WHEN ledger_transactions.transaction_type = 'expense' AND ledger_transactions.account_id = accounts.id THEN -ledger_transactions.amount_minor
				WHEN ledger_transactions.transaction_type = 'transfer' AND ledger_transactions.account_id = accounts.id THEN -ledger_transactions.amount_minor
				WHEN ledger_transactions.transaction_type = 'transfer' AND ledger_transactions.transfer_account_id = accounts.id THEN ledger_transactions.amount_minor
				ELSE 0 END), 0) AS balance_minor`).
		Joins("LEFT JOIN ledger_transactions ON ledger_transactions.account_id = accounts.id OR ledger_transactions.transfer_account_id = accounts.id").
		Group("accounts.id").
		Order("accounts.archived ASC, accounts.name ASC").
		Scan(&accounts).Error
	return accounts, err
}

func CreateTransaction(transaction *Transaction) error {
	if err := normalizeAndValidateTransaction(transaction); err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := validateTransactionAccounts(tx, transaction); err != nil {
			return err
		}
		if err := applyCategoryRule(tx, transaction); err != nil {
			return err
		}
		return tx.Create(transaction).Error
	})
}

func GetTransactions(accountID uint, transactionType, fromDate, toDate string, limit, offset int) ([]Transaction, error) {
	query := db.Order("date DESC, id DESC")
	if accountID > 0 {
		query = query.Where("account_id = ? OR transfer_account_id = ?", accountID, accountID)
	}
	if transactionType != "" {
		if transactionType != "income" && transactionType != "expense" && transactionType != "transfer" {
			return nil, ErrInvalidTxn
		}
		query = query.Where("transaction_type = ?", transactionType)
	}
	if fromDate != "" {
		if !validDate(fromDate) {
			return nil, ErrInvalidTxn
		}
		query = query.Where("date >= ?", fromDate)
	}
	if toDate != "" {
		if !validDate(toDate) {
			return nil, ErrInvalidTxn
		}
		query = query.Where("date <= ?", toDate)
	}
	if limit < 1 || limit > 100 || offset < 0 || offset > 10000 {
		return nil, ErrInvalidTxn
	}

	transactions := make([]Transaction, 0)
	if err := query.Limit(limit).Offset(offset).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func GetTransactionByID(id int) (*Transaction, error) {
	var transaction Transaction
	if err := db.First(&transaction, id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func UpdateTransaction(id int, updated *Transaction) error {
	if err := normalizeAndValidateTransaction(updated); err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		var existing Transaction
		if err := tx.First(&existing, id).Error; err != nil {
			return err
		}
		if err := validateTransactionAccounts(tx, updated); err != nil {
			return err
		}
		if err := applyCategoryRule(tx, updated); err != nil {
			return err
		}
		return tx.Model(&existing).Select(
			"AccountID", "TransferAccountID", "CategoryID", "CategoryRuleID", "Type", "AmountMinor", "Currency",
			"Payee", "Description", "Date",
		).Updates(updated).Error
	})
}

func DeleteTransaction(id int) error {
	result := db.Delete(&Transaction{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func GetCategoryRules() ([]CategoryRule, error) {
	rules := make([]CategoryRule, 0)
	if err := db.Order("priority DESC, id ASC").Find(&rules).Error; err != nil {
		return nil, err
	}
	return rules, nil
}

func CreateCategoryRule(rule *CategoryRule) error {
	rule.Name = strings.TrimSpace(rule.Name)
	rule.MatchValue = strings.TrimSpace(rule.MatchValue)
	if rule.Name == "" || len(rule.Name) > 120 || rule.MatchValue == "" || len(rule.MatchValue) > 180 || rule.CategoryID == 0 {
		return ErrInvalidRule
	}
	var existing category.Category
	if err := db.First(&existing, rule.CategoryID).Error; err != nil {
		return fmt.Errorf("rule category: %w", err)
	}
	return db.Create(rule).Error
}

func UpdateCategoryRule(id int, updated *CategoryRule) error {
	updated.Name = strings.TrimSpace(updated.Name)
	updated.MatchValue = strings.TrimSpace(updated.MatchValue)
	if updated.Name == "" || len(updated.Name) > 120 || updated.MatchValue == "" || len(updated.MatchValue) > 180 || updated.CategoryID == 0 {
		return ErrInvalidRule
	}
	return db.Transaction(func(tx *gorm.DB) error {
		var existing CategoryRule
		if err := tx.First(&existing, id).Error; err != nil {
			return err
		}
		var existingCategory category.Category
		if err := tx.First(&existingCategory, updated.CategoryID).Error; err != nil {
			return fmt.Errorf("rule category: %w", err)
		}
		return tx.Model(&existing).Select("Name", "MatchValue", "CategoryID", "Priority", "Enabled").Updates(updated).Error
	})
}

func DeleteCategoryRule(id int) error {
	result := db.Delete(&CategoryRule{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func normalizeAndValidateTransaction(transaction *Transaction) error {
	transaction.Type = strings.ToLower(strings.TrimSpace(transaction.Type))
	transaction.Currency = strings.ToUpper(strings.TrimSpace(transaction.Currency))
	transaction.Payee = strings.TrimSpace(transaction.Payee)
	transaction.Description = strings.TrimSpace(transaction.Description)
	transaction.CategoryRuleID = nil
	if transaction.AccountID == 0 || transaction.AmountMinor <= 0 ||
		!validCurrency(transaction.Currency) || !validDate(transaction.Date) ||
		(transaction.Type != "income" && transaction.Type != "expense" && transaction.Type != "transfer") {
		return ErrInvalidTxn
	}
	if transaction.Type == "transfer" {
		if transaction.TransferAccountID == nil || *transaction.TransferAccountID == transaction.AccountID || transaction.CategoryID != nil {
			return ErrInvalidTransfer
		}
	} else if transaction.TransferAccountID != nil {
		return ErrInvalidTxnTarget
	}
	return nil
}

func applyCategoryRule(tx *gorm.DB, transaction *Transaction) error {
	if transaction.CategoryID != nil || transaction.Type == "transfer" {
		if transaction.CategoryID == nil {
			return nil
		}
		var existing category.Category
		if err := tx.First(&existing, *transaction.CategoryID).Error; err != nil {
			return fmt.Errorf("transaction category: %w", err)
		}
		return nil
	}
	if transaction.Payee == "" {
		return nil
	}

	rules, err := GetCategoryRulesFrom(tx)
	if err != nil {
		return err
	}
	for _, rule := range rules {
		if strings.Contains(strings.ToLower(transaction.Payee), strings.ToLower(rule.MatchValue)) {
			transaction.CategoryID = &rule.CategoryID
			transaction.CategoryRuleID = &rule.ID
			return nil
		}
	}
	return nil
}

func GetCategoryRulesFrom(tx *gorm.DB) ([]CategoryRule, error) {
	rules := make([]CategoryRule, 0)
	if err := tx.Where("enabled = ?", true).Order("priority DESC, id ASC").Find(&rules).Error; err != nil {
		return nil, err
	}
	return rules, nil
}

func validateTransactionAccounts(tx *gorm.DB, transaction *Transaction) error {
	var source Account
	if err := tx.First(&source, transaction.AccountID).Error; err != nil {
		return fmt.Errorf("transaction account: %w", err)
	}
	if source.Currency != transaction.Currency {
		return fmt.Errorf("%w: transaction currency must match account currency", ErrInvalidTxn)
	}
	if transaction.TransferAccountID == nil {
		return nil
	}
	var destination Account
	if err := tx.First(&destination, *transaction.TransferAccountID).Error; err != nil {
		return fmt.Errorf("transfer account: %w", err)
	}
	if source.Currency != destination.Currency {
		return ErrInvalidTransfer
	}
	return nil
}

func validAccountType(accountType string) bool {
	_, ok := accountTypes[accountType]
	return ok
}

func validCurrency(currency string) bool {
	if !currencyCodePattern.MatchString(currency) {
		return false
	}
	_, ok := currencyMinorUnitDigits[currency]
	return ok
}

func SupportedCurrency(currency string) bool {
	return validCurrency(strings.ToUpper(strings.TrimSpace(currency)))
}

func validDate(value string) bool {
	parsed, err := time.Parse("2006-01-02", value)
	return err == nil && parsed.Format("2006-01-02") == value
}
