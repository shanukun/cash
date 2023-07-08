# cash

 Cash is an in‐memory key‐value store that may be used as a distributed cache, and it's also concurrency safe.
- gRPC is used to construct APIs for adding/replacing key/value, getting value using a key, adding key/value pairs to the list and map, deleting a specific key and deleting all keys.
- It uses Red-Black Tree for storing a variety of abstract data types, including lists, maps, and strings.

Note: Implementing another Data Structure to store data types is quite straightforward.


## Installing

```
go install github.com/shanukun/cash@latest
```

### Options

```
Usage of cash:
  -addr string
    	address (default ":8001")
  -clu int
    	cleanup after expiration (min) (default 3)
  -exp int
    	expiration (min) (default 7)
```



## API

You can find the [proto file here](https://github.com/shanukun/cash/blob/master/cash_proto/cash.proto).
For examples checkout [demo file](https://github.com/shanukun/cash/blob/master/demo/demo.go)


### Set

Set supplied value at key. If key already exists, value is updated.


```go
func (c Cache) Set(ctx context.Context, item *pb.String) (*pb.Response, error)
```

### Get

Get value stored at key.

```go
func (c Cache) Get(ctx context.Context, args *pb.Key) (*pb.String, error)
```

### LPush

Add all of the supplied values to the front of the list stored at key. If key does not exists, new empty list will be created.

```go
func (c Cache) LPush(ctx context.Context, item *pb.String) (*pb.Response, error)
```

### RPush

Add all of the supplied values to the back of the list stored at key. If key does not exists, new empty list will be created.

```go
func (c Cache) RPush(ctx context.Context, item *pb.String) (*pb.Response, error)
```

### Get List

Get list stored at key. 

```go
func (c Cache) GetList(ctx context.Context, args *pb.Key) (*pb.List, error)
```

### HMSet

Set supplied fields to their respective values to the HashMap stored at key. If key does not exists, new HashMap will be created.

```go
func (c Cache) HMSet(ctx context.Context, item *pb.HashMapItem) (*pb.Response, error)
```

### Get HashMap

Get HasmMap stored at key. Function will return a List.
```
Example: 
HashMap = {key1: value1, key2: value2} 
List = [key1, value1, key2, value2]
```

```go
func (c Cache) GetHashMap(ctx context.Context, args *pb.Key) (*pb.List, error)
```

### DeleteKey

Delete key along with value stored.

```go
func (c Cache) DeleteKey(ctx context.Context, args *pb.Key) (*pb.Response, error)
```

### DeleteAll

Delete all the keys at once.

```go
func (c Cache) DeleteAll(ctx context.Context, in *empty.Empty) (*pb.Response, error)
```
