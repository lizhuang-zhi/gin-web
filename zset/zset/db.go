package zset

import (
	"github.com/go-redis/redis"
)

type DefaultDB struct {
	client *redis.Client
}

// ZAdd wraps the ZADD command.
func (db *DefaultDB) ZAdd(key string, score int, member string) *redis.IntCmd {
	return db.client.ZAdd(key, redis.Z{Score: float64(score), Member: member})
}

// ZRem wraps the ZREM command.
func (db *DefaultDB) ZRem(key string, members ...string) *redis.IntCmd {
	return db.client.ZRem(key, members)
}

// Do sends a command to the Redis server and returns the received reply.
func (db *DefaultDB) Do(args ...interface{}) *redis.Cmd {
	return db.client.Do(args...)
}
