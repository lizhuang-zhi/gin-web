package main

import (
	"booking-app/zset/zset"
	"fmt"
)

func main() {
	// 添加玩家排行榜
	gameRankService := zset.NewGameRankService("gin_web_game_rank")
	gameRankService.AddOrUpdatePlayer("leo", 6000)
	gameRankService.AddOrUpdatePlayer("tim", 1100)
	gameRankService.AddOrUpdatePlayer("joe", 3100)
	gameRankService.AddOrUpdatePlayer("shuo", 3900)
	gameRankService.AddOrUpdatePlayer("chen", 1400)

	// 获取玩家排行榜
	rank, err := gameRankService.GetPlayerRank("player1223")
	if err != nil {
		panic(err)
	}
	fmt.Println("rank:", rank)

	// 获取前n名玩家
	topN, err := gameRankService.GetTopPlayers(23)
	if err != nil {
		panic(err)
	}
	fmt.Println("topN:", topN)
}
