package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/redis/go-redis/v9"
)

// 从文件加载 Lua 脚本
var deductScript *redis.Script

func init() {
	// 获取脚本绝对路径
	scriptPath := filepath.Join("lua", "stock.lua")

	// 读取脚本内容
	scriptContent, err := ioutil.ReadFile(scriptPath)
	if err != nil {
		panic(fmt.Sprintf("加载Lua脚本失败: %v", err))
	}

	// 创建 Redis 脚本对象
	deductScript = redis.NewScript(string(scriptContent))
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx := context.Background()
	rdb.Set(ctx, "gin-web:stock", 100, 0)

	// 使用通用类型处理
	result, err := deductScript.Run(ctx, rdb, []string{"gin-web:stock"}, 5, 10).Result()
	if err != nil {
		// 处理真正的Redis错误
		fmt.Println("操作失败:", err)
		return
	}

	switch v := result.(type) {
	case int64:
		fmt.Printf("直接返回数值: %d\n", v)
	case []interface{}:
		if v[0].(int64) == 1 {
			fmt.Printf("剩余库存: %d\n", v[1].(int64))
		} else {
			fmt.Printf("当前库存%d, 库存不足", v[1].(int64))
		}
	case string:
		fmt.Println("错误信息:", v)
	default:
		fmt.Println("未知返回类型")
	}
}
