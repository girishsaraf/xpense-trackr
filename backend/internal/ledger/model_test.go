package ledger

import (
	"errors"
	"testing"
	"time"

	"github.com/girishsaraf/xpense-trackr/backend/internal/category"
	"github.com/girishsaraf/xpense-trackr/backend/internal/testutil"
	"gorm.io/gorm"
)

func setupLedgerDB(t *testing.T) *gorm.DB {
	t.Helper()
	database := testutil.NewDB(t, &Account{}, &Transaction{}, &CategoryRule{}, &ScheduledFlow{}, &ScheduledOccurrence{}, &category.Category{})
	SetDB(database)
	return database
}

func TestLedgerTransactionsUpdateExactAccountBalances(t *testing.T) {
	database := setupLedgerDB(t)
	checking := &Account{Name: "Checking", Type: "checking", Currency: "usd", OpeningBalanceMinor: 10000}
	savings := &Account{Name: "Savings", Type: "savings", Currency: "USD", OpeningBalanceMinor: 100}
	if err := CreateAccount(checking); err != nil {
		t.Fatalf("CreateAccount(checking) error = %v", err)
	}
	if err := CreateAccount(savings); err != nil {
		t.Fatalf("CreateAccount(savings) error = %v", err)
	}
	food := &category.Category{Name: "Food"}
	if err := database.Create(food).Error; err != nil {
		t.Fatalf("create category: %v", err)
	}

	expenseCategoryID := food.ID
	transactions := []Transaction{
		{AccountID: checking.ID, CategoryID: &expenseCategoryID, Type: "expense", AmountMinor: 500, Currency: "usd", Date: "2026-10-12"},
		{AccountID: checking.ID, Type: "income", AmountMinor: 1500, Currency: "USD", Date: "2026-10-12"},
		{AccountID: checking.ID, TransferAccountID: &savings.ID, Type: "transfer", AmountMinor: 200, Currency: "USD", Date: "2026-10-13"},
	}
	for i := range transactions {
		if err := CreateTransaction(&transactions[i]); err != nil {
			t.Fatalf("CreateTransaction(%d) error = %v", i, err)
		}
	}

	accounts, err := GetAllAccounts()
	if err != nil {
		t.Fatalf("GetAllAccounts() error = %v", err)
	}
	if len(accounts) != 2 {
		t.Fatalf("GetAllAccounts() returned %d accounts, want 2", len(accounts))
	}
	balances := map[uint]int64{}
	for _, account := range accounts {
		balances[account.ID] = account.BalanceMinor
	}
	if balances[checking.ID] != 10800 {
		t.Errorf("checking balance = %d, want 10800", balances[checking.ID])
	}
	if balances[savings.ID] != 300 {
		t.Errorf("savings balance = %d, want 300", balances[savings.ID])
	}

	results, err := GetTransactions(checking.ID, "transfer", "", "", 50, 0)
	if err != nil {
		t.Fatalf("GetTransactions() error = %v", err)
	}
	if len(results) != 1 || results[0].ID != transactions[2].ID {
		t.Errorf("filtered transactions = %+v, want the account transfer", results)
	}
}

func TestCreateTransferRejectsDifferentCurrencies(t *testing.T) {
	setupLedgerDB(t)
	source := &Account{Name: "USD", Type: "checking", Currency: "USD"}
	target := &Account{Name: "EUR", Type: "checking", Currency: "EUR"}
	if err := CreateAccount(source); err != nil {
		t.Fatalf("CreateAccount(source) error = %v", err)
	}
	if err := CreateAccount(target); err != nil {
		t.Fatalf("CreateAccount(target) error = %v", err)
	}

	err := CreateTransaction(&Transaction{
		AccountID: source.ID, TransferAccountID: &target.ID, Type: "transfer",
		AmountMinor: 100, Currency: "USD", Date: "2026-10-12",
	})
	if !errors.Is(err, ErrInvalidTransfer) {
		t.Errorf("CreateTransaction() error = %v, want %v", err, ErrInvalidTransfer)
	}
}

