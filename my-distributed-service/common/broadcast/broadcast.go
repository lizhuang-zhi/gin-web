package broadcast

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type Broadcaster struct {
	client *redis.Client
}

func NewBroadcaster(client *redis.Client) *Broadcaster {
	return &Broadcaster{
		client: client,
	}
}

func (b *Broadcaster) Publish(channel string, message interface{}) error {
	return b.client.Publish(ctx, channel, message).Err()
}

func (b *Broadcaster) Subscribe(channel string) *redis.PubSub {
	return b.client.Subscribe(ctx, channel)
}
