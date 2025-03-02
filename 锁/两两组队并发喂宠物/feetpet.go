package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

const (
	totalTeams    = 25 // 共25组队伍
	feedPerUser   = 20 // 每人投喂20次
	redisPoolSize = 50 // Redis连接池大小
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

// 初始化Redis连接
func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: redisPoolSize, // 重要：确保连接池足够
	})
}

// 创建队伍宠物数据
func createPets() {
	for i := 1; i <= totalTeams; i++ {
		key := getPetKey(i)
		if err := rdb.Set(ctx, key, 0, 0).Err(); err != nil {
			panic(fmt.Sprintf("初始化宠物失败: %v", err))
		}
	}
}

// 投喂操作
func feedPet(teamID, userID int, wg *sync.WaitGroup) {
	defer wg.Done()

	petKey := getPetKey(teamID)
	logKey := getLogKey(teamID)

	for i := 0; i < feedPerUser; i++ {
		// 原子操作增加养成值
		newVal, err := rdb.IncrBy(ctx, petKey, 1).Result()
		if err != nil {
			fmt.Printf("队伍%d-用户%d投喂失败: %v\n", teamID, userID, err)
			continue
		}

		// 记录日志（异步提升性能）
		go func(iteration int) {
			logMsg := fmt.Sprintf("用户%d第%d次投喂，当前值:%d",
				userID, iteration+1, newVal)
			rdb.RPush(ctx, logKey, logMsg)
		}(i)
	}
}

func getPetKey(teamID int) string {
	return fmt.Sprintf("pet:team%d:value", teamID)
}

func getLogKey(teamID int) string {
	return fmt.Sprintf("pet:team%d:logs", teamID)
}

func main() {
	initRedis()
	createPets()

	var wg sync.WaitGroup

	// 为每个队伍启动两个协程（两人组队）
	for teamID := 1; teamID <= totalTeams; teamID++ {
		wg.Add(2)
		go feedPet(teamID, 1, &wg) // 用户1
		go feedPet(teamID, 2, &wg) // 用户2
	}

	wg.Wait()
	fmt.Println("所有投喂操作完成！")

	// 验证结果
	for i := 1; i <= totalTeams; i++ {
		val, _ := rdb.Get(ctx, getPetKey(i)).Int()
		fmt.Printf("队伍%d最终养成值: %d\n", i, val)
	}
}
