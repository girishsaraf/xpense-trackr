package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func userRequest(t *testing.T, handler http.HandlerFunc, body string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	return recorder
}

func TestLogin(t *testing.T) {
	recorder := userRequest(t, Login, `{"email":"user@example.com","password":"password"}`)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var response map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if _, ok := response["token"]; !ok {
		t.Errorf("response = %+v, want a token field", response)
	}
}

func TestLoginInvalidBody(t *testing.T) {
	recorder := userRequest(t, Login, `{`)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestRegister(t *testing.T) {
	recorder := userRequest(t, Register, `{"username":"user","email":"user@example.com","password":"password"}`)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusCreated)
	}
	if body := recorder.Body.String(); body != "User account created successfully" {
		t.Errorf("body = %q", body)
	}
}

func TestRegisterInvalidBody(t *testing.T) {
	recorder := userRequest(t, Register, `{`)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}
