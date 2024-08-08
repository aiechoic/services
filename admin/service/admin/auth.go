package admin

import (
	"github.com/aiechoic/services/admin/internal/auth"
	"github.com/aiechoic/services/admin/internal/model"
	"github.com/aiechoic/services/database/gorm"
	"github.com/aiechoic/services/database/redis"
	"github.com/aiechoic/services/ioc"
	"github.com/aiechoic/services/openapi"
)

var SecurityScheme = "AdminAuth"

var SecuritySchemes = openapi.SecuritySchemes{
	SecurityScheme: &openapi.SecurityScheme{
		Type: "apiKey",
		In:   "header",
		Name: SecurityScheme,
	},
	// Add more security schemes here
}

func NewAuth(c *ioc.Container) *auth.Auth[Admin] {
	gormDB := gorm.GetGormDB(c)
	rds := redis.GetRedis(c)
	authDB, err := model.NewAuthDB[Admin](gormDB, rds)
	if err != nil {
		panic(err)
	}
	return auth.NewAuth[Admin](authDB, SecurityScheme, SecurityScheme)
}
