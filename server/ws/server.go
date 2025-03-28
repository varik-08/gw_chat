package ws

import "github.com/varik-08/gw_chat/config"

func RunWsHubServer(app *config.App) *Hub {
	hub := NewHub(app)
	go hub.Run()

	return hub
}
