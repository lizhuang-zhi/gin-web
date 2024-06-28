package mongodb

import (
	"booking-app/micro-service/cluster/common"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient() {
	// 连接到 MongoDB
	mongoClient, err := connectMongoDB()
	if err != nil {
		panic(err)
	}
	common.MongoClient = mongoClient
}

func connectMongoDB() (*mongo.Client, error) {
	// 设置 MongoDB 客户端选项
	clientOptions := options.Client().ApplyURI(common.Config.MongoDB.Host)

	// 连接到 MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// 检查连接
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
