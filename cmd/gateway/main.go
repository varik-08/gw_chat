package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	kgo "github.com/segmentio/kafka-go"

	mqk "github.com/varik-08/gw_chat/internal/mq/kafka"
	otelinit "github.com/varik-08/gw_chat/internal/otel"
	jwtrs "github.com/varik-08/gw_chat/internal/security/jwt"
)

func main() {
	ctx := context.Background()
	shutdown, err := otelinit.Init(ctx, "gateway")
	if err != nil {
		log.Fatalf("otel: %v", err)
	}
	defer func() { _ = shutdown(context.Background()) }()

	jwksURL := os.Getenv("AUTH_JWKS_URL")
	if jwksURL == "" {
		jwksURL = "http://auth:8082/.well-known/jwks.json"
	}
	jwks, err := jwtrs.NewRemoteJWKSProvider(ctx, jwksURL)
	if err != nil {
		log.Fatalf("jwks: %v", err)
	}

	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "kafka:9092"
	}
	prod := mqk.NewProducer([]string{brokers}, "gateway")
	cons := mqk.NewConsumer([]string{brokers}, "gateway-fanout", "evt.message.posted")
	go func() {
		_ = cons.Start(ctx, func(ctx context.Context, m kgo.Message) error {
			// TODO: fanout в подключенных WS клиентов
			return nil
		})
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	mux.HandleFunc("/ws/typing", func(w http.ResponseWriter, r *http.Request) {
		authz := r.Header.Get("Authorization")
		if authz == "" {
			http.Error(w, "no auth", http.StatusUnauthorized)
			return
		}
		token := authz[len("Bearer "):]
		claims, err := jwtrs.ValidateRS256(token, jwks)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		chatID := r.URL.Query().Get("chat_id")
		isTyping := r.URL.Query().Get("typing") == "1"
		key := []byte(chatID)
		value := []byte(`{"msg_type":"evt.chat.typing","chat_id":` + chatID + `,"user_id":` + strconv.FormatInt(claims.UserID, 10) + `,"is_typing":` + strconv.FormatBool(isTyping) + `}`)
		_ = prod.Publish(r.Context(), "evt.chat.typing", key, value, map[string]string{"msg_type": "evt.chat.typing"})
		w.WriteHeader(http.StatusAccepted)
	})

	srv := &http.Server{Addr: ":8080", Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	log.Printf("gateway http :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("http: %v", err)
	}
}

