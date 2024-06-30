package rpc

import (
	pb "booking-app/micro-service/protobuf/gen-pb"
	"context"
	"log"
	"time"
)

// 获取公告通过ID
func GetNotice(client pb.NoticeServiceClient) {
	resp, err := client.GetNotice(context.Background(), &pb.GetNoticeRequest{Id: 1})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Response: %v", resp)
}

// 获取所有公告
func GetNotices(client pb.NoticeServiceClient) {
	// 获取所有公告
	resp, err := client.GetNotices(context.Background(), &pb.GetNoticesRequest{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Response: %v", resp)
}

// 修改公告通过ID
func UpdateNoticeById(client pb.NoticeServiceClient) {
	resp, err := client.UpdateNoticeById(context.Background(), &pb.UpdateNoticeRequest{
		Id:      1,
		Title:   "更新公告",
		Content: "版本更新日志",
		SubType: 5,
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Response: %v", resp)
}

// 启动广播接收
func BroadcastPlayerNotify(broadcaseClient pb.BroadcastServiceClient) {
	// 启动广播接收
	stream, err := broadcaseClient.BroadcastPlayerNotify(context.Background(), &pb.BroadcastRequest{})
	if err != nil {
		log.Fatalf("could not start broadcast: %v", err)
	}

	for {
		// 接收广播消息
		broadcast, err := stream.Recv()
		if err != nil {
			log.Fatalf("could not receive broadcast: %v", err)
		}

		log.Println("Received broadcast: ", broadcast.OsTypes)

		broadcastOsTypes := broadcast.OsTypes

		// 过滤广播消息
		if len(broadcastOsTypes) > 0 {
			for _, osType := range broadcastOsTypes {
				if osType == "android" {
					log.Printf("Received broadcast: %v", broadcast.Packets[0].Message)
				}

				if osType == "ios" {
					log.Printf("Received broadcast: %v", broadcast.Packets[0].Message)
				}
			}
		}

		time.Sleep(1 * time.Second)
	}
}

// 获取用户通过ID
func GetUser(client pb.UserServiceClient) {
	resp, err := client.GetUser(context.Background(), &pb.GetUserRequest{Id: "12"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Response: %v", resp)
}

// 创建用户
func CreateUser(client pb.UserServiceClient) {
	newUser := &pb.User{
		Id:    "12",
		Name:  "Jim",
		Email: "19223@qq.com",
	}

	resp, err := client.CreateUser(context.Background(), &pb.CreateUserRequest{
		User: newUser,
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Response: %v", resp)
}
