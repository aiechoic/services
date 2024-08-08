package dingtalk

import (
	"github.com/aiechoic/services/ioc"
	"github.com/spf13/viper"
	"log"
)

const (
	ConfigKey = "dingtalk"
)

var defaultConfigData = []byte(
	`# DingTalk config

# Note: this file will be watched and reloaded 
# automatically, so you can change the config 
# without restarting the service.

# DingTalk access token
access_token: "your_access_token"
# DingTalk secret
secret: "your_secret"
`)

var providers = ioc.NewProviders[*Client]()

func GetProvider(configSection string) *ioc.Provider[*Client] {
	return providers.GetProvider(configSection, func(c *ioc.Container) (*Client, error) {
		client := NewClient()
		err := c.UnmarshalAndWatchConfig(configSection, defaultConfigData, func(v *viper.Viper) {
			var cfg Config
			err := v.Unmarshal(&cfg)
			if err != nil {
				log.Println(err)
			}
			client.UpdateConfig(&cfg)
		})
		if err != nil {
			return nil, err
		}
		return client, nil
	})
}

func GetClient(c *ioc.Container) *Client {
	return GetProvider(ConfigKey).MustGet(c)
}

type Config struct {
	AccessToken string `mapstructure:"access_token"`
	Secret      string `mapstructure:"secret"`
}
