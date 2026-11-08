package recurring

import (
	"log"
	"time"

	"github.com/girishsaraf/xpense-trackr/backend/internal/expense"
	"github.com/girishsaraf/xpense-trackr/backend/internal/investments"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var db *gorm.DB

// SetDB sets the db object for the budget package.
func SetDB(databaseConfig *gorm.DB) {
	db = databaseConfig
}

func HandleRecurringExpense() {
	if err := FetchRecurringExpense(); err != nil {
		log.Printf("process recurring expenses: %v", err)
	}
}

func HandleRecurringInvestment() {
	if err := FetchRecurringInvestment(); err != nil {
		log.Printf("process recurring investments: %v", err)
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
		return lastRunDate.Before(time.Now().AddDate(0, 0, -7))
	case "biweekly":
		return lastRunDate.Before(time.Now().AddDate(0, 0, -14))
	case "monthly":
		now := time.Now()
		return lastRunDate.Year() < now.Year() || lastRunDate.Year() == now.Year() && lastRunDate.Month() < now.Month()
	default:
		// Unsupported frequency
		return false
	}
}

func getCurrDayMonth() string {
	today := time.Now()
	return "-" + today.Format("01/02")
}

func insertExpense(currExpense RecurringExpense) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var current RecurringExpense
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&current, currExpense.ID).Error; err != nil {
			return err
		}
		if !checkIfExpenseOrInvestmentIsDue(current.Frequency, current.StartDate) {
			return nil
		}

		// Insert expense into the Expense table
		newExpense := expense.Expense{
			CategoryID:  current.CategoryID,
			Amount:      current.Amount,
			Description: current.Description + getCurrDayMonth(),
			Date:        time.Now().Format("2006-01-02"),
		}
		if err := tx.Create(&newExpense).Error; err != nil {
			return err
		}

		current.StartDate = time.Now()
		return tx.Save(&current).Error
	})
}

func insertInvestment(investment RecurringInvestment) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var current RecurringInvestment
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&current, investment.ID).Error; err != nil {
			return err
		}
		if !checkIfExpenseOrInvestmentIsDue(current.Frequency, current.StartDate) {
			return nil
		}

		// Insert investment into the Investment table
		newInvestment := investments.Investment{
			Amount:         current.Amount,
			InvestmentType: current.InvestmentType,
			Date:           time.Now().Format("2006-01-02"),
		}
		if err := tx.Create(&newInvestment).Error; err != nil {
			return err
		}

		current.StartDate = time.Now()
		return tx.Save(&current).Error
	})
}
