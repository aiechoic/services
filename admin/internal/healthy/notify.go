package healthy

import (
	"bytes"
	"github.com/aiechoic/services/ioc"
	"github.com/aiechoic/services/ioc/healthy"
	"github.com/aiechoic/services/message/notifier"
	"github.com/aiechoic/services/message/notifier/dingtalk"
	"github.com/spf13/viper"
	"log"
	"sync"
	"text/template"
	"time"
)

var DefaultConfigSection ConfigSection = "healthy-check"

type ConfigSection string

func CheckHealthyWithNotifier(c *ioc.Container, n notifier.Notifier, cfgSection ConfigSection) error {
	var cfg *Config
	mu := sync.Mutex{}
	err := c.UnmarshalAndWatchConfig(string(cfgSection), defaultConfigData, func(v *viper.Viper) {
		mu.Lock()
		defer mu.Unlock()
		err := v.Unmarshal(&cfg)
		if err != nil {
			log.Printf("failed to unmarshal config: %v\n", err)
		}
		log.Printf("loaded healthy config \"%s\" success\n", cfgSection)
	})
	if err != nil {
		return err
	}
	c.RunHealthCheck(1*time.Minute, 5*time.Second, func(errs []*healthy.Error) {
		mu.Lock()
		defer mu.Unlock()
		if errs != nil {
			var _errs []*healthy.Error
			for _, err := range errs {
				if err.Level.Number() >= cfg.NotifyLevel.Number() {
					_errs = append(_errs, err)
				}
			}
			if len(_errs) == 0 {
				return
			}
			msg, err := healthyErrorToMarkdownText(errs)
			if err != nil {
				log.Printf("failed to generate markdown text: %v", err)
				return
			}
			err = n.Notify(msg)
			if err != nil {
				log.Printf("failed to notify: %v", err)
			}
		}
	})
	return nil
}

func CheckHealthy(c *ioc.Container) error {
	dingding := dingtalk.GetClient(c)
	return CheckHealthyWithNotifier(c, dingding, DefaultConfigSection)
}

var tpl = func() *template.Template {
	t := template.New("healthy")
	t.Parse(`
# Health Check
{{- range . }}

### Error Level: {{ .Level }}

- **Message**: {{ .Msg }}
{{ end -}}
`)
	return t
}()

func healthyErrorToMarkdownText(errs []*healthy.Error) (string, error) {
	buf := bytes.NewBuffer(nil)

	err := tpl.Execute(buf, errs)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
