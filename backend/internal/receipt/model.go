package receipt

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/girishsaraf/xpense-trackr/backend/internal/ledger"
	"gorm.io/gorm"
)

const MaxFileSize = 10 << 20

var (
	ErrInvalidReceipt = errors.New("receipt must be a PDF, PNG, or JPEG linked to a transaction and no larger than 10 MB")
	db                *gorm.DB
	storageDir        = "./data/receipts"
)

type Receipt struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	TransactionID uint   `gorm:"not null;index" json:"transaction_id"`
	FileName      string `gorm:"size:255;not null" json:"file_name"`
	MediaType     string `gorm:"size:80;not null" json:"media_type"`
	SizeBytes     int64  `gorm:"not null" json:"size_bytes"`
	StorageKey    string `gorm:"size:64;not null;uniqueIndex" json:"-"`
	CreatedAt     int64  `gorm:"autoCreateTime" json:"created_at"`
}

func SetDB(database *gorm.DB) {
	db = database
}

func SetStorageDir(directory string) {
	storageDir = filepath.Clean(directory)
}

func GetReceipts(transactionID uint) ([]Receipt, error) {
	receipts := make([]Receipt, 0)
	query := db.Order("created_at DESC, id DESC")
	if transactionID > 0 {
		query = query.Where("transaction_id = ?", transactionID)
	}
	if err := query.Find(&receipts).Error; err != nil {
		return nil, err
	}
	return receipts, nil
}

func GetReceipt(id uint) (*Receipt, error) {
	var receipt Receipt
	if err := db.First(&receipt, id).Error; err != nil {
		return nil, err
	}
	return &receipt, nil
}

func DeleteReceipt(id uint) error {
	var item Receipt
	if err := db.First(&item, id).Error; err != nil {
		return err
	}
	path, err := filePath(item.StorageKey)
	if err != nil {
		return err
	}
	trashPath := path + ".deleting"
	if err := os.Rename(path, trashPath); err != nil {
		return fmt.Errorf("stage receipt file for deletion: %w", err)
	}
	if err := db.Delete(&item).Error; err != nil {
		if restoreErr := os.Rename(trashPath, path); restoreErr != nil {
			return fmt.Errorf("delete receipt metadata: %v; restore receipt file: %w", err, restoreErr)
		}
		return err
	}
	if err := os.Remove(trashPath); err != nil {
		return fmt.Errorf("remove receipt file after deleting metadata: %w", err)
	}
	return nil
}

func filePath(key string) (string, error) {
	if len(key) != 64 {
		return "", ErrInvalidReceipt
	}
	for _, character := range key {
		if !(character >= '0' && character <= '9') && !(character >= 'a' && character <= 'f') {
			return "", ErrInvalidReceipt
		}
	}
	return filepath.Join(storageDir, key), nil
}

func validateTransaction(transactionID uint) error {
	var transaction ledger.Transaction
	if err := db.First(&transaction, transactionID).Error; err != nil {
		return fmt.Errorf("receipt transaction: %w", err)
	}
	return nil
}
