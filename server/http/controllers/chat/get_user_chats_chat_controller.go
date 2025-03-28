package chat

import (
	"encoding/json"
	"net/http"

	"github.com/varik-08/gw_chat/server/http/controllers"
)

func GetUserChatsHandler(w http.ResponseWriter, r *http.Request) {
	app := controllers.GetAppFromContext(r.Context())

	currentUser := r.Context().Value("userID").(int)

	chats, err := app.Services.ChatService.GetChatsByUserID(currentUser)
	if err != nil {
		http.Error(w, "Failed to get user chats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"chats": chats,
	})
}
