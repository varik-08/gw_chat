package auth

import (
	"encoding/json"
	"net/http"

	"github.com/varik-08/gw_chat/internal/entities/users"
	"github.com/varik-08/gw_chat/server/http/controllers"
)

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	app := controllers.GetAppFromContext(r.Context())

	var credentials users.Credentials

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	idUser, err := app.Services.UserService.CreateUser(credentials)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": idUser})
}
