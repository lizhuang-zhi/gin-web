package zset

import (
	"github.com/go-redis/redis"
)

// GameRankService 代表游戏排行榜服务
type GameRankService struct {
	db      *DefaultDB
	rankKey string
}

// NewGameRankService 创建新的游戏排行榜服务实例
func NewGameRankService(rankKey string) *GameRankService {
	// 初始化Redis客户端
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // 如果没有的话，默认为空
		DB:       0,  // 默认为0号数据库
	})

	// 创建 defaultDB 实例
	db := &DefaultDB{client: redisClient}

	return &GameRankService{db: db, rankKey: rankKey}
}

// AddOrUpdatePlayer 添加或更新玩家分数
func (s *GameRankService) AddOrUpdatePlayer(player string, score int) error {
	redisIntCmd := s.db.ZAdd(s.rankKey, score, player)
	return redisIntCmd.Err()
}

// RemovePlayer 移除玩家
func (s *GameRankService) RemovePlayer(player string) error {
	redisIntCmd := s.db.ZRem(s.rankKey, player)
	return redisIntCmd.Err()
}

// GetPlayerRank 获取玩家排名
func (s *GameRankService) GetPlayerRank(player string) (int, error) {
	return s.db.Do("ZREVRANK", s.rankKey, player).Int()
}

// GetTopPlayers 获取前n名玩家
func (s *GameRankService) GetTopPlayers(n int) ([]string, error) {
	result, err := s.db.client.ZRevRange(s.rankKey, 0, int64(n-1)).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}
