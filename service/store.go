package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/shanukun/cash/cash_proto"
)

var (
	ErrNoKey      = errors.New("No key found")
	ErrKeyExpired = errors.New("Key expired")
)

func (c *cache) Set(ctx context.Context, item *pb.Data) (*pb.Data, error) {
	var expiration int64
	duration, _ := time.ParseDuration(item.Expiration)
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}
	c.mu.Lock()
	c.data[item.Key] = Data{
		Object:     item.Value,
		Expiration: expiration,
	}
	c.mu.Unlock()
	return item, nil
}

func (c *cache) Get(ctx context.Context, args *pb.Key) (*pb.Data, error) {
	key := args.Key
	c.mu.RLock()
	value, exists := c.data[key]
	if !exists {
		c.mu.RUnlock()
		return nil, ErrNoKey
	}

	if value.(Data).Expiration > 0 {
		if time.Now().UnixNano() > value.(Data).Expiration {
			c.mu.RUnlock()
			return nil, ErrKeyExpired
		}
	}
	c.mu.RUnlock()
	return &pb.Data{
		Key:        key,
		Value:      value.(Data).Object.(string),
		Expiration: time.Unix(0, value.(Data).Expiration).String(),
	}, nil
}

func (c *cache) GetByPrefix(ctx context.Context, args *pb.Key) (*pb.AllData, error) {
	key := args.Key
	c.mu.RLock()
	defer c.mu.RUnlock()
	var data []*pb.Data
	now := time.Now().UnixNano()
	for k, v := range c.data {
		if v.(Data).Expiration > 0 {
			if now > v.(Data).Expiration {
				continue
			}
		}

		if strings.Contains(k.(string), key) {
			data = append(data, &pb.Data{
				Key:        k.(string),
				Value:      v.(Data).Object.(string),
				Expiration: time.Unix(0, v.(Data).Expiration).String(),
			})
		}
	}
	if len(data) < 1 {
		return nil, ErrNoKey
	}

	return &pb.AllData{
		Data: data,
	}, nil
}

func (c *cache) GetAllData(ctx context.Context, in *empty.Empty) (*pb.AllData, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var data []*pb.Data
	now := time.Now().UnixNano()
	for k, v := range c.data {
		if v.(Data).Expiration > 0 {
			if now > v.(Data).Expiration {
				continue
			}
		}

		data = append(data, &pb.Data{
			Key:        k.(string),
			Value:      v.(Data).Object.(string),
			Expiration: time.Unix(0, v.(Data).Expiration).String(),
		})
	}

	if len(data) < 1 {
		return nil, ErrNoKey
	}

	return &pb.AllData{
		Data: data,
	}, nil
}

func (c *cache) DeleteKey(ctx context.Context, args *pb.Key) (*pb.Response, error) {
	c.mu.Lock()
	c.delete(args.Key)
	c.mu.Unlock()
	return &pb.Response{
		Response: true,
	}, nil
}

func (c *cache) DeleteAll(ctx context.Context, in *empty.Empty) (*pb.Response, error) {
	c.mu.Lock()
	c.data = map[interface{}]interface{}{}
	c.mu.Unlock()
	return &pb.Response{
		Response: true,
	}, nil
}
