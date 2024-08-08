package model_test

import (
	"context"
	"github.com/aiechoic/services/admin/internal/model"
	"github.com/aiechoic/services/database/gorm"
	"github.com/aiechoic/services/database/redis"
	"github.com/aiechoic/services/ioc"
	"github.com/stretchr/testify/assert"
	"testing"
)

var configPath = "../../../configs"

type testUser struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Email    string `json:"email" gorm:"uniqueIndex,size:50"`
	Password string `json:"password" gorm:"size:60"`
	Token    string `json:"token" gorm:"uniqueIndex,size:64"`
}

func (t *testUser) GetEmail() string {
	return t.Email
}

func (t *testUser) SetEmail(email string) {
	t.Email = email
}

func (t *testUser) SetToken(token string) {
	t.Token = token
}

func (t *testUser) GetToken() string {
	return t.Token
}

func (t *testUser) GetHashedPassword() string {
	return t.Password
}

func (t *testUser) SetHashedPassword(password string) {
	t.Password = password
}

func setupTestAuthDB[T any](t *testing.T) (*model.AuthDB[T], func()) {
	c := ioc.NewContainer()
	err := c.LoadConfig(configPath, ioc.ConfigEnvTest)
	if err != nil {
		t.Fatal(err)
	}
	db := gorm.GetGormDB(c)
	rds := redis.GetRedis(c)
	userDB, err := model.NewAuthDB[T](db, rds)
	if err != nil {
		panic(err)
	}
	deferFunc := func() {
		err = userDB.DropTable()
		assert.NoError(t, err)
	}
	return userDB, deferFunc
}

func TestAuthDB_Create(t *testing.T) {
	userDB, closer := setupTestAuthDB[testUser](t)
	defer closer()
	user := &testUser{
		Email:    "test@example.com",
		Password: "password",
	}

	err := userDB.Create(user)
	assert.NoError(t, err)

	createdUser, err := userDB.GetByColumn("email", user.Email)
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, user.Email, createdUser.Email)
}

func TestAuthDB_GetByToken(t *testing.T) {
	userDB, closer := setupTestAuthDB[testUser](t)
	defer closer()

	password := "password"
	user := &testUser{
		Email:    "test@example.com",
		Password: password,
	}

	err := userDB.Create(user)
	assert.NoError(t, err)
	assert.Empty(t, user.Token)

	loginUser, err := userDB.Login(user.Email, password)
	assert.NoError(t, err)
	assert.NotEmpty(t, loginUser.Token)

	readUser, err := userDB.GetByToken(context.Background(), loginUser.Token)
	assert.NoError(t, err)
	assert.NotNil(t, readUser)
	assert.Equal(t, loginUser.Token, readUser.Token)
}

func TestAuthDB_Login(t *testing.T) {
	userDB, closer := setupTestAuthDB[testUser](t)
	defer closer()

	password := "password"
	user := &testUser{
		Email:    "test@example.com",
		Password: password,
	}

	err := userDB.Create(user)
	assert.NoError(t, err)

	loggedInUser, err := userDB.Login(user.Email, password)
	assert.NoError(t, err)
	assert.NotNil(t, loggedInUser)
	assert.NotEmpty(t, loggedInUser.Token)

	_, err = userDB.Login(user.Email, "wrong_password")
	assert.ErrorIs(t, err, model.ErrorIncorrectEmailOrPassword)
}

func TestAuthDB_Logout(t *testing.T) {
	userDB, closer := setupTestAuthDB[testUser](t)
	defer closer()

	password := "password"
	user := &testUser{
		Email:    "test@example.com",
		Password: password,
	}

	err := userDB.Create(user)
	assert.NoError(t, err)

	user, err = userDB.Login(user.Email, password)
	assert.NoError(t, err)

	loginUser, err := userDB.GetByToken(context.Background(), user.Token)
	assert.NoError(t, err)
	assert.NotNil(t, loginUser)

	err = userDB.Logout(user.Token)
	assert.NoError(t, err)

	loggedOutUser, err := userDB.GetByToken(context.Background(), user.Token)
	assert.NoError(t, err)
	assert.Nil(t, loggedOutUser)
}

func TestAuthDB_ResetPassword(t *testing.T) {
	userDB, closer := setupTestAuthDB[testUser](t)
	defer closer()

	// Create an user
	user := &testUser{
		Email:    "test@example.com",
		Password: "old_password",
	}
	err := userDB.Create(user)
	assert.NoError(t, err)

	// Reset the password
	newPassword := "new_password"
	err = userDB.ResetPassword(user.Email, newPassword)
	assert.NoError(t, err)

	// Verify the password has been updated
	updatedUser, err := userDB.Login(user.Email, newPassword)
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)

	// Verify the old password no longer works
	_, err = userDB.Login(user.Email, "old_password")
	assert.ErrorIs(t, err, model.ErrorIncorrectEmailOrPassword)
}

func TestAuthDB_Save(t *testing.T) {
	userDB, closer := setupTestAuthDB[testUser](t)
	defer closer()

	// Create a user
	user := &testUser{
		Email:    "test@example.com",
		Password: "password",
	}
	err := userDB.Create(user)
	assert.NoError(t, err)

	// Login to generate token and cache
	loggedInUser, err := userDB.Login(user.Email, "password")
	assert.NoError(t, err)
	assert.NotNil(t, loggedInUser)

	// Read by token to get old data
	oldUser, err := userDB.GetByToken(context.Background(), loggedInUser.Token)
	assert.NoError(t, err)
	assert.NotNil(t, oldUser)

	// Update user data and save
	oldUser.Password = "new_password"
	updatedUser, err := userDB.Save(oldUser)
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)

	// Read by token to get new data
	newUser, err := userDB.GetByToken(context.Background(), loggedInUser.Token)
	assert.NoError(t, err)
	assert.NotNil(t, newUser)

	// Verify the cache is cleared and new data is cached
	assert.Equal(t, oldUser.Password, newUser.Password)
}
