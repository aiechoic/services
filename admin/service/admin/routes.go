package admin

import (
	"github.com/aiechoic/services/admin/config"
	"github.com/aiechoic/services/admin/service/user"
	"github.com/aiechoic/services/database/gorm"
	"github.com/aiechoic/services/database/redis"
	"github.com/aiechoic/services/email/verify"
	"github.com/aiechoic/services/gins"
	"github.com/aiechoic/services/ioc"
)

func registerSuperAdmin(c *ioc.Container, adminDB *DB) {
	total, err := adminDB.Count()
	if err != nil {
		panic(err)
	}
	if total == 0 {
		cfg := config.Provider.MustGet(c)
		err = adminDB.Create(&Admin{
			Email:    cfg.Email,
			Password: cfg.Password,
		})
		if err != nil {
			panic(err)
		}
	}
}

func NewService(c *ioc.Container) *gins.Service {
	es := verify.GetGenerator(c)
	db := gorm.GetGormDB(c)
	rds := redis.GetRedis(c)
	adminDB, err := NewDB(db, rds)
	if err != nil {
		panic(err)
	}
	registerSuperAdmin(c, adminDB)
	userDB, err := user.NewDB(db, rds)
	if err != nil {
		panic(err)
	}
	auth := NewAuth(c)
	a := NewHandlers(es, adminDB, userDB, auth)
	return &gins.Service{
		Tag:         "Admin",
		Description: "Admin service",
		Path:        "/admin",
		Routes: []gins.Route{
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
			gins.Route{
				Method:  "POST",
				Path:    "/count_user",
				Handler: a.CountUser(),
			}.Use(auth),
			gins.Route{
				Method:  "POST",
				Path:    "/find_users",
				Handler: a.FindUsers(),
			}.Use(auth),
			gins.Route{
				Method:  "POST",
				Path:    "/set_user_vip_time",
				Handler: a.SetUserVipTime(),
			}.Use(auth),
		},
	}
}
