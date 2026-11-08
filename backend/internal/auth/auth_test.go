package auth

import "testing"

func TestGenerateJWTToken(t *testing.T) {
	if token := GenerateJWTToken("user@example.com"); token == "" {
		t.Error("GenerateJWTToken() returned an empty token")
	}
}

func TestAuthenticateUser(t *testing.T) {
	token, err := AuthenticateUser("user@example.com", "password")
	if err != nil {
		t.Errorf("AuthenticateUser() error = %v", err)
	}
	if token != "dummy_jwt_token" && token != "" {
		t.Errorf("AuthenticateUser() token = %q", token)
	}
}

func TestRegisterUser(t *testing.T) {
	if err := RegisterUser("user", "user@example.com", "password"); err != nil {
		t.Errorf("RegisterUser() error = %v", err)
	}
}
