package service

import (
	pb "github.com/shanukun/cash/cash_proto"
	"runtime"
	"sync"
	"time"
)

type Data struct {
	Object     interface{}
	Expiration int64
}

type Cache struct {
	*cache
}

type worker struct {
	Interval time.Duration
	stop     chan bool
}

type cache struct {
	defaultExpiration time.Duration
	mu                sync.RWMutex
	data              map[interface{}]interface{}
	worker            *worker
	pb.UnimplementedCacheServiceServer
}

func NewCacheService(defaultExpiration, cleanupInterval time.Duration) *Cache {
	data := make(map[interface{}]interface{})
	return newCacheWithWorker(defaultExpiration, cleanupInterval, data)
}

func newCacheWithWorker(defaultExpiration time.Duration, cleanupInterval time.Duration, data map[interface{}]interface{}) *Cache {
	c := newCache(defaultExpiration, data)
	C := &Cache{c}
	if cleanupInterval > 0 {
		runWorker(c, cleanupInterval)
		runtime.SetFinalizer(C, stopWorker)
	}
	return C
}

func newCache(defaultExpiration time.Duration, data map[interface{}]interface{}) *cache {
	c := &cache{
		defaultExpiration: defaultExpiration,
		data:              data,
	}
	return c
}

func stopWorker(c *Cache) {
	c.worker.stop <- true
}

func runWorker(c *cache, cleanupInterval time.Duration) {
	w := &worker{
		Interval: cleanupInterval,
		stop:     make(chan bool),
	}
	c.worker = w
	go w.Run(c)
}

func (w *worker) Run(c *cache) {
	ticker := time.NewTicker(w.Interval)
	for {
		select {
		case <-ticker.C:
			c.deleteExpired()
		case <-w.stop:
			ticker.Stop()
			return
		}
	}
}

func (c *cache) delete(k interface{}) (interface{}, bool) {
	delete(c.data, k)
	return nil, false
}

func (c *cache) deleteExpired() {
	now := time.Now().UnixNano()
	c.mu.Lock()
	for k, v := range c.data {
		if v.(Data).Expiration > 0 && now > v.(Data).Expiration {
			c.delete(k)
		}
	}

	c.mu.Unlock()
}
