package queue

import (
	"context"
	"errors"
	"github.com/aiechoic/services/internal/encoding"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisQueue[T any] struct {
	c   *redis.Client
	s   encoding.Serializer
	key string
}

func NewRedisQueue[T any](c *redis.Client, s encoding.Serializer, key string) *RedisQueue[T] {
	return &RedisQueue[T]{
		c:   c,
		s:   s,
		key: key,
	}
}

func (r *RedisQueue[T]) Push(ctx context.Context, message *T, expire time.Duration) error {
	ex := NewExpiringMessage(message, expire)
	data, err := r.s.Serialize(ex)
	if err != nil {
		return err
	}
	return r.c.LPush(ctx, r.key, data).Err()
}

func (r *RedisQueue[T]) Pop(ctx context.Context) (*T, error) {
	timeout := 15 * time.Second
	result, err := r.c.BRPop(ctx, timeout, r.key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	var ex ExpiringMessage[T]
	err = r.s.Deserialize([]byte(result[1]), &ex)
	if err != nil {
		return nil, err
	}
	if ex.Expired() {
		return r.Pop(ctx)
	}
	return ex.Data, nil
}

func (r *RedisQueue[T]) Len(ctx context.Context) (int64, error) {
	return r.c.LLen(ctx, r.key).Result()
}
