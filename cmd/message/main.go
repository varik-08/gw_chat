package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/grpc"

	conf "github.com/varik-08/gw_chat/config"
	"github.com/varik-08/gw_chat/internal/entities/message"
	"github.com/varik-08/gw_chat/internal/grpc/gen/messagev1"
	mqk "github.com/varik-08/gw_chat/internal/mq/kafka"
	otelinit "github.com/varik-08/gw_chat/internal/otel"
)

type messageServer struct {
	messagev1.UnimplementedMessageServiceServer
	domain   *message.Service
	producer *mqk.Producer
}

func (s *messageServer) CreateMessage(ctx context.Context, req *messagev1.CreateMessageRequest) (*messagev1.CreateMessageResponse, error) {
	id, err := s.domain.CreateChatMessage(&message.DTO{ChatID: int(req.GetChatId()), UserID: int(req.GetSenderId()), Content: req.GetText()})
	if err != nil {
		return nil, err
	}
	// Publish evt.message.posted (немедленно; дублируется outbox'ом — допускаем для демо)
	evt := map[string]any{
		"msg_type":  "evt.message.posted",
		"chat_id":   req.GetChatId(),
		"message_id": id,
		"sender_id": req.GetSenderId(),
		"text":      req.GetText(),
		"ts":        time.Now().Unix(),
	}
	b, _ := json.Marshal(evt)
	_ = s.producer.Publish(ctx, "evt.message.posted", []byte(strconv.FormatInt(req.GetChatId(), 10)), b, map[string]string{"msg_type": "evt.message.posted"})

	return &messagev1.CreateMessageResponse{MessageId: int64(id), CreatedUnix: time.Now().Unix()}, nil
}

func (s *messageServer) GetMessagesByChatID(ctx context.Context, req *messagev1.GetMessagesByChatIDRequest) (*messagev1.GetMessagesByChatIDResponse, error) {
	items, err := s.domain.GetMessagesByChatID(int(req.GetChatId()), 0)
	if err != nil {
		return nil, err
	}
	out := make([]*messagev1.MessageItem, 0, len(items))
	for _, m := range items {
		out = append(out, &messagev1.MessageItem{
			Id:          int64(m.ID),
			ChatId:      int64(m.ChatID),
			SenderId:    int64(m.UserID),
			Text:        m.Content,
			CreatedUnix: m.CreatedAt.Unix(),
		})
	}
	return &messagev1.GetMessagesByChatIDResponse{Messages: out}, nil
}

func main() {
	cfg, err := conf.GetConfig()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	shutdown, err := otelinit.Init(context.Background(), "message-service")
	if err != nil {
		log.Fatalf("otel: %v", err)
	}
	defer func() { _ = shutdown(context.Background()) }()

	db := conf.InitDB(cfg.DB)
	repos := conf.NewRepository()
	services := conf.NewService(cfg, repos)

	prod := mqk.NewProducer(cfg.Kafka.Brokers, "message-service")

	// healthz HTTP
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
		_ = http.ListenAndServe(":9101-health", mux)
	}()

	// outbox worker
	go func() {
		pub := message.NewOutboxPublisher(db, prod)
		_ = pub.Run(context.Background())
	}()

	lis, err := net.Listen("tcp", ":9101")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	srv := grpc.NewServer()
	messagev1.RegisterMessageServiceServer(srv, &messageServer{domain: services.MessageService, producer: prod})
	log.Printf("message gRPC :9101")
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

