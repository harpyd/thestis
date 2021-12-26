package performance

import "sync"

type Context struct {
	mu sync.RWMutex

	store map[string]interface{}
}

func (c *Context) Store(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[key] = value
}

func (c *Context) Load(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.Unlock()

	value, ok := c.store[key]

	return value, ok
}
