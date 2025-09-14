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

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	kgo "github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authv1 "github.com/varik-08/gw_chat/internal/grpc/gen/authv1"
	chatv1 "github.com/varik-08/gw_chat/internal/grpc/gen/chatv1"
	messagev1 "github.com/varik-08/gw_chat/internal/grpc/gen/messagev1"
	userv1 "github.com/varik-08/gw_chat/internal/grpc/gen/userv1"
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

	// gRPC clients
	authAddr := os.Getenv("GRPC_AUTH_ADDR"); if authAddr == "" { authAddr = "auth:9099" }
	userAddr := os.Getenv("GRPC_USER_ADDR"); if userAddr == "" { userAddr = "user:9102" }
	chatAddr := os.Getenv("GRPC_CHAT_ADDR"); if chatAddr == "" { chatAddr = "chat:9100" }
	msgAddr := os.Getenv("GRPC_MESSAGE_ADDR"); if msgAddr == "" { msgAddr = "message:9101" }

	dial := func(addr string) *grpc.ClientConn {
		cc, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil { log.Fatalf("grpc dial %s: %v", addr, err) }
		return cc
	}
	ccAuth := dial(authAddr); defer ccAuth.Close()
	ccUser := dial(userAddr); defer ccUser.Close()
	ccChat := dial(chatAddr); defer ccChat.Close()
	ccMsg := dial(msgAddr); defer ccMsg.Close()
	authCli := authv1.NewAuthServiceClient(ccAuth)
	userCli := userv1.NewUserServiceClient(ccUser)
	chatCli := chatv1.NewChatServiceClient(ccChat)
	msgCli := messagev1.NewMessageServiceClient(ccMsg)

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

	r := mux.NewRouter()
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }).Methods("GET")

	r.HandleFunc("/ws/typing", func(w http.ResponseWriter, r *http.Request) {
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
	}).Methods("POST")

	// WS subscribe endpoint: /ws?chat_id=123 with Bearer token
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
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
	}).Methods("GET")

	// REST: Auth
	r.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		var in struct{ Username, Password string }
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil { http.Error(w, "bad json", http.StatusBadRequest); return }
		resp, err := authCli.Login(r.Context(), &authv1.LoginRequest{Username: in.Username, Password: in.Password})
		if err != nil { http.Error(w, err.Error(), http.StatusUnauthorized); return }
		_ = json.NewEncoder(w).Encode(resp)
	}).Methods("POST")
	r.HandleFunc("/api/auth/refresh", func(w http.ResponseWriter, r *http.Request) {
		var in struct{ RefreshToken string `json:"refresh_token"` }
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil { http.Error(w, "bad json", http.StatusBadRequest); return }
		resp, err := authCli.Refresh(r.Context(), &authv1.RefreshRequest{RefreshToken: in.RefreshToken})
		if err != nil { http.Error(w, err.Error(), http.StatusUnauthorized); return }
		_ = json.NewEncoder(w).Encode(resp)
	}).Methods("POST")

	// JWT middleware
	requireJWT := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authz := r.Header.Get("Authorization")
			if authz == "" { http.Error(w, "no auth", http.StatusUnauthorized); return }
			token := authz[len("Bearer "):]
			_, err := jwtrs.ValidateRS256(token, jwks)
			if err != nil { http.Error(w, "invalid token", http.StatusUnauthorized); return }
			next(w, r)
		}
	}

	// REST: Users
	r.HandleFunc("/api/users", requireJWT(func(w http.ResponseWriter, r *http.Request) {
		resp, err := userCli.GetUsers(r.Context(), &userv1.GetUsersRequest{})
		if err != nil { http.Error(w, err.Error(), http.StatusBadGateway); return }
		_ = json.NewEncoder(w).Encode(resp)
	})).Methods("GET")
	r.HandleFunc("/api/users/{id}", requireJWT(func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["id"]
		id, _ := strconv.ParseInt(idStr, 10, 64)
		resp, err := userCli.GetUserByID(r.Context(), &userv1.GetUserByIDRequest{Id: id})
		if err != nil { http.Error(w, err.Error(), http.StatusBadGateway); return }
		_ = json.NewEncoder(w).Encode(resp)
	})).Methods("GET")
	r.HandleFunc("/api/users/{id}/password", requireJWT(func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["id"]
		id, _ := strconv.ParseInt(idStr, 10, 64)
		var in struct{ OldPassword, NewPassword string }
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil { http.Error(w, "bad json", http.StatusBadRequest); return }
		_, err := userCli.UpdatePassword(r.Context(), &userv1.UpdatePasswordRequest{Id: id, OldPassword: in.OldPassword, NewPassword: in.NewPassword})
		if err != nil { http.Error(w, err.Error(), http.StatusBadGateway); return }
		w.WriteHeader(http.StatusNoContent)
	})).Methods("POST")

	// REST: Chats
	r.HandleFunc("/api/chats", requireJWT(func(w http.ResponseWriter, r *http.Request) {
		var in struct{ Title string; IsPublic bool; MemberIds []int64 }
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil { http.Error(w, "bad json", http.StatusBadRequest); return }
		resp, err := chatCli.CreateChat(r.Context(), &chatv1.CreateChatRequest{Title: in.Title, IsPublic: in.IsPublic, MemberIds: in.MemberIds})
		if err != nil { http.Error(w, err.Error(), http.StatusBadGateway); return }
		_ = json.NewEncoder(w).Encode(resp)
	})).Methods("POST")
	r.HandleFunc("/api/chats/{id}/members", requireJWT(func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["id"]
		id, _ := strconv.ParseInt(idStr, 10, 64)
		var in struct{ MemberIds []int64 }
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil { http.Error(w, "bad json", http.StatusBadRequest); return }
		resp, err := chatCli.AddMembers(r.Context(), &chatv1.AddMembersRequest{ChatId: id, MemberIds: in.MemberIds})
		if err != nil { http.Error(w, err.Error(), http.StatusBadGateway); return }
		_ = json.NewEncoder(w).Encode(resp)
	})).Methods("POST")
	r.HandleFunc("/api/users/{id}/chats", requireJWT(func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["id"]
		id, _ := strconv.ParseInt(idStr, 10, 64)
		resp, err := chatCli.GetUserChats(r.Context(), &chatv1.GetUserChatsRequest{UserId: id})
		if err != nil { http.Error(w, err.Error(), http.StatusBadGateway); return }
		_ = json.NewEncoder(w).Encode(resp)
	})).Methods("GET")

	// REST: Messages
	r.HandleFunc("/api/chats/{id}/messages", requireJWT(func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["id"]
		id, _ := strconv.ParseInt(idStr, 10, 64)
		resp, err := msgCli.GetMessagesByChatID(r.Context(), &messagev1.GetMessagesByChatIDRequest{ChatId: id})
		if err != nil { http.Error(w, err.Error(), http.StatusBadGateway); return }
		_ = json.NewEncoder(w).Encode(resp)
	})).Methods("GET")
	r.HandleFunc("/api/messages", requireJWT(func(w http.ResponseWriter, r *http.Request) {
		var in struct{ ChatId, SenderId int64; Text string }
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil { http.Error(w, "bad json", http.StatusBadRequest); return }
		resp, err := msgCli.CreateMessage(r.Context(), &messagev1.CreateMessageRequest{ChatId: in.ChatId, SenderId: in.SenderId, Text: in.Text})
		if err != nil { http.Error(w, err.Error(), http.StatusBadGateway); return }
		_ = json.NewEncoder(w).Encode(resp)
	})).Methods("POST")

	srv := &http.Server{Addr: ":8080", Handler: r, ReadHeaderTimeout: 5 * time.Second}
	log.Printf("gateway http :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("http: %v", err)
	}
}

