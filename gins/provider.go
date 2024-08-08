package gins

import (
	"github.com/aiechoic/services/ioc"
)

type ConfigSection string

var providers = ioc.NewProviders[*Server]()

func GetServerByConfig(configSection ConfigSection, c *ioc.Container) *Server {
	pvd := providers.GetProvider(string(configSection), func(c *ioc.Container) (*Server, error) {
		var cfg Config
		err := c.UnmarshalConfig(string(configSection), &cfg, defaultConfigData)
		if err != nil {
			return nil, err
		}
		return cfg.NewServer(), nil
	})
	return pvd.MustGet(c)
}

func GetServer(c *ioc.Container) *Server {
	return GetServerByConfig(DefaultConfigSection, c)
}
