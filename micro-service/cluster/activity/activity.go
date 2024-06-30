package activity

import (
	"booking-app/micro-service/cluster/activity/manager"
	"booking-app/micro-service/cluster/activity/router"
	"booking-app/micro-service/cluster/common"
	"booking-app/micro-service/cluster/common/commandpb"
	"booking-app/micro-service/cluster/common/logger"
	pb "booking-app/micro-service/protobuf/gen-pb"
	"net"
	"sync"

	"google.golang.org/grpc"
)

func Start(wg *sync.WaitGroup) {
	noticeService := manager.NewNoticeService()
	boardcastService := manager.NewBoardcastService()

	wg.Add(2)

	go StartHTTPServer(wg)
	go StartGRPCServer(wg, noticeService, boardcastService)
}

// HTTP server
func StartHTTPServer(wg *sync.WaitGroup) {
	defer wg.Done()

	r := router.Routers()
	logger.Info("Activity HTTP Service is running on port " + common.Config.Activity.TCPPort)

	r.Run(":" + common.Config.Activity.TCPPort)
}

// GRPC server
func StartGRPCServer(wg *sync.WaitGroup, n *manager.NoticeService, b *manager.BroadcastService) {
	defer wg.Done()

	listen, err := net.Listen("tcp", ":"+common.Config.Activity.GRPCPort)
	if err != nil {
		logger.Errorf("grpc failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNoticeServiceServer(grpcServer, n)
	pb.RegisterBroadcastServiceServer(grpcServer, b)
	logger.Info("Activity GRPC Service is running on port " + common.Config.Activity.GRPCPort)
	if err := grpcServer.Serve(listen); err != nil {
		logger.Errorf("Failed to serve gRPC: %v", err)
	}

	// 注册rpc方法
	n.RegisterHandler(commandpb.Command_ActivityGetNotice, n.GetNotice)
}
