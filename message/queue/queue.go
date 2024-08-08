package queue

import (
	"context"
	"time"
)

type Queue[T any] interface {
	Push(ctx context.Context, message *T, expire time.Duration) error
	Pop(ctx context.Context) (*T, error)
	Len(ctx context.Context) (int64, error)
}
