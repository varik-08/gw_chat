package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	kgo "github.com/segmentio/kafka-go"

	mqk "github.com/varik-08/gw_chat/internal/mq/kafka"
	otelinit "github.com/varik-08/gw_chat/internal/otel"
	jwtrs "github.com/varik-08/gw_chat/internal/security/jwt"
)

type hub struct {
	mu          sync.RWMutex
	chatClients map[string]map[*websocket.Conn]struct{}
	upgrader    websocket.Upgrader
}

func newHub() *hub {
	return &hub{
		chatClients: make(map[string]map[*websocket.Conn]struct{}),
		upgrader:    websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
	}
}

func (h *hub) addClient(chatID string, c *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	m := h.chatClients[chatID]
	if m == nil {
		m = make(map[*websocket.Conn]struct{})
	}
	m[c] = struct{}{}
	h.chatClients[chatID] = m
}

func (h *hub) removeClient(chatID string, c *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if m, ok := h.chatClients[chatID]; ok {
		delete(m, c)
		if len(m) == 0 {
			delete(h.chatClients, chatID)
		}
	}
}

func (h *hub) broadcast(chatID string, payload []byte) {
	h.mu.RLock()
	conns := h.chatClients[chatID]
	h.mu.RUnlock()
	for c := range conns {
		_ = c.WriteMessage(websocket.TextMessage, payload)
	}
}

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

	h := newHub()

	// Kafka fanout consumers
	consMsg := mqk.NewConsumer([]string{brokers}, "gateway-fanout", "evt.message.posted")
	go func() {
		_ = consMsg.Start(ctx, func(ctx context.Context, m kgo.Message) error {
			var env map[string]any
			if err := json.Unmarshal(m.Value, &env); err != nil { return nil }
			chatID := ""
			if v, ok := env["chat_id"].(float64); ok { chatID = strconv.FormatInt(int64(v), 10) }
			if v, ok := env["chat_id"].(string); ok { chatID = v }
			if chatID == "" { return nil }
			h.broadcast(chatID, m.Value)
			return nil
		})
	}()

	consTyping := mqk.NewConsumer([]string{brokers}, "gateway-fanout", "evt.chat.typing")
	go func() {
		_ = consTyping.Start(ctx, func(ctx context.Context, m kgo.Message) error {
			var env map[string]any
			if err := json.Unmarshal(m.Value, &env); err != nil { return nil }
			chatID := ""
			if v, ok := env["chat_id"].(float64); ok { chatID = strconv.FormatInt(int64(v), 10) }
			if v, ok := env["chat_id"].(string); ok { chatID = v }
			if chatID == "" { return nil }
			h.broadcast(chatID, m.Value)
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

	// WS subscribe endpoint: /ws?chat_id=123 with Bearer token
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		authz := r.Header.Get("Authorization")
		if authz == "" { http.Error(w, "no auth", http.StatusUnauthorized); return }
		token := authz[len("Bearer "):]
		_, err := jwtrs.ValidateRS256(token, jwks)
		if err != nil { http.Error(w, "invalid token", http.StatusUnauthorized); return }
		chatID := r.URL.Query().Get("chat_id")
		if chatID == "" { http.Error(w, "chat_id required", http.StatusBadRequest); return }
		c, err := h.upgrader.Upgrade(w, r, nil)
		if err != nil { return }
		h.addClient(chatID, c)
		defer func() {
			h.removeClient(chatID, c)
			_ = c.Close()
		}()
		c.SetReadLimit(64)
		c.SetReadDeadline(time.Now().Add(60 * time.Second))
		c.SetPongHandler(func(string) error { c.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil })
		for {
			if _, _, err := c.ReadMessage(); err != nil { break }
		}
	})

	srv := &http.Server{Addr: ":8080", Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	log.Printf("gateway http :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("http: %v", err)
	}
}

