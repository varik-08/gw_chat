package user

import (
	"encoding/json"
	"net/http"

	"github.com/varik-08/gw_chat/server/http/controllers"
	"github.com/varik-08/gw_chat/server/http/middlewares"
)

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	app := controllers.GetAppFromContext(r.Context())

	var req ChangePasswordRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	currentUser := r.Context().Value(middlewares.UserIDKey).(int)

	err = app.Services.UserService.UpdatePassword(currentUser, req.NewPassword, req.OldPassword)
	if err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
