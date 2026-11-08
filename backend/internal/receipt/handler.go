package receipt

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetReceiptsHandler(w http.ResponseWriter, r *http.Request) {
	var transactionID uint
	if raw := r.URL.Query().Get("transaction_id"); raw != "" {
		value, err := strconv.ParseUint(raw, 10, 64)
		if err != nil || value == 0 {
			http.Error(w, "Invalid transaction_id", http.StatusBadRequest)
			return
		}
		transactionID = uint(value)
	}
	receipts, err := GetReceipts(transactionID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, receipts)
}

func CreateReceiptHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxFileSize+(1<<20))
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		http.Error(w, "Upload must be a multipart form with a file no larger than 10 MB", http.StatusBadRequest)
		return
	}
	transactionIDValue, err := strconv.ParseUint(r.FormValue("transaction_id"), 10, 64)
	if err != nil || transactionIDValue == 0 {
		http.Error(w, "A valid transaction_id is required", http.StatusBadRequest)
		return
	}
	if err := validateTransaction(uint(transactionIDValue)); err != nil {
		writeError(w, err)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "A receipt file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(io.LimitReader(file, MaxFileSize+1))
	if err != nil {
		http.Error(w, "Unable to read receipt file", http.StatusBadRequest)
		return
	}
	if len(content) == 0 || len(content) > MaxFileSize {
		http.Error(w, "Receipt file must be between 1 byte and 10 MB", http.StatusBadRequest)
		return
	}
	mediaType := http.DetectContentType(content[:min(len(content), 512)])
	if mediaType != "application/pdf" && mediaType != "image/png" && mediaType != "image/jpeg" {
		http.Error(w, "Only PDF, PNG, and JPEG receipts are supported", http.StatusBadRequest)
		return
	}
	fileName := filepath.Base(strings.TrimSpace(header.Filename))
	if fileName == "." || fileName == string(filepath.Separator) || len(fileName) > 255 {
		http.Error(w, "Invalid receipt filename", http.StatusBadRequest)
		return
	}
	key, err := randomStorageKey()
	if err != nil {
		writeError(w, err)
		return
	}
	path, err := filePath(key)
	if err != nil {
		writeError(w, err)
		return
	}
	if err := os.MkdirAll(storageDir, 0700); err != nil {
		writeError(w, fmt.Errorf("create receipt storage: %w", err))
		return
	}
	if err := os.WriteFile(path, content, 0600); err != nil {
		writeError(w, fmt.Errorf("store receipt file: %w", err))
		return
	}

	item := Receipt{
		TransactionID: uint(transactionIDValue), FileName: fileName,
		MediaType: mediaType, SizeBytes: int64(len(content)), StorageKey: key,
	}
	if err := db.Create(&item).Error; err != nil {
		_ = os.Remove(path)
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func GetReceiptFileHandler(w http.ResponseWriter, r *http.Request) {
	item, ok := receiptByID(w, r)
	if !ok {
		return
	}
	path, err := filePath(item.StorageKey)
	if err != nil {
		writeError(w, err)
		return
	}
	file, err := os.Open(path)
	if err != nil {
		writeError(w, fmt.Errorf("open receipt file: %w", err))
		return
	}
	defer file.Close()
	disposition := mime.FormatMediaType("attachment", map[string]string{"filename": item.FileName})
	w.Header().Set("Content-Type", item.MediaType)
	w.Header().Set("Content-Disposition", disposition)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	http.ServeContent(w, r, item.FileName, time.Unix(item.CreatedAt, 0), file)
}

func DeleteReceiptHandler(w http.ResponseWriter, r *http.Request) {
	item, ok := receiptByID(w, r)
	if !ok {
		return
	}
	if err := DeleteReceipt(item.ID); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func receiptByID(w http.ResponseWriter, r *http.Request) (*Receipt, bool) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil || id == 0 {
		http.Error(w, "Invalid receipt ID", http.StatusBadRequest)
		return nil, false
	}
	item, err := GetReceipt(uint(id))
	if err != nil {
		writeError(w, err)
		return nil, false
	}
	return item, true
}

func randomStorageKey() (string, error) {
	value := make([]byte, 32)
	if _, err := rand.Read(value); err != nil {
		return "", fmt.Errorf("generate receipt storage key: %w", err)
	}
	return hex.EncodeToString(value), nil
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		http.Error(w, "Transaction or receipt not found", http.StatusNotFound)
	case errors.Is(err, ErrInvalidReceipt):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		log.Printf("receipt request failed: %v", err)
		http.Error(w, "Unable to complete receipt request", http.StatusInternalServerError)
	}
}

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("encode receipt response: %v", err)
	}
}
