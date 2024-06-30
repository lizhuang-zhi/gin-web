package rpc

import (
	pb "booking-app/micro-service/protobuf/gen-pb"

	"google.golang.org/grpc"
)

func GetActivityNoticeClient(conn *grpc.ClientConn) pb.NoticeServiceClient {
	return pb.NewNoticeServiceClient(conn)
}

func GetActivityBroadcastClient(conn *grpc.ClientConn) pb.BroadcastServiceClient {
	return pb.NewBroadcastServiceClient(conn)
}

func GetLobbyUserClient(conn *grpc.ClientConn) pb.UserServiceClient {
	return pb.NewUserServiceClient(conn)
}
