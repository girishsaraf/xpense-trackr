package internal

import (
	"github.com/girishsaraf/xpense-trackr/backend/internal/budget"
	"github.com/girishsaraf/xpense-trackr/backend/internal/category"
	"github.com/girishsaraf/xpense-trackr/backend/internal/expense"
	"github.com/girishsaraf/xpense-trackr/backend/internal/investments"
	"github.com/girishsaraf/xpense-trackr/backend/internal/recurring"
	"github.com/girishsaraf/xpense-trackr/backend/internal/user"
	"github.com/gorilla/mux"
)

// SetRoutes initializes the routes for the API endpoints.
func SetRoutes(router *mux.Router) {
	// User Authentication
	authRouter := router.PathPrefix("/api/auth").Subrouter()
	authRouter.HandleFunc("/login", user.Login).Methods("POST")
	authRouter.HandleFunc("/register", user.Register).Methods("POST")

	// Expense Management
	expenseRouter := router.PathPrefix("/api/expenses").Subrouter()
	expenseRouter.HandleFunc("", expense.GetAllExpensesHandler).Methods("GET")
	expenseRouter.HandleFunc("", expense.CreateExpenseHandler).Methods("POST")
	expenseRouter.HandleFunc("/{id:[0-9]+}", expense.GetExpenseByIDHandler).Methods("GET")
	expenseRouter.HandleFunc("/{id:[0-9]+}", expense.UpdateExpenseHandler).Methods("PUT")
	expenseRouter.HandleFunc("/{id:[0-9]+}", expense.DeleteExpenseHandler).Methods("DELETE")
	expenseRouter.HandleFunc("/cumulative", expense.CumulativeExpensesByMonthHandler).Methods("GET")

	// Category Management
	categoryRouter := router.PathPrefix("/api/categories").Subrouter()
	categoryRouter.HandleFunc("", category.GetAllCategoriesHandler).Methods("GET")
	categoryRouter.HandleFunc("", category.CreateCategoryHandler).Methods("POST")
	categoryRouter.HandleFunc("/{id}", category.GetCategoryByIDHandler).Methods("GET")
	categoryRouter.HandleFunc("/{id}", category.UpdateCategoryHandler).Methods("PUT")
	categoryRouter.HandleFunc("/{id}", category.DeleteCategoryHandler).Methods("DELETE")

	// Budget Management
	budgetRouter := router.PathPrefix("/api/budgets").Subrouter()
	budgetRouter.HandleFunc("", budget.GetAllBudgetsHandler).Methods("GET")
	budgetRouter.HandleFunc("", budget.CreateBudgetHandler).Methods("POST")
	budgetRouter.HandleFunc("/{id}", budget.GetBudgetByIDHandler).Methods("GET")
	budgetRouter.HandleFunc("/{id}", budget.UpdateBudgetHandler).Methods("PUT")
	budgetRouter.HandleFunc("/{id}", budget.DeleteBudgetHandler).Methods("DELETE")

	// Investment Management
	investmentRouter := router.PathPrefix("/api/investments").Subrouter()
	investmentRouter.HandleFunc("", investments.GetAllInvestmentsHandler).Methods("GET")
	investmentRouter.HandleFunc("", investments.CreateInvestmentHandler).Methods("POST")
	investmentRouter.HandleFunc("/{id:[0-9]+}", investments.GetInvestmentByIDHandler).Methods("GET")
	investmentRouter.HandleFunc("/{id:[0-9]+}", investments.UpdateInvestmentHandler).Methods("PUT")
	investmentRouter.HandleFunc("/{id:[0-9]+}", investments.DeleteInvestmentHandler).Methods("DELETE")
	investmentRouter.HandleFunc("/cumulative", investments.CumulativeInvestmentsByMonthHandler).Methods("GET")

	// Recurring Entry Management
	recurringRouter := router.PathPrefix("/api/recurring").Subrouter()
	recurringRouter.HandleFunc("/expenses", recurring.GetAllRecurringExpensesHandler).Methods("GET")
	recurringRouter.HandleFunc("/expenses", recurring.CreateRecurringExpenseHandler).Methods("POST")
	recurringRouter.HandleFunc("/expenses/{id:[0-9]+}", recurring.GetRecurringExpenseByIDHandler).Methods("GET")
	recurringRouter.HandleFunc("/expenses/{id:[0-9]+}", recurring.UpdateRecurringExpenseHandler).Methods("PUT")
	recurringRouter.HandleFunc("/expenses/{id:[0-9]+}", recurring.DeleteRecurringExpenseHandler).Methods("DELETE")

	recurringRouter.HandleFunc("/investments", recurring.GetAllRecurringInvestmentsHandler).Methods("GET")
	recurringRouter.HandleFunc("/investments", recurring.CreateRecurringInvestmentHandler).Methods("POST")
	recurringRouter.HandleFunc("/investments/{id:[0-9]+}", recurring.GetRecurringInvestmentByIDHandler).Methods("GET")
	recurringRouter.HandleFunc("/investments/{id:[0-9]+}", recurring.UpdateRecurringInvestmentHandler).Methods("PUT")
	recurringRouter.HandleFunc("/investments/{id:[0-9]+}", recurring.DeleteRecurringInvestmentHandler).Methods("DELETE")
}
