package verify

import (
	"github.com/aiechoic/services/database/redis"
	"github.com/aiechoic/services/email"
	"github.com/aiechoic/services/encoding"
	"github.com/aiechoic/services/ioc"
	"github.com/aiechoic/services/rate"
	"github.com/spf13/viper"
	"log"
	"sync"
)

var DefaultConfigSection ConfigSection = "verify-code-generator"

var DefaultRedisCodeKey RedisCodeKey = "verify-code:"

type ConfigSection string

type RedisCodeKey string

var generatorProviders = ioc.NewProviders[*Generator]()

func GetGeneratorProvider(
	configSection ConfigSection,
	pusherConfig email.PusherConfigSection,
	redisConfig redis.ConfigSection,
	redisPusherKey email.RedisQueueKey,
	redisCodeKey RedisCodeKey,
) *ioc.Provider[*Generator] {
	return generatorProviders.GetProvider(string(configSection), func(c *ioc.Container) (*Generator, error) {
		pusher := email.GetPusherProvider(pusherConfig, redisConfig, redisPusherKey).MustGet(c)
		rds := redis.GetProvider(redisConfig).MustGet(c)
		storage := rate.NewRedisStorage[string](rds, encoding.JSONSerializer, string(redisCodeKey))
		generator := NewGenerator(storage, pusher)
		mu := &sync.Mutex{}
		err := c.UnmarshalAndWatchConfig(string(configSection), defaultConfigData, func(v *viper.Viper) {
			mu.Lock()
			defer mu.Unlock()
			var cfg Config
			err := v.Unmarshal(&cfg)
			if err != nil {
				log.Println(err)
				return
			}
			err = generator.UpdateConfig(cfg.ToOptions())
			if err != nil {
				log.Printf("update generator config \"%s\"error: %v\n", configSection, err)
				return
			}
			log.Printf("loaded generator config \"%s\"success\n", configSection)
		})
		if err != nil {
			return nil, err
		}
		return generator, nil
	})
}

func GetGenerator(c *ioc.Container) *Generator {
	return GetGeneratorProvider(
		DefaultConfigSection,
		email.DefaultPusherConfigSection,
		redis.DefaultConfigSection,
		email.DefaultRedisQueueKey,
		DefaultRedisCodeKey,
	).MustGet(c)
}