func TestCreateTransactionRequiresPositiveExactAmountAndValidDate(t *testing.T) {
	setupLedgerDB(t)
	account := &Account{Name: "Cash", Type: "cash", Currency: "USD"}
	if err := CreateAccount(account); err != nil {
		t.Fatalf("CreateAccount() error = %v", err)
	}

	for _, transaction := range []Transaction{
		{AccountID: account.ID, Type: "expense", AmountMinor: 0, Currency: "USD", Date: "2026-10-12"},
		{AccountID: account.ID, Type: "expense", AmountMinor: 100, Currency: "USD", Date: "10/12/2026"},
		{AccountID: account.ID, Type: "expense", AmountMinor: 100, Currency: "XYZ", Date: "2026-10-12"},
	} {
		if err := CreateTransaction(&transaction); !errors.Is(err, ErrInvalidTxn) {
			t.Errorf("CreateTransaction(%+v) error = %v, want %v", transaction, err, ErrInvalidTxn)
		}
	}

}

func TestCreateAccountRejectsUnsupportedCurrency(t *testing.T) {
	setupLedgerDB(t)
	err := CreateAccount(&Account{Name: "Unknown", Type: "checking", Currency: "XYZ"})
	if !errors.Is(err, ErrInvalidAccount) {
		t.Errorf("CreateAccount() error = %v, want %v", err, ErrInvalidAccount)
	}
}

func TestGetTransactionsBoundsPagination(t *testing.T) {
	setupLedgerDB(t)
	if _, err := GetTransactions(0, "", "", "", 101, 0); !errors.Is(err, ErrInvalidTxn) {
		t.Errorf("GetTransactions() error = %v, want %v", err, ErrInvalidTxn)
	}
	if _, err := GetTransactions(0, "", "", "", 50, 10001); !errors.Is(err, ErrInvalidTxn) {
		t.Errorf("GetTransactions() error = %v, want %v", err, ErrInvalidTxn)
	}
}

func TestCategoryRulesApplyByPriorityAndPreserveExplicitCategory(t *testing.T) {
	database := setupLedgerDB(t)
	account := &Account{Name: "Checking", Type: "checking", Currency: "USD"}
	food := &category.Category{Name: "Food"}
	transport := &category.Category{Name: "Transport"}
	for _, record := range []interface{}{account, food, transport} {
		if err := database.Create(record).Error; err != nil {
			t.Fatalf("create fixture: %v", err)
		}
	}

	rules := []CategoryRule{
		{Name: "Broad", MatchValue: "market", CategoryID: food.ID, Priority: 1, Enabled: true},
		{Name: "Specific", MatchValue: "corner market", CategoryID: transport.ID, Priority: 2, Enabled: true},
	}
	for i := range rules {
		if err := CreateCategoryRule(&rules[i]); err != nil {
			t.Fatalf("CreateCategoryRule() error = %v", err)
		}
	}

	matched := Transaction{
		AccountID: account.ID, Type: "expense", AmountMinor: 1250,
		Currency: "USD", Date: "2026-10-12", Payee: "Corner Market #4",
	}
	if err := CreateTransaction(&matched); err != nil {
		t.Fatalf("CreateTransaction(matched) error = %v", err)
	}
	if matched.CategoryID == nil || *matched.CategoryID != transport.ID || matched.CategoryRuleID == nil || *matched.CategoryRuleID != rules[1].ID {
		t.Errorf("matched transaction category/rule = %v/%v, want %d/%d", matched.CategoryID, matched.CategoryRuleID, transport.ID, rules[1].ID)
	}

	explicitCategory := food.ID
	explicit := Transaction{
		AccountID: account.ID, CategoryID: &explicitCategory, Type: "expense", AmountMinor: 1250,
		Currency: "USD", Date: "2026-10-12", Payee: "Corner Market #4",
	}
	if err := CreateTransaction(&explicit); err != nil {
		t.Fatalf("CreateTransaction(explicit) error = %v", err)
	}
	if explicit.CategoryID == nil || *explicit.CategoryID != food.ID || explicit.CategoryRuleID != nil {
		t.Errorf("explicit category/rule = %v/%v, want %d/nil", explicit.CategoryID, explicit.CategoryRuleID, food.ID)
	}
}

