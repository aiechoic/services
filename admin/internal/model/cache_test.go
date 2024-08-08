package model_test

import (
	"context"
	"github.com/aiechoic/services/admin/internal/model"
	"github.com/aiechoic/services/database/redis"
	"github.com/aiechoic/services/internal/encoding"
	"github.com/aiechoic/services/ioc"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestStruct struct {
	Field1 string
	Field2 int
}

func setupTestCache(t *testing.T) (*model.Caches[TestStruct], func()) {
	c := ioc.NewContainer()
	err := c.LoadConfig(configPath, ioc.ConfigEnvTest)
	if err != nil {
		t.Fatal(err)
	}
	rds := redis.GetRedis(c)
	cache := model.NewCaches[TestStruct](rds, "test_table", encoding.JSONSerializer)
	deferFunc := func() {
		err = cache.DropTable(context.Background())
		assert.NoError(t, err)
	}
	return cache, deferFunc
}

func TestCaches_Set(t *testing.T) {
	cache, closer := setupTestCache(t)
	defer closer()
	testValue := &TestStruct{
		Field1: "example",
		Field2: 42,
	}

	ctx := context.Background()

	err := cache.Set(ctx, "column", "value", testValue)
	assert.NoError(t, err)

	result, err := cache.Get(ctx, "column", "value")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testValue.Field1, result.Field1)
	assert.Equal(t, testValue.Field2, result.Field2)
}

func TestCaches_Drop(t *testing.T) {
	cache, closer := setupTestCache(t)
	defer closer()
	testValue := &TestStruct{
		Field1: "example",
		Field2: 42,
	}
	ctx := context.Background()
	err := cache.Set(ctx, "column", "value", testValue)
	assert.NoError(t, err)

	err = cache.Drop(ctx, "column", "value")
	assert.NoError(t, err)

	result, err := cache.Get(ctx, "column", "value")
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestCaches_DropColumn(t *testing.T) {
	cache, closer := setupTestCache(t)
	defer closer()
	testValue := &TestStruct{
		Field1: "example",
		Field2: 42,
	}
	ctx := context.Background()
	err := cache.Set(ctx, "column1", "value1", testValue)
	assert.NoError(t, err)
	err = cache.Set(ctx, "column1", "value2", testValue)
	assert.NoError(t, err)
	err = cache.Set(ctx, "column2", "value1", testValue)
	assert.NoError(t, err)

	err = cache.DropColumn(ctx, "column1")
	assert.NoError(t, err)

	result1, err := cache.Get(ctx, "column1", "value1")
	assert.NoError(t, err)
	assert.Nil(t, result1)

	result2, err := cache.Get(ctx, "column1", "value2")
	assert.NoError(t, err)
	assert.Nil(t, result2)

	result3, err := cache.Get(ctx, "column2", "value1")
	assert.NoError(t, err)
	assert.NotNil(t, result3)
}

func TestCaches_DropTable(t *testing.T) {
	cache, closer := setupTestCache(t)
	defer closer()
	testValue := &TestStruct{
		Field1: "example",
		Field2: 42,
	}
	ctx := context.Background()
	err := cache.Set(ctx, "column1", "value1", testValue)
	assert.NoError(t, err)
	err = cache.Set(ctx, "column2", "value2", testValue)
	assert.NoError(t, err)

	err = cache.DropTable(ctx)
	assert.NoError(t, err)

	result1, err := cache.Get(ctx, "column1", "value1")
	assert.NoError(t, err)
	assert.Nil(t, result1)

	result2, err := cache.Get(ctx, "column2", "value2")
	assert.NoError(t, err)
	assert.Nil(t, result2)
}

func TestCaches_CountTable(t *testing.T) {
	cache, closer := setupTestCache(t)
	defer closer()
	testValue := &TestStruct{
		Field1: "example",
		Field2: 42,
	}
	ctx := context.Background()
	err := cache.Set(ctx, "column1", "value1", testValue)
	assert.NoError(t, err)
	err = cache.Set(ctx, "column2", "value2", testValue)
	assert.NoError(t, err)

	count, err := cache.CountTable(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

func TestCaches_CountColumn(t *testing.T) {
	cache, closer := setupTestCache(t)
	defer closer()
	testValue := &TestStruct{
		Field1: "example",
		Field2: 42,
	}
	ctx := context.Background()
	err := cache.Set(ctx, "column1", "value1", testValue)
	assert.NoError(t, err)
	err = cache.Set(ctx, "column1", "value2", testValue)
	assert.NoError(t, err)
	err = cache.Set(ctx, "column2", "value1", testValue)
	assert.NoError(t, err)

	count, err := cache.CountColumn(ctx, "column1")
	assert.NoError(t, err)
	assert.Equal(t, 2, count)

	count, err = cache.CountColumn(ctx, "column2")
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}
