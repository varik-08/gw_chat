package main

import (
    "context"
    "log"
    "net"
    "net/http"

    "github.com/prometheus/client_golang/prometheus/promhttp"
    "google.golang.org/grpc"

    conf "github.com/varik-08/gw_chat/config"
    "github.com/varik-08/gw_chat/internal/entities/chat"
    "github.com/varik-08/gw_chat/internal/grpc/gen/chatv1"
    otelinit "github.com/varik-08/gw_chat/internal/otel"
)

type chatServer struct {
    chatv1.UnimplementedChatServiceServer
    domain *chat.Service
}

func (s *chatServer) CreateChat(ctx context.Context, req *chatv1.CreateChatRequest) (*chatv1.CreateChatResponse, error) {
    id, err := s.domain.CreateChat(&chat.DTO{OwnerID: 0, Name: req.GetTitle(), IsPublic: req.GetIsPublic()})
    if err != nil {
        return nil, err
    }
    return &chatv1.CreateChatResponse{ChatId: int64(id)}, nil
}

func (s *chatServer) AddMembers(ctx context.Context, req *chatv1.AddMembersRequest) (*chatv1.AddMembersResponse, error) {
    var added int32
    for _, uid := range req.GetMemberIds() {
        if err := s.domain.AddMember(&chat.MemberDTO{ChatID: int(req.GetChatId()), UserID: int(uid)}, 0); err == nil {
            added++
        }
    }
    return &chatv1.AddMembersResponse{ChatId: req.GetChatId(), AddedCount: added}, nil
}

func (s *chatServer) GetUserChats(ctx context.Context, req *chatv1.GetUserChatsRequest) (*chatv1.GetUserChatsResponse, error) {
    cs, err := s.domain.GetChatsByUserID(int(req.GetUserId()))
    if err != nil {
        return nil, err
    }
    out := make([]*chatv1.Chat, 0, len(cs))
    for _, c := range cs {
        out = append(out, &chatv1.Chat{Id: int64(c.ID), Title: c.Name, IsPublic: c.IsPublic})
    }
    return &chatv1.GetUserChatsResponse{Chats: out}, nil
}

func (s *chatServer) GetMembersIDsByChatID(ctx context.Context, req *chatv1.GetMembersIDsByChatIDRequest) (*chatv1.GetMembersIDsByChatIDResponse, error) {
    // Минимальный ответ: можно расширить репозиторием позже
    return &chatv1.GetMembersIDsByChatIDResponse{MemberIds: []int64{}}, nil
}

func main() {
    cfg, err := conf.GetConfig()
    if err != nil {
        log.Fatalf("config: %v", err)
    }
    shutdown, err := otelinit.Init(context.Background(), "chat-service")
    if err != nil {
        log.Fatalf("otel: %v", err)
    }
    defer func() { _ = shutdown(context.Background()) }()

    _ = conf.InitDB(cfg.DB)
    repos := conf.NewRepository()
    services := conf.NewService(cfg, repos)

    // healthz/readyz/metrics HTTP
    go func() {
        mux := http.NewServeMux()
        mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
        mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
        mux.Handle("/metrics", promhttp.Handler())
        _ = http.ListenAndServe(":9100-health", mux)
    }()

    lis, err := net.Listen("tcp", ":9100")
    if err != nil {
        log.Fatalf("listen: %v", err)
    }
    srv := grpc.NewServer()
    chatv1.RegisterChatServiceServer(srv, &chatServer{domain: services.ChatService})
    log.Printf("chat gRPC :9100")
    if err := srv.Serve(lis); err != nil {
        log.Fatalf("serve: %v", err)
    }
}

