package main

import (
	"context"
	"crypto/rsa"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc"

	conf "github.com/varik-08/gw_chat/config"
	"github.com/varik-08/gw_chat/internal/entities/auth"
	"github.com/varik-08/gw_chat/internal/entities/users"
	"github.com/varik-08/gw_chat/internal/grpc/gen/authv1"
	"github.com/varik-08/gw_chat/internal/otel"
	jwtrs "github.com/varik-08/gw_chat/internal/security/jwt"
)

type authServer struct {
	authv1.UnimplementedAuthServiceServer
	domain *auth.Service
}

func (s *authServer) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	tkn, err := s.domain.Login(users.Credentials{Username: req.GetUsername(), Password: req.GetPassword()})
	if err != nil {
		return nil, err
	}
	return &authv1.LoginResponse{Tokens: &authv1.TokenPair{
		AccessToken:        tkn.AccessToken,
		RefreshToken:       tkn.RefreshToken,
		AccessExpiresUnix:  tkn.ExpiresAt.Unix(),
		RefreshExpiresUnix: tkn.ExpiresAt.Add(7 * 24 * time.Hour).Unix(),
	}}, nil
}

func (s *authServer) Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.RefreshResponse, error) {
	r, err := s.domain.RefreshToken(&auth.Token{RefreshToken: req.GetRefreshToken()})
	if err != nil {
		return nil, err
	}
	return &authv1.RefreshResponse{Tokens: &authv1.TokenPair{
		AccessToken:  r.AccessToken,
		RefreshToken: r.RefreshToken,
	}}, nil
}

func main() {
	// Init config
	cfg, err := conf.GetConfig()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	ctx := context.Background()
	shutdown, err := otel.Init(ctx, "auth-service")
	if err != nil {
		log.Fatalf("otel init error: %v", err)
	}
	defer func() { _ = shutdown(context.Background()) }()

	// RSA keys (for demo: generate on startup or load from files)
	var priv *rsa.PrivateKey
	var kid = "auth-key-1"
	if pem := os.Getenv("AUTH_RSA_PRIVATE_PEM"); pem != "" {
		p, err := jwtrs.ParseRSAPrivateKeyPEM([]byte(pem))
		if err != nil {
			log.Fatalf("invalid AUTH_RSA_PRIVATE_PEM: %v", err)
		}
		priv = p
	} else {
		kp, err := jwtrs.GenerateRSAKey(2048, kid)
		if err != nil {
			log.Fatalf("rsa gen error: %v", err)
		}
		priv = kp.PrivateKey
	}
	pub := &priv.PublicKey

	// JWKS HTTP server
	jwks, err := jwtrs.NewJWKSProvider(pub, kid)
	if err != nil {
		log.Fatalf("jwks init error: %v", err)
	}
	mux := http.NewServeMux()
	mux.Handle("/.well-known/jwks.json", jwks.Handler())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	httpSrv := &http.Server{
		Addr:              net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Wire domain service
	repos := conf.NewRepository()
	services := conf.NewService(cfg, repos)

	// gRPC server
	grpcLis, err := net.Listen("tcp", ":9099")
	if err != nil {
		log.Fatalf("grpc listen: %v", err)
	}
	grpcSrv := grpc.NewServer()
	authv1.RegisterAuthServiceServer(grpcSrv, &authServer{domain: services.AuthService})
	go func() {
		if err := grpcSrv.Serve(grpcLis); err != nil {
			log.Printf("grpc serve stopped: %v", err)
		}
	}()

	log.Printf("auth-service http jwks on %s, grpc on %s", httpSrv.Addr, ":9099")
	if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("http serve: %v", err)
	}
}

