package admin

import (
	"github.com/aiechoic/services/admin/internal/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Admin struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Email    string `json:"email" gorm:"uniqueIndex,size:50"`
	Password string `json:"-" gorm:"size:60"`
	Token    string `json:"token" gorm:"uniqueIndex,size:64"`
}

type DB struct {
	*model.AuthDB[Admin]
}

func NewDB(db *gorm.DB, rds *redis.Client) (*DB, error) {
	authDB, err := model.NewAuthDB[Admin](db, rds)
	if err != nil {
		return nil, err
	}
	return &DB{
		AuthDB: authDB,
	}, nil
}

func (a *Admin) GetEmail() string {
	return a.Email
}

func (a *Admin) SetEmail(email string) {
	a.Email = email
}

func (a *Admin) SetToken(token string) {
	a.Token = token
}

func (a *Admin) GetToken() string {
	return a.Token
}

func (a *Admin) GetHashedPassword() string {
	return a.Password
}

func (a *Admin) SetHashedPassword(password string) {
	a.Password = password
}
