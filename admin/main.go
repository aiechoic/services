package main

import (
	"context"
	"github.com/aiechoic/services/admin/internal/healthy"
	"github.com/aiechoic/services/admin/service/admin"
	"github.com/aiechoic/services/admin/service/email"
	"github.com/aiechoic/services/admin/service/user"
	email2 "github.com/aiechoic/services/email"
	"github.com/aiechoic/services/gins"
	"github.com/aiechoic/services/gins/docs/redoc"
	"github.com/aiechoic/services/gins/docs/swagger"
	"github.com/aiechoic/services/ioc"
	"log"
)

func main() {
	configDir := "configs"
	env := ioc.ConfigEnvTest

	c := ioc.NewContainer()
	defer c.Close()
	err := c.LoadConfig(configDir, env)
	if err != nil {
		log.Fatal(err)
	}
	err = healthy.CheckHealthy(c)
	if err != nil {
		log.Fatal(err)
	}

	server := gins.GetServer(c)

	userService := user.NewService(c)
	adminService := admin.NewService(c)
	emailService := email.NewService(c)

	server.SetSecuritySchemes(user.SecuritySchemes)
	server.SetSecuritySchemes(admin.SecuritySchemes)

	server.Register(
		emailService,
		userService,
		adminService,
	)

	redoc.ServeAPI(server)
	swagger.ServeAPI(server)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sender := email2.GetSender(c)

	// start checking email message queue, and send email
	go func() {
		sender.Run(ctx, func(err error) {
			log.Println(err)
		})
	}()

	server.Run(ctx)
}
