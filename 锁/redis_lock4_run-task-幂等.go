package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

/*
	我这里模拟的场景: 分布式环境下, 同时删除数据库中的同一份id数据

	实际上，多pod同时执行相同的任务，是否需要分布式锁，取决于任务的性质。
	如果任务是幂等的，也就是说无论执行多少次，产生的结果都是一样的，那么这个任务就不需要分布式锁。
	举个例子，比如数据库的清理任务，它的目标是删除一天前的所有日志记录。无论这个任务执行一次还是执行多次，结果都是一样的：一天前的所有日志记录都被删除了。因此，这个任务可以在多个pod中同时执行，而不需要使用分布式锁。
	然而，如果任务不是幂等的，就可能需要使用分布式锁了。
	比如说，你有一个任务是向用户发送电子邮件。如果这个任务在多个pod中同时执行，用户就会收到多封相同的电子邮件。这显然不是我们想要的结果。为了防止这种情况发生，我们就需要在发送电子邮件的任务开始时获取一个分布式锁，
	确保一次只有一个pod可以执行这个任务。
	又比如，你有一个任务是从共享队列中取出一项并处理。如果多个pod同时从队列中取出相同的项，就可能会处理两次。为了防止这种情况，我们就需要使用分布式锁来保证一次只有一个pod可以取出和处理队列中的项。
	这就是为什么在某些情况下，我们需要在多pod环境中使用分布式锁。当然，实际情况可能会更复杂，你需要根据你的具体需求和业务逻辑来判断是否需要使用分布式锁。
*/

var mongoClient *mongo.Client

func main() {
	mongoClient = NewMongoClient()

	/*
		添加数据
	*/
	// InsertFunc()

	/*
		删除数据: 验证分布式环境下的任务执行
	*/
	DeleteFunc()

	time.Sleep(2 * time.Second)
}

func InsertFunc() {
	addData()
}

func addData() {
	// 添加数据
	collection := mongoClient.Database("distributed-lock").Collection("sync")
	_, err := collection.InsertOne(context.Background(), bson.M{"id": 1, "title": "test1", "content": "content1"})
	if err != nil {
		fmt.Printf("添加数据失败: %v", err)
		return
	}
	fmt.Println("添加数据成功")
}

func DeleteFunc() {
	// 启动3个goroutine执行任务, 模拟分布式环境下的任务执行(可以简单理解为三个pod,同时执行相同的任务)
	for i := 0; i < 3; i++ {
		go task(i + 1)
	}
}

// 任务
func task(id int) {
	fmt.Printf("任务[%v]执行开始>>\n", id)

	deleteID := 1 // 删除id为1的数据

	collection := mongoClient.Database("distributed-lock").Collection("sync")
	result, err := collection.DeleteOne(context.Background(), bson.M{"id": deleteID})
	if err != nil {
		fmt.Printf("任务[%v], 查询数据失败: %v", id, err)
		return
	}
	fmt.Printf("任务[%v], 删除数据成功, 删除数量: %d\n", id, result.DeletedCount)

	fmt.Printf("任务[%v]执行结束!!\n", id)
}

func NewMongoClient() *mongo.Client {
	// 连接到 MongoDB
	mongoClient, err := connectMongoDB()
	if err != nil {
		panic(err)
	}
	return mongoClient
}

func connectMongoDB() (*mongo.Client, error) {
	// 设置 MongoDB 客户端选项
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")

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
