package backup

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/girishsaraf/xpense-trackr/backend/config"
	"gorm.io/gorm"
)

const (
	archiveVersion = 1
	maxArchiveSize = int64(2 << 30)
	maxReceiptSize = int64(10 << 20)
	maxEntries     = 10002
)

var ErrInvalidArchive = errors.New("invalid xpense-trackr backup archive")

type manifest struct {
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
}

// Create writes a consistent SQLite snapshot and its receipt files to a ZIP.
func Create(databasePath, dataDir, archivePath string) error {
	if databasePath == "" || dataDir == "" || archivePath == "" {
		return errors.New("database, data, and archive paths are required")
	}
	databaseAbsolute, err := filepath.Abs(databasePath)
	if err != nil {
		return fmt.Errorf("resolve database path: %w", err)
	}
	archiveAbsolute, err := filepath.Abs(archivePath)
	if err != nil {
		return fmt.Errorf("resolve archive path: %w", err)
	}
	if databaseAbsolute == archiveAbsolute {
		return errors.New("backup archive must not replace the database")
	}
	if err := os.MkdirAll(filepath.Dir(archivePath), 0700); err != nil {
		return fmt.Errorf("create backup directory: %w", err)
	}
	workDir, err := os.MkdirTemp("", "xpense-backup-*")
	if err != nil {
		return fmt.Errorf("create backup workspace: %w", err)
	}
	defer os.RemoveAll(workDir)

	snapshotPath := filepath.Join(workDir, "database.sqlite")
	db, err := config.OpenDatabase(config.DatabaseConfig{Driver: "sqlite", Path: databasePath})
	if err != nil {
		return err
	}
	if err := db.Exec("VACUUM INTO ?", snapshotPath).Error; err != nil {
		closeDatabase(db)
		return fmt.Errorf("create consistent SQLite snapshot: %w", err)
	}
	if err := closeDatabase(db); err != nil {
		return err
	}

	tempArchive, err := os.CreateTemp(filepath.Dir(archivePath), ".xpense-backup-*.zip")
	if err != nil {
		return fmt.Errorf("create backup archive: %w", err)
	}
	tempPath := tempArchive.Name()
	if err := tempArchive.Close(); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("close temporary backup: %w", err)
	}
	defer os.Remove(tempPath)

	file, err := os.OpenFile(tempPath, os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("open temporary backup: %w", err)
	}
	zipWriter := zip.NewWriter(file)
	writeErr := writeArchive(zipWriter, snapshotPath, filepath.Join(dataDir, "receipts"))
	closeZipErr := zipWriter.Close()
	closeFileErr := file.Close()
	if writeErr != nil {
		return writeErr
	}
	if closeZipErr != nil {
		return fmt.Errorf("finalize backup archive: %w", closeZipErr)
	}
	if closeFileErr != nil {
		return fmt.Errorf("close backup archive: %w", closeFileErr)
	}
	if err := os.Rename(tempPath, archivePath); err != nil {
		return fmt.Errorf("publish backup archive: %w", err)
	}
	return nil
}

