package healthy

import (
	"context"
	"fmt"
	"github.com/aiechoic/services/ioc"
	"github.com/aiechoic/services/ioc/healthy"
	"time"
)

type Engine struct {
	start bool
}

func (e *Engine) Start() {
	if !e.start {
		fmt.Printf("start engine\n")
		e.start = true
	}
}

func (e *Engine) Stop() {
	if e.start {
		fmt.Printf("stop engine\n")
		e.start = false
	}
}

func (e *Engine) Toggle() {
	if e.start {
		e.Stop()
	} else {
		e.Start()
	}
}

var ErrorEngineNotStart = &healthy.Error{Level: healthy.LError, Msg: "Engine is not start"}

func (e *Engine) HealthyCheck() *healthy.Error {
	if !e.start {
		return ErrorEngineNotStart
	}
	return nil
}

var engineProvider = ioc.NewProvider(func(c *ioc.Container) (*Engine, error) {
	engine := &Engine{}
	c.OnClose(func() error {
		fmt.Printf("close engine\n")
		engine.Stop()
		return nil
	})
	c.OnHealthCheck(engine.HealthyCheck)
	fmt.Println("create new engine")
	return engine, nil
})

func Example() {
	healthyTicker := 1 * time.Second
	checkTimeout := 5 * time.Second

	c := ioc.NewContainer()
	defer c.Close()

	engine := engineProvider.MustGet(c)
	engine.Start()

	c.RunHealthCheck(healthyTicker, checkTimeout, func(errs []*healthy.Error) {
		if len(errs) == 0 {
			fmt.Printf("Healthy: OK\n")
		} else {
			for _, err := range errs {
				fmt.Printf("Healthy ERROR: %s\n", err.Msg)
			}
		}
	})

	toggleTicker := time.NewTicker(2 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for {
		select {
		case <-toggleTicker.C:
			engine.Toggle()
		case <-ctx.Done():
			return
		}
	}

	// Output:
	// create new engine
	// start engine
	// Healthy: OK
	// stop engine
	// Healthy ERROR: Engine is not start
	// Healthy ERROR: Engine is not start
	// start engine
	// Healthy: OK
	// close engine
	// stop engine
}
