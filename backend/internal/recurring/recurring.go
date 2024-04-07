package recurring

import (
	"github.com/girishsaraf/xpense-trackr/backend/internal/expense"
	"github.com/girishsaraf/xpense-trackr/backend/internal/investments"
	"gorm.io/gorm"
	"time"
)

var db *gorm.DB

// SetDB sets the db object for the budget package.
func SetDB(databaseConfig *gorm.DB) {
	db = databaseConfig
}

func HandleRecurringExpense() {
	err := FetchRecurringExpense()
	if err != nil {
		return
	}
}

func HandleRecurringInvestment() {
	err := FetchRecurringInvestment()
	if err != nil {
		return
	}
}

func FetchRecurringExpense() error {
	// Get recurring expenses due for today
	recurringExpenses, err := getRecurringExpensesForDate(time.Now())
	if err != nil {
		return err
	}

	// Insert recurring expenses into the database
	for _, pendingExpense := range recurringExpenses {
		err := insertExpense(pendingExpense)
		if err != nil {
			return err
		}
	}

	return nil
}

func FetchRecurringInvestment() error {
	// Get recurring investments due for today
	recurringInvestments, err := getRecurringInvestmentsForDate(time.Now())
	if err != nil {
		return err
	}

	// Insert recurring investments into the database
	for _, investment := range recurringInvestments {
		err := insertInvestment(investment)
		if err != nil {
			return err
		}
	}

	return nil
}

func getRecurringExpensesForDate(date time.Time) ([]RecurringExpense, error) {
	var recurringExpenses []RecurringExpense
	// Query recurring_expense table for expenses due for the given date
	err := db.Where("start_date <= ? AND (end_date IS NULL OR end_date >= ?)", date, date).Find(&recurringExpenses).Error
	if err != nil {
		return nil, err
	}
	return recurringExpenses, nil
}

func getRecurringInvestmentsForDate(date time.Time) ([]RecurringInvestment, error) {
	var recurringInvestments []RecurringInvestment
	// Query recurring_investment table for investments due for the given date
	err := db.Where("start_date <= ? AND (end_date IS NULL OR end_date >= ?)", date, date).Find(&recurringInvestments).Error
	if err != nil {
		return nil, err
	}
	return recurringInvestments, nil
}

func checkIfExpenseOrInvestmentIsDue(frequency string, lastRunDate time.Time) bool {
	switch frequency {
	case "daily":
		// Check if last run date is before today
		return lastRunDate.Before(time.Now().Truncate(24 * time.Hour))
	case "weekly":
		// Check if last run date is before the start of this week
		startOfWeek := time.Now().Truncate(24 * 7 * time.Hour)
		return lastRunDate.Before(startOfWeek)
	case "biweekly":
		// Check if last run date is before the start of this biweekly
		startOfBiweekly := time.Now().Truncate(24 * 14 * time.Hour)
		return lastRunDate.Before(startOfBiweekly)
	case "monthly":
		// Check if last run date is before the start of this month
		startOfMonth := time.Now().Truncate(24 * 30 * time.Hour) // Approximation for a month
		return lastRunDate.Before(startOfMonth)
	default:
		// Unsupported frequency
		return false
	}
}

func getCurrDayMonth() string {
	today := time.Now()
	return "-" + today.Format("DD/MM")
}

func insertExpense(currExpense RecurringExpense) error {
	if checkIfExpenseOrInvestmentIsDue(currExpense.Frequency, currExpense.StartDate) {
		// Insert expense into the Expense table
		newExpense := expense.Expense{
			CategoryID:  currExpense.CategoryID,
			Amount:      currExpense.Amount,
			Description: currExpense.Description + getCurrDayMonth(),
			Date:        time.Now().String(),
		}
		err := db.Create(&newExpense).Error
		if err != nil {
			return err
		}

		// Update the recurring expense's start date to today's date
		currExpense.StartDate = time.Now()
		err = db.Save(&currExpense).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func insertInvestment(investment RecurringInvestment) error {

	if checkIfExpenseOrInvestmentIsDue(investment.Frequency, investment.StartDate) {
		// Insert investment into the Investment table
		newInvestment := investments.Investment{
			Amount:         investment.Amount,
			InvestmentType: investment.InvestmentType,
			Date:           time.Now().String(),
		}
		err := db.Create(&newInvestment).Error
		if err != nil {
			return err
		}

		// Update the recurring expense's start date to today's date
		investment.StartDate = time.Now()
		err = db.Save(&investment).Error
		if err != nil {
			return err
		}
	}
	return nil
}
