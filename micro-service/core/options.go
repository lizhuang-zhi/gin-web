package core

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Options struct {
	Logger      *log.Logger // 日志组件
	MongoClient *mongo.Client
}

func NewOptions() *Options {
	// 创建日志组件
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	// 连接到 MongoDB
	mongoClient, err := connectMongoDB()
	if err != nil {
		logger.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	return &Options{
		Logger:      logger,
		MongoClient: mongoClient,
	}
}

func connectMongoDB() (*mongo.Client, error) {
	// 设置 MongoDB 客户端选项
	clientOptions := options.Client().ApplyURI(Config.MongoDB.Host)

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
