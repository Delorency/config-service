package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "main/main/api/proto"
)

func (s *GRPCServer) GetConfig(ctx context.Context, req *pb.GetRequest) (*pb.Config, error) {
	config, err := s.repo.GetCurrent(ctx, req.ServiceId)
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
