package auth

import (
	"errors"
	"testing"
)

func TestAuthenticateUserReportsNotImplemented(t *testing.T) {
	if _, err := AuthenticateUser("user@example.com", "password"); !errors.Is(err, ErrNotImplemented) {
		t.Errorf("AuthenticateUser() error = %v, want %v", err, ErrNotImplemented)
	}
}

func TestRegisterUserReportsNotImplemented(t *testing.T) {
	if err := RegisterUser("user", "user@example.com", "password"); !errors.Is(err, ErrNotImplemented) {
		t.Errorf("RegisterUser() error = %v, want %v", err, ErrNotImplemented)
	}
}
