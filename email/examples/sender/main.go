package main

import (
	"context"
	"github.com/aiechoic/services/email"
	"github.com/aiechoic/services/ioc"
	"log"
)

func main() {

	c := ioc.NewContainer()
	defer c.Close()

	err := c.LoadConfig("./configs", ioc.ConfigEnvTest)
	if err != nil {
		panic(err)
	}

	sender := email.GetSender(c)

	ctx := context.Background()

	// start checking email message queue, and send email
	sender.Run(ctx, func(err error) {
		log.Println(err)
	})
}
