package main

import (
	"booking-app/micro-service/core"
	"booking-app/micro-service/manager"
	pb "booking-app/micro-service/protobuf/gen-pb"
	"booking-app/micro-service/router"
	"fmt"
	"net"
	"sync"

	"github.com/spf13/pflag"
	"google.golang.org/grpc"
)

func main() {
	// 读取配置
	LoadConfig()

	// 启动服务
	var wg sync.WaitGroup

	options := core.NewOptions()
	noticeService := manager.NewNoticeService(options)

	wg.Add(2)

	go StartHTTPServer(&wg, options)
	go StartGRPCServer(&wg, noticeService)

	wg.Wait()
}

func LoadConfig() {
	// 读取命令行参数
	configFilePath := pflag.StringP("config", "c", "./configs/local/config.yaml", "config file path")
	pflag.Parse()

	// 初始化配置
	err := core.InitConfig(*configFilePath)
	if err != nil {
		panic(fmt.Errorf("init config err:%v", err))
	}
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
