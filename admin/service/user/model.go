package user

import (
	"errors"
	"github.com/aiechoic/services/admin/internal/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	ErrorUserEmailNotExist = errors.New("user email not exist")
)

type User struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	Email     string `json:"email" gorm:"uniqueIndex,size:50"`
	Password  string `json:"-" gorm:"size:60"`
	Token     string `json:"token" gorm:"uniqueIndex,size:64"`
	VipExpire int64  `json:"vip_expire"`
}

type DB struct {
	*model.AuthDB[User]
}

func NewDB(db *gorm.DB, rds *redis.Client) (*DB, error) {
	authDB, err := model.NewAuthDB[User](db, rds)
	if err != nil {
		return nil, err
	}
	return &DB{
		AuthDB: authDB,
	}, nil
}

func (a *DB) SetVipExpire(email string, expire int64) (*User, error) {
	user, err := a.GetByColumn("email", email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrorUserEmailNotExist
	}
	user.VipExpire = expire
	// Save method is from AuthDB, will automatically remove token cache
	return a.Save(user)
}

func (a *User) GetEmail() string {
	return a.Email
}

func (a *User) SetEmail(email string) {
	a.Email = email
}

func (a *User) SetToken(token string) {
	a.Token = token
}

func (a *User) GetToken() string {
	return a.Token
}

func (a *User) GetHashedPassword() string {
	return a.Password
}

func (a *User) SetHashedPassword(password string) {
	a.Password = password
}
