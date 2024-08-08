package ioc

import (
	"fmt"
	"sync"
)

type injector interface {
	new(c *Container) (any, error)
}

type Provider[T any] struct {
	f  func(c *Container) (T, error)
	mu sync.Mutex
}

func (f *Provider[T]) new(c *Container) (any, error) {
	return f.f(c)
}

func (f *Provider[T]) Get(c *Container) (t T, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	ins, ok := c.get(f)
	if ok {
		return ins.(T), nil
	} else {
		ins, err = f.new(c)
		if err != nil {
			return t, err
		}
		if ins == nil {
			return t, fmt.Errorf("ioc: provider %T returned nil", f)
		}
		c.set(f, ins)
		return ins.(T), nil
	}
}

func (f *Provider[T]) MustGet(c *Container) T {
	t, err := f.Get(c)
	if err != nil {
		panic(err)
	}
	return t
}

func (f *Provider[T]) GetNew(c *Container) (T, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	t, err := f.new(c)
	if err != nil {
		var zero T
		return zero, err
	}
	return t.(T), nil
}

func (f *Provider[T]) MustGetNew(c *Container) T {
	t, err := f.GetNew(c)
	if err != nil {
		panic(err)
	}
	return t
}

func (f *Provider[T]) IsSet(c *Container) bool {
	_, ok := c.get(f)
	return ok
}

func (f *Provider[T]) Set(c *Container, ins T) {
	c.set(f, ins)
}

func (f *Provider[T]) Refresh(c *Container) (ins T, err error) {
	newIns, err := f.GetNew(c)
	if err != nil {
		return newIns, err
	}
	f.Set(c, newIns)
	return newIns, nil
}

func (f *Provider[T]) MustRefresh(c *Container) T {
	ins, err := f.Refresh(c)
	if err != nil {
		panic(err)
	}
	return ins
}

func NewProvider[T any](new func(c *Container) (T, error)) *Provider[T] {
	return &Provider[T]{f: new}
}

type Providers[T any] struct {
	ps map[string]*Provider[T]
	mu sync.Mutex
}

func NewProviders[T any]() *Providers[T] {
	return &Providers[T]{
		ps: map[string]*Provider[T]{},
	}
}

func (r *Providers[T]) GetProvider(name string, new func(c *Container) (T, error)) *Provider[T] {
	r.mu.Lock()
	defer r.mu.Unlock()
	pvd := r.ps[name]
	if pvd == nil {
		pvd = NewProvider(new)
		r.ps[name] = pvd
	}
	return pvd
}
