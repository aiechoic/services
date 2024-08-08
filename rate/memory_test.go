package rate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage_Set(t *testing.T) {
	storage := NewMemoryStorage[string](time.Minute)
	v := "value1"

	err := storage.Set("key1", &v, time.Second*2)
	assert.NoError(t, err)

	data, err := storage.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", *data)
}

func TestMemoryStorage_Get(t *testing.T) {
	storage := NewMemoryStorage[string](time.Minute)

	v := "value1"

	err := storage.Set("key1", &v, time.Second*2)
	assert.NoError(t, err)

	time.Sleep(time.Second * 3)
	data, err := storage.Get("key1")
	assert.NoError(t, err)
	assert.Nil(t, data)
}

func TestMemoryStorage_Del(t *testing.T) {
	storage := NewMemoryStorage[string](time.Minute)
	v := "value1"

	err := storage.Set("key1", &v, time.Second*2)
	assert.NoError(t, err)

	err = storage.Del("key1")
	assert.NoError(t, err)

	data, err := storage.Get("key1")
	assert.NoError(t, err)
	assert.Nil(t, data)
}

func TestMemoryStorage_Cleanup(t *testing.T) {
	storage := NewMemoryStorage[string](time.Second * 1)

	v := "value1"

	err := storage.Set("key1", &v, time.Second*1)
	assert.NoError(t, err)

	time.Sleep(time.Second * 2)
	data, err := storage.Get("key1")
	assert.NoError(t, err)
	assert.Nil(t, data)
}
