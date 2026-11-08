package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestWithFrontendMapsVersionedAPIPath(t *testing.T) {
	api := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/categories" {
			t.Errorf("API path = %q, want %q", r.URL.Path, "/api/categories")
		}
		w.WriteHeader(http.StatusNoContent)
	})
	handler := withFrontend(api, t.TempDir())
	request := httptest.NewRequest(http.MethodGet, "/api/v1/categories", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", response.Code, http.StatusNoContent)
	}
}

func TestWithFrontendServesClientRoute(t *testing.T) {
	webDir := t.TempDir()
	index := []byte("<html>app shell</html>")
	if err := os.WriteFile(filepath.Join(webDir, "index.html"), index, 0600); err != nil {
		t.Fatalf("write index: %v", err)
	}
	handler := withFrontend(http.NotFoundHandler(), webDir)
	request := httptest.NewRequest(http.MethodGet, "/budgets", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if response.Body.String() != string(index) {
		t.Errorf("body = %q, want %q", response.Body.String(), index)
	}
}
