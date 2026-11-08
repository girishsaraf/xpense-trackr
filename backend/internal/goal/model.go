package goal

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	db                *gorm.DB
	ErrRecordNotFound = gorm.ErrRecordNotFound
	ErrInvalidGoal    = errors.New("goal name and a positive target amount are required")
	ErrInvalidDate    = errors.New("target_date must use YYYY-MM-DD format")
)

type Goal struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"size:120;not null" json:"name"`
	TargetAmount  float64   `gorm:"not null" json:"target_amount"`
	CurrentAmount float64   `gorm:"not null;default:0" json:"current_amount"`
	TargetDate    string    `gorm:"size:10" json:"target_date,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func SetDB(database *gorm.DB) {
	db = database
}

func validate(goal *Goal) error {
	if strings.TrimSpace(goal.Name) == "" || goal.TargetAmount <= 0 || goal.CurrentAmount < 0 {
		return ErrInvalidGoal
	}
	if goal.TargetDate != "" {
		if _, err := time.Parse("2006-01-02", goal.TargetDate); err != nil {
			return ErrInvalidDate
		}
	}
	return nil
}

func GetAllGoals() ([]Goal, error) {
	goals := make([]Goal, 0)
	if err := db.Order("target_date IS NULL, target_date ASC, id ASC").Find(&goals).Error; err != nil {
		return nil, err
	}
	return goals, nil
}

func GetGoalByID(id int) (*Goal, error) {
	var goal Goal
	if err := db.First(&goal, id).Error; err != nil {
		return nil, err
	}
	return &goal, nil
}

func CreateGoal(goal *Goal) error {
	if err := validate(goal); err != nil {
		return err
	}
	return db.Create(goal).Error
}

func UpdateGoal(id int, updated *Goal) error {
	if err := validate(updated); err != nil {
		return err
	}
	var existing Goal
	if err := db.First(&existing, id).Error; err != nil {
		return err
	}
	return db.Model(&existing).Select("Name", "TargetAmount", "CurrentAmount", "TargetDate").Updates(updated).Error
}

func DeleteGoal(id int) error {
	result := db.Delete(&Goal{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
