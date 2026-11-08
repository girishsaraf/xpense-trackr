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
	Currency         string                 `json:"currency"`
	AssetsMinor      int64                  `json:"assets_minor"`
	LiabilitiesMinor int64                  `json:"liabilities_minor"`
	NetWorthMinor    int64                  `json:"net_worth_minor"`
	PortfolioMinor   int64                  `json:"portfolio_minor"`
	Accounts         []NetWorthAccount      `json:"accounts"`
	ManualItems      []NetWorthItem         `json:"manual_items"`
	Holdings         []PortfolioHoldingItem `json:"holdings"`
}

type NetWorthItem struct {
	AssetID    uint   `json:"asset_id"`
	Name       string `json:"name"`
	Kind       string `json:"kind"`
	ValueMinor int64  `json:"value_minor"`
	AsOf       string `json:"as_of"`
	Source     string `json:"source"`
}

type PortfolioHoldingItem struct {
	HoldingID      uint   `json:"holding_id"`
	Name           string `json:"name"`
	Ticker         string `json:"ticker"`
	AssetClass     string `json:"asset_class"`
	Quantity       string `json:"quantity"`
	CostBasisMinor int64  `json:"cost_basis_minor"`
	ValueMinor     int64  `json:"value_minor"`
	AsOf           string `json:"as_of"`
	Source         string `json:"source"`
}

type NetWorthReport struct {
	AsOf       string             `json:"as_of"`
	Currencies []NetWorthCurrency `json:"currencies"`
}

type NetWorthHistoryEntry struct {
	AsOf       string             `json:"as_of"`
	Currencies []NetWorthCurrency `json:"currencies"`
}

type NetWorthHistoryReport struct {
	Months  int                    `json:"months"`
	Entries []NetWorthHistoryEntry `json:"entries"`
}

type netWorthAccountRow struct {
	ID                  uint   `gorm:"column:id"`
	Name                string `gorm:"column:name"`
	Type                string `gorm:"column:account_type"`
	Currency            string `gorm:"column:currency"`
	OpeningBalanceMinor int64  `gorm:"column:opening_balance_minor"`
	BalanceMinor        int64  `gorm:"column:balance_minor"`
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
	asOf = asOf.UTC()
	accounts, err := getAccountsAtDate(asOf.Format("2006-01-02"))
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
				ManualItems: make([]NetWorthItem, 0), Holdings: make([]PortfolioHoldingItem, 0),
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
	manualAssets, err := GetManualAssetsAt(report.AsOf)
	if err != nil {
		return nil, err
	}
	for _, asset := range manualAssets {
		if !asset.HasValuation {
			continue
		}
		index, found := byCurrency[asset.Currency]
		if !found {
			index = len(report.Currencies)
			byCurrency[asset.Currency] = index
			report.Currencies = append(report.Currencies, NetWorthCurrency{
				Currency: asset.Currency, Accounts: make([]NetWorthAccount, 0),
				ManualItems: make([]NetWorthItem, 0), Holdings: make([]PortfolioHoldingItem, 0),
			})
		}
		currencyReport := &report.Currencies[index]
		if asset.Kind == "asset" {
			currencyReport.AssetsMinor += asset.LatestValueMinor
			currencyReport.NetWorthMinor += asset.LatestValueMinor
		} else {
			currencyReport.LiabilitiesMinor += asset.LatestValueMinor
			currencyReport.NetWorthMinor -= asset.LatestValueMinor
		}
		currencyReport.ManualItems = append(currencyReport.ManualItems, NetWorthItem{
			AssetID: asset.ID, Name: asset.Name, Kind: asset.Kind, ValueMinor: asset.LatestValueMinor,
			AsOf: asset.ValuationDate, Source: asset.ValuationSource,
		})
	}
	holdings, err := GetPortfolioHoldingsAt(report.AsOf)
	if err != nil {
		return nil, err
	}
	for _, holding := range holdings {
		if !holding.HasValuation {
			continue
		}
		index, found := byCurrency[holding.Currency]
		if !found {
			index = len(report.Currencies)
			byCurrency[holding.Currency] = index
			report.Currencies = append(report.Currencies, NetWorthCurrency{
				Currency: holding.Currency, Accounts: make([]NetWorthAccount, 0),
				ManualItems: make([]NetWorthItem, 0), Holdings: make([]PortfolioHoldingItem, 0),
			})
		}
		currencyReport := &report.Currencies[index]
		currencyReport.AssetsMinor += holding.LatestValueMinor
		currencyReport.NetWorthMinor += holding.LatestValueMinor
		currencyReport.PortfolioMinor += holding.LatestValueMinor
		currencyReport.Holdings = append(currencyReport.Holdings, PortfolioHoldingItem{
			HoldingID: holding.ID, Name: holding.Name, Ticker: holding.Ticker,
			AssetClass: holding.AssetClass, Quantity: holding.Quantity,
			CostBasisMinor: holding.CostBasisMinor, ValueMinor: holding.LatestValueMinor,
			AsOf: holding.ValuationDate, Source: holding.ValuationSource,
		})
	}
	sort.Slice(report.Currencies, func(i, j int) bool {
		return report.Currencies[i].Currency < report.Currencies[j].Currency
	})
	return report, nil
}

