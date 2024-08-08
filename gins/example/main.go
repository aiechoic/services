package main

import (
	"context"
	"github.com/aiechoic/services/gins"
	"github.com/aiechoic/services/gins/docs/redoc"
	"github.com/aiechoic/services/gins/docs/swagger"
	"github.com/aiechoic/services/gins/example/user"
	"github.com/aiechoic/services/ioc"
)

func main() {
	secret := "secret"
	c := ioc.NewContainer()
	defer c.Close()
	err := c.LoadConfig("./configs", ioc.ConfigEnvTest)
	if err != nil {
		panic(err)
	}

	server := gins.GetServer(c)

	userService := user.NewService(secret)

	server.SetSecuritySchemes(user.SecuritySchemes)

	server.Register(userService)

	redoc.ServeAPI(server)

	swagger.ServeAPI(server)

	server.Run(context.Background())
}
