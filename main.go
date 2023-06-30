package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/shanukun/cash/cash_proto"
	service "github.com/shanukun/cash/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	address string
	expire  int
	cleanup int
)

func parseFlags() {
	flag.StringVar(&address, "addr", ":8001", "address")
	flag.IntVar(&expire, "exp", 7, "expiration (min)")
	flag.IntVar(&cleanup, "clu", 3, "cleanup after expiration (min)")
	flag.Parse()
}

func main() {
	parseFlags()

	opts := []grpc.ServerOption{
		grpc.MaxConcurrentStreams(100),
	}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterCacheServiceServer(
		grpcServer,
		service.NewCacheService(time.Duration(expire)*time.Minute,
			time.Duration(cleanup)*time.Minute))

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("start error %v", err)
	}
	fmt.Println("server running on:", address)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("grpc server failed: %v\n", err)
	}
}
