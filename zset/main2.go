package main

import (
	"booking-app/zset/zset"
	"fmt"
	"strconv"
)

func main() {
	gameRankService := zset.NewGameRankService("gin_web_game_rank")

	// 添加玩家排行榜
	for i := 1; i <= 10000; i++ {
		player := "player" + strconv.Itoa(i) // 每个玩家名都是唯一的
		score := i                           // 这里假设分数就是i，你可以按需要改变
		err := gameRankService.AddOrUpdatePlayer(player, score)
		if err != nil {
			fmt.Printf("Error when adding player: %v\n", err)
			return
		}
	}
}
