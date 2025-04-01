package auth

import (
	"encoding/json"
	"net/http"

	"github.com/varik-08/gw_chat/server/http/controllers"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	app := controllers.GetAppFromContext(r.Context())

	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "не удалось декодировать тело запроса", http.StatusUnauthorized)
		return
	}

	token, err := app.Services.AuthService.Login(credentials)
	if err != nil {
		http.Error(w, "ошибка авторизации", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"userId":       token.UserID,
		"username":     token.Username,
		"accessToken":  token.AccessToken,
		"refreshToken": token.RefreshToken,
		"expiresAt":    token.ExpiresAt.Format("2006-01-02T15:04:05Z"),
	})
}
