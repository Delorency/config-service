package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "main/api/proto"
)

func (s *GRPCServer) GetConfig(ctx context.Context, req *pb.GetRequest) (*pb.Config, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	if req.ServiceId == "" {
		return nil, status.Error(codes.InvalidArgument, "service_id is required")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	config, err := s.repo.GetActualVersion(ctx, req.ServiceId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get config: %v", err)
	}

	if config == nil {
		return nil, status.Errorf(codes.NotFound, "config not found for service: %s", req.ServiceId)
	}

	return &pb.Config{
		ServiceId: config.ServiceID,
		Data:      config.Data,
		Version:   int32(config.Version),
	}, nil
}
