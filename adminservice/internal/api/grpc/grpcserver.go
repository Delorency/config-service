package grpc

import (
	"context"
	"fmt"
	"log"
	"main/internal/repo"
	"main/internal/watcher"
	"net"
	"sync"
	"time"

	pb "main/api/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	pb.UnimplementedConfigServiceServer

	repo    *repo.ConfigRepository
	watcher *watcher.WatcherManager

	grpcServer   *grpc.Server
	listener     net.Listener
	host         string
	port         int
	mu           sync.RWMutex
	isRunning    bool
	shutdownDone chan struct{}
	healthServer *health.Server
	log          *log.Logger
}

func NewGRPCServer(
	repo *repo.ConfigRepository,
	watcher *watcher.WatcherManager,
	log *log.Logger,
) *GRPCServer {
	return &GRPCServer{
		repo:    repo,
		watcher: watcher,
		log:     log,
	}
}

func (s *GRPCServer) Start(host string, port int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return fmt.Errorf("server is already running")
	}

	s.host = host
	s.port = port
	s.shutdownDone = make(chan struct{})

	keepaliveParams := keepalive.ServerParameters{
		MaxConnectionIdle:     5 * time.Minute,
		MaxConnectionAge:      10 * time.Minute,
		MaxConnectionAgeGrace: 5 * time.Second,
		Time:                  30 * time.Second,
		Timeout:               10 * time.Second,
	}

	s.grpcServer = grpc.NewServer(
		grpc.MaxRecvMsgSize(4<<20),
		grpc.MaxSendMsgSize(4<<20),
		grpc.ConnectionTimeout(120*time.Second),
		grpc.KeepaliveParams(keepaliveParams),
		grpc.ChainUnaryInterceptor(
			s.loggingInterceptor(),
			s.recoveryInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			s.streamLoggingInterceptor(),
			s.streamRecoveryInterceptor(),
		),
	)

	pb.RegisterConfigServiceServer(s.grpcServer, s)

	s.healthServer = health.NewServer()
	grpc_health_v1.RegisterHealthServer(s.grpcServer, s.healthServer)
	s.healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	reflection.Register(s.grpcServer)

	addr := fmt.Sprintf("%s:%d", host, port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	s.listener = listener
	s.isRunning = true

	go func() {
		defer close(s.shutdownDone)
		s.log.Printf("gRPC server listening on %s", addr)
		if err := s.grpcServer.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			s.log.Printf("gRPC server error: %v", err)
		}
	}()

	return nil
}

func (s *GRPCServer) GracefulStop(timeout time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return nil
	}

	s.log.Print("gRPC server: starting graceful shutdown...")

	if s.healthServer != nil {
		s.healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	}

	if s.watcher != nil {
		s.log.Print("Stopping watcher manager...")
		s.watcher.Stop()
	}

	done := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		s.log.Print("gRPC server: graceful shutdown completed")
	case <-time.After(timeout):
		s.log.Printf("gRPC server: graceful shutdown timeout after %v, forcing stop...", timeout)
		s.grpcServer.Stop()
	}

	if s.listener != nil {
		s.listener.Close()
	}

	select {
	case <-s.shutdownDone:
		s.log.Print("Serve goroutine stopped")
	case <-time.After(5 * time.Second):
		s.log.Print("Serve goroutine stop timeout")
	}

	s.isRunning = false
	s.log.Print("gRPC server: stopped")

	return nil
}

func (s *GRPCServer) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return
	}

	s.log.Print("gRPC server: forcing stop...")

	if s.watcher != nil {
		s.watcher.Stop()
	}

	if s.grpcServer != nil {
		s.grpcServer.Stop()
	}

	if s.listener != nil {
		s.listener.Close()
	}

	select {
	case <-s.shutdownDone:
		s.log.Print("Serve goroutine stopped")
	case <-time.After(5 * time.Second):
		s.log.Print("Serve goroutine stop timeout")
	}

	s.isRunning = false
	s.log.Print("gRPC server: stopped")
}

func (s *GRPCServer) loggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			s.log.Printf("[gRPC] Method: %s, Duration: %v, Error: %v\n", info.FullMethod, duration, err)
		} else {
			s.log.Printf("[gRPC] Method: %s, Duration: %v, Success\n", info.FullMethod, duration)
		}

		return resp, err
	}
}

func (s *GRPCServer) recoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				s.log.Printf("[gRPC] Panic recovered: %v, Method: %s\n", r, info.FullMethod)
				err = fmt.Errorf("internal server error")
			}
		}()
		return handler(ctx, req)
	}
}

func (s *GRPCServer) streamLoggingInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		err := handler(srv, ss)
		duration := time.Since(start)

		if err != nil {
			s.log.Printf("[gRPC Stream] Method: %s, Duration: %v, Error: %v\n", info.FullMethod, duration, err)
		} else {
			s.log.Printf("[gRPC Stream] Method: %s, Duration: %v, Success\n", info.FullMethod, duration)
		}

		return err
	}
}

func (s *GRPCServer) streamRecoveryInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				s.log.Printf("[gRPC Stream] Panic recovered: %v, Method: %s\n", r, info.FullMethod)
				err = fmt.Errorf("internal server error")
			}
		}()
		return handler(srv, ss)
	}
}

func (s *GRPCServer) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isRunning
}

func (s *GRPCServer) GetAddr() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.listener != nil {
		return s.listener.Addr().String()
	}
	return ""
}

func (s *GRPCServer) CheckHealth(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	if !s.IsRunning() {
		return &pb.HealthResponse{Status: "NOT_SERVING"}, nil
	}

	if s.watcher == nil {
		return &pb.HealthResponse{Status: "NOT_SERVING"}, nil
	}

	return &pb.HealthResponse{Status: "SERVING"}, nil
}
