package model

import (
	"booking-app/micro-service/cluster/common/core"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// 公告数据
type Notice struct {
	ID      int    `json:"id" bson:"id"`
	Title   string `json:"title" bson:"title"`
	SubType int    `json:"sub_type" bson:"sub_type"`
	Content string `json:"content" bson:"content"`
}

// 内存中的公告数据
var GlobalNotice = make(map[int]*Notice)

// 查询公告数据
func QueryNoticeData(c *gin.Context, opts *core.Options) ([]*Notice, error) {
	collection := opts.MongoClient.Database("micro-service-activity").Collection("notice")
	cursor, err := collection.Find(c.Request.Context(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(c.Request.Context())

	var notices []*Notice
	for cursor.Next(c.Request.Context()) {
		var notice Notice
		err := cursor.Decode(&notice)
		if err != nil {
			return nil, err
		}
		notices = append(notices, &notice)
	}

	return notices, nil
}

// 新增公告数据
func InsertNoticeData(c *gin.Context, data *Notice, opts *core.Options) error {
	collection := opts.MongoClient.Database("micro-service-activity").Collection("notice")
	document := bson.M{
		"id":       data.ID,
		"title":    data.Title,
		"sub_type": data.SubType,
		"content":  data.Content,
	}
	_, err := collection.InsertOne(c.Request.Context(), document)
	if err != nil {
		return err
	}

	return nil
}
