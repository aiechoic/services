package redis_test

import (
	"context"
	"fmt"
	"github.com/aiechoic/services/database/redis"
	"github.com/aiechoic/services/ioc"
)

func ExampleGetRedis() {
	c := ioc.NewContainer()
	defer c.Close()

	err := c.LoadConfig("../../configs", ioc.ConfigEnvTest)
	if err != nil {
		panic(err)
	}

	client := redis.GetRedis(c)

	err = client.Ping(context.Background()).Err()

	fmt.Printf("%v\n", err == nil)

	// Output: true
}
