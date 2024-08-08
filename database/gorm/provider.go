package gorm

import (
	"github.com/aiechoic/services/ioc"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"log"
)

const DefaultConfigSection = "gorm"

var providers = ioc.NewProviders[*gorm.DB]()

type ConfigSection string

func GetProvider(section ConfigSection) *ioc.Provider[*gorm.DB] {
	return providers.GetProvider(string(section), func(c *ioc.Container) (*gorm.DB, error) {
		var cfg Config
		err := c.UnmarshalConfig(string(section), &cfg, defaultConfigData)
		if err != nil {
			return nil, err
		}
		gdb, closer := cfg.Connect()
		c.OnClose(closer)

		err = c.WatchConfig(string(section), func(v *viper.Viper) {
			var newCfg Config
			err := v.Unmarshal(&newCfg)
			if err != nil {
				log.Println(err)
				return
			}
			if cfg.LogLevel != newCfg.LogLevel {
				gdb.Logger = newCfg.NewLogger()
			}
			log.Printf("gorm config reloaded: %+v\n", newCfg)
		})
		return gdb, nil
	})
}

func GetGormDB(c *ioc.Container) *gorm.DB {
	return GetProvider(DefaultConfigSection).MustGet(c)
}
