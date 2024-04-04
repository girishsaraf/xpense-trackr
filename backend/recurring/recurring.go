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
	for _, expense := range recurringExpenses {
		err := insertExpense(expense)
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

func getCurrDayMonth() string {
	today := time.Now()
	return "-" + today.Format("DD/MM")
}

func insertExpense(currExpense RecurringExpense) error {
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
	return nil
}

func insertInvestment(investment RecurringInvestment) error {
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
	return nil
}
