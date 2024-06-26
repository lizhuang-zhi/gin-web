package activity

import (
	"booking-app/micro-service/cluster/activity/manager"
	"booking-app/micro-service/cluster/activity/router"
	"booking-app/micro-service/cluster/common/core"
	"booking-app/micro-service/core/server"
	pb "booking-app/micro-service/protobuf/gen-pb"
	"net"
	"sync"

	"google.golang.org/grpc"
)

type Server struct {
	server.Base
}

func NewServerInstance() *Server {
	return &Server{}
}

func (s *Server) Start() error {
	// 启动服务
	var wg sync.WaitGroup

	options := core.NewOptions()
	noticeService := manager.NewNoticeService(options)

	wg.Add(2)

	go StartHTTPServer(&wg, options)
	go StartGRPCServer(&wg, noticeService)

	wg.Wait()

	return nil
}

// HTTP server
func StartHTTPServer(wg *sync.WaitGroup, opts *core.Options) {
	defer wg.Done()

	r := router.Routers(opts)
	opts.Logger.Println("Activity HTTP Service is running on port " + core.Config.System.TCPPort)

	r.Run(":" + core.Config.System.TCPPort)
}

// GRPC server
func StartGRPCServer(wg *sync.WaitGroup, noticeService *manager.NoticeService) {
	defer wg.Done()

	listen, err := net.Listen("tcp", ":"+core.Config.System.RPCPort)
	if err != nil {
		noticeService.GetOptions().Logger.Fatalf("grpc failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNoticeServiceServer(grpcServer, noticeService)
	noticeService.GetOptions().Logger.Println("Activity GRPC Service is running on port " + core.Config.System.RPCPort)
	if err := grpcServer.Serve(listen); err != nil {
		noticeService.GetOptions().Logger.Fatalf("Failed to serve gRPC: %v", err)
	}
}
