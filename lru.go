package lru

import (
	. "github.com/golang/groupcache/lru"
)

type LruCache interface {
	Add(Key, interface{})
	Get(Key) (interface{}, bool)
	Remove(Key)
	RemoveOldest()
	Len() int
}

type get struct {
	key        Key
	returnChan chan interface{}
}

type add struct {
	key        Key
	value      interface{}
	returnChan chan struct{}
}

type remove struct {
	key        Key
	returnChan chan struct{}
}

type removeOldest struct {
	returnChan chan struct{}
}

type lruCache struct {
	cache            *Cache
	getChan          chan *get
	addChan          chan *add
	removeChan       chan *remove
	removeOldestChan chan *removeOldest
}

func (l *lruCache) Add(key Key, value interface{}) {
	returnChan := make(chan struct{})
	l.addChan <- &add{
		key:        key,
		value:      value,
		returnChan: returnChan}

	<-returnChan
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	returnChan := make(chan interface{})
	l.getChan <- &get{
		key:        key,
		returnChan: returnChan,
	}

	v, ok := <-returnChan
	return v, ok
}

func (l *lruCache) Remove(key Key) {
	returnChan := make(chan struct{})

	l.removeChan <- &remove{
		key:        key,
		returnChan: returnChan,
	}

	<-returnChan
}

func (l *lruCache) RemoveOldest() {
	returnChan := make(chan struct{})

	l.removeOldestChan <- &removeOldest{returnChan: returnChan}

	<-returnChan
}

func (l *lruCache) Len() int {
	return l.cache.Len()
}

func NewLruCache(maxEntries int) LruCache {
	cache := New(maxEntries)

	threadSafeCache := &lruCache{
		cache:            cache,
		getChan:          make(chan *get),
		addChan:          make(chan *add),
		removeChan:       make(chan *remove),
		removeOldestChan: make(chan *removeOldest),
	}

	go func() {
		for {
			select {
			case get := <-threadSafeCache.getChan:
				returnChan := get.returnChan
				if value, ok := threadSafeCache.cache.Get(get.key); ok {
					returnChan <- value
					continue
				}
				close(returnChan)

			case add := <-threadSafeCache.addChan:
				returnChan := add.returnChan
				threadSafeCache.cache.Add(add.key, add.value)
				close(returnChan)

			case remove := <-threadSafeCache.removeChan:
				returnChan := remove.returnChan
				threadSafeCache.cache.Remove(remove.key)

				close(returnChan)

			case removeOldest := <-threadSafeCache.removeOldestChan:
				returnChan := removeOldest.returnChan
				threadSafeCache.cache.RemoveOldest()
				close(returnChan)

			}
		}
	}()

	return threadSafeCache
}
