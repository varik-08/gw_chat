package routes

import (
	"github.com/gorilla/mux"
	"github.com/varik-08/gw_chat/config"
	"github.com/varik-08/gw_chat/server/http/controllers/chat"
	"github.com/varik-08/gw_chat/server/http/middlewares"
)

func SetupChatRoutes(router *mux.Router, app *config.App) {
	chatRouter := router.PathPrefix("/chats").Subrouter()
	chatRouter.Use(middlewares.AuthMiddleware(app.Config.SecretKey), middlewares.AppMiddleware(app))

	chatRouter.HandleFunc("", chat.CreateChatHandler).Methods("POST")
	chatRouter.HandleFunc("", chat.GetUserChatsHandler).Methods("GET")
	chatRouter.HandleFunc("/{id}/messages", chat.GetChatMessagesHandler).Methods("GET")
	chatRouter.HandleFunc("/add-member", chat.AddMemberHandler).Methods("POST")
}
