package rate

import (
	"sync"
	"time"
)

type MemoryStorage[T any] struct {
	data   map[string]*T
	expire map[string]time.Time
	mu     sync.Mutex
}

func NewMemoryStorage[T any](cleanUpInterval time.Duration) *MemoryStorage[T] {
	s := &MemoryStorage[T]{
		data:   make(map[string]*T),
		expire: make(map[string]time.Time),
	}
	go s.runCleanUp(cleanUpInterval)
	return s
}

func (m *MemoryStorage[T]) runCleanUp(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case now := <-ticker.C:
			func() {
				m.mu.Lock()
				defer m.mu.Unlock()
				for name, expire := range m.expire {
					if now.After(expire) {
						delete(m.data, name)
						delete(m.expire, name)
					}
				}
			}()
		}
	}
}

func (m *MemoryStorage[T]) Set(name string, data *T, expire time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[name] = data
	m.expire[name] = time.Now().Add(expire)
	return nil
}

func (m *MemoryStorage[T]) Get(name string) (data *T, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	expire, ok := m.expire[name]
	if !ok {
		return data, nil
	}
	if time.Now().After(expire) {
		delete(m.data, name)
		delete(m.expire, name)
		return data, nil
	}
	return m.data[name], nil
}

func (m *MemoryStorage[T]) Del(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, name)
	delete(m.expire, name)
	return nil
}