func GetNetWorthHistory(months int, now time.Time) (*NetWorthHistoryReport, error) {
	if months < 1 || months > 36 {
		return nil, ErrInvalidTxn
	}
	now = now.UTC()
	currentDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	currentMonth := time.Date(currentDate.Year(), currentDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	firstMonth := currentMonth.AddDate(0, -(months - 1), 0)
	report := &NetWorthHistoryReport{Months: months, Entries: make([]NetWorthHistoryEntry, 0, months)}
	for index := 0; index < months-1; index++ {
		month := firstMonth.AddDate(0, index, 0)
		monthEnd := month.AddDate(0, 1, -1)
		snapshot, err := GetNetWorthReport(monthEnd)
		if err != nil {
			return nil, err
		}
		report.Entries = append(report.Entries, NetWorthHistoryEntry{
			AsOf: snapshot.AsOf, Currencies: snapshot.Currencies,
		})
	}
	current, err := GetNetWorthReport(currentDate)
	if err != nil {
		return nil, err
	}
	report.Entries = append(report.Entries, NetWorthHistoryEntry{
		AsOf: current.AsOf, Currencies: current.Currencies,
	})
	return report, nil
}

func getAccountsAtDate(asOf string) ([]Account, error) {
	asOfDate, err := time.Parse("2006-01-02", asOf)
	if err != nil {
		return nil, ErrInvalidTxn
	}
	rows := make([]netWorthAccountRow, 0)
	err = db.Table("accounts").
		Select(`accounts.id, accounts.name, accounts.account_type, accounts.currency,
			accounts.opening_balance_minor,
			accounts.opening_balance_minor + COALESCE(SUM(CASE
				WHEN ledger_transactions.transaction_type = 'income' AND ledger_transactions.account_id = accounts.id THEN ledger_transactions.amount_minor
				WHEN ledger_transactions.transaction_type = 'expense' AND ledger_transactions.account_id = accounts.id THEN -ledger_transactions.amount_minor
				WHEN ledger_transactions.transaction_type = 'transfer' AND ledger_transactions.account_id = accounts.id THEN -ledger_transactions.amount_minor
				WHEN ledger_transactions.transaction_type = 'transfer' AND ledger_transactions.transfer_account_id = accounts.id THEN ledger_transactions.amount_minor
				ELSE 0 END), 0) AS balance_minor`).
		Joins(`LEFT JOIN ledger_transactions
			ON (ledger_transactions.account_id = accounts.id OR ledger_transactions.transfer_account_id = accounts.id)
			AND ledger_transactions.date <= ?`, asOf).
		Where("accounts.created_at < ?", asOfDate.AddDate(0, 0, 1)).
		Group("accounts.id").
		Order("accounts.name ASC, accounts.id ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	accounts := make([]Account, 0, len(rows))
	for _, row := range rows {
		accounts = append(accounts, Account{
			ID: row.ID, Name: row.Name, Type: row.Type, Currency: row.Currency,
			OpeningBalanceMinor: row.OpeningBalanceMinor, BalanceMinor: row.BalanceMinor,
		})
	}
	return accounts, nil
}
