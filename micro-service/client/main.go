package main

import (
	"log"

	"booking-app/micro-service/client/rpc"
	pb "booking-app/micro-service/protobuf/gen-pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 执行rpc名称
var ExcuteRpcName = map[string]bool{
	// 公告
	"GetNotice":        false,
	"GetNotices":       false,
	"UpdateNoticeById": false,

	// 广播
	"BroadcastPlayerNotify": true,
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:13002", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close() // 关闭连接

	// 建立连接
	client := rpc.GetActivityNoticeClient(conn)
	broadcaseClient := rpc.GetActivityBroadcastClient(conn)

	// 映射各rpc方法
	var excuteNoticeMap = map[string]func(client pb.NoticeServiceClient){
		"GetNotice":        rpc.GetNotice,
		"GetNotices":       rpc.GetNotices,
		"UpdateNoticeById": rpc.UpdateNoticeById,
	}
	var excuteBroadcastMap = map[string]func(client pb.BroadcastServiceClient){
		"BroadcastPlayerNotify": rpc.BroadcastPlayerNotify,
	}

	// 执行rpc方法
	for key, excute := range excuteNoticeMap {
		if !ExcuteRpcName[key] {
			continue
		}
		excute(client)
	}
	for key, excute := range excuteBroadcastMap {
		if !ExcuteRpcName[key] {
			continue
		}
		excute(broadcaseClient)
	}
}