func TestCashflowForecastUsesClampedMonthlyBillsAndOneTimeIncome(t *testing.T) {
	setupLedgerDB(t)
	account := &Account{Name: "Checking", Type: "checking", Currency: "USD", OpeningBalanceMinor: 10000}
	if err := CreateAccount(account); err != nil {
		t.Fatalf("CreateAccount() error = %v", err)
	}
	bill := &ScheduledFlow{
		AccountID: account.ID, Label: "Rent", Type: "expense", AmountMinor: 1000,
		Currency: "USD", Frequency: "monthly", NextDate: "2026-01-31", Enabled: true,
	}
	income := &ScheduledFlow{
		AccountID: account.ID, Label: "Bonus", Type: "income", AmountMinor: 3000,
		Currency: "USD", Frequency: "once", NextDate: "2026-02-15", Enabled: true,
	}
	for _, flow := range []*ScheduledFlow{bill, income} {
		if err := CreateScheduledFlow(flow); err != nil {
			t.Fatalf("CreateScheduledFlow() error = %v", err)
		}
	}

	forecast, err := BuildCashflowForecast(time.Date(2026, 1, 30, 12, 0, 0, 0, time.UTC), 30)
	if err != nil {
		t.Fatalf("BuildCashflowForecast() error = %v", err)
	}
	if len(forecast.Events) != 3 {
		t.Fatalf("forecast events = %+v, want Jan 31/Feb 28 bills and Feb 15 income", forecast.Events)
	}
	if forecast.Events[0].Date != "2026-01-31" || forecast.Events[1].Date != "2026-02-15" || forecast.Events[2].Date != "2026-02-28" {
		t.Errorf("forecast event dates = %v, want Jan 31, Feb 15, and Feb 28", []string{forecast.Events[0].Date, forecast.Events[1].Date, forecast.Events[2].Date})
	}
	if len(forecast.Accounts) != 1 || forecast.Accounts[0].ProjectedBalanceMinor != 11000 {
		t.Errorf("projected account balance = %+v, want 11000 minor units", forecast.Accounts)
	}
}

func TestScheduledFlowRequiresAccountCurrencyMatch(t *testing.T) {
	setupLedgerDB(t)
	account := &Account{Name: "Euro", Type: "checking", Currency: "EUR"}
	if err := CreateAccount(account); err != nil {
		t.Fatalf("CreateAccount() error = %v", err)
	}
	err := CreateScheduledFlow(&ScheduledFlow{
		AccountID: account.ID, Label: "Salary", Type: "income", AmountMinor: 1000,
		Currency: "USD", Frequency: "monthly", NextDate: "2026-01-01",
	})
	if !errors.Is(err, ErrInvalidSchedule) {
		t.Errorf("CreateScheduledFlow() error = %v, want %v", err, ErrInvalidSchedule)
	}
}

