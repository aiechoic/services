package email

import (
	"context"
	"github.com/aiechoic/services/database/redis"
	"github.com/aiechoic/services/internal/encoding"
	"github.com/aiechoic/services/ioc"
	"github.com/aiechoic/services/ioc/healthy"
	"github.com/aiechoic/services/message/queue"
	"github.com/spf13/viper"
	"log"
	"sync"
)

const (
	DefaultSenderConfigSection SenderConfigSection = "email-sender"
	DefaultPusherConfigSection PusherConfigSection = "email-pusher"
	DefaultRedisQueueKey                           = "email-verify"
)

var MsgSerializer = encoding.JSONSerializer

var (
	queueProviders  = ioc.NewProviders[queue.Queue[Msg]]()
	senderProviders = ioc.NewProviders[*Sender]()
	pusherProviders = ioc.NewProviders[*Pusher]()
)

type SenderConfigSection string
type PusherConfigSection string
type RedisQueueKey string

func getRedisQueueProvider(redisConfig redis.ConfigSection, redisKey RedisQueueKey) *ioc.Provider[queue.Queue[Msg]] {
	return queueProviders.GetProvider(string(redisKey), func(c *ioc.Container) (queue.Queue[Msg], error) {
		var rds = redis.GetProvider(redisConfig).MustGet(c)
		return queue.NewRedisQueue[Msg](rds, MsgSerializer, string(redisKey)), nil
	})
}

func GetSenderProvider(
	senderConfig SenderConfigSection,
	redisConfig redis.ConfigSection,
	redisKey RedisQueueKey,
) *ioc.Provider[*Sender] {
	return senderProviders.GetProvider(string(senderConfig), func(c *ioc.Container) (*Sender, error) {
		rdsQueue := getRedisQueueProvider(redisConfig, redisKey).MustGet(c)
		sender := NewSender(rdsQueue)
		err := c.UnmarshalAndWatchConfig(string(senderConfig), defaultSenderConfigData, func(v *viper.Viper) {
			var cfg SenderConfig
			err := v.Unmarshal(&cfg)
			if err != nil {
				log.Println(err)
				return
			}
			err = sender.UpdateConfig(&cfg)
			if err != nil {
				log.Printf("update sender config \"%s\"error: %v\n", senderConfig, err)
				return
			}
			log.Printf("loaded sender config \"%s\"success\n", senderConfig)
		})
		if err != nil {
			return nil, err
		}
		return sender, nil
	})
}

func GetSender(c *ioc.Container) *Sender {
	return GetSenderProvider(
		DefaultSenderConfigSection,
		redis.DefaultConfigSection,
		DefaultRedisQueueKey,
	).MustGet(c)
}

func GetPusherProvider(
	pusherConfig PusherConfigSection,
	redisConfig redis.ConfigSection,
	redisKey RedisQueueKey,
) *ioc.Provider[*Pusher] {
	return pusherProviders.GetProvider(string(pusherConfig), func(c *ioc.Container) (*Pusher, error) {
		var rdsQueue = getRedisQueueProvider(redisConfig, redisKey).MustGet(c)
		var cfg = &PusherConfig{}
		var mu = &sync.Mutex{}
		err := c.UnmarshalAndWatchConfig(string(pusherConfig), defaultPusherConfigData, func(v *viper.Viper) {
			mu.Lock()
			defer mu.Unlock()
			err := v.Unmarshal(cfg)
			if err != nil {
				log.Println(err)
				return
			}
			log.Printf("loaded pusher config \"%s\"success\n", pusherConfig)
		})
		if err != nil {
			return nil, err
		}
		c.OnHealthCheck(func() *healthy.Error {
			mu.Lock()
			defer mu.Unlock()
			n, err := rdsQueue.Len(context.Background())
			if err != nil {
				return &healthy.Error{
					Level: healthy.LError,
					Msg:   err.Error(),
				}
			}
			return cfg.GetError(n)
		})
		return NewPusher(rdsQueue), nil
	})
}

func GetPusher(c *ioc.Container) *Pusher {
	return GetPusherProvider(
		DefaultPusherConfigSection,
		redis.DefaultConfigSection,
		DefaultRedisQueueKey,
	).MustGet(c)
}
