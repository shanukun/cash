package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/shanukun/cash/cash_proto"
	dt "github.com/shanukun/cash/datatypes"
	"github.com/shanukun/cash/ds"
)

var (
	ErrNoKey      = errors.New("No key found")
	ErrKeyExpired = errors.New("Key expired")
)

func getExpiration(expiration string) int64 {
	var exp int64
	duration, _ := time.ParseDuration(expiration)
	if duration > 0 {
		exp = time.Now().Add(duration).UnixNano()
	}

	return exp
}

func isExpired(expiration int64) bool {
	if expiration > 0 {
		if time.Now().UnixNano() > expiration {
			return true
		}
	}
	return false
}

type keyReport struct {
	val       dt.AnyT
	exists    bool
	typeMatch bool
}

func genKeyReport(c *cache, key string, t int) *keyReport {
	typeMatch := true
	p, exists := c.store.Find(key)
	if exists {
		switch t {
		case 0:
			_, typeMatch = p.(*dt.StringT)
		case 1:
			_, typeMatch = p.(*dt.ListT)
		case 2:
			_, typeMatch = p.(*dt.HashMapT)
		}

	}
	return &keyReport{
		val:       p,
		exists:    exists,
		typeMatch: typeMatch,
	}
}

func (c *cache) Set(ctx context.Context, item *pb.String) (*pb.Response, error) {
	expiration := getExpiration(item.Expiration)
	c.mu.Lock()

	c.expList[item.Key] = expiration

	kr := genKeyReport(c, item.Key, 0)
	if !kr.exists {
		stringData := &dt.StringT{
			Data:       item.Value,
			Expiration: expiration,
		}
		anyT := dt.AnyT(stringData)

		c.store.Insert(item.Key, anyT)
	} else if kr.typeMatch {
		stringValue := (kr.val).(*dt.StringT)
		stringValue.Data = item.Value
		stringValue.Expiration = getExpiration(item.Expiration)
	}
	c.mu.Unlock()
	return &pb.Response{
		Response: true,
	}, nil
}

func (c *cache) Get(ctx context.Context, args *pb.Key) (*pb.String, error) {
	key := args.Key
	c.mu.RLock()
	kr := genKeyReport(c, key, 0)
	if !kr.exists || !kr.typeMatch {
		c.mu.RUnlock()
		return nil, ErrNoKey
	}

	stringValue := (kr.val).(*dt.StringT)

	if isExpired(stringValue.Expiration) {
		c.mu.RUnlock()
		return nil, ErrKeyExpired
	}

	c.mu.RUnlock()

	return &pb.String{
		Key:        key,
		Value:      stringValue.Data,
		Expiration: time.Unix(0, stringValue.Expiration).String(),
	}, nil
}

func (c *cache) DeleteKey(ctx context.Context, args *pb.Key) (*pb.Response, error) {
	c.mu.Lock()
	c.store.Delete(args.Key)
	c.mu.Unlock()

	return &pb.Response{
		Response: true,
	}, nil
}

func (c *cache) LPush(ctx context.Context, item *pb.String) (*pb.Response, error) {
	expiration := getExpiration(item.Expiration)
	key := item.Key

	c.mu.RLock()
	kr := genKeyReport(c, key, 1)
	if !kr.exists {
		c.expList[item.Key] = expiration

		newList := &dt.ListT{
			Data:       []string{item.Value},
			Expiration: expiration,
		}
		anyT := dt.AnyT(newList)

		c.store.Insert(item.Key, anyT)
	} else if kr.typeMatch {
		list := (kr.val).(*dt.ListT)

		if isExpired(list.Expiration) {
			c.mu.RUnlock()
			return nil, ErrKeyExpired
		}

		list.Data = append([]string{item.Value}, list.Data...)
	}
	c.mu.RUnlock()

	return &pb.Response{
		Response: true,
	}, nil
}

func (c *cache) RPush(ctx context.Context, item *pb.String) (*pb.Response, error) {
	expiration := getExpiration(item.Expiration)
	key := item.Key
	c.mu.RLock()
	kr := genKeyReport(c, key, 1)
	if !kr.exists {
		c.expList[item.Key] = expiration

		newList := &dt.ListT{
			Data:       []string{item.Value},
			Expiration: expiration,
		}
		anyT := dt.AnyT(newList)

		c.store.Insert(item.Key, anyT)
	} else if kr.typeMatch {
		list := (kr.val).(*dt.ListT)

		if isExpired(list.Expiration) {
			c.mu.RUnlock()
			return nil, ErrKeyExpired
		}

		list.Data = append(list.Data, item.Value)
	}
	c.mu.RUnlock()

	return &pb.Response{
		Response: true,
	}, nil
}

func (c *cache) GetList(ctx context.Context, args *pb.Key) (*pb.List, error) {
	key := args.Key
	c.mu.RLock()
	kr := genKeyReport(c, key, 1)
	if !kr.exists || !kr.typeMatch {
		c.mu.RUnlock()
		return nil, ErrNoKey
	}

	list := (kr.val).(*dt.ListT)
	if isExpired(list.Expiration) {
		c.mu.RUnlock()
		return nil, ErrKeyExpired
	}

	c.mu.RUnlock()

	return &pb.List{
		Key:        args.Key,
		List:       list.Data,
		Expiration: time.Unix(0, list.Expiration).String(),
	}, nil
}

func (c *cache) HMSet(ctx context.Context, item *pb.HashMapItem) (*pb.Response, error) {
	expiration := getExpiration(item.Expiration)
	key := item.Key
	c.mu.RLock()

	kr := genKeyReport(c, key, 2)
	if !kr.exists {
		c.expList[item.Key] = expiration
		newHashMap := &dt.HashMapT{
			Data: map[string]string{
				item.Field: item.Value,
			},
			Expiration: expiration,
		}
		anyT := dt.AnyT(newHashMap)

		c.store.Insert(item.Key, anyT)
	} else if kr.typeMatch {
		hashMap := (kr.val).(*dt.HashMapT)

		if isExpired(hashMap.Expiration) {
			c.mu.RUnlock()
			return nil, ErrKeyExpired
		}

		hashMap.Data[item.Field] = item.Value

	}
	c.mu.RUnlock()

	return &pb.Response{
		Response: true,
	}, nil
}

func (c *cache) GetHashMap(ctx context.Context, args *pb.Key) (*pb.List, error) {
	c.mu.RLock()
	kr := genKeyReport(c, args.Key, 2)
	if !kr.exists || !kr.typeMatch {
		c.mu.RUnlock()
		return nil, ErrNoKey
	}

	hashMap := (kr.val).(*dt.HashMapT)
	if isExpired(hashMap.Expiration) {
		c.mu.RUnlock()
		return nil, ErrKeyExpired
	}

	var list []string
	for k, v := range hashMap.Data {
		list = append(list, k)
		list = append(list, v)
	}

	c.mu.RUnlock()

	return &pb.List{
		Key:        args.Key,
		List:       list,
		Expiration: time.Unix(0, hashMap.Expiration).String(),
	}, nil
}

func (c *cache) DeleteAll(ctx context.Context, in *empty.Empty) (*pb.Response, error) {
	c.store = ds.InitRBTree()
	return &pb.Response{
		Response: true,
	}, nil
}
