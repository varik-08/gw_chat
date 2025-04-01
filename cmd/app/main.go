package main

import (
	"log"

	"github.com/varik-08/gw_chat/config"
	"github.com/varik-08/gw_chat/server/http"
)

func main() {
	app, err := config.GetApp()
	if err != nil {
		log.Fatal(err)
	}

	http.RunServer(app)
}
