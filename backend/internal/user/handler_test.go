package user

import (
	"bytes"
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

func TestLoginReportsAuthenticationUnavailable(t *testing.T) {
	recorder := userRequest(t, Login, `{"email":"user@example.com","password":"password"}`)

	if recorder.Code != http.StatusNotImplemented {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusNotImplemented)
	}
}

func TestLoginInvalidBody(t *testing.T) {
	recorder := userRequest(t, Login, `{`)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestRegisterReportsAuthenticationUnavailable(t *testing.T) {
	recorder := userRequest(t, Register, `{"username":"user","email":"user@example.com","password":"password"}`)

	if recorder.Code != http.StatusNotImplemented {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusNotImplemented)
	}
}

func TestRegisterInvalidBody(t *testing.T) {
	recorder := userRequest(t, Register, `{`)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}
