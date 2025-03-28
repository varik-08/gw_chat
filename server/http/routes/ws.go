package routes

import (
	"github.com/gorilla/mux"
	"github.com/varik-08/gw_chat/config"
	"github.com/varik-08/gw_chat/server/http/middlewares"
	"github.com/varik-08/gw_chat/server/ws"
	"log"
	"net/http"
)

func SetupWSRoutes(router *mux.Router, app *config.App) {
	wsHub := ws.RunWsHubServer(app)

	log.Println("Run ws server")

	wsRouter := router.PathPrefix("/ws").Subrouter()
	wsRouter.Use(middlewares.WsAuthMiddleware(app.Config.SecretKey))

	wsRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		ws.WebSocketHandler(wsHub, w, r)
	})
}
