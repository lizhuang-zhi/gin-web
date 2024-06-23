package main

import (
	"context"
	"log"

	pb "booking-app/micro-service/protobuf/gen-pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:13002", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close() // 关闭连接

	// 建立连接
	client := getActivityNoticeClient(conn)

	resp, err := client.GetNotice(context.Background(), &pb.GetNoticeRequest{Id: 1})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Response: %v", resp)
}

func getActivityNoticeClient(conn *grpc.ClientConn) pb.NoticeServiceClient {
	return pb.NewNoticeServiceClient(conn)
}
