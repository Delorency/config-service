package grpc

import (
	"fmt"
	pb "main/main/api/proto"
)

func (s *GRPCServer) WatchConfig(req *pb.WatchRequest, stream pb.ConfigService_WatchConfigServer) error {
	serviceID := req.ServiceId
	subscriberID := fmt.Sprintf("%s-%d", serviceID, stream.Context().Value("client_id"))

	sub, unsubscribe := s.watcher.Subscribe(serviceID, subscriberID, int(req.CurrentVersion))
	defer unsubscribe()

	current, err := s.repo.GetCurrent(stream.Context(), serviceID)
	if err != nil {
		return err
	}

	if current != nil && int(req.CurrentVersion) < current.Version {
		if err := stream.Send(&pb.Config{
			ServiceId: current.ServiceID,
			Data:      current.Data,
			Version:   int32(current.Version),
		}); err != nil {
			return err
		}
	}

	for {
		select {
		case <-stream.Context().Done():
			return nil

		case update := <-sub.UpdateChan:
			if err := stream.Send(&pb.Config{
				ServiceId: serviceID,
				Data:      update.Data,
				Version:   int32(update.Version),
			}); err != nil {
				return err
			}
		}
	}
}
