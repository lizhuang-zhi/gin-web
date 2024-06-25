package main

import (
	"booking-app/micro-service/core"
	"booking-app/micro-service/manager"
	pb "booking-app/micro-service/protobuf/gen-pb"
	"booking-app/micro-service/router"
	"net"
	"sync"

	"google.golang.org/grpc"
)

func main() {
	var wg sync.WaitGroup

	options := core.NewOptions()
	noticeService := manager.NewNoticeService(options)

	wg.Add(2)

	go startHTTPServer(&wg, options)
	go startGRPCServer(&wg, noticeService)

	wg.Wait()
}

// HTTP server
func startHTTPServer(wg *sync.WaitGroup, opts *core.Options) {
	defer wg.Done()

	r := router.Routers(opts)
	opts.Logger.Println("Activity HTTP Service is running on port 13001...")

	r.Run(":13001")
}

// GRPC server
func startGRPCServer(wg *sync.WaitGroup, noticeService *manager.NoticeService) {
	defer wg.Done()

	listen, err := net.Listen("tcp", ":13002")
	if err != nil {
		noticeService.GetOptions().Logger.Fatalf("grpc failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNoticeServiceServer(grpcServer, noticeService)
	noticeService.GetOptions().Logger.Println("Activity GRPC Service is running on port 13002...")
	if err := grpcServer.Serve(listen); err != nil {
		noticeService.GetOptions().Logger.Fatalf("Failed to serve gRPC: %v", err)
	}
}
