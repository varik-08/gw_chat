package chat

import (
	"encoding/json"
	"net/http"

	"github.com/varik-08/gw_chat/internal/entities/chat"
	"github.com/varik-08/gw_chat/server/http/controllers"
	"github.com/varik-08/gw_chat/server/http/middlewares"
)

func AddMemberHandler(w http.ResponseWriter, r *http.Request) {
	app := controllers.GetAppFromContext(r.Context())

	var memberDTO chat.MemberDTO

	err := json.NewDecoder(r.Body).Decode(&memberDTO)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	currentUser := r.Context().Value(middlewares.UserIDKey).(int)

	err = app.Services.ChatService.AddMember(&memberDTO, currentUser)
	if err != nil {
		http.Error(w, "Failed to create chat member", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
