package user

import (
	"encoding/json"
	"net/http"

	"github.com/varik-08/gw_chat/server/http/controllers"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	app := controllers.GetAppFromContext(r.Context())

	users, err := app.Services.UserService.GetUsers()
	if err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"users": users,
	})
}
