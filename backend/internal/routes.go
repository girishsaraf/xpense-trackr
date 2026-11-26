package internal

import "github.com/gorilla/mux"

// SetRoutes registers the public HTTP surface by domain.
func SetRoutes(router *mux.Router) {
	registerAuthRoutes(router)
	registerLedgerRoutes(router)
	registerReceiptRoutes(router)
	registerPlanningRoutes(router)
	registerLegacyRoutes(router)
}
