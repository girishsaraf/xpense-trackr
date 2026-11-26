package internal

import (
	"github.com/girishsaraf/xpense-trackr/backend/internal/receipt"
	"github.com/gorilla/mux"
)

func registerReceiptRoutes(router *mux.Router) {
	receiptRouter := router.PathPrefix("/api/receipts").Subrouter()
	receiptRouter.HandleFunc("", receipt.GetReceiptsHandler).Methods("GET")
	receiptRouter.HandleFunc("", receipt.CreateReceiptHandler).Methods("POST")
	receiptRouter.HandleFunc("/{id:[0-9]+}", receipt.DeleteReceiptHandler).Methods("DELETE")
	receiptRouter.HandleFunc("/{id:[0-9]+}/file", receipt.GetReceiptFileHandler).Methods("GET")
}
