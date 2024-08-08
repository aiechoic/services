package ioc

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	ConfigEnvTest ConfigEnv = "test"
	ConfigEnvDev  ConfigEnv = "dev"
	ConfigEnvProd ConfigEnv = "prod"
)

type ConfigEnv string

type Vipers struct {
	dir      string
	env      ConfigEnv
	encoders map[string]*viper.Viper
	mu       sync.Mutex
}

func NewVipers(dir string, env ConfigEnv) (*Vipers, error) {

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	encoders := make(map[string]*viper.Viper)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filename := file.Name()
		ext := filepath.Ext(filename)
		suffix := fmt.Sprintf(".%s%s", env, ext)
		if !strings.HasSuffix(filename, suffix) {
			continue
		}
		subName := strings.TrimSuffix(filename, suffix)
		if encoders[subName] != nil {
			return nil, fmt.Errorf("duplicate config file: \"%s\"", subName)
		}
		v := viper.New()
		switch ext {
		case ".json":
			v.SetConfigType("json")
		case ".yaml", ".yml":
			v.SetConfigType("yaml")
		case ".toml":
			v.SetConfigType("toml")
		default:
			return nil, fmt.Errorf("unsupported config file type: %s", ext)
		}
		v.AutomaticEnv()
		v.SetEnvPrefix(subName)
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		v.SetConfigFile(filepath.Join(dir, filename))
		err = v.ReadInConfig()
		if err != nil {
			return nil, fmt.Errorf("error reading config file \"%s\": %w", filename, err)
		}
		encoders[subName] = v
	}

	return &Vipers{
		dir:      dir,
		env:      env,
		encoders: encoders,
	}, nil
}

func (c *Vipers) Unmarshal(name string, v any, def []byte) error {
	vp, err := c.GetOrCreateViper(name, def)
	if err != nil {
		return err
	}
	err = vp.Unmarshal(v)
	if err != nil {
		file := vp.ConfigFileUsed()
		return fmt.Errorf("unmarshalling config \"%s\": %w", file, err)
	}
	return nil
}

func (c *Vipers) GetOrCreateViper(name string, defaultContent []byte) (*viper.Viper, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	vp, ok := c.encoders[name]
	if !ok {
		ct := c.getContentType(defaultContent)
		if ct == "" {
			return nil, fmt.Errorf("cannot determine the format of the default data")
		}
		filename := fmt.Sprintf("%s.%s.%s", name, c.env, ct)
		filename = filepath.Join(c.dir, filename)
		err := os.WriteFile(filename, defaultContent, 0644)
		if err != nil {
			return nil, fmt.Errorf("error writing default config file \"%s\": %w", filename, err)
		}
		log.Printf("default config file created: %s\n", filename)
		vp = viper.New()
		vp.AutomaticEnv()
		vp.SetEnvPrefix(name)
		vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		vp.SetConfigFile(filename)
		err = vp.ReadInConfig()
		if err != nil {
			return nil, fmt.Errorf("error reading config file \"%s\": %w", filename, err)
		}
		c.encoders[name] = vp
	}
	return vp, nil
}

func (c *Vipers) WatchConfig(name string, callback func(v *viper.Viper)) error {
	vp, ok := c.encoders[name]
	if !ok {
		return fmt.Errorf("config %s not found", name)
	}

	vp.WatchConfig()
	vp.OnConfigChange(func(in fsnotify.Event) {
		if in.Op&fsnotify.Write == fsnotify.Write {
			callback(vp)
		}
	})
	return nil
}

func (c *Vipers) UnmarshalAndWatch(name string, defaultContent []byte, callback func(v *viper.Viper)) error {
	vp, err := c.GetOrCreateViper(name, defaultContent)
	if err != nil {
		return err
	}
	callback(vp)
	return c.WatchConfig(name, callback)
}

func (c *Vipers) getContentType(content []byte) string {
	dt := map[string]interface{}{}
	err := json.Unmarshal(content, &dt)
	if err == nil {
		return "json"
	}
	err = yaml.Unmarshal(content, &dt)
	if err == nil {
		return "yaml"
	}
	err = toml.Unmarshal(content, &dt)
	if err == nil {
		return "toml"
	}
	return ""
}
