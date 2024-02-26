package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Credentials struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
}

type Client struct {
	*redis.Client
}

type RedisClient interface {
	redis.Client
	Check(ctx context.Context) error
}

func (r *Client) Check(ctx context.Context) error {
	return r.Ping(ctx).Err()
}

func New(ctx context.Context, creds *Credentials) (*Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         creds.Host,
		Password:     creds.Password,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	})

	if ping := redisClient.Ping(ctx); ping.Err() != nil {
		return nil, ping.Err()
	}

	return &Client{
		redisClient,
	}, nil
}
