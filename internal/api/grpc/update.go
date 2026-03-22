package grpc

import (
	"context"
	"main/internal/models"
	"main/internal/watcher"
	pb "main/main/api/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GRPCServer) UpdateConfig(ctx context.Context, req *pb.UpdateRequest) (*pb.Config, error) {
	config := &models.Config{
		ServiceID: req.ServiceId,
		Data:      req.Data,
	}

	if err := s.repo.Update(ctx, config); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update: %v", err)
	}

	updated, _ := s.repo.GetCurrent(ctx, req.ServiceId)

	s.watcher.Notify(req.ServiceId, &watcher.ConfigUpdate{
		ServiceID: req.ServiceId,
		Data:      updated.Data,
		Version:   updated.Version,
	})

	return &pb.Config{
		ServiceId: updated.ServiceID,
		Data:      updated.Data,
		Version:   int32(updated.Version),
	}, nil
}
