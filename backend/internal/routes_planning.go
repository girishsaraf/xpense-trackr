package internal

import (
	"github.com/girishsaraf/xpense-trackr/backend/internal/planning"
	"github.com/gorilla/mux"
)

func registerPlanningRoutes(router *mux.Router) {
	planningRouter := router.PathPrefix("/api/planning").Subrouter()
	planningRouter.HandleFunc("/goals", planning.GetGoalsHandler).Methods("GET")
	planningRouter.HandleFunc("/goals", planning.CreateGoalHandler).Methods("POST")
	planningRouter.HandleFunc("/goals/{id:[0-9]+}", planning.DeleteGoalHandler).Methods("DELETE")
	planningRouter.HandleFunc("/goals/{id:[0-9]+}/contributions", planning.GetContributionsHandler).Methods("GET")
	planningRouter.HandleFunc("/goals/{id:[0-9]+}/contributions", planning.CreateContributionHandler).Methods("POST")
	planningRouter.HandleFunc("/envelopes", planning.GetEnvelopesHandler).Methods("GET")
	planningRouter.HandleFunc("/envelopes", planning.CreateEnvelopeHandler).Methods("POST")
	planningRouter.HandleFunc("/envelopes/{id:[0-9]+}", planning.UpdateEnvelopeHandler).Methods("PUT")
	planningRouter.HandleFunc("/envelopes/{id:[0-9]+}", planning.DeleteEnvelopeHandler).Methods("DELETE")
}
