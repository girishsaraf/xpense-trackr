package internal

import (
	"github.com/girishsaraf/xpense-trackr/backend/internal/ledger"
	"github.com/gorilla/mux"
)

func registerLedgerRoutes(router *mux.Router) {
	accountRouter := router.PathPrefix("/api/accounts").Subrouter()
	accountRouter.HandleFunc("", ledger.GetAccountsHandler).Methods("GET")
	accountRouter.HandleFunc("", ledger.CreateAccountHandler).Methods("POST")

	transactionRouter := router.PathPrefix("/api/transactions").Subrouter()
	transactionRouter.HandleFunc("", ledger.GetTransactionsHandler).Methods("GET")
	transactionRouter.HandleFunc("", ledger.CreateTransactionHandler).Methods("POST")
	transactionRouter.HandleFunc("/export.csv", ledger.ExportTransactionsCSVHandler).Methods("GET")
	transactionRouter.HandleFunc("/import.csv", ledger.ImportTransactionsCSVHandler).Methods("POST")
	transactionRouter.HandleFunc("/{id:[0-9]+}", ledger.GetTransactionHandler).Methods("GET")
	transactionRouter.HandleFunc("/{id:[0-9]+}", ledger.UpdateTransactionHandler).Methods("PUT")
	transactionRouter.HandleFunc("/{id:[0-9]+}", ledger.DeleteTransactionHandler).Methods("DELETE")

	categoryRuleRouter := router.PathPrefix("/api/category-rules").Subrouter()
	categoryRuleRouter.HandleFunc("", ledger.GetCategoryRulesHandler).Methods("GET")
	categoryRuleRouter.HandleFunc("", ledger.CreateCategoryRuleHandler).Methods("POST")
	categoryRuleRouter.HandleFunc("/{id:[0-9]+}", ledger.UpdateCategoryRuleHandler).Methods("PUT")
	categoryRuleRouter.HandleFunc("/{id:[0-9]+}", ledger.DeleteCategoryRuleHandler).Methods("DELETE")

	scheduledFlowRouter := router.PathPrefix("/api/scheduled-flows").Subrouter()
	scheduledFlowRouter.HandleFunc("", ledger.GetScheduledFlowsHandler).Methods("GET")
	scheduledFlowRouter.HandleFunc("", ledger.CreateScheduledFlowHandler).Methods("POST")
	scheduledFlowRouter.HandleFunc("/{id:[0-9]+}", ledger.UpdateScheduledFlowHandler).Methods("PUT")
	scheduledFlowRouter.HandleFunc("/{id:[0-9]+}", ledger.DeleteScheduledFlowHandler).Methods("DELETE")
	scheduledFlowRouter.HandleFunc("/{id:[0-9]+}/occurrences", ledger.ProcessScheduledOccurrenceHandler).Methods("POST")

	router.HandleFunc("/api/bills", ledger.GetBillsHandler).Methods("GET")
	router.HandleFunc("/api/forecast/cashflow", ledger.CashflowForecastHandler).Methods("GET")
	router.HandleFunc("/api/reports/cashflow", ledger.CashflowReportHandler).Methods("GET")
	router.HandleFunc("/api/reports/cashflow.csv", ledger.CashflowCSVHandler).Methods("GET")
	router.HandleFunc("/api/reports/insights", ledger.FinancialInsightsHandler).Methods("GET")
	router.HandleFunc("/api/net-worth", ledger.NetWorthHandler).Methods("GET")
	router.HandleFunc("/api/net-worth/history", ledger.NetWorthHistoryHandler).Methods("GET")

	assetRouter := router.PathPrefix("/api/assets").Subrouter()
	assetRouter.HandleFunc("", ledger.GetManualAssetsHandler).Methods("GET")
	assetRouter.HandleFunc("", ledger.CreateManualAssetHandler).Methods("POST")
	assetRouter.HandleFunc("/{id:[0-9]+}", ledger.UpdateManualAssetHandler).Methods("PUT")
	assetRouter.HandleFunc("/{id:[0-9]+}", ledger.DeleteManualAssetHandler).Methods("DELETE")
	assetRouter.HandleFunc("/{id:[0-9]+}/valuations", ledger.GetAssetValuationsHandler).Methods("GET")
	assetRouter.HandleFunc("/{id:[0-9]+}/valuations", ledger.CreateAssetValuationHandler).Methods("POST")

	portfolioRouter := router.PathPrefix("/api/portfolio/holdings").Subrouter()
	portfolioRouter.HandleFunc("", ledger.GetPortfolioHoldingsHandler).Methods("GET")
	portfolioRouter.HandleFunc("", ledger.CreatePortfolioHoldingHandler).Methods("POST")
	portfolioRouter.HandleFunc("/{id:[0-9]+}", ledger.UpdatePortfolioHoldingHandler).Methods("PUT")
	portfolioRouter.HandleFunc("/{id:[0-9]+}", ledger.DeletePortfolioHoldingHandler).Methods("DELETE")
	portfolioRouter.HandleFunc("/{id:[0-9]+}/valuations", ledger.GetHoldingValuationsHandler).Methods("GET")
	portfolioRouter.HandleFunc("/{id:[0-9]+}/valuations", ledger.CreateHoldingValuationHandler).Methods("POST")
}
