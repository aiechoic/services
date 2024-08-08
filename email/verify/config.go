package verify

import (
	"time"
)

var defaultConfigData = []byte(
	`# email verification code generator config

# Note: this file will be watched and reloaded 
# automatically, so you can change the config 
# without restarting the service.

# random chars used to generate code
random_chars: "0123456789"
# length of the generated code
code_length: 6
# cache expire time in seconds
cache_expire_in_seconds: 60
# rate limit every n times in m seconds
rate_limit_every_in_seconds: m
rate_limit_allow_n: n
`)

var defaultTemplateData = []byte(
	`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Verification Code</title>
</head>
<body>
    <p>Please enter the code to verify your account.</p>
    <p><strong>{{.code}}</strong></p>
    <p>The code will expire in {{.expireIn}} minutes.</p>
</body>
</html>`)

type Config struct {
	RandomChars             string `mapstructure:"random_chars"`
	CodeLength              int    `mapstructure:"code_length"`
	CacheExpireInSeconds    int    `mapstructure:"cache_expire_in_seconds"`
	RateLimitEveryInSeconds int    `mapstructure:"rate_limit_every_in_seconds"`
	RateLimitAllowN         int    `mapstructure:"rate_limit_allow_n"`
}

func (c *Config) ToOptions() *GeneratorOptions {
	return &GeneratorOptions{
		RandomChars:     c.RandomChars,
		CodeLength:      c.CodeLength,
		CacheExpire:     time.Duration(c.CacheExpireInSeconds) * time.Second,
		RateLimitEvery:  time.Duration(c.RateLimitEveryInSeconds) * time.Second,
		RateLimitAllowN: c.RateLimitAllowN,
	}
}
