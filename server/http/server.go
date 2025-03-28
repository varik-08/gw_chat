package http

import (
	"github.com/rs/cors"
	"github.com/varik-08/gw_chat/server/http/routes"
	"log"
	"net/http"
	"time"

	"github.com/varik-08/gw_chat/config"
)

func RunServer(app *config.App) {
	log.Println("Run server on", app.Config.HTTP.Host, ":", app.Config.HTTP.Port)

	mux := routes.SetupRoutes(app)

	var allowedOrigins []string
	if app.Config.AppEnv == "production" {
		allowedOrigins = []string{"https://yourapp.com"}
	} else {
		allowedOrigins = []string{"http://localhost:3000", "http://192.168.0.220:3000"}
	}

	handler := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(mux)

	server := &http.Server{
		Addr:         app.Config.HTTP.Host + ":" + app.Config.HTTP.Port,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
