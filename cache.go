package cache

import (
	"sync"
	"time"
)

type Cache struct {
	cacheMap map[string]Value
	mutex    sync.Mutex
}

type Value struct {
	deadline *time.Time
	str      string
}

func NewCache() Cache {
	return Cache{cacheMap: map[string]Value{}}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, ok := c.cacheMap[key]; !ok {
		return "", false
	} else if c.cacheMap[key].deadline != nil && c.cacheMap[key].deadline.Before(time.Now()) {
		return "", false
	}
	return c.cacheMap[key].str, true
}

func (c *Cache) Put(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cacheMap[key] = Value{str: value, deadline: nil}
}

func (c *Cache) Keys() []string {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	var keys []string
	time := time.Now()
	for key, value := range c.cacheMap {
		if value.deadline != nil && value.deadline.Before(time) {
			continue
		}
		keys = append(keys, key)
	}
	return keys
}

func (c *Cache) PutTill(key, value string, deadline time.Time) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cacheMap[key] = Value{&deadline, value}
}
