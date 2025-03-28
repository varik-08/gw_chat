package routes

import (
	"github.com/gorilla/mux"
	"github.com/varik-08/gw_chat/config"
	"github.com/varik-08/gw_chat/server/http/controllers/auth"
	"github.com/varik-08/gw_chat/server/http/middlewares"
)

func SetupAuthRoutes(router *mux.Router, app *config.App) {
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.Use(middlewares.AppMiddleware(app))

	authRouter.HandleFunc("/registration", auth.RegistrationHandler).Methods("POST")
	authRouter.HandleFunc("/login", auth.LoginHandler).Methods("POST")
	authRouter.HandleFunc("/refresh", auth.RefreshHandler).Methods("POST")
}
