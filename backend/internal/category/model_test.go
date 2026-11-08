package category

import (
	"errors"
	"testing"

	"github.com/girishsaraf/xpense-trackr/backend/internal/testutil"
	"gorm.io/gorm"
)

func newCategoryModelDB(t *testing.T) *gorm.DB {
	t.Helper()
	database := testutil.NewDB(t, &Category{})
	SetDB(database)
	return database
}

func TestCategoryCRUD(t *testing.T) {
	newCategoryModelDB(t)

	first := &Category{Name: "Food"}
	second := &Category{Name: "Travel"}
	if err := CreateCategory(first); err != nil {
		t.Fatalf("CreateCategory() error = %v", err)
	}
	if err := CreateCategory(second); err != nil {
		t.Fatalf("CreateCategory() error = %v", err)
	}

	categories, err := GetAllCategories()
	if err != nil {
		t.Fatalf("GetAllCategories() error = %v", err)
	}
	if len(categories) != 2 {
		t.Fatalf("GetAllCategories() returned %d categories, want 2", len(categories))
	}

	got, err := GetCategoryByID(int(first.ID))
	if err != nil {
		t.Fatalf("GetCategoryByID() error = %v", err)
	}
	if got.Name != first.Name {
		t.Errorf("GetCategoryByID() = %+v, want %+v", got, first)
	}

	if err := UpdateCategory(int(first.ID), &Category{Name: "Groceries"}); err != nil {
		t.Fatalf("UpdateCategory() error = %v", err)
	}
	got, err = GetCategoryByID(int(first.ID))
	if err != nil {
		t.Fatalf("GetCategoryByID() after update error = %v", err)
	}
	if got.Name != "Groceries" {
		t.Errorf("updated category name = %q, want %q", got.Name, "Groceries")
	}

	if err := DeleteCategory(int(second.ID)); err != nil {
		t.Fatalf("DeleteCategory() error = %v", err)
	}
	if _, err := GetCategoryByID(int(second.ID)); !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("GetCategoryByID() after delete error = %v, want ErrRecordNotFound", err)
	}
}

func TestGetCategoryByIDNotFound(t *testing.T) {
	newCategoryModelDB(t)

	if _, err := GetCategoryByID(404); !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("GetCategoryByID() error = %v, want ErrRecordNotFound", err)
	}
}

func TestUpdateCategoryNotFound(t *testing.T) {
	newCategoryModelDB(t)

	if err := UpdateCategory(404, &Category{Name: "Missing"}); !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("UpdateCategory() error = %v, want ErrRecordNotFound", err)
	}
}

func TestDeleteCategoryNotFound(t *testing.T) {
	newCategoryModelDB(t)

	if err := DeleteCategory(404); !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("DeleteCategory() error = %v, want ErrRecordNotFound", err)
	}
}
