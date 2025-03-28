package routes

import (
	"github.com/gorilla/mux"
	"github.com/varik-08/gw_chat/config"
)

func SetupRoutes(app *config.App) *mux.Router {
	router := mux.NewRouter()

	SetupAuthRoutes(router, app)
	SetupChatRoutes(router, app)
	SetupUserRoutes(router, app)

	SetupWSRoutes(router, app)

	return router
}
