package category

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/girishsaraf/xpense-trackr/backend/internal/testutil"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func newCategoryHandlerDB(t *testing.T) *gorm.DB {
	t.Helper()
	database := testutil.NewDB(t, &Category{})
	SetDB(database)
	return database
}

func newCategoryTestRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/categories", GetAllCategoriesHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/categories", CreateCategoryHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/categories/{id}", GetCategoryByIDHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/categories/{id}", UpdateCategoryHandler).Methods(http.MethodPut)
	router.HandleFunc("/api/categories/{id}", DeleteCategoryHandler).Methods(http.MethodDelete)
	return router
}

func categoryRequest(t *testing.T, method, path string, body []byte) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(method, path, bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	newCategoryTestRouter().ServeHTTP(recorder, request)
	return recorder
}

func TestGetAllCategoriesHandler(t *testing.T) {
	database := newCategoryHandlerDB(t)
	if err := database.Create(&Category{Name: "Food"}).Error; err != nil {
		t.Fatalf("failed to seed category: %v", err)
	}

	recorder := categoryRequest(t, http.MethodGet, "/api/categories", nil)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var categories []Category
	if err := json.NewDecoder(recorder.Body).Decode(&categories); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(categories) != 1 || categories[0].Name != "Food" {
		t.Errorf("categories = %+v", categories)
	}
}

func TestCreateCategoryHandler(t *testing.T) {
	newCategoryHandlerDB(t)

	recorder := categoryRequest(t, http.MethodPost, "/api/categories", []byte(`{"name":"Utilities"}`))

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body = %s", recorder.Code, http.StatusCreated, recorder.Body.String())
	}
	var created Category
	if err := json.NewDecoder(recorder.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if created.ID == 0 || created.Name != "Utilities" {
		t.Errorf("created category = %+v", created)
	}
}

func TestCreateCategoryHandlerInvalidBody(t *testing.T) {
	newCategoryHandlerDB(t)

	recorder := categoryRequest(t, http.MethodPost, "/api/categories", []byte(`{`))

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestGetCategoryByIDHandler(t *testing.T) {
	database := newCategoryHandlerDB(t)
	seeded := &Category{Name: "Health"}
	if err := database.Create(seeded).Error; err != nil {
		t.Fatalf("failed to seed category: %v", err)
	}

	recorder := categoryRequest(t, http.MethodGet, fmt.Sprintf("/api/categories/%d", seeded.ID), nil)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var got Category
	if err := json.NewDecoder(recorder.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got.Name != seeded.Name {
		t.Errorf("response = %+v, want %+v", got, seeded)
	}
}

func TestGetCategoryByIDHandlerInvalidID(t *testing.T) {
	newCategoryHandlerDB(t)
	request := mux.SetURLVars(httptest.NewRequest(http.MethodGet, "/api/categories/abc", nil), map[string]string{"id": "abc"})
	recorder := httptest.NewRecorder()

	GetCategoryByIDHandler(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestGetCategoryByIDHandlerNotFound(t *testing.T) {
	newCategoryHandlerDB(t)

	recorder := categoryRequest(t, http.MethodGet, "/api/categories/999", nil)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusNotFound)
	}
}

func TestUpdateCategoryHandler(t *testing.T) {
	database := newCategoryHandlerDB(t)
	seeded := &Category{Name: "Old"}
	if err := database.Create(seeded).Error; err != nil {
		t.Fatalf("failed to seed category: %v", err)
	}

	recorder := categoryRequest(t, http.MethodPut, fmt.Sprintf("/api/categories/%d", seeded.ID), []byte(`{"name":"New"}`))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var updated Category
	if err := json.NewDecoder(recorder.Body).Decode(&updated); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if updated.Name != "New" {
		t.Errorf("updated category = %+v", updated)
	}
}

func TestUpdateCategoryHandlerNotFound(t *testing.T) {
	newCategoryHandlerDB(t)

	recorder := categoryRequest(t, http.MethodPut, "/api/categories/999", []byte(`{"name":"Missing"}`))

	if recorder.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusNotFound)
	}
}

func TestDeleteCategoryHandler(t *testing.T) {
	database := newCategoryHandlerDB(t)
	seeded := &Category{Name: "Delete"}
	if err := database.Create(seeded).Error; err != nil {
		t.Fatalf("failed to seed category: %v", err)
	}

	recorder := categoryRequest(t, http.MethodDelete, fmt.Sprintf("/api/categories/%d", seeded.ID), nil)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNoContent)
	}
}

func TestDeleteCategoryHandlerNotFound(t *testing.T) {
	newCategoryHandlerDB(t)

	recorder := categoryRequest(t, http.MethodDelete, "/api/categories/999", nil)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusNotFound)
	}
}
