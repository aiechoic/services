package rate

import (
	"errors"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var ErrLimitExceed = errors.New("rate limit exceed")

type Storage[T any] interface {
	Set(name string, data *T, expire time.Duration) error
	Get(name string) (data *T, err error)
	Del(name string) error
}

type Limiter[T any] struct {
	storage       Storage[T]
	memory        *MemoryStorage[rate.Limiter]
	mu            sync.Mutex
	rate          rate.Limit
	limiterExpire time.Duration
	burst         int
}

func NewRateLimiter[T any](storage Storage[T], every time.Duration, allowN int) *Limiter[T] {
	rls := &Limiter[T]{
		storage:       storage,
		memory:        NewMemoryStorage[rate.Limiter](every),
		rate:          rate.Every(every / time.Duration(allowN)),
		limiterExpire: every,
		burst:         allowN,
	}
	return rls
}

func (r *Limiter[T]) getOrCreateLimiter(name string) (*rate.Limiter, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	limiter, err := r.memory.Get(name)
	if err != nil {
		return nil, err
	}
	if limiter != nil {
		return limiter, nil
	}
	limiter = rate.NewLimiter(r.rate, r.burst)
	err = r.memory.Set(name, limiter, r.limiterExpire)
	if err != nil {
		return nil, err
	}
	return limiter, nil
}

func (r *Limiter[T]) GetWaiteTime(name string) time.Duration {
	r.mu.Lock()
	defer r.mu.Unlock()
	limiter, err := r.memory.Get(name)
	if err != nil {
		return 0
	}
	if limiter == nil {
		return 0
	}
	reservation := limiter.Reserve()
	defer reservation.Cancel()
	return reservation.Delay()
}

func (r *Limiter[T]) Set(name string, data *T, expire time.Duration) error {
	limiter, err := r.getOrCreateLimiter(name)
	if err != nil {
		return err
	}
	if !limiter.Allow() {
		return ErrLimitExceed
	}
	err = r.storage.Set(name, data, expire)
	if err != nil {
		return err
	}
	return nil
}

func (r *Limiter[T]) Get(name string) (data *T, err error) {
	return r.storage.Get(name)
}

func (r *Limiter[T]) Del(name string) error {
	err := r.DelData(name)
	if err != nil {
		return err
	}
	err = r.DelLimiter(name)
	return err
}

func (r *Limiter[T]) DelData(name string) error {
	err := r.storage.Del(name)
	return err
}

func (r *Limiter[T]) DelLimiter(name string) error {
	err := r.memory.Del(name)
	return err
}
