package chat

import (
	"encoding/json"
	"net/http"

	"github.com/varik-08/gw_chat/internal/entities/chat"
	"github.com/varik-08/gw_chat/server/http/controllers"
)

func CreateChatHandler(w http.ResponseWriter, r *http.Request) {
	app := controllers.GetAppFromContext(r.Context())

	var chatDTO chat.ChatDTO

	err := json.NewDecoder(r.Body).Decode(&chatDTO)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	chatID, err := app.Services.ChatService.CreateChat(&chatDTO)
	if err != nil {
		http.Error(w, "Failed to create chat", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{
		"id": chatID,
	})
}
