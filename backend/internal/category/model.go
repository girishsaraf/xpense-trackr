// internal/category/model.go

package category

import (
	"errors"
	"gorm.io/gorm"
)

var db *gorm.DB

// SetDB sets the db object for the category package.
func SetDB(database *gorm.DB) {
	db = database
}

// ErrRecordNotFound is an error indicating that a record was not found.
var ErrRecordNotFound = gorm.ErrRecordNotFound

// Category represents the category model.
type Category struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

// GetAllCategories retrieves all categories from the database.
func GetAllCategories() ([]Category, error) {
	var categories []Category
	if err := db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetCategoryByID retrieves a specific category by its ID from the database.
func GetCategoryByID(id int) (*Category, error) {
	var category Category
	if err := db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &category, nil
}

// CreateCategory creates a new category in the database.
func CreateCategory(category *Category) error {
	if err := db.Create(&category).Error; err != nil {
		return err
	}
	return nil
}

// UpdateCategory updates an existing category in the database.
func UpdateCategory(id int, updatedCategory *Category) error {
	if err := db.Model(&Category{}).Where("id = ?", id).Updates(updatedCategory).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRecordNotFound
		}
		return err
	}
	return nil
}

// DeleteCategory deletes an existing category from the database.
func DeleteCategory(id int) error {
	if err := db.Delete(&Category{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRecordNotFound
		}
		return err
	}
	return nil
}
