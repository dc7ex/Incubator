package context

import (
	"sync"
)

type (
	// Map defines a generic map of type `map[string]interface{}`.
	Map map[string]interface{}
)

type (
	Context interface {
		// Get retrieves data from the context.
		Get(key string) interface{}

		// Set saves data in the context.
		Set(key string, val interface{})
	}
)

type context struct {
	store Map
	lock  sync.RWMutex
}

func (c *context) Get(key string) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.store[key]
}

func (c *context) Set(key string, val interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.store == nil {
		c.store = make(Map)
	}
	c.store[key] = val
}
