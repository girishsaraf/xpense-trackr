package auth

import "errors"

var ErrNotImplemented = errors.New("local account authentication is not implemented")

// AuthenticateUser reports that account authentication is not yet available.
func AuthenticateUser(email, password string) (string, error) {
	return "", ErrNotImplemented
}

// RegisterUser reports that account registration is not yet available.
func RegisterUser(username, email, password string) error {
	return ErrNotImplemented
}
