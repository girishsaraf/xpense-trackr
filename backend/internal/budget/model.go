// internal/model/budget.go

package budget

import (
	"errors"
	"gorm.io/gorm"
)

var db *gorm.DB

// SetDB sets the db object for the budget package.
func SetDB(databaseConfig *gorm.DB) {
	db = databaseConfig
}

// ErrRecordNotFound is an error indicating that a record was not found.
var ErrRecordNotFound = errors.New("record not found")

// Budget represents the budget model.
type Budget struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	UserID     uint    `json:"user_id"`
	CategoryID uint    `json:"category_id"`
	Amount     float64 `json:"amount"`
}

// GetAllBudgets retrieves all budgets from the database.
func GetAllBudgets() ([]Budget, error) {
	var budgets []Budget
	if err := db.Find(&budgets).Error; err != nil {
		return nil, err
	}
	return budgets, nil
}

// GetBudgetByID retrieves a specific budget by its ID from the database.
func GetBudgetByID(id int) (*Budget, error) {
	var budget Budget
	if err := db.First(&budget, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &budget, nil
}

// CreateBudget creates a new budget in the database.
func CreateBudget(budget *Budget) error {
	if err := db.Create(&budget).Error; err != nil {
		return err
	}
	return nil
}

// UpdateBudget updates an existing budget in the database.
func UpdateBudget(id int, updatedBudget *Budget) error {
	if err := db.Model(&Budget{}).Where("id = ?", id).Updates(updatedBudget).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRecordNotFound
		}
		return err
	}
	return nil
}

// DeleteBudget deletes an existing budget from the database.
func DeleteBudget(id int) error {
	if err := db.Delete(&Budget{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRecordNotFound
		}
		return err
	}
	return nil
}
