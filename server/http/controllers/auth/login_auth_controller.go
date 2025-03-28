package auth

import (
	"encoding/json"
	"github.com/varik-08/gw_chat/server/http/controllers"
	"net/http"
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
		"user_id":       token.UserID,
		"username":      token.Username,
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
		"expires_at":    token.ExpiresAt.Format("2006-01-02T15:04:05Z"),
	})
}
