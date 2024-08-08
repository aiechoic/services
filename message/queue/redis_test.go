package queue_test

import (
	"context"
	"github.com/aiechoic/services/database/redis"
	"github.com/aiechoic/services/internal/encoding"
	"github.com/aiechoic/services/ioc"
	"github.com/aiechoic/services/message/queue"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func setupRedisQueue(key string, t *testing.T) *queue.RedisQueue[string] {
	c := ioc.NewContainer()
	err := c.LoadConfig("../../configs", ioc.ConfigEnvTest)
	if err != nil {
		t.Fatal(err)
	}
	client := redis.GetRedis(c)
	return queue.NewRedisQueue[string](client, encoding.JSONSerializer, key)
}

func TestRedisQueue_Pop(t *testing.T) {
	key := "test_key"
	expire := time.Second
	redisQueue := setupRedisQueue(key, t)

	message := "test_message"
	err := redisQueue.Push(context.Background(), &message, expire)
	assert.NoError(t, err)
	result, err := redisQueue.Pop(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, &message, result)

	err = redisQueue.Push(context.Background(), &message, expire)
	assert.NoError(t, err)
	err = redisQueue.Push(context.Background(), &message, expire)
	assert.NoError(t, err)
	// Wait for the message to expire
	time.Sleep(2 * time.Second)

	result, err = redisQueue.Pop(context.Background())
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestRedisQueue_Len(t *testing.T) {
	key := "test_key"
	expire := time.Minute
	redisQueue := setupRedisQueue(key, t)
	message := "test_message"
	err := redisQueue.Push(context.Background(), &message, expire)
	assert.Nil(t, err)
	length, err := redisQueue.Len(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, int64(1), length)
}
