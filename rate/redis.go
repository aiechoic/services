package rate

import (
	"context"
	"errors"
	"github.com/aiechoic/services/internal/encoding"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisStorage[T any] struct {
	client     *redis.Client
	serializer encoding.Serializer
	key        string
}

func (r *RedisStorage[T]) Set(name string, data *T, expire time.Duration) error {
	ctx := context.Background()
	val, err := r.serializer.Serialize(data)
	if err != nil {
		return err
	}
	err = r.client.Set(ctx, r.key+name, val, expire).Err()
	return err
}

func (r *RedisStorage[T]) Get(name string) (data *T, err error) {
	ctx := context.Background()
	val, err := r.client.Get(ctx, r.key+name).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	var v T
	err = r.serializer.Deserialize(val, &v)
	if err != nil {
		return nil, err
	}
	return &v, err
}

func (r *RedisStorage[T]) Del(name string) error {
	ctx := context.Background()
	err := r.client.Del(ctx, r.key+name).Err()
	return err
}

func NewRedisStorage[T any](client *redis.Client, serializer encoding.Serializer, key string) *RedisStorage[T] {
	return &RedisStorage[T]{
		client:     client,
		serializer: serializer,
		key:        key,
	}
}
