package user

import (
	"context"
	"github.com/aiechoic/services/database/gorm"
	"github.com/aiechoic/services/database/redis"
	"github.com/aiechoic/services/ioc"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupTestDB(t *testing.T) (*DB, func(), error) {
	c := ioc.NewContainer()
	err := c.LoadConfig("../../../configs", ioc.ConfigEnvTest)
	if err != nil {
		t.Fatal(err)
	}
	db := gorm.GetGormDB(c)
	rds := redis.GetRedis(c)
	userDB, err := NewDB(db, rds)
	if err != nil {
		return nil, nil, err
	}
	deferFunc := func() {
		err = userDB.DropTable()
		assert.NoError(t, err)
	}
	return userDB, deferFunc, nil
}

func TestDB_SetVipExpire(t *testing.T) {
	userDB, deferFunc, err := setupTestDB(t)
	defer deferFunc()
	assert.NoError(t, err)

	// Create a user
	user := &User{
		Email:    "test@example.com",
		Password: "password",
		Token:    "test_token",
	}
	err = userDB.Create(user)
	assert.NoError(t, err)

	// Login to generate token and cache
	loggedInUser, err := userDB.Login(user.Email, "password")
	assert.NoError(t, err)
	assert.NotNil(t, loggedInUser)

	// Read by token to get old data
	oldUser, err := userDB.GetByToken(context.Background(), loggedInUser.Token)
	assert.NoError(t, err)
	assert.NotNil(t, oldUser)

	// Set VIP expire, drop old data by token from cache
	newExpire := int64(1234567890)
	updatedUser, err := userDB.SetVipExpire(user.Email, newExpire)
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.Equal(t, newExpire, updatedUser.VipExpire)

	// Read by token to get new data
	newUser, err := userDB.GetByToken(context.Background(), loggedInUser.Token)
	assert.NoError(t, err)
	assert.NotNil(t, newUser)

	// Compare old and new VipExpire
	assert.NotEqual(t, oldUser.VipExpire, newUser.VipExpire)
	assert.Equal(t, newExpire, newUser.VipExpire)
}
