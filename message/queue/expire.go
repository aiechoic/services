package queue

import "time"

type ExpiringMessage[T any] struct {
	Data      *T    `json:"d"`
	ExpiresAt int64 `json:"ex"`
}

func NewExpiringMessage[T any](data *T, expiration time.Duration) *ExpiringMessage[T] {
	var ex int64
	if expiration > 0 {
		ex = time.Now().Add(expiration).Unix()
	} else {
		ex = -1
	}
	return &ExpiringMessage[T]{
		Data:      data,
		ExpiresAt: ex,
	}
}

func (e *ExpiringMessage[T]) Expired() bool {
	if e.ExpiresAt < 0 {
		return false
	}
	return e.ExpiresAt < time.Now().Unix()
}
