package manager

import (
	"booking-app/micro-service/cluster/common/logger"
	"time"

	pb "booking-app/micro-service/protobuf/gen-pb"
)

type BroadcastService struct {
	pb.UnimplementedBroadcastServiceServer
}

func NewBoardcastService() *BroadcastService {
	return &BroadcastService{}
}

func (s *BroadcastService) BroadcastPlayerNotify(req *pb.BroadcastRequest, stream pb.BroadcastService_BroadcastPlayerNotifyServer) error {
	for {
		// 模拟广播消息
		packet := &pb.NotifyPacket{
			Message: "Broadcast message " + time.Now().Format(time.RFC3339),
		}
		broadcast := &pb.PlayerBroadcast{
			Delay: 1000, // 延时 单位：ms
			Packets: []*pb.NotifyPacket{
				packet,
			},
			OsTypes: []string{"android", "ios"},
		}

		if err := stream.Send(broadcast); err != nil {
			logger.Infof("Error sending broadcast message: %v", err)
			return err
		}

		logger.Infof("Sent broadcast: %v", packet.Message)
		time.Sleep(1 * time.Second)
	}
}
