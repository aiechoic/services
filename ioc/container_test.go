package ioc_test

import (
	"context"
	"errors"
	"github.com/aiechoic/services/ioc"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewContainer(t *testing.T) {
	c := ioc.NewContainer()

	provider := ioc.NewProvider(func(c *ioc.Container) (*string, error) {
		msg := "test instance"
		return &msg, nil
	})

	instance := provider.MustGet(c)
	assert.Equal(t, "test instance", *instance)

	// Test caching behavior
	cachedInstance := provider.MustGet(c)
	assert.True(t, instance == cachedInstance)

	// Test New instance
	newInstance := provider.MustGetNew(c)
	assert.False(t, instance == newInstance)
	assert.Equal(t, "test instance", *newInstance)

	provider.Set(c, newInstance)

	cachedInstance = provider.MustGet(c)
	assert.True(t, newInstance == cachedInstance)
}

func TestCloseWithContext(t *testing.T) {
	c := ioc.NewContainer()

	// Mock closers
	closer1Called := false
	closer1 := func() error {
		closer1Called = true
		return nil
	}

	closer2Called := false
	closer2 := func() error {
		closer2Called = true
		return errors.New("closer2 error")
	}

	closer3Called := false
	closer3 := func() error {
		time.Sleep(4 * time.Second) // Simulate a long-running closer
		closer3Called = true
		return nil
	}

	c.OnClose(closer1)
	c.OnClose(closer2)
	c.OnClose(closer3)

	// Call CloseWithContext with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := c.CloseWithContext(ctx)

	// Verify that all closers are called
	assert.True(t, closer1Called, "closer1 should be called")
	assert.True(t, closer2Called, "closer2 should be called")

	// Verify that the error is returned correctly
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "closer2 error")

	// Verify that the long-running closer is canceled
	assert.False(t, closer3Called, "closer3 should not be called")
	assert.Contains(t, err.Error(), "context deadline exceeded")
}
