package rate

import (
	"github.com/aiechoic/services/database/redis"
	"github.com/aiechoic/services/internal/encoding"
	"github.com/aiechoic/services/ioc"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func setupRedis[T any](c *ioc.Container, t *testing.T) *RedisStorage[T] {
	err := c.LoadConfig("../configs", ioc.ConfigEnvTest)
	if err != nil {
		t.Fatal(err)
	}
	rds := redis.GetRedis(c)
	return NewRedisStorage[T](rds, encoding.JSONSerializer, "test:")
}

func TestRedisStorage_Set(t *testing.T) {
	container := ioc.NewContainer()
	storage := setupRedis[string](container, t)

	v := "value1"
	err := storage.Set("key1", &v, time.Second*2)
	assert.NoError(t, err)
}

func TestRedisStorage_Get(t *testing.T) {
	container := ioc.NewContainer()
	storage := setupRedis[string](container, t)

	v := "value1"
	err := storage.Set("key1", &v, time.Second*2)
	assert.NoError(t, err)

	data, err := storage.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", *data)
}

func TestRedisStorage_Del(t *testing.T) {
	container := ioc.NewContainer()
	storage := setupRedis[string](container, t)

	v := "value1"
	err := storage.Set("key1", &v, time.Second*2)
	assert.NoError(t, err)

	err = storage.Del("key1")
	assert.NoError(t, err)

	data, err := storage.Get("key1")
	assert.NoError(t, err)
	assert.Nil(t, data)
}
