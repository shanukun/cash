package service

import (
	pb "github.com/shanukun/cash/cash_proto"
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
