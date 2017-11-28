package main

import (
	"sync"
)

type RegularMap struct {
	sync.RWMutex
	internal map[interface{}]interface{}
}

func NewRegularMap() *RegularMap {
	return &RegularMap{
		internal: make(map[interface{}]interface{}),
	}
}

func (rm *RegularMap) Load(key interface{}) (value interface{}, ok bool) {
	rm.RLock()
	result, ok := rm.internal[key]
	rm.RUnlock()
	return result, ok
}

func (rm *RegularMap) Delete(key interface{}) {
	rm.Lock()
	delete(rm.internal, key)
	rm.Unlock()
}

func (rm *RegularMap) Store(key, value interface{}) {
	rm.Lock()
	rm.internal[key] = value
	rm.Unlock()
}
