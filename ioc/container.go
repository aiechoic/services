package ioc

import (
	"context"
	"fmt"
	"github.com/aiechoic/services/ioc/healthy"
	"github.com/spf13/viper"
	"runtime"
	"sync"
	"time"
)

type withPkgFunc[T any] struct {
	pkg string
	f   func() T
}

type Container struct {
	instances      map[injector]any
	closers        []*withPkgFunc[error]
	healthCheckers []*withPkgFunc[*healthy.Error]
	vipers         *Vipers
	cancel         context.CancelFunc
	mu             sync.Mutex
}

func NewContainer() *Container {
	return &Container{
		instances: map[injector]any{},
	}
}

func (c *Container) LoadConfig(dir string, env ConfigEnv) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	config, err := NewVipers(dir, env)
	if err != nil {
		return err
	}
	c.vipers = config
	return nil
}

func (c *Container) UnmarshalConfig(name string, v any, defaultContent []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.vipers == nil {
		return fmt.Errorf("config not loaded")
	}
	return c.vipers.Unmarshal(name, v, defaultContent)
}

func (c *Container) WatchConfig(name string, handler func(v *viper.Viper)) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.vipers == nil {
		return fmt.Errorf("config not loaded")
	}
	return c.vipers.WatchConfig(name, handler)
}

func (c *Container) UnmarshalAndWatchConfig(name string, defaultContent []byte, handler func(v *viper.Viper)) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.vipers == nil {
		return fmt.Errorf("config not loaded")
	}
	return c.vipers.UnmarshalAndWatch(name, defaultContent, handler)
}

func (c *Container) get(p injector) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	i, ok := c.instances[p]
	return i, ok
}

func (c *Container) set(p injector, ins any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.instances[p] = ins
}

func (c *Container) OnClose(closer func() error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closers = append(c.closers, &withPkgFunc[error]{
		pkg: getCallerLocation(2),
		f:   closer,
	})
}

func getCallerLocation(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	var location string
	if ok {
		location = fmt.Sprintf("%s:%d", file, line)
	} else {
		location = "unknown location"
	}
	return location
}

func (c *Container) OnHealthCheck(checker func() *healthy.Error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.healthCheckers = append(c.healthCheckers, &withPkgFunc[*healthy.Error]{
		pkg: getCallerLocation(2),
		f:   checker,
	})
}

func (c *Container) RunHealthCheck(ticker, timeout time.Duration, handler func(errs []*healthy.Error)) {
	t := time.NewTicker(ticker)
	ctx, cancel := context.WithCancel(context.Background())
	c.mu.Lock()
	c.cancel = cancel
	c.mu.Unlock()
	go func() {
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				func() {
					var subCtx context.Context
					if timeout == 0 {
						subCtx = ctx
					} else {
						var subCancel context.CancelFunc
						subCtx, subCancel = context.WithTimeout(ctx, timeout)
						defer subCancel()
					}
					errs := c.CheckHealth(subCtx)
					handler(errs)
				}()
			}
		}
	}()
}

func (c *Container) CheckHealth(ctx context.Context) []*healthy.Error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if context is done
	select {
	case <-ctx.Done():
		return nil
	default:
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errors []*healthy.Error

	for _, checker := range c.healthCheckers {
		wg.Add(1)
		go func(checker *withPkgFunc[*healthy.Error]) {
			defer wg.Done()
			errChan := make(chan *healthy.Error, 1)
			go func() {
				errChan <- checker.f()
			}()

			var err *healthy.Error
			select {
			case err = <-errChan:
			case <-ctx.Done():
				err = &healthy.Error{
					Level: healthy.LError,
					Msg:   fmt.Sprintf("%s: %s", checker.pkg, ctx.Err()),
				}
			}
			if err != nil {
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
			}
		}(checker)
	}

	wg.Wait()
	return errors
}

type MultiError []error

func (m MultiError) Error() string {
	var s string
	for _, e := range m {
		s += e.Error() + "\n"
	}
	return s
}

func (c *Container) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.CloseWithContext(ctx)
}

func (c *Container) CloseWithContext(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cancel != nil {
		c.cancel()
	}

	var errs []error
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, clo := range c.closers {
		wg.Add(1)
		go func(clo *withPkgFunc[error]) {
			defer wg.Done()
			errChan := make(chan error, 1)
			go func() {
				errChan <- clo.f()
			}()

			var err error
			select {
			case err = <-errChan:
			case <-ctx.Done():
				err = ctx.Err()
			}
			if err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("%s: %w", clo.pkg, err))
				mu.Unlock()
			}
		}(clo)
	}

	wg.Wait()

	if len(errs) > 0 {
		return MultiError(errs)
	}
	return nil
}
