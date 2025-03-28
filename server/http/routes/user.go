package routes

import (
	"github.com/gorilla/mux"
	"github.com/varik-08/gw_chat/config"
	"github.com/varik-08/gw_chat/server/http/controllers/user"
	"github.com/varik-08/gw_chat/server/http/middlewares"
)

func SetupUserRoutes(router *mux.Router, app *config.App) {
	chatRouter := router.PathPrefix("/users").Subrouter()
	chatRouter.Use(middlewares.AuthMiddleware(app.Config.SecretKey), middlewares.AppMiddleware(app))

	chatRouter.HandleFunc("/change-password", user.UpdatePasswordHandler).Methods("PUT")
}
