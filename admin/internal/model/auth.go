package model

import (
	"context"
	"errors"
	"github.com/aiechoic/services/internal/encoding"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var Serializer encoding.Serializer = encoding.GobSerializer

var (
	ErrorIncorrectEmailOrPassword = errors.New("incorrect email or password")
	ErrorEmailAlreadyExists       = errors.New("email already exists")
	ErrorEmailNotRegistered       = errors.New("email not registered")
)

type UserEntity interface {
	GetEmail() string
	SetEmail(email string)
	SetToken(token string)
	GetToken() string
	GetHashedPassword() string
	SetHashedPassword(password string)
}

type AuthDB[T any] struct {
	*BaseDB[T]
	cache *Caches[T]
}

func NewAuthDB[T any](db *gorm.DB, cache *redis.Client) (*AuthDB[T], error) {
	bdb, err := NewBaseDB[T](db)
	if err != nil {
		return nil, err
	}
	table, err := GetTableName[T](db)
	if err != nil {
		return nil, err
	}
	c := NewCaches[T](cache, table, Serializer)
	return &AuthDB[T]{BaseDB: bdb, cache: c}, nil
}

func (a *AuthDB[T]) newUser() (*T, UserEntity) {
	var u T
	var user = &u
	return user, any(user).(UserEntity)
}

func (a *AuthDB[T]) Create(user *T) error {
	userInter := any(user).(UserEntity)
	if userInter.GetEmail() == "" {
		return errors.New("email is empty")
	}
	if userInter.GetHashedPassword() == "" {
		return errors.New("password is empty")
	}
	if _user, _ := a.BaseDB.GetByColumn("email", userInter.GetEmail()); _user != nil {
		return ErrorEmailAlreadyExists
	}
	pass, err := bcrypt.GenerateFromPassword([]byte(userInter.GetHashedPassword()), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	userInter.SetHashedPassword(string(pass))
	userInter.SetToken("")
	return a.BaseDB.Create(user)
}

func (a *AuthDB[T]) GetByToken(ctx context.Context, token string) (*T, error) {
	if token == "" {
		return nil, errors.New("token is empty")
	}
	user, err := a.cache.Get(ctx, "token", token)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}
	user, err = a.BaseDB.GetByColumn("token", token)
	if err != nil {
		return nil, err
	}
	if user != nil {
		err = a.cache.Set(ctx, "token", token, user)
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}

func (a *AuthDB[T]) Login(email, password string) (*T, error) {
	user, err := a.BaseDB.GetByColumn("email", email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrorIncorrectEmailOrPassword
	}
	userInter := any(user).(UserEntity)
	if bcrypt.CompareHashAndPassword([]byte(userInter.GetHashedPassword()), []byte(password)) != nil {
		return nil, ErrorIncorrectEmailOrPassword
	}
	token := userInter.GetToken()
	if token != "" {
		// drop old token from cache
		err = a.cache.Drop(context.Background(), "token", token)
		if err != nil {
			return nil, err
		}
	}
	token = uuid.NewString()
	user, userInter = a.newUser()
	userInter.SetToken(token)
	user, err = a.BaseDB.PartialUpdateByColumn(user, "email", email)
	if err != nil {
		return nil, err
	}
	// set new token to cache
	err = a.cache.Set(context.Background(), "token", token, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (a *AuthDB[T]) Logout(token string) error {
	if token == "" {
		return nil
	}
	t, err := a.GetByColumn("token", token)
	if err != nil {
		return nil
	}
	if t == nil {
		return nil
	}
	userInter := any(t).(UserEntity)
	userInter.SetToken("")
	_, err = a.Save(t)
	if err != nil {
		return err
	}
	err = a.cache.Drop(context.Background(), "token", token)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return err
	}
	return nil
}

func (a *AuthDB[T]) ResetPassword(email, password string) error {
	user, err := a.BaseDB.GetByColumn("email", email)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrorEmailNotRegistered
	}
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	userInter := any(user).(UserEntity)
	err = a.cache.Drop(context.Background(), "token", userInter.GetToken())
	if err != nil {
		return err
	}
	userInter.SetHashedPassword(string(pass))
	userInter.SetToken("")
	_, err = a.BaseDB.FullUpdateByColumn(user, "email", email)
	return err
}

// Save saves the entity and clears the cache
func (a *AuthDB[T]) Save(entity *T) (*T, error) {
	userInter := any(entity).(UserEntity)
	email := userInter.GetEmail()
	var err error
	entity, err = a.BaseDB.FullUpdateByColumn(entity, "email", email)
	if err != nil {
		return nil, err
	}
	token := userInter.GetToken()
	if token != "" {
		err = a.cache.Drop(context.Background(), "token", token)
		if err != nil {
			return nil, err
		}
	}
	return entity, nil
}

// DropTable drops the table and cache, used for testing, be careful to use it
func (a *AuthDB[T]) DropTable() error {
	err := a.BaseDB.DropTable()
	if err != nil {
		return err
	}
	return a.cache.DropTable(context.Background())
}
