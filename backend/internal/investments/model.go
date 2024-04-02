// internal/model/Investment.go

package investments

import (
	"errors"
	"gorm.io/gorm"
)

var db *gorm.DB

// SetDB sets the db object for the Investment package.
func SetDB(databaseConfig *gorm.DB) {
	db = databaseConfig
}

// ErrRecordNotFound is an error indicating that a record was not found.
var ErrRecordNotFound = errors.New("record not found")

// Investment represents the Investment model.
type Investment struct {
	ID             uint    `gorm:"primaryKey" json:"id"`
	Amount         float64 `json:"amount,string,omitempty"`
	InvestmentType string  `json:"investment_type"`
	Date           string  `json:"date"`
}

// GetAllInvestments retrieves all Investments from the database.
func GetAllInvestments() ([]Investment, error) {
	var Investments []Investment
	if err := db.Find(&Investments).Error; err != nil {
		return nil, err
	}
	return Investments, nil
}

// GetInvestmentByID retrieves a specific Investment by its ID from the database.
func GetInvestmentByID(id int) (*Investment, error) {
	var Investment Investment
	if err := db.First(&Investment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &Investment, nil
}

// CreateInvestment creates a new Investment in the database.
func CreateInvestment(Investment *Investment) error {
	if err := db.Create(&Investment).Error; err != nil {
		return err
	}
	return nil
}

// UpdateInvestment updates an existing Investment in the database.
func UpdateInvestment(id int, updatedInvestment *Investment) error {
	if err := db.Model(&Investment{}).Where("id = ?", id).Updates(updatedInvestment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRecordNotFound
		}
		return err
	}
	return nil
}

// DeleteInvestment deletes an existing Investment from the database.
func DeleteInvestment(id int) error {
	if err := db.Delete(&Investment{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRecordNotFound
		}
		return err
	}
	return nil
}
