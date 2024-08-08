package rate

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLimiter_Set(t *testing.T) {
	storage := NewMemoryStorage[string](time.Minute)
	rls := NewRateLimiter[string](storage, time.Second, 1)

	v := "value1"

	err := rls.Set("key1", &v, time.Second*2)
	assert.NoError(t, err)

	data, err := storage.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", *data)
}

func TestLimiter_Get(t *testing.T) {
	storage := NewMemoryStorage[string](time.Minute)
	rls := NewRateLimiter[string](storage, time.Second, 1)
	v := "value1"
	err := rls.Set("key1", &v, time.Second*2)
	assert.NoError(t, err)

	data, err := rls.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", *data)
}

func TestLimiter_Del(t *testing.T) {
	storage := NewMemoryStorage[string](time.Minute)
	rls := NewRateLimiter[string](storage, time.Second, 1)
	v := "value1"
	err := storage.Set("key1", &v, time.Second*2)
	assert.NoError(t, err)

	err = rls.Del("key1")
	assert.NoError(t, err)

	data, err := storage.Get("key1")
	assert.NoError(t, err)
	assert.Nil(t, data)
}

func TestLimiter_RateLimit(t *testing.T) {
	storage := NewMemoryStorage[string](time.Minute)
	rls := NewRateLimiter[string](storage, time.Second, 1)
	v := "value1"
	err := rls.Set("key1", &v, time.Second*2)
	assert.NoError(t, err)
	v2 := "value2"
	err = rls.Set("key1", &v2, time.Second*2)
	assert.Error(t, err)
	assert.Equal(t, ErrLimitExceed, err)
}

func TestLimiter_GetWaiteTime(t *testing.T) {
	storage := NewMemoryStorage[string](time.Minute)
	rls := NewRateLimiter[string](storage, time.Second, 1)

	// Set a value in Limiter
	v := "value1"
	err := rls.Set("key1", &v, time.Second*2)
	assert.NoError(t, err)

	// Get the wait time
	waitTime := rls.GetWaiteTime("key1")
	need := time.Second
	assert.Equal(t, waitTime, need)

	time.Sleep(500 * time.Millisecond)
	// Get the wait time
	waitTime = rls.GetWaiteTime("key1")
	round := 50 * time.Millisecond
	waitTime = waitTime/round*round + round
	need = 500 * time.Millisecond
	assert.Equal(t, waitTime, need)
}
