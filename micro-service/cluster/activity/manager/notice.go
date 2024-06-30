package manager

import (
	"booking-app/micro-service/cluster/activity/model"
	"booking-app/micro-service/cluster/common"
	"booking-app/micro-service/cluster/common/commandpb"
	pb "booking-app/micro-service/protobuf/gen-pb"
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
)

// TODO: *pb.GetNoticeRequest和 *pb.GetNoticeResponse 还没有抽象,所以目前只有一个rpc注册
type MsgHandler func(context.Context, *pb.GetNoticeRequest) (*pb.GetNoticeResponse, error)

type NoticeService struct {
	pb.UnimplementedNoticeServiceServer

	sync.RWMutex
	regHandler map[commandpb.Command]MsgHandler
}

func NewNoticeService() *NoticeService {
	return &NoticeService{
		regHandler: make(map[commandpb.Command]MsgHandler),
	}
}

func (s *NoticeService) RegisterHandler(comd commandpb.Command, msgHandler MsgHandler) {
	s.Lock()
	defer s.Unlock()

	s.regHandler[comd] = msgHandler
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

func (s *NoticeService) GetNotices(ctx context.Context, req *pb.GetNoticesRequest) (*pb.GetNoticesResponse, error) {
	collection := common.MongoClient.Database("micro-service-activity").Collection("notice")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notices []*pb.Notice
	for cursor.Next(ctx) {
		var notice pb.Notice
		err := cursor.Decode(&notice)
		if err != nil {
			return nil, err
		}
		notices = append(notices, &notice)
	}

	return &pb.GetNoticesResponse{
		Notices: notices,
	}, nil
}

func (s *NoticeService) UpdateNoticeById(ctx context.Context, req *pb.UpdateNoticeRequest) (*pb.UpdateNoticeResponse, error) {
	notice_id := req.Id
	notice_title := req.Title
	notice_content := req.Content
	notice_sub_type := req.SubType

	for id, notice := range model.GlobalNotice {
		int64_id := int64(id)
		if int64_id == notice_id {
			notice.Title = notice_title
			notice.Content = notice_content
			notice.SubType = int(notice_sub_type)
			return &pb.UpdateNoticeResponse{
				Success: true,
			}, nil
		}
	}

	return &pb.UpdateNoticeResponse{
		Success: false,
	}, nil
}
