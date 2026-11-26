package internal

import (
	"github.com/girishsaraf/xpense-trackr/backend/internal/user"
	"github.com/gorilla/mux"
)

func registerAuthRoutes(router *mux.Router) {
	authRouter := router.PathPrefix("/api/auth").Subrouter()
	authRouter.HandleFunc("/login", user.Login).Methods("POST")
	authRouter.HandleFunc("/register", user.Register).Methods("POST")
}
