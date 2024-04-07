package recurring

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

// ErrRecordNotFound is an error indicating that a record was not found.
var ErrRecordNotFound = gorm.ErrRecordNotFound

type RecurringExpense struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CategoryID  uint      `json:"category_id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Frequency   string    `json:"frequency"`
}

type RecurringInvestment struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Amount         float64   `json:"amount"`
	InvestmentType string    `json:"investment_type"`
	Description    string    `json:"description"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Frequency      string    `json:"frequency"`
}

// GetAllRecurringExpenses retrieves all expenses from the database.
func GetAllRecurringExpenses() ([]RecurringExpense, error) {
	var expenses []RecurringExpense
	if err := db.Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

// GetRecurringExpenseByID retrieves a specific expense by its ID from the database.
func GetRecurringExpenseByID(id int) (*RecurringExpense, error) {
	var expense RecurringExpense
	if err := db.First(&expense, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &expense, nil
}

// CreateRecurringExpense creates a new expense in the database.
func CreateRecurringExpense(expense *RecurringExpense) error {
	if err := db.Create(&expense).Error; err != nil {
		return err
	}
	return nil
}

// UpdateRecurringExpense updates an existing expense in the database.
func UpdateRecurringExpense(id int, updatedExpense *RecurringExpense) error {
	if err := db.Model(&RecurringExpense{}).Where("id = ?", id).Updates(updatedExpense).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRecordNotFound
		}
		return err
	}
	return nil
}

// DeleteRecurringExpense deletes an existing expense from the database.
func DeleteRecurringExpense(id int) error {
	if err := db.Delete(&RecurringExpense{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRecordNotFound
		}
		return err
	}
	return nil
}

// GetAllRecurringInvestments retrieves all expenses from the database.
func GetAllRecurringInvestments() ([]RecurringInvestment, error) {
	var expenses []RecurringInvestment
	if err := db.Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

// GetRecurringInvestmentByID retrieves a specific expense by its ID from the database.
func GetRecurringInvestmentByID(id int) (*RecurringInvestment, error) {
	var expense RecurringInvestment
	if err := db.First(&expense, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &expense, nil
}

// CreateRecurringInvestment creates a new expense in the database.
func CreateRecurringInvestment(expense *RecurringInvestment) error {
	if err := db.Create(&expense).Error; err != nil {
		return err
	}
	return nil
}

// UpdateRecurringInvestment updates an existing expense in the database.
func UpdateRecurringInvestment(id int, updatedInvestment *RecurringInvestment) error {
	if err := db.Model(&RecurringInvestment{}).Where("id = ?", id).Updates(updatedInvestment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRecordNotFound
		}
		return err
	}
	return nil
}

// DeleteRecurringInvestment deletes an existing expense from the database.
func DeleteRecurringInvestment(id int) error {
	if err := db.Delete(&RecurringInvestment{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRecordNotFound
		}
		return err
	}
	return nil
}
