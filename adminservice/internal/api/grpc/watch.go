package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	pb "main/api/proto"
	"main/internal/watcher"
)

func (s *GRPCServer) Watch(req *pb.WatchRequest, stream pb.ConfigService_WatchServer) error {
	if req == nil || req.ServiceId == "" {
		return status.Error(codes.InvalidArgument, "service_id is required")
	}

	subscriberID := s.generateSubscriberID(stream.Context(), req.ServiceId)

	sub, unsubscribe := s.watcher.Subscribe(req.ServiceId, subscriberID, int(req.CurrentVersion))
	defer unsubscribe()

	if err := s.sendCurrentVersionIfNeeded(stream, req.ServiceId, int(req.CurrentVersion)); err != nil {
		return err
	}

	return s.waitForUpdates(stream.Context(), stream, req.ServiceId, sub)
}

func (s *GRPCServer) generateSubscriberID(ctx context.Context, serviceID string) string {
	var clientID string
	if peer, ok := peer.FromContext(ctx); ok {
		clientID = peer.Addr.String()
	} else {
		clientID = "unknown"
	}

	return fmt.Sprintf("%s-%s-%d", serviceID, clientID, time.Now().UnixNano())
}

func (s *GRPCServer) sendCurrentVersionIfNeeded(stream pb.ConfigService_WatchServer, serviceID string, currentVersion int) error {
	current, err := s.repo.GetActualVersion(stream.Context(), serviceID)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to get current config: %v", err)
	}

	if current != nil && currentVersion < current.Version {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		default:
			if err := stream.Send(&pb.Config{
				ServiceId: current.ServiceID,
				Data:      current.Data,
				Version:   int32(current.Version),
			}); err != nil {
				return status.Errorf(codes.Unavailable, "failed to send current config: %v", err)
			}
		}
	}

	return nil
}

func (s *GRPCServer) waitForUpdates(
	ctx context.Context,
	stream pb.ConfigService_WatchServer,
	serviceID string,
	sub *watcher.Subscriber,
) error {
	heartbeat := time.NewTicker(30 * time.Second)
	defer heartbeat.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-heartbeat.C:
			if err := stream.Send(&pb.Config{
				ServiceId: serviceID,
				Data:      []byte(`{"type":"heartbeat"}`),
				Version:   0,
			}); err != nil {
				return status.Errorf(codes.Unavailable, "heartbeat failed: %v", err)
			}

		case update, ok := <-sub.UpdateChan:
			if !ok {
				return status.Error(codes.Canceled, "subscription cancelled")
			}

			sendCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			err := stream.Send(&pb.Config{
				ServiceId: serviceID,
				Data:      update.Data,
				Version:   int32(update.Version),
			})
			cancel()

			if err != nil {
				if sendCtx.Err() == context.DeadlineExceeded {
					return status.Error(codes.DeadlineExceeded, "send timeout")
				}
				return status.Errorf(codes.Unavailable, "failed to send update: %v", err)
			}
		}
	}
}
