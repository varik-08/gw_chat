package auth

import (
	"encoding/json"
	"net/http"

	"github.com/varik-08/gw_chat/internal/entities/auth"
	"github.com/varik-08/gw_chat/server/http/controllers"
)

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	app := controllers.GetAppFromContext(r.Context())

	var oldToken auth.Token

	// Получаем refresh token из запроса
	err := json.NewDecoder(r.Body).Decode(&oldToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	newToken, err := app.Services.AuthService.RefreshToken(&oldToken)
	if err != nil {
		http.Error(w, "ошибка авторизации", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"accessToken":  newToken.AccessToken,
		"refreshToken": newToken.RefreshToken,
	})
}
