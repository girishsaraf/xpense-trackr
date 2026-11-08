package ledger

import (
	"math"
	"sort"
	"strconv"
	"time"
)

type FinancialInsight struct {
	Kind                string  `json:"kind"`
	Title               string  `json:"title"`
	Explanation         string  `json:"explanation"`
	Currency            string  `json:"currency"`
	CategoryName        string  `json:"category_name,omitempty"`
	CurrentAmountMinor  int64   `json:"current_amount_minor"`
	PreviousAmountMinor int64   `json:"previous_amount_minor"`
	ChangePercent       float64 `json:"change_percent,omitempty"`
}

type FinancialInsightsReport struct {
	From         string             `json:"from"`
	To           string             `json:"to"`
	PreviousFrom string             `json:"previous_from"`
	PreviousTo   string             `json:"previous_to"`
	Insights     []FinancialInsight `json:"insights"`
}

func GetFinancialInsights(from, to string) (*FinancialInsightsReport, error) {
	fromDate, err := time.Parse("2006-01-02", from)
	if err != nil {
		return nil, ErrInvalidTxn
	}
	toDate, err := time.Parse("2006-01-02", to)
	if err != nil || toDate.Before(fromDate) || toDate.After(fromDate.AddDate(1, 0, 0)) {
		return nil, ErrInvalidTxn
	}
	current, err := GetCashflowReport(from, to)
	if err != nil {
		return nil, err
	}
	periodDays := int(toDate.Sub(fromDate).Hours()/24) + 1
	previousTo := fromDate.AddDate(0, 0, -1)
	previousFrom := previousTo.AddDate(0, 0, -(periodDays - 1))
	previous, err := GetCashflowReport(previousFrom.Format("2006-01-02"), previousTo.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	report := &FinancialInsightsReport{
		From: from, To: to,
		PreviousFrom: previousFrom.Format("2006-01-02"),
		PreviousTo:   previousTo.Format("2006-01-02"),
		Insights:     make([]FinancialInsight, 0),
	}
	previousByCurrency := make(map[string]CashflowCurrencyReport, len(previous.Currencies))
	currentByCurrency := make(map[string]CashflowCurrencyReport, len(current.Currencies))
	currencySet := make(map[string]struct{}, len(current.Currencies)+len(previous.Currencies))
	for _, currency := range previous.Currencies {
		previousByCurrency[currency.Currency] = currency
		currencySet[currency.Currency] = struct{}{}
	}
	for _, currentCurrency := range current.Currencies {
		currentByCurrency[currentCurrency.Currency] = currentCurrency
		currencySet[currentCurrency.Currency] = struct{}{}
	}
	currencies := make([]string, 0, len(currencySet))
	for currency := range currencySet {
		currencies = append(currencies, currency)
	}
	sort.Strings(currencies)
	for _, currency := range currencies {
		currentCurrency := currentByCurrency[currency]
		previousCurrency := previousByCurrency[currency]
		if insight, ok := spendingInsight(currency, currentCurrency.ExpenseMinor, previousCurrency.ExpenseMinor); ok {
			report.Insights = append(report.Insights, insight)
		}
		categoryInsights := categorySpendingInsights(currency, currentCurrency.Categories, previousCurrency.Categories)
		report.Insights = append(report.Insights, categoryInsights...)
	}
	return report, nil
}

func spendingInsight(currency string, current, previous int64) (FinancialInsight, bool) {
	if current == 0 && previous == 0 {
		return FinancialInsight{}, false
	}
	if previous == 0 {
		return FinancialInsight{
			Kind: "spending_started", Title: "Spending appeared in " + currency,
			Explanation: "You recorded " + formatAmountForCSV(current, currency) + " in this period; there was no recorded spending in the previous period.",
			Currency:    currency, CurrentAmountMinor: current,
		}, current >= 500
	}
	change := current - previous
	percent := math.Round(float64(change)/float64(previous)*1000) / 10
	if math.Abs(percent) < 10 || absMinor(change) < 500 {
		return FinancialInsight{}, false
	}
	title := ""
	if change > 0 {
		title = "Spending increased by " + formatPercent(math.Abs(percent)) + " in " + currency
	} else if current == 0 {
		title = "Spending stopped in " + currency
	} else {
		title = "Spending decreased by " + formatPercent(math.Abs(percent)) + " in " + currency
	}
	return FinancialInsight{
		Kind:  "spending_change",
		Title: title,
		Explanation: "Recorded spending changed from " + formatAmountForCSV(previous, currency) +
			" to " + formatAmountForCSV(current, currency) + " compared with the previous period.",
		Currency: currency, CurrentAmountMinor: current, PreviousAmountMinor: previous, ChangePercent: percent,
	}, true
}

func categorySpendingInsights(currency string, current, previous []CashflowCategoryTotal) []FinancialInsight {
	type categoryAmounts struct {
		name     string
		current  int64
		previous int64
	}
	amounts := make(map[string]*categoryAmounts, len(current)+len(previous))
	keyFor := func(item CashflowCategoryTotal) string {
		if item.CategoryID == nil {
			return "uncategorized"
		}
		return formatID(*item.CategoryID)
	}
	for _, item := range current {
		key := keyFor(item)
		if amounts[key] == nil {
			amounts[key] = &categoryAmounts{name: item.CategoryName}
		}
		amounts[key].current = item.AmountMinor
	}
	for _, item := range previous {
		key := keyFor(item)
		if amounts[key] == nil {
			amounts[key] = &categoryAmounts{name: item.CategoryName}
		}
		amounts[key].previous = item.AmountMinor
	}
	insights := make([]FinancialInsight, 0, len(amounts))
	for _, item := range amounts {
		if item.name == "" {
			item.name = "Uncategorized"
		}
		if item.previous == 0 && item.current == 0 {
			continue
		}
		if item.previous == 0 {
			if item.current < 500 {
				continue
			}
			insights = append(insights, FinancialInsight{
				Kind: "category_spending_started", Title: "New spending in " + item.name,
				Explanation: "You recorded " + formatAmountForCSV(item.current, currency) +
					" in " + item.name + "; there was no recorded spending in the previous period.",
				Currency: currency, CategoryName: item.name, CurrentAmountMinor: item.current,
			})
			continue
		}
		change := item.current - item.previous
		percent := math.Round(float64(change)/float64(item.previous)*1000) / 10
		if math.Abs(percent) < 20 || absMinor(change) < 500 {
			continue
		}
		title := ""
		if change > 0 {
			title = item.name + " spending increased by " + formatPercent(math.Abs(percent))
		} else if item.current == 0 {
			title = "Spending stopped in " + item.name
		} else {
			title = item.name + " spending decreased by " + formatPercent(math.Abs(percent))
		}
		insights = append(insights, FinancialInsight{
			Kind:  "category_spending_change",
			Title: title,
			Explanation: "Recorded spending in " + item.name + " changed from " +
				formatAmountForCSV(item.previous, currency) + " to " +
				formatAmountForCSV(item.current, currency) + ".",
			Currency: currency, CategoryName: item.name, CurrentAmountMinor: item.current,
			PreviousAmountMinor: item.previous, ChangePercent: percent,
		})
	}
	sort.Slice(insights, func(i, j int) bool {
		leftChange := absMinor(insights[i].CurrentAmountMinor - insights[i].PreviousAmountMinor)
		rightChange := absMinor(insights[j].CurrentAmountMinor - insights[j].PreviousAmountMinor)
		if leftChange != rightChange {
			return leftChange > rightChange
		}
		return insights[i].CategoryName < insights[j].CategoryName
	})
	if len(insights) > 3 {
		insights = insights[:3]
	}
	return insights
}

func absMinor(value int64) int64 {
	if value < 0 {
		return -value
	}
	return value
}

func formatPercent(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64) + "%"
}

func formatID(value uint) string {
	return strconv.FormatUint(uint64(value), 10)
}
