package main

import "sync"

type InMemoryPlayerStore struct {
	store map[string]int
	mutex sync.RWMutex
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	return i.store[name]
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.mutex.Lock()
	i.store[name]++
	i.mutex.Unlock()
}

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{
		map[string]int{}, 
		sync.RWMutex{},
	}
}