// Restore validates the archive before replacing the SQLite database and receipts.
// The caller must stop the application before restoring.
func Restore(archivePath, databasePath, dataDir string) error {
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return fmt.Errorf("create application data directory: %w", err)
	}
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return fmt.Errorf("%w: open archive: %v", ErrInvalidArchive, err)
	}
	defer reader.Close()
	if len(reader.File) == 0 || len(reader.File) > maxEntries {
		return ErrInvalidArchive
	}

	workDir, err := os.MkdirTemp(dataDir, ".restore-*")
	if err != nil {
		return fmt.Errorf("create restore workspace: %w", err)
	}
	defer os.RemoveAll(workDir)

	dbSnapshot := filepath.Join(workDir, "database.sqlite")
	receiptsSnapshot := filepath.Join(workDir, "receipts")
	if err := os.Mkdir(receiptsSnapshot, 0700); err != nil {
		return fmt.Errorf("create receipt restore workspace: %w", err)
	}
	var (
		foundManifest bool
		foundDatabase bool
		totalSize     int64
	)
	for _, entry := range reader.File {
		if entry.FileInfo().IsDir() {
			if entry.Name != "receipts/" {
				return ErrInvalidArchive
			}
			continue
		}
		totalSize += int64(entry.UncompressedSize64)
		if totalSize > maxArchiveSize {
			return fmt.Errorf("%w: archive exceeds size limit", ErrInvalidArchive)
		}
		switch {
		case entry.Name == "manifest.json":
			if foundManifest || entry.UncompressedSize64 > 4096 {
				return ErrInvalidArchive
			}
			foundManifest = true
			if err := readManifest(entry); err != nil {
				return err
			}
		case entry.Name == "database.sqlite":
			if foundDatabase || entry.UncompressedSize64 == 0 {
				return ErrInvalidArchive
			}
			foundDatabase = true
			if err := extract(entry, dbSnapshot, maxArchiveSize); err != nil {
				return err
			}
		case strings.HasPrefix(entry.Name, "receipts/"):
			key := strings.TrimPrefix(entry.Name, "receipts/")
			if !validReceiptKey(key) || entry.UncompressedSize64 > uint64(maxReceiptSize) {
				return ErrInvalidArchive
			}
			if err := extract(entry, filepath.Join(receiptsSnapshot, key), maxReceiptSize); err != nil {
				return err
			}
		default:
			return ErrInvalidArchive
		}
	}
	if !foundManifest || !foundDatabase {
		return ErrInvalidArchive
	}
	if err := verifyDatabase(dbSnapshot); err != nil {
		return fmt.Errorf("%w: database integrity check failed: %v", ErrInvalidArchive, err)
	}
	if err := os.MkdirAll(filepath.Dir(databasePath), 0700); err != nil {
		return fmt.Errorf("create database directory: %w", err)
	}
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return fmt.Errorf("create application data directory: %w", err)
	}

	timestamp := time.Now().UTC().Format("20060102T150405.000000000Z")
	safetyArchive := filepath.Join(dataDir, "pre-restore-"+timestamp+".zip")
	if _, err := os.Stat(databasePath); err == nil {
		if err := Create(databasePath, dataDir, safetyArchive); err != nil {
			return fmt.Errorf("create pre-restore safety backup: %w", err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("inspect existing database: %w", err)
	}

	dbTemp, err := os.CreateTemp(filepath.Dir(databasePath), ".restore-database-*")
	if err != nil {
		return fmt.Errorf("stage restored database: %w", err)
	}
	dbTempPath := dbTemp.Name()
	if _, err := copyFile(dbTemp, dbSnapshot); err != nil {
		dbTemp.Close()
		os.Remove(dbTempPath)
		return fmt.Errorf("copy restored database: %w", err)
	}
	if err := dbTemp.Close(); err != nil {
		os.Remove(dbTempPath)
		return fmt.Errorf("close restored database: %w", err)
	}
	defer os.Remove(dbTempPath)

	receiptPath := filepath.Join(dataDir, "receipts")
	oldReceipts := filepath.Join(dataDir, "receipts.pre-restore-"+timestamp)
	hadReceipts := false
	if _, err := os.Stat(receiptPath); err == nil {
		if err := os.Rename(receiptPath, oldReceipts); err != nil {
			return fmt.Errorf("preserve existing receipts: %w", err)
		}
		hadReceipts = true
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("inspect existing receipts: %w", err)
	}
	if err := os.Rename(receiptsSnapshot, receiptPath); err != nil {
		if hadReceipts {
			_ = os.Rename(oldReceipts, receiptPath)
		}
		return fmt.Errorf("install restored receipts: %w", err)
	}

	movedSidecars := make([]string, 0, 2)
	rollbackSidecars := func() {
		for i := len(movedSidecars) - 1; i >= 0; i-- {
			_ = os.Rename(movedSidecars[i]+".pre-restore-"+timestamp, movedSidecars[i])
		}
	}
	for _, suffix := range []string{"-wal", "-shm"} {
		sidecar := databasePath + suffix
		if _, err := os.Stat(sidecar); err == nil {
			if err := os.Rename(sidecar, sidecar+".pre-restore-"+timestamp); err != nil {
				rollbackSidecars()
				_ = os.RemoveAll(receiptPath)
				if hadReceipts {
					_ = os.Rename(oldReceipts, receiptPath)
				}
				return fmt.Errorf("preserve SQLite sidecar %s: %w", suffix, err)
			}
			movedSidecars = append(movedSidecars, sidecar)
		} else if !errors.Is(err, os.ErrNotExist) {
			rollbackSidecars()
			_ = os.RemoveAll(receiptPath)
			if hadReceipts {
				_ = os.Rename(oldReceipts, receiptPath)
			}
			return fmt.Errorf("inspect SQLite sidecar %s: %w", suffix, err)
		}
	}
	if err := os.Rename(dbTempPath, databasePath); err != nil {
		rollbackSidecars()
		_ = os.RemoveAll(receiptPath)
		if hadReceipts {
			_ = os.Rename(oldReceipts, receiptPath)
		}
		return fmt.Errorf("install restored database: %w", err)
	}
	return nil
}

func writeArchive(writer *zip.Writer, snapshotPath, receiptsPath string) error {
	manifestData, err := json.Marshal(manifest{Version: archiveVersion, CreatedAt: time.Now().UTC()})
	if err != nil {
		return fmt.Errorf("encode backup manifest: %w", err)
	}
	manifestFile, err := writer.Create("manifest.json")
	if err != nil {
		return fmt.Errorf("add backup manifest: %w", err)
	}
	if _, err := manifestFile.Write(manifestData); err != nil {
		return fmt.Errorf("write backup manifest: %w", err)
	}
	if err := addFile(writer, "database.sqlite", snapshotPath); err != nil {
		return err
	}
	info, err := os.Stat(receiptsPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("inspect receipt directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("receipt storage is not a directory")
	}
	entries, err := os.ReadDir(receiptsPath)
	if err != nil {
		return fmt.Errorf("read receipt directory: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() || entry.Type()&os.ModeSymlink != 0 || !validReceiptKey(entry.Name()) {
			return fmt.Errorf("unexpected file in receipt directory")
		}
		if err := addFile(writer, "receipts/"+entry.Name(), filepath.Join(receiptsPath, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}

func addFile(writer *zip.Writer, name, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %s for backup: %w", name, err)
	}
	defer file.Close()
	entry, err := writer.Create(name)
	if err != nil {
		return fmt.Errorf("add %s to backup: %w", name, err)
	}
	if _, err := io.Copy(entry, file); err != nil {
		return fmt.Errorf("write %s to backup: %w", name, err)
	}
	return nil
}

func readManifest(entry *zip.File) error {
	file, err := entry.Open()
	if err != nil {
		return fmt.Errorf("%w: read manifest: %v", ErrInvalidArchive, err)
	}
	defer file.Close()
	var value manifest
	if err := json.NewDecoder(io.LimitReader(file, 4096)).Decode(&value); err != nil {
		return fmt.Errorf("%w: invalid manifest: %v", ErrInvalidArchive, err)
	}
	if value.Version != archiveVersion || value.CreatedAt.IsZero() {
		return ErrInvalidArchive
	}
	return nil
}

func extract(entry *zip.File, destination string, limit int64) error {
	source, err := entry.Open()
	if err != nil {
		return fmt.Errorf("%w: open archive entry: %v", ErrInvalidArchive, err)
	}
	defer source.Close()
	out, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return fmt.Errorf("%w: create extracted file: %v", ErrInvalidArchive, err)
	}
	written, copyErr := io.Copy(out, io.LimitReader(source, limit+1))
	closeErr := out.Close()
	if copyErr != nil {
		return fmt.Errorf("%w: extract archive entry: %v", ErrInvalidArchive, copyErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close extracted file: %w", closeErr)
	}
	if written > limit || uint64(written) != entry.UncompressedSize64 {
		return ErrInvalidArchive
	}
	return nil
}

func verifyDatabase(path string) error {
	db, err := config.OpenDatabase(config.DatabaseConfig{Driver: "sqlite", Path: path})
	if err != nil {
		return err
	}
	var result string
	if err := db.Raw("PRAGMA integrity_check").Scan(&result).Error; err != nil {
		closeDatabase(db)
		return err
	}
	if result != "ok" {
		closeDatabase(db)
		return fmt.Errorf("integrity_check returned %q", result)
	}
	return closeDatabase(db)
}

func closeDatabase(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func copyFile(destination *os.File, sourcePath string) (int64, error) {
	source, err := os.Open(sourcePath)
	if err != nil {
		return 0, err
	}
	defer source.Close()
	return io.Copy(destination, source)
}

func validReceiptKey(key string) bool {
	if len(key) != 64 {
		return false
	}
	for _, character := range key {
		if (character < '0' || character > '9') && (character < 'a' || character > 'f') {
			return false
		}
	}
	return true
}
