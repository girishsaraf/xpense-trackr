package recurring

import (
	"net/http"
	"testing"

	"github.com/girishsaraf/xpense-trackr/backend/internal/expense"
	"github.com/girishsaraf/xpense-trackr/backend/internal/investments"
	"github.com/girishsaraf/xpense-trackr/backend/internal/testutil"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func newRecurringTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	database := testutil.NewDB(t,
		&RecurringExpense{},
		&RecurringInvestment{},
		&expense.Expense{},
		&investments.Investment{},
	)
	SetDB(database)
	return database
}

func newRecurringTestRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/recurring/expenses", GetAllRecurringExpensesHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/recurring/expenses", CreateRecurringExpenseHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/recurring/expenses/{id}", GetRecurringExpenseByIDHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/recurring/expenses/{id}", UpdateRecurringExpenseHandler).Methods(http.MethodPut)
	router.HandleFunc("/api/recurring/expenses/{id}", DeleteRecurringExpenseHandler).Methods(http.MethodDelete)

	router.HandleFunc("/api/recurring/investments", GetAllRecurringInvestmentsHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/recurring/investments", CreateRecurringInvestmentHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/recurring/investments/{id}", GetRecurringInvestmentByIDHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/recurring/investments/{id}", UpdateRecurringInvestmentHandler).Methods(http.MethodPut)
	router.HandleFunc("/api/recurring/investments/{id}", DeleteRecurringInvestmentHandler).Methods(http.MethodDelete)
	return router
}
