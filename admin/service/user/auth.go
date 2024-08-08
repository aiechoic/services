package user

import (
	"github.com/aiechoic/services/admin/internal/auth"
	"github.com/aiechoic/services/admin/internal/model"
	"github.com/aiechoic/services/database/gorm"
	"github.com/aiechoic/services/database/redis"
	"github.com/aiechoic/services/ioc"
	"github.com/aiechoic/services/openapi"
)

var SecurityScheme = "UserAuth"

var SecuritySchemes = openapi.SecuritySchemes{
	SecurityScheme: &openapi.SecurityScheme{
		Type: "apiKey",
		In:   "header",
		Name: SecurityScheme,
	},
}

func NewAuth(c *ioc.Container) *auth.Auth[User] {
	authDB, err := model.NewAuthDB[User](gorm.GetGormDB(c), redis.GetRedis(c))
	if err != nil {
		panic(err)
	}
	return auth.NewAuth[User](authDB, SecurityScheme, SecurityScheme)
}
