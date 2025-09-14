package main

import (
    "context"
    "log"
    "net"
    "net/http"

    "google.golang.org/grpc"

    conf "github.com/varik-08/gw_chat/config"
    "github.com/varik-08/gw_chat/internal/entities/users"
    "github.com/varik-08/gw_chat/internal/grpc/gen/userv1"
    otelinit "github.com/varik-08/gw_chat/internal/otel"
)

type userServer struct {
    userv1.UnimplementedUserServiceServer
    domain *users.UserService
}

func (s *userServer) GetUsers(ctx context.Context, _ *userv1.GetUsersRequest) (*userv1.GetUsersResponse, error) {
    us, err := s.domain.GetUsers()
    if err != nil { return nil, err }
    out := make([]*userv1.UserItem, 0, len(us))
    for _, u := range us { out = append(out, &userv1.UserItem{Id: int64(u.ID), Username: u.Username}) }
    return &userv1.GetUsersResponse{Users: out}, nil
}

func (s *userServer) GetUserByID(ctx context.Context, req *userv1.GetUserByIDRequest) (*userv1.GetUserByIDResponse, error) {
    u, err := s.domain.UserRepository.GetUserByID(int(req.GetId()))
    if err != nil { return nil, err }
    return &userv1.GetUserByIDResponse{User: &userv1.UserItem{Id: int64(u.ID), Username: u.Username}}, nil
}

func (s *userServer) UpdatePassword(ctx context.Context, req *userv1.UpdatePasswordRequest) (*userv1.UpdatePasswordResponse, error) {
    if err := s.domain.UpdatePassword(int(req.GetId()), req.GetNewPassword(), req.GetOldPassword()); err != nil { return nil, err }
    return &userv1.UpdatePasswordResponse{}, nil
}

func main() {
    cfg, err := conf.GetConfig()
    if err != nil { log.Fatalf("config: %v", err) }
    shutdown, err := otelinit.Init(context.Background(), "user-service")
    if err != nil { log.Fatalf("otel: %v", err) }
    defer func() { _ = shutdown(context.Background()) }()

    _ = conf.InitDB(cfg.DB)
    repos := conf.NewRepository()
    services := conf.NewService(cfg, repos)

    // healthz HTTP
    go func() {
        mux := http.NewServeMux()
        mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
        _ = http.ListenAndServe(":9102-health", mux)
    }()

    lis, err := net.Listen("tcp", ":9102")
    if err != nil { log.Fatalf("listen: %v", err) }
    srv := grpc.NewServer()
    userv1.RegisterUserServiceServer(srv, &userServer{domain: services.UserService})
    log.Printf("user gRPC :9102")
    if err := srv.Serve(lis); err != nil { log.Fatalf("serve: %v", err) }
}

