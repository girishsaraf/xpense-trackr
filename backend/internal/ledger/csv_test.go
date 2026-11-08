package ledger

import (
	"errors"
	"fmt"
	"testing"
)

func TestParseAmountMinorUsesCurrencyPrecision(t *testing.T) {
	tests := []struct {
		value    string
		currency string
		want     int64
		wantErr  bool
	}{
		{value: "12.34", currency: "USD", want: 1234},
		{value: "100", currency: "JPY", want: 100},
		{value: "0.001", currency: "KWD", want: 1},
		{value: "1.234", currency: "USD", wantErr: true},
		{value: "1.0", currency: "JPY", wantErr: true},
		{value: "0.00", currency: "USD", wantErr: true},
	}
	for _, test := range tests {
		got, err := ParseAmountMinor(test.value, test.currency)
		if test.wantErr {
			if err == nil {
				t.Errorf("ParseAmountMinor(%q, %q) error = nil, want error", test.value, test.currency)
			}
			continue
		}
		if err != nil || got != test.want {
			t.Errorf("ParseAmountMinor(%q, %q) = %d, %v; want %d, nil", test.value, test.currency, got, err, test.want)
		}
	}
}

func TestCSVImportPreviewCommitAndRepeatAreIdempotent(t *testing.T) {
	database := setupLedgerDB(t)
	account := &Account{Name: "Checking", Type: "checking", Currency: "USD"}
	if err := CreateAccount(account); err != nil {
		t.Fatalf("CreateAccount() error = %v", err)
	}
	content := []byte(fmt.Sprintf(
		"account_id,type,amount,currency,date,payee,description,category_id,transfer_account_id\n"+
			"%d,expense,12.34,USD,2026-10-12,Market,Groceries,,\n"+
			"%d,expense,12.34,USD,2026-10-12,Market,Groceries,,\n",
		account.ID, account.ID,
	))

	preview, err := ImportTransactionsCSV(content, false)
	if err != nil {
		t.Fatalf("ImportTransactionsCSV(preview) error = %v", err)
	}
	if preview.Total != 2 || preview.Ready != 1 || preview.Duplicates != 1 || preview.Imported != 0 {
		t.Errorf("preview report = %+v", preview)
	}
	var count int64
	if err := database.Model(&Transaction{}).Count(&count).Error; err != nil || count != 0 {
		t.Fatalf("preview transaction count = %d, %v; want 0", count, err)
	}

	committed, err := ImportTransactionsCSV(content, true)
	if err != nil {
		t.Fatalf("ImportTransactionsCSV(commit) error = %v", err)
	}
	if committed.Imported != 1 || committed.Ready != 0 || committed.Duplicates != 1 {
		t.Errorf("commit report = %+v", committed)
	}
	repeated, err := ImportTransactionsCSV(content, true)
	if err != nil {
		t.Fatalf("ImportTransactionsCSV(repeat) error = %v", err)
	}
	if repeated.Imported != 0 || repeated.Duplicates != 2 {
		t.Errorf("repeat report = %+v, want both rows detected as duplicates", repeated)
	}
}

func TestCSVImportInvalidRowPreventsBatchCommit(t *testing.T) {
	database := setupLedgerDB(t)
	account := &Account{Name: "Checking", Type: "checking", Currency: "USD"}
	if err := CreateAccount(account); err != nil {
		t.Fatalf("CreateAccount() error = %v", err)
	}
	content := []byte(fmt.Sprintf(
		"account_id,type,amount,currency,date,payee\n"+
			"%d,expense,10.00,USD,2026-10-12,Market\n"+
			"%d,expense,5.00,USD,not-a-date,Shop\n",
		account.ID, account.ID,
	))
	report, err := ImportTransactionsCSV(content, true)
	if err != nil {
		t.Fatalf("ImportTransactionsCSV() error = %v", err)
	}
	if report.Invalid != 1 || report.Ready != 1 || report.Imported != 0 {
		t.Errorf("invalid batch report = %+v", report)
	}
	var count int64
	if err := database.Model(&Transaction{}).Count(&count).Error; err != nil || count != 0 {
		t.Fatalf("invalid batch transaction count = %d, %v; want 0", count, err)
	}
	if _, err := ImportTransactionsCSV([]byte("date,type\n"), false); !errors.Is(err, ErrInvalidCSV) {
		t.Errorf("invalid header error = %v, want %v", err, ErrInvalidCSV)
	}
}
