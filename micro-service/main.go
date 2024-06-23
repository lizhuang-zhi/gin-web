package main

import (
	"booking-app/micro-service/model"
	pb "booking-app/micro-service/protobuf/gen-pb"
	"booking-app/micro-service/router"
	"context"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

type NoticeService struct {
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

	wg.Add(2)

	go startHTTPServer(&wg)
	go startGRPCServer(&wg)

	wg.Wait()
}

// HTTP server
func startHTTPServer(wg *sync.WaitGroup) {
	defer wg.Done()

	r := router.Routers()
	log.Println("Activity HTTP Service is running on port 13001...")
	r.Run(":13001")
}

// GRPC server
func startGRPCServer(wg *sync.WaitGroup) {
	defer wg.Done()

	listen, err := net.Listen("tcp", ":13002")
	if err != nil {
		log.Fatalf("grpc failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNoticeServiceServer(grpcServer, &NoticeService{})
	log.Println("Activity GRPC Service is running on port 13002...")
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
