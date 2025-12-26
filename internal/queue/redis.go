package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	redis *redis.Client
}

func New(addr string) *Client {
	return &Client{redis: redis.NewClient(&redis.Options{Addr: addr})}
}

func (c *Client) Close() error {
	return c.redis.Close()
}

func (c *Client) EnsureGroup(ctx context.Context, stream, group string) error {
	if err := c.redis.XGroupCreateMkStream(ctx, stream, group, "$" ).Err(); err != nil {
		if err.Error() == "BUSYGROUP Consumer Group name already exists" {
			return nil
		}
		return fmt.Errorf("create group: %w", err)
	}
	return nil
}

func (c *Client) AddJob(ctx context.Context, stream string, values map[string]any) (string, error) {
	id, err := c.redis.XAdd(ctx, &redis.XAddArgs{Stream: stream, Values: values}).Result()
	if err != nil {
		return "", fmt.Errorf("xadd: %w", err)
	}
	return id, nil
}

func (c *Client) ReadGroup(ctx context.Context, stream, group, consumer string, count int64, block time.Duration) ([]redis.XMessage, error) {
	res, err := c.redis.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: consumer,
		Streams:  []string{stream, ">"},
		Count:    count,
		Block:    block,
	}).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("xreadgroup: %w", err)
	}
	if len(res) == 0 {
		return nil, nil
	}
	return res[0].Messages, nil
}

func (c *Client) Ack(ctx context.Context, stream, group string, ids ...string) error {
	if len(ids) == 0 {
		return nil
	}
	if _, err := c.redis.XAck(ctx, stream, group, ids...).Result(); err != nil {
		return fmt.Errorf("xack: %w", err)
	}
	return nil
}