func TestBillInsightAndConfirmedOccurrenceArePriceAwareAndIdempotent(t *testing.T) {
	database := setupLedgerDB(t)
	account := &Account{Name: "Checking", Type: "checking", Currency: "USD", OpeningBalanceMinor: 100000}
	food := &category.Category{Name: "Utilities"}
	for _, record := range []interface{}{account, food} {
		if err := database.Create(record).Error; err != nil {
			t.Fatalf("create fixture: %v", err)
		}
	}
	categoryID := food.ID
	for _, transaction := range []Transaction{
		{AccountID: account.ID, CategoryID: &categoryID, Type: "expense", AmountMinor: 10000, Currency: "USD", Date: "2026-01-10", Payee: "Internet"},
		{AccountID: account.ID, CategoryID: &categoryID, Type: "expense", AmountMinor: 12000, Currency: "USD", Date: "2026-02-10", Payee: "Internet"},
	} {
		if err := database.Create(&transaction).Error; err != nil {
			t.Fatalf("create bill history: %v", err)
		}
	}
	flow := &ScheduledFlow{
		AccountID: account.ID, CategoryID: &categoryID, Label: "Internet", Type: "expense",
		AmountMinor: 12000, Currency: "USD", Frequency: "monthly", NextDate: "2026-03-10",
		ReminderDays: 3, Enabled: true,
	}
	if err := CreateScheduledFlow(flow); err != nil {
		t.Fatalf("CreateScheduledFlow() error = %v", err)
	}
	insights, err := GetBillInsights(time.Date(2026, 3, 8, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("GetBillInsights() error = %v", err)
	}
	if len(insights) != 1 || insights[0].Status != "due_soon" ||
		!insights[0].HasPriceChange || insights[0].ChangePercent != 20 {
		t.Fatalf("bill insight = %+v, want due-soon with 20%% price increase", insights)
	}

	request := ProcessOccurrenceRequest{ScheduledDate: "2026-03-10", Action: "record", AmountMinor: 13000}
	first, err := ProcessScheduledOccurrence(flow.ID, request)
	if err != nil {
		t.Fatalf("ProcessScheduledOccurrence() error = %v", err)
	}
	retry, err := ProcessScheduledOccurrence(flow.ID, request)
	if err != nil {
		t.Fatalf("idempotent ProcessScheduledOccurrence() error = %v", err)
	}
	if first.ID != retry.ID || first.TransactionID == nil || *first.TransactionID != *retry.TransactionID {
		t.Fatalf("idempotent occurrence = %+v, retry = %+v", first, retry)
	}
	_, err = ProcessScheduledOccurrence(flow.ID, ProcessOccurrenceRequest{
		ScheduledDate: request.ScheduledDate, Action: "record", AmountMinor: 14000,
	})
	if !errors.Is(err, ErrOccurrenceExists) {
		t.Fatalf("conflicting retry error = %v, want %v", err, ErrOccurrenceExists)
	}
	var recorded Transaction
	if err := database.First(&recorded, *first.TransactionID).Error; err != nil {
		t.Fatalf("fetch recorded transaction: %v", err)
	}
	if recorded.AmountMinor != 13000 || recorded.CategoryID == nil || *recorded.CategoryID != categoryID {
		t.Fatalf("recorded transaction = %+v", recorded)
	}
	var advanced ScheduledFlow
	if err := database.First(&advanced, flow.ID).Error; err != nil {
		t.Fatal(err)
	}
	if advanced.NextDate != "2026-04-10" {
		t.Errorf("next occurrence date = %q, want 2026-04-10", advanced.NextDate)
	}
}

func TestCashflowAndNetWorthReportsRemainCurrencyAware(t *testing.T) {
	database := setupLedgerDB(t)
	checking := &Account{Name: "Checking", Type: "checking", Currency: "USD", OpeningBalanceMinor: 10000}
	card := &Account{Name: "Card", Type: "credit_card", Currency: "USD"}
	food := &category.Category{Name: "Food"}
	for _, record := range []interface{}{checking, card, food} {
		if err := database.Create(record).Error; err != nil {
			t.Fatalf("create fixture: %v", err)
		}
	}
	categoryID := food.ID
	for _, transaction := range []Transaction{
		{AccountID: checking.ID, Type: "income", AmountMinor: 5000, Currency: "USD", Date: "2026-10-12"},
		{AccountID: checking.ID, CategoryID: &categoryID, Type: "expense", AmountMinor: 1500, Currency: "USD", Date: "2026-10-13"},
		{AccountID: card.ID, Type: "expense", AmountMinor: 2000, Currency: "USD", Date: "2026-10-13"},
		{AccountID: checking.ID, TransferAccountID: &card.ID, Type: "transfer", AmountMinor: 100, Currency: "USD", Date: "2026-10-14"},
	} {
		if err := CreateTransaction(&transaction); err != nil {
			t.Fatalf("CreateTransaction() error = %v", err)
		}
	}

	report, err := GetCashflowReport("2026-10-01", "2026-10-31")
	if err != nil {
		t.Fatalf("GetCashflowReport() error = %v", err)
	}
	if len(report.Currencies) != 1 || report.Currencies[0].IncomeMinor != 5000 ||
		report.Currencies[0].ExpenseMinor != 3500 || report.Currencies[0].NetMinor != 1500 {
		t.Errorf("cashflow report = %+v, want USD income 5000, expense 3500, net 1500", report)
	}
	if len(report.Currencies[0].Categories) != 2 {
		t.Errorf("cashflow expense categories = %+v, want categorized and uncategorized totals", report.Currencies[0].Categories)
	}

	netWorth, err := GetNetWorthReport(time.Date(2026, 10, 31, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("GetNetWorthReport() error = %v", err)
	}
	if len(netWorth.Currencies) != 1 || netWorth.Currencies[0].AssetsMinor != 13400 ||
		netWorth.Currencies[0].LiabilitiesMinor != 1900 || netWorth.Currencies[0].NetWorthMinor != 11500 {
		t.Errorf("net worth report = %+v, want USD assets 13400, liabilities 1900, net 11500", netWorth)
	}
	if _, err := GetCashflowReport("2026-11-01", "2026-10-01"); !errors.Is(err, ErrInvalidTxn) {
		t.Errorf("reversed report period error = %v, want %v", err, ErrInvalidTxn)
	}
}
