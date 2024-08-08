package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"net"
	"os"
)

var defaultConfigData = []byte(
	`# Redis configuration file
#
# configure "redis.conf" for add user and password:
# user yourusername on +@all ~* >somepassword

# use debug mode
debug: true
# Redis server
host: "localhost"
# Redis port
port: 6379
# Redis database
db: 9
# Username for Redis server
username: ""
# Password for Redis server
password: ""
`)

type Config struct {
	Debug    bool   `mapstructure:"debug"`
	Host     string `mapstructure:"host"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func (r *Config) NewClient() (*redis.Client, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", r.Host, 6379),
		Password: r.Password,
		Username: r.Username,
		DB:       r.DB,
	})
	err := c.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	if r.Debug {
		c.AddHook(redisHook{
			logger: log.New(os.Stderr, "redis: ", log.LstdFlags|log.Llongfile),
		})
	}
	return c, nil
}

type redisHook struct {
	logger *log.Logger
}

func (h redisHook) printf(format string, v ...interface{}) {
	_ = h.logger.Output(2, fmt.Sprintf(format, v...))
}

func (h redisHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		h.logger.Printf("Dialing to %s:%s\n", network, addr)
		return next(ctx, network, addr)
	}
}

func (h redisHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		err := next(ctx, cmd)
		h.printf("command: %s, err: %v\n", cmd, err)
		return err
	}
}

func (h redisHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		err := next(ctx, cmds)
		h.printf("pipeline: %v, err: %v\n", cmds, err)
		return err
	}
}
