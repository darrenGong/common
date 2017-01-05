package common

import "sync"

type ConnMap struct {
	m map[string]interface{}
	sync.Mutex
}

func NewConnMap() *ConnMap {
	return &ConnMap{m: make(map[string]interface{})}
}

func (c *ConnMap) AddConn(k string, v interface{}) {
	c.Lock()
	defer c.Unlock()

	c.m[k] = v
}

func (c *ConnMap) GetConn(k string) interface{} {
	c.Lock()
	defer c.Unlock()

	if v, ok := c.m[k]; ok {
		return v
	}

	return nil
}

func (c *ConnMap) DelConn(k string) {
	c.Lock()
	defer c.Unlock()

	delete(c.m, k)
}
