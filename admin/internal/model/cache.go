package model

import (
	"context"
	"errors"
	"github.com/aiechoic/services/internal/encoding"
	"github.com/redis/go-redis/v9"
	"strings"
)

type Caches[T any] struct {
	rds        *redis.Client
	table      string
	serializer encoding.Serializer
}

func NewCaches[T any](rds *redis.Client, table string, serializer encoding.Serializer) *Caches[T] {
	return &Caches[T]{rds: rds, table: table, serializer: serializer}
}

func (c *Caches[T]) field(uniqueColumn, value string) string {
	return uniqueColumn + ":" + value
}

func (c *Caches[T]) Set(ctx context.Context, uniqueColumn, value string, v *T) error {
	field := c.field(uniqueColumn, value)
	data, err := c.serializer.Serialize(v)
	if err != nil {
		return err
	}
	return c.rds.HSet(ctx, c.table, field, data).Err()
}

func (c *Caches[T]) Get(ctx context.Context, uniqueColumn, value string) (*T, error) {
	field := c.field(uniqueColumn, value)
	data, err := c.rds.HGet(ctx, c.table, field).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	var v T
	err = c.serializer.Deserialize(data, &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (c *Caches[T]) Drop(ctx context.Context, uniqueColumn, value string) error {
	field := c.field(uniqueColumn, value)
	err := c.rds.HDel(ctx, c.table, field).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return err
	}
	return nil
}

func (c *Caches[T]) DropColumn(ctx context.Context, uniqueColumn string) error {
	fields, err := c.rds.HKeys(ctx, c.table).Result()
	if err != nil {
		return err
	}
	for _, field := range fields {
		if strings.HasPrefix(field, uniqueColumn+":") {
			if err := c.rds.HDel(ctx, c.table, field).Err(); err != nil {
				return err
			}
		}
	}
	return nil
}

// DropTable drops the table, used for testing
func (c *Caches[T]) DropTable(ctx context.Context) error {
	err := c.rds.Del(ctx, c.table).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return err
	}
	return nil
}

func (c *Caches[T]) CountTable(ctx context.Context) (int64, error) {
	return c.rds.HLen(ctx, c.table).Result()
}

func (c *Caches[T]) CountColumn(ctx context.Context, uniqueColumn string) (int, error) {
	fields, err := c.rds.HKeys(ctx, c.table).Result()
	if err != nil {
		return 0, err
	}
	count := 0
	for _, field := range fields {
		if strings.HasPrefix(field, uniqueColumn+":") {
			count++
		}
	}
	return count, nil
}
