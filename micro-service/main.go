package main

import (
	"booking-app/micro-service/core"
	"booking-app/micro-service/model"
	pb "booking-app/micro-service/protobuf/gen-pb"
	"booking-app/micro-service/router"
	"context"
	"net"
	"sync"

	"google.golang.org/grpc"
)

type NoticeService struct {
	opts *core.Options
	pb.UnimplementedNoticeServiceServer
}

func (s *NoticeService) GetNotice(ctx context.Context, req *pb.GetNoticeRequest) (*pb.GetNoticeResponse, error) {
	notice_id := req.Id

	resp_notice := &pb.Notice{}

	for id, notice := range model.GlobalNotice {
		int64_id := int64(id)
		if int64_id == notice_id {
			resp_notice = &pb.Notice{
				Id:      int64_id,
				Title:   notice.Title,
				Content: notice.Content,
				SubType: int64(notice.SubType),
			}
		}
	}

	return &pb.GetNoticeResponse{
		Notice: resp_notice,
	}, nil
}

func main() {
	var wg sync.WaitGroup

	options := core.NewOptions()
	noticeService := &NoticeService{
		opts: options,
	}

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
func startGRPCServer(wg *sync.WaitGroup, noticeService *NoticeService) {
	defer wg.Done()

	listen, err := net.Listen("tcp", ":13002")
	if err != nil {
		noticeService.opts.Logger.Fatalf("grpc failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNoticeServiceServer(grpcServer, noticeService)
	noticeService.opts.Logger.Println("Activity GRPC Service is running on port 13002...")
	if err := grpcServer.Serve(listen); err != nil {
		noticeService.opts.Logger.Fatalf("Failed to serve gRPC: %v", err)
	}
}
