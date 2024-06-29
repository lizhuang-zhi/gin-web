package activity

import (
	"booking-app/micro-service/cluster/activity/manager"
	"booking-app/micro-service/cluster/activity/router"
	"booking-app/micro-service/cluster/common"
	"booking-app/micro-service/cluster/common/logger"
	pb "booking-app/micro-service/protobuf/gen-pb"
	"net"
	"sync"

	"google.golang.org/grpc"
)

func Start() error {
	// 启动服务
	var wg sync.WaitGroup

	noticeService := manager.NewNoticeService()
	boardcastService := manager.NewBoardcastService()

	wg.Add(2)

	go StartHTTPServer(&wg)
	go StartGRPCServer(&wg, noticeService, boardcastService)

	wg.Wait()

	return nil
}

// HTTP server
func StartHTTPServer(wg *sync.WaitGroup) {
	defer wg.Done()

	r := router.Routers()
	logger.Info("Activity HTTP Service is running on port " + common.Config.System.TCPPort)

	r.Run(":" + common.Config.System.TCPPort)
}

// GRPC server
func StartGRPCServer(wg *sync.WaitGroup, noticeService *manager.NoticeService, boardcastService *manager.BroadcastService) {
	defer wg.Done()

	listen, err := net.Listen("tcp", ":"+common.Config.System.RPCPort)
	if err != nil {
		logger.Errorf("grpc failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNoticeServiceServer(grpcServer, noticeService)
	pb.RegisterBroadcastServiceServer(grpcServer, boardcastService)
	logger.Info("Activity GRPC Service is running on port " + common.Config.System.RPCPort)
	if err := grpcServer.Serve(listen); err != nil {
		logger.Errorf("Failed to serve gRPC: %v", err)
	}
}
