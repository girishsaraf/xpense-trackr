package trends

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/girishsaraf/xpense-trackr/backend/internal/category"
	"github.com/girishsaraf/xpense-trackr/backend/internal/expense"
	"github.com/girishsaraf/xpense-trackr/backend/internal/testutil"
)

func TestSpendingHandlerAggregatesExpensesByMonthAndCategory(t *testing.T) {
	database := testutil.NewDB(t, &category.Category{}, &expense.Expense{})
	SetDB(database)
	food := category.Category{Name: "Food"}
	travel := category.Category{Name: "Travel"}
	if err := database.Create(&food).Error; err != nil {
		t.Fatal(err)
	}
	if err := database.Create(&travel).Error; err != nil {
		t.Fatal(err)
	}
	thisMonth := time.Now().Format("2006-01") + "-05"
	lastMonth := time.Now().AddDate(0, -1, 0).Format("2006-01") + "-08"
	seed := []expense.Expense{
		{CategoryID: food.ID, Amount: 10, Date: thisMonth},
		{CategoryID: food.ID, Amount: 15, Date: thisMonth},
		{CategoryID: travel.ID, Amount: 8, Date: lastMonth},
	}
	if err := database.Create(&seed).Error; err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/trends/spending?months=3", nil)
	response := httptest.NewRecorder()
	SpendingHandler(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %s", response.Code, http.StatusOK, response.Body)
	}

	var points []SpendingPoint
	if err := json.NewDecoder(response.Body).Decode(&points); err != nil {
		t.Fatal(err)
	}
	if len(points) != 2 {
		t.Fatalf("got %d trend points, want 2: %+v", len(points), points)
	}
	if points[0].Month != time.Now().AddDate(0, -1, 0).Format("2006-01") || points[0].Category != "Travel" || points[0].Amount != 8 {
		t.Errorf("last month point = %+v", points[0])
	}
	if points[1].Month != time.Now().Format("2006-01") || points[1].Category != "Food" || points[1].Amount != 25 {
		t.Errorf("current month point = %+v", points[1])
	}
}

func TestSpendingHandlerRejectsInvalidMonthRange(t *testing.T) {
	database := testutil.NewDB(t, &category.Category{}, &expense.Expense{})
	SetDB(database)
	for _, months := range []string{"0", "25", "invalid"} {
		t.Run(fmt.Sprintf("months_%s", months), func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/api/trends/spending?months="+months, nil)
			response := httptest.NewRecorder()
			SpendingHandler(response, request)
			if response.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want %d", response.Code, http.StatusBadRequest)
			}
		})
	}
}
