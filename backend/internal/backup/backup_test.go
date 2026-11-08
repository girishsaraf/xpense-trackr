package backup

import (
	"archive/zip"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/girishsaraf/xpense-trackr/backend/config"
)

func TestCreateAndRestoreRoundTrip(t *testing.T) {
	root := t.TempDir()
	dataDir := filepath.Join(root, "data")
	dbPath := filepath.Join(dataDir, "app.sqlite")
	archivePath := filepath.Join(root, "backup.zip")
	if err := os.MkdirAll(filepath.Join(dataDir, "receipts"), 0700); err != nil {
		t.Fatal(err)
	}

	db, err := config.OpenDatabase(config.DatabaseConfig{Driver: "sqlite", Path: dbPath})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("CREATE TABLE sample (value TEXT NOT NULL)").Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("INSERT INTO sample (value) VALUES (?)", "before").Error; err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	if err := sqlDB.Close(); err != nil {
		t.Fatal(err)
	}
	receiptKey := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	receiptPath := filepath.Join(dataDir, "receipts", receiptKey)
	if err := os.WriteFile(receiptPath, []byte("receipt before"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := Create(dbPath, dataDir, archivePath); err != nil {
		t.Fatalf("create backup: %v", err)
	}

	db, err = config.OpenDatabase(config.DatabaseConfig{Driver: "sqlite", Path: dbPath})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("UPDATE sample SET value = ?", "after").Error; err != nil {
		t.Fatal(err)
	}
	sqlDB, err = db.DB()
	if err != nil {
		t.Fatal(err)
	}
	if err := sqlDB.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(receiptPath, []byte("receipt after"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := Restore(archivePath, dbPath, dataDir); err != nil {
		t.Fatalf("restore backup: %v", err)
	}
	db, err = config.OpenDatabase(config.DatabaseConfig{Driver: "sqlite", Path: dbPath})
	if err != nil {
		t.Fatal(err)
	}
	var restoredValue string
	if err := db.Raw("SELECT value FROM sample").Scan(&restoredValue).Error; err != nil {
		t.Fatal(err)
	}
	sqlDB, err = db.DB()
	if err != nil {
		t.Fatal(err)
	}
	if err := sqlDB.Close(); err != nil {
		t.Fatal(err)
	}
	if restoredValue != "before" {
		t.Fatalf("restored database value = %q, want before", restoredValue)
	}
	restoredReceipt, err := os.ReadFile(receiptPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(restoredReceipt) != "receipt before" {
		t.Fatalf("restored receipt = %q, want original content", restoredReceipt)
	}
}

func TestRestoreRejectsTraversalArchive(t *testing.T) {
	root := t.TempDir()
	archivePath := filepath.Join(root, "bad.zip")
	file, err := os.Create(archivePath)
	if err != nil {
		t.Fatal(err)
	}
	writer := zip.NewWriter(file)
	manifestFile, err := writer.Create("manifest.json")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := manifestFile.Write([]byte(`{"version":1,"created_at":"2026-11-08T00:00:00Z"}`)); err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Create("receipts/../../outside"); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	err = Restore(archivePath, filepath.Join(root, "app.sqlite"), filepath.Join(root, "data"))
	if !errors.Is(err, ErrInvalidArchive) {
		t.Fatalf("restore error = %v, want invalid archive", err)
	}
	if _, err := os.Stat(filepath.Join(root, "outside")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("traversal destination exists or could not be checked: %v", err)
	}
}
