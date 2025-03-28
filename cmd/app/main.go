package main

import (
	"github.com/varik-08/gw_chat/server/http"
	"log"

	"github.com/varik-08/gw_chat/config"
)

func main() {
	app, err := config.GetApp()
	if err != nil {
		log.Fatal(err)
	}

	http.RunServer(app)
}
