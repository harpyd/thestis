package performance

import "sync"

type Environment struct {
	mu    sync.RWMutex
	store map[string]interface{}
}

const defaultEnvStoreInitialSize = 10

func newEnvironment() *Environment {
	return &Environment{
		store: make(map[string]interface{}, defaultEnvStoreInitialSize),
	}
}

func (c *Environment) Store(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[key] = value
}

func (c *Environment) Load(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.store[key]

	return value, ok
}
