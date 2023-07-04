package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	pb "github.com/shanukun/cash/cash_proto"
)

var (
	address string
	conn    *grpc.ClientConn
	err     error
)

func main() {

	// Get address from flag
	flag.StringVar(&address, "addr", "127.0.0.1:8001", "Setress on which you want to run server")
	flag.Parse()

    // Connecting to server
	conn, err = grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

    // Creating new client
	c := pb.NewCacheServiceClient(conn)
	ctx := context.Background()

	// Set key
	keyVal1 := &pb.String{
		Key:        "author",
		Value:      "Brandon Sanderson",
		Expiration: "-1m",
	}

	keyVal2 := &pb.String{
		Key:        "book",
		Value:      "Mistborn",
		Expiration: "20s",
	}

	keyVal3 := &pb.String{
		Key:        "planet",
		Value:      "Roshar",
		Expiration: "2min10s",
	}

	c.Set(ctx, keyVal1)
	c.Set(ctx, keyVal2)

	addKeyRes, err := c.Set(ctx, keyVal3)
	if err != nil {
		log.Fatalf("Error when calling Set: %s", err)
	}
	fmt.Println("Response from server for adding a key", addKeyRes)

	// Checking for race condition
	for i := 0; i < 50; i++ {
		go c.Set(ctx, &pb.String{
			Key:        strconv.Itoa(i),
			Value:      "Value of i is ",
			Expiration: strconv.Itoa(i),
		})
	}

	// Get key
	keyGet := &pb.Key{
		Key: "book",
	}

	getKeyRes, err := c.Get(ctx, keyGet)
	if err != nil {
		log.Fatalf("Error when calling Get: Key: %s %s", keyGet.Key, err)
	}
	fmt.Println("Response from server for getting a key", getKeyRes)


	// String item for List
	keyVal4 := &pb.String{
		Key:        "tbr",
		Value:      "mistborn",
		Expiration: "2min10s",
	}

	keyVal5 := &pb.String{
		Key:        "tbr",
		Value:      "stormlight",
		Expiration: "2min10s",
	}

    // Left push
	c.LPush(ctx, keyVal4)
    // Right push
	c.RPush(ctx, keyVal5)

	keyGetList := &pb.Key{
		Key: "tbr",
	}

    // Get complete list
	list, err := c.GetList(ctx, keyGetList)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(list.List)

	// Delete Key
	deleteKeyRes, err := c.DeleteKey(ctx, keyGet)
	if err != nil {
		log.Fatalf("Error when calling DeleteKey: %s", err)
	}
	fmt.Println("Response from server after deleting a key", deleteKeyRes)

    // Set Key, Field and Value for HashMap
	keyVal6 := &pb.HashMapItem{
		Key:        "mybooks",
		Field:      "fantasy",
		Value:      "stormlight",
		Expiration: "2min10s",
	}

	keyVal7 := &pb.HashMapItem{
		Key:        "mybooks",
		Field:      "scifi",
		Value:      "snow crash",
		Expiration: "2min10s",
	}

    // Set Hash Map values
	c.HMSet(ctx, keyVal6)
	c.HMSet(ctx, keyVal7)

	keyGetHashMap := &pb.Key{
		Key: "mybooks",
	}

    // Get HashMap as a List
	hmlist, err := c.GetHashMap(ctx, keyGetHashMap)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(hmlist.List)


    // Deleting all keys
	c.DeleteAll(ctx, &empty.Empty{})
}
