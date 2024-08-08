package queue

import (
	"context"
	"time"
)

type MemoryQueue[T any] struct {
	messages []*ExpiringMessage[T]
}

func NewMemoryQueue[T any]() *MemoryQueue[T] {
	return &MemoryQueue[T]{}
}

func (e *MemoryQueue[T]) Push(ctx context.Context, message *T, expire time.Duration) error {
	e.messages = append(e.messages, NewExpiringMessage(message, expire))
	return nil
}

func (e *MemoryQueue[T]) Pop(ctx context.Context) (*T, error) {
	if len(e.messages) == 0 {
		return nil, nil
	}
	message := e.messages[0]
	e.messages = e.messages[1:]
	if message.Expired() {
		return e.Pop(ctx)
	}
	return message.Data, nil
}

func (e *MemoryQueue[T]) Len(ctx context.Context) (int64, error) {
	return int64(len(e.messages)), nil
}
