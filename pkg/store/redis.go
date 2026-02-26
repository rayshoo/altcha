package store

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisStore(url string, ttlMinutes int) (*RedisStore, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisStore{
		client: client,
		ttl:    time.Duration(ttlMinutes) * time.Minute,
	}, nil
}

func (s *RedisStore) Exists(token string) (bool, error) {
	n, err := s.client.Exists(context.Background(), "altcha:"+token).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (s *RedisStore) Add(token string) error {
	return s.client.Set(context.Background(), "altcha:"+token, "1", s.ttl).Err()
}

func (s *RedisStore) Close() error { return s.client.Close() }
