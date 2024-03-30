// internal/expense/model.go

package expense

import (
	"errors"
	"gorm.io/gorm"
)

var db *gorm.DB

// SetDB sets the db object for the expense package.
func SetDB(databaseConfig *gorm.DB) {
	db = databaseConfig
}

// ErrRecordNotFound is an error indicating that a record was not found.
var ErrRecordNotFound = gorm.ErrRecordNotFound

// Expense represents the expense model.
type Expense struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	UserID      uint    `json:"user_id"`
	CategoryID  uint    `json:"category_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

// GetAllExpenses retrieves all expenses from the database.
func GetAllExpenses() ([]Expense, error) {
	var expenses []Expense
	if err := db.Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

// GetExpenseByID retrieves a specific expense by its ID from the database.
func GetExpenseByID(id int) (*Expense, error) {
	var expense Expense
	if err := db.First(&expense, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &expense, nil
}

// CreateExpense creates a new expense in the database.
func CreateExpense(expense *Expense) error {
	if err := db.Create(&expense).Error; err != nil {
		return err
	}
	return nil
}

// UpdateExpense updates an existing expense in the database.
func UpdateExpense(id int, updatedExpense *Expense) error {
	if err := db.Model(&Expense{}).Where("id = ?", id).Updates(updatedExpense).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRecordNotFound
		}
		return err
	}
	return nil
}

// DeleteExpense deletes an existing expense from the database.
func DeleteExpense(id int) error {
	if err := db.Delete(&Expense{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRecordNotFound
		}
		return err
	}
	return nil
}
