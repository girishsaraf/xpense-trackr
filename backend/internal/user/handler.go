// internal/auth/handler.go

package user

import (
	"encoding/json"
	"net/http"

	"github.com/girishsaraf/xpense-trackr/backend/internal/auth"
)

// LoginRequest represents the request body for user login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest represents the request body for user registration.
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login handles the user login request.
func Login(w http.ResponseWriter, r *http.Request) {
	// Decode request body
	var loginRequest LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Authenticate user
	token, err := auth.AuthenticateUser(loginRequest.Email, loginRequest.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Respond with token
	err = json.NewEncoder(w).Encode(map[string]string{"token": token})
	if err != nil {
		return
	}
}

// Register handles the user registration request.
func Register(w http.ResponseWriter, r *http.Request) {
	// Decode request body
	var registerRequest RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Create new user account
	if err := auth.RegisterUser(registerRequest.Username, registerRequest.Email, registerRequest.Password); err != nil {
		http.Error(w, "Failed to create user account", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusCreated)
	_, err := w.Write([]byte("User account created successfully"))
	if err != nil {
		return
	}
}
