package service

import (
	pb "github.com/shanukun/cash/cash_proto"
	ds "github.com/shanukun/cash/ds"
	"runtime"
	"sync"
	"time"
)

type Cache struct {
	*cache
}

type worker struct {
	Interval time.Duration
	stop     chan bool
}

type cache struct {
	defaultExpiration time.Duration
	expList           map[string]int64
	mu                sync.RWMutex
	store             *ds.RBTree
	worker            *worker
	pb.UnimplementedCacheServiceServer
}

func NewCacheService(defaultExpiration, cleanupInterval time.Duration) *Cache {
	store := ds.InitRBTree()
	return newCacheWithWorker(defaultExpiration, cleanupInterval, store)
}

func newCacheWithWorker(defaultExpiration time.Duration, cleanupInterval time.Duration, store *ds.RBTree) *Cache {
	c := newCache(defaultExpiration, store)
	C := &Cache{c}
	if cleanupInterval > 0 {
		runWorker(c, cleanupInterval)
		runtime.SetFinalizer(C, stopWorker)
	}
	return C
}

func newCache(defaultExpiration time.Duration, store *ds.RBTree) *cache {
	c := &cache{
		defaultExpiration: defaultExpiration,
		store:             store,
		expList:           make(map[string]int64),
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

func (c *cache) deleteExpired() {
	now := time.Now().UnixNano()
	c.mu.Lock()
	for k, v := range c.expList {
		if v > 0 && now > v {
			c.store.Delete(k)
			delete(c.expList, k)
		}
	}
	c.mu.Unlock()
}
