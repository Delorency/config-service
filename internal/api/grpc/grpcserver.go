package grpc

import (
	"fmt"
	"main/internal/core/validator"
	"main/internal/repo"
	"main/internal/watcher"
	"net"
	"sync"
	"time"

	pb "main/main/api/proto"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	pb.UnimplementedConfigServiceServer

	repo    *repo.ConfigRepository
	watcher *watcher.WatcherManager

	grpcServer   *grpc.Server
	listener     net.Listener
	host         string
	port         int
	mu           sync.Mutex
	isRunning    bool
	shutdownDone chan struct{}
}

func NewGRPCServer(
	repo *repo.ConfigRepository,
	watcher *watcher.WatcherManager,
	validator *validator.Validator,
) *GRPCServer {
	return &GRPCServer{
		repo:    repo,
		watcher: watcher,
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

	s.grpcServer = grpc.NewServer(
		grpc.MaxRecvMsgSize(4<<20), // 4MB
		grpc.MaxSendMsgSize(4<<20), // 4MB
		grpc.ConnectionTimeout(120*time.Second),
	)

	pb.RegisterConfigServiceServer(s.grpcServer, s)

	addr := fmt.Sprintf("%s:%d", host, port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	s.listener = listener
	s.isRunning = true

	go func() {
		defer close(s.shutdownDone)

		if err := s.grpcServer.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			fmt.Printf("gRPC server error: %v\n", err)
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

	fmt.Println("gRPC server: starting graceful shutdown...")

	// Канал для сигнала завершения
	done := make(chan struct{})

	go func() {
		s.grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("gRPC server: graceful shutdown completed")
	case <-time.After(timeout):
		fmt.Printf("gRPC server: graceful shutdown timeout after %v, forcing stop...\n", timeout)
		s.grpcServer.Stop()
	}

	if s.listener != nil {
		s.listener.Close()
	}

	<-s.shutdownDone

	s.isRunning = false
	fmt.Println("gRPC server: stopped")

	return nil
}

func (s *GRPCServer) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return
	}

	fmt.Println("gRPC server: forcing stop...")

	if s.grpcServer != nil {
		s.grpcServer.Stop()
	}

	if s.listener != nil {
		s.listener.Close()
	}

	<-s.shutdownDone
	s.isRunning = false

	fmt.Println("gRPC server: stopped")
}

func (s *GRPCServer) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.isRunning
}

func (s *GRPCServer) GetAddr() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.listener != nil {
		return s.listener.Addr().String()
	}
	return ""
}
