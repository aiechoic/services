package config

import (
	"github.com/aiechoic/services/ioc"
)

var defaultConfigData = []byte(
	`# Amin console config

# Email of default admin
email: "your_email"
# Password of default admin
password: "your_password"
`)

var DefaultSection Section = "amin-console"

type Section string

type Admin struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

var Provider = ioc.NewProvider(func(c *ioc.Container) (*Admin, error) {
	var cfg Admin
	err := c.UnmarshalConfig(string(DefaultSection), &cfg, defaultConfigData)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
})

func GetAdmin(c *ioc.Container) (*Admin, error) {
	return Provider.Get(c)
}
