package ledger

import (
	"sort"
	"time"
)

type CashflowCategoryTotal struct {
	CategoryID   *uint  `json:"category_id,omitempty"`
	CategoryName string `json:"category_name"`
	AmountMinor  int64  `json:"amount_minor"`
}

type CashflowCurrencyReport struct {
	Currency     string                  `json:"currency"`
	IncomeMinor  int64                   `json:"income_minor"`
	ExpenseMinor int64                   `json:"expense_minor"`
	NetMinor     int64                   `json:"net_minor"`
	Categories   []CashflowCategoryTotal `json:"categories"`
}

type CashflowReport struct {
	From       string                   `json:"from"`
	To         string                   `json:"to"`
	Currencies []CashflowCurrencyReport `json:"currencies"`
}

type NetWorthAccount struct {
	AccountID    uint   `json:"account_id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	BalanceMinor int64  `json:"balance_minor"`
}

type NetWorthCurrency struct {
	Currency         string            `json:"currency"`
	AssetsMinor      int64             `json:"assets_minor"`
	LiabilitiesMinor int64             `json:"liabilities_minor"`
	NetWorthMinor    int64             `json:"net_worth_minor"`
	Accounts         []NetWorthAccount `json:"accounts"`
}

type NetWorthReport struct {
	AsOf       string             `json:"as_of"`
	Currencies []NetWorthCurrency `json:"currencies"`
}

type cashflowAggregateRow struct {
	Currency     string `gorm:"column:currency"`
	Type         string `gorm:"column:transaction_type"`
	CategoryID   *uint  `gorm:"column:category_id"`
	CategoryName string `gorm:"column:category_name"`
	AmountMinor  int64  `gorm:"column:amount_minor"`
}

func GetCashflowReport(from, to string) (*CashflowReport, error) {
	fromDate, err := time.Parse("2006-01-02", from)
	if err != nil {
		return nil, ErrInvalidTxn
	}
	toDate, err := time.Parse("2006-01-02", to)
	if err != nil || toDate.Before(fromDate) || toDate.After(fromDate.AddDate(1, 0, 0)) {
		return nil, ErrInvalidTxn
	}
	rows := make([]cashflowAggregateRow, 0)
	if err := db.Model(&Transaction{}).
		Select(`ledger_transactions.currency, ledger_transactions.transaction_type,
			ledger_transactions.category_id, COALESCE(categories.name, '') AS category_name,
			SUM(ledger_transactions.amount_minor) AS amount_minor`).
		Joins("LEFT JOIN categories ON categories.id = ledger_transactions.category_id").
		Where("ledger_transactions.date >= ? AND ledger_transactions.date <= ?", from, to).
		Where("ledger_transactions.transaction_type IN ?", []string{"income", "expense"}).
		Group("ledger_transactions.currency, ledger_transactions.transaction_type, ledger_transactions.category_id, categories.name").
		Order("ledger_transactions.currency ASC, ledger_transactions.transaction_type ASC, categories.name ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	report := &CashflowReport{
		From: from, To: to, Currencies: make([]CashflowCurrencyReport, 0),
	}
	byCurrency := make(map[string]int, len(rows))
	for _, row := range rows {
		index, found := byCurrency[row.Currency]
		if !found {
			index = len(report.Currencies)
			byCurrency[row.Currency] = index
			report.Currencies = append(report.Currencies, CashflowCurrencyReport{
				Currency: row.Currency, Categories: make([]CashflowCategoryTotal, 0),
			})
		}
		currencyReport := &report.Currencies[index]
		if row.Type == "income" {
			currencyReport.IncomeMinor += row.AmountMinor
		} else {
			currencyReport.ExpenseMinor += row.AmountMinor
			currencyReport.Categories = append(currencyReport.Categories, CashflowCategoryTotal{
				CategoryID: row.CategoryID, CategoryName: row.CategoryName, AmountMinor: row.AmountMinor,
			})
		}
	}
	for i := range report.Currencies {
		report.Currencies[i].NetMinor = report.Currencies[i].IncomeMinor - report.Currencies[i].ExpenseMinor
	}
	return report, nil
}

func GetNetWorthReport(asOf time.Time) (*NetWorthReport, error) {
	accounts, err := GetAllAccounts()
	if err != nil {
		return nil, err
	}
	report := &NetWorthReport{
		AsOf: asOf.Format("2006-01-02"), Currencies: make([]NetWorthCurrency, 0),
	}
	byCurrency := make(map[string]int, len(accounts))
	for _, account := range accounts {
		index, found := byCurrency[account.Currency]
		if !found {
			index = len(report.Currencies)
			byCurrency[account.Currency] = index
			report.Currencies = append(report.Currencies, NetWorthCurrency{
				Currency: account.Currency, Accounts: make([]NetWorthAccount, 0),
			})
		}
		balance := account.BalanceMinor
		currencyReport := &report.Currencies[index]
		currencyReport.NetWorthMinor += balance
		if balance > 0 {
			currencyReport.AssetsMinor += balance
		} else {
			currencyReport.LiabilitiesMinor += -balance
		}
		currencyReport.Accounts = append(currencyReport.Accounts, NetWorthAccount{
			AccountID: account.ID, Name: account.Name, Type: account.Type, BalanceMinor: balance,
		})
	}
	sort.Slice(report.Currencies, func(i, j int) bool {
		return report.Currencies[i].Currency < report.Currencies[j].Currency
	})
	return report, nil
}
