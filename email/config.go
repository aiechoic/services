package email

import (
	"fmt"
	"github.com/aiechoic/services/ioc/healthy"
)

var defaultSenderConfigData = []byte(
	`# email sender config

# Note: this file will be watched and reloaded 
# automatically, so you can change the config 
# without restarting the service.

# title of the email
title: "Your Company"
# email address of the sender
from: "sender@example.com"
# smtp host
host: "smtp.example.com"
# smtp port
port: "465"
# smtp password
password: "password"
# email template file
template: "template/email/verify.gohtml"
# email subject
subject: "Email Verification"
`)

var defaultPusherConfigData = []byte(
	`# email pusher config

# Note: this file will be watched and reloaded 
# automatically, so you can change the config 
# without restarting the service.

# healthy check levels of email queue length
levels:
  fatal: 50
  error: 20
  warn: 10
  info: 5
  #debug: 0
`)

type PusherConfig struct {
	Levels map[healthy.Level]int64 `mapstructure:"levels"`
}

func (c *PusherConfig) GetError(queueLength int64) *healthy.Error {
	var levels = []healthy.Level{healthy.LFatal, healthy.LError, healthy.LWarn, healthy.LInfo, healthy.LDebug}
	for _, level := range levels {
		if n, ok := c.Levels[level]; ok && queueLength >= n {
			return &healthy.Error{
				Level: level,
				Msg:   fmt.Sprintf("email queue length: %d", queueLength),
			}
		}
	}
	return nil
}

type SenderConfig struct {
	Title    string `mapstructure:"title"`
	From     string `mapstructure:"from"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Template string `mapstructure:"template"`
	Subject  string `mapstructure:"subject"`
}
