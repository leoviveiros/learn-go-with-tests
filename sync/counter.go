package sync

import "sync"

type Counter struct {
	value int
	mutex sync.Mutex
}

func NewCounter() *Counter {
    return &Counter{}
}

func (c *Counter) Inc() {
	c.mutex.Lock()
	c.value++
	c.mutex.Unlock()
}

func (c *Counter) Value() int {
	return c.value
}