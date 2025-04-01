package chat

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/varik-08/gw_chat/server/http/controllers"
	"github.com/varik-08/gw_chat/server/http/middlewares"
)

func GetChatMessagesHandler(w http.ResponseWriter, r *http.Request) {
	app := controllers.GetAppFromContext(r.Context())

	currentUser := r.Context().Value(middlewares.UserIDKey).(int)

	vars := mux.Vars(r)
	chatIDStr := vars["id"]

	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	messages, err := app.Services.MessageService.GetMessagesByChatID(chatID, currentUser)
	if err != nil {
		http.Error(w, "Failed to get chat messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"messages": messages,
	})
}
