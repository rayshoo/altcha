package store

import (
	"context"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client redis.Cmdable
	ttl    time.Duration
	closer func() error
}

func NewRedisStore(url string, cluster bool, ttlMinutes int) (*RedisStore, error) {
	ttl := time.Duration(ttlMinutes) * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if cluster || strings.Contains(url, ",") {
		return newClusterStore(ctx, url, ttl)
	}
	return newSingleStore(ctx, url, ttl)
}

func newSingleStore(ctx context.Context, url string, ttl time.Duration) (*RedisStore, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisStore{client: client, ttl: ttl, closer: client.Close}, nil
}

func newClusterStore(ctx context.Context, url string, ttl time.Duration) (*RedisStore, error) {
	var addrs []string
	for _, a := range strings.Split(url, ",") {
		a = strings.TrimSpace(a)
		a = strings.TrimPrefix(a, "redis://")
		addrs = append(addrs, a)
	}

	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: addrs,
	})
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisStore{client: client, ttl: ttl, closer: client.Close}, nil
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

func (s *RedisStore) Ping() error {
	return s.client.Ping(context.Background()).Err()
}

func (s *RedisStore) Close() error { return s.closer() }
