package lobby

import (
	"booking-app/micro-service/cluster/common"
	"booking-app/micro-service/cluster/common/logger"
	"booking-app/micro-service/cluster/lobby/manager"
	"net"
	"sync"

	pb "booking-app/micro-service/protobuf/gen-pb"

	"google.golang.org/grpc"
)

func Start(wg *sync.WaitGroup) {
	userService := manager.NewUserService()

	wg.Add(1)

	go StartGRPCServer(wg, userService)
}

func StartGRPCServer(wg *sync.WaitGroup, u *manager.UserServcie) {
	defer wg.Done()

	listen, err := net.Listen("tcp", ":"+common.Config.Lobby.GRPCPort)
	if err != nil {
		logger.Errorf("grpc failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, u)
	logger.Info("Lobby GRPC Service is running on port " + common.Config.Lobby.GRPCPort)
	if err := grpcServer.Serve(listen); err != nil {
		logger.Errorf("Failed to serve gRPC: %v", err)
	}
}
