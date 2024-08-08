package main

import (
	"context"
	"github.com/aiechoic/services/email"
	"github.com/aiechoic/services/ioc"
	"log"
	"time"
)

// Before running this example, you need to start a message queue server.
// run the email/examples/sender/main.go first
func main() {
	yourEmail := ""

	c := ioc.NewContainer()
	defer c.Close()

	err := c.LoadConfig("./configs", ioc.ConfigEnvTest)
	if err != nil {
		panic(err)
	}

	pusher := email.GetPusher(c)

	// push a message to email message queue
	err = pusher.Push(context.Background(), yourEmail, map[string]interface{}{
		"code":     123456,
		"expireIn": 5,
	}, 5*time.Second)

	log.Fatal(err)
}
