package user

import (
	"github.com/aiechoic/services/database/gorm"
	"github.com/aiechoic/services/database/redis"
	"github.com/aiechoic/services/email/verify"
	"github.com/aiechoic/services/gins"
	"github.com/aiechoic/services/ioc"
)

func NewService(c *ioc.Container) *gins.Service {
	es := verify.GetGenerator(c)
	db := gorm.GetGormDB(c)
	rds := redis.GetRedis(c)
	userDB, err := NewDB(db, rds)
	if err != nil {
		panic(err)
	}
	auth := NewAuth(c)
	a := NewHandlers(es, userDB, auth)
	return &gins.Service{
		Tag:         "User",
		Description: "User service",
		Path:        "/user",
		Routes: []gins.Route{
			{
				Method:  "POST",
				Path:    "/register",
				Handler: a.Register(),
			},
			{
				Method:  "POST",
				Path:    "/login",
				Handler: a.Login(),
			},
			gins.Route{
				Method:  "POST",
				Path:    "/logout",
				Handler: a.Logout(),
			}.Use(auth),
			gins.Route{
				Method:  "POST",
				Path:    "/reset_password",
				Handler: a.ResetPassword(),
			},
			gins.Route{
				Method:  "POST",
				Path:    "/profile",
				Handler: a.GetProfile(),
			}.Use(auth),
		},
	}
}
