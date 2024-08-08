package healthy

import "github.com/aiechoic/services/ioc/healthy"

var defaultConfigData = []byte(
	`# Health Check Config

# Note: this file will be watched and reloaded 
# automatically, so you can change the config 
# without restarting the service.

# NotifyLevel is the level of error that will be notified
# Available levels are: "debug", "info", "warn", "error", "fatal"
notify_level: "debug"
`)

type Config struct {
	NotifyLevel healthy.Level `mapstructure:"notify_level"`
}
