package redis

import (
	"fmt"
	"github.com/aiechoic/services/ioc"
	"github.com/redis/go-redis/v9"
)

const (
	DefaultConfigSection ConfigSection = "redis"
)

var providers = ioc.NewProviders[*redis.Client]()

type ConfigSection string

func GetProvider(section ConfigSection) *ioc.Provider[*redis.Client] {
	return providers.GetProvider(string(section), func(c *ioc.Container) (*redis.Client, error) {
		var cfg Config
		err := c.UnmarshalConfig(string(section), &cfg, defaultConfigData)
		if err != nil {
			return nil, err
		}
		client, err := cfg.NewClient()
		if err != nil {
			return nil, err
		}
		c.OnClose(client.Close)
		return client, nil
	})
}

func GetRedis(c *ioc.Container) *redis.Client {
	config := DefaultConfigSection
	client, err := GetProvider(config).Get(c)
	if err != nil {
		panic(fmt.Errorf("get redis client \"%s\": %w", config, err))
	}
	return client
}
