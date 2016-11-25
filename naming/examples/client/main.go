package main

import (
	"flag"
	"fmt"
	"time"

	naming "gomicro/naming/etcd"
	"gomicro/naming/examples/pb"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	service = flag.String("service", "hello_service", "service name")
	addr    = flag.String("address", "http://127.0.0.1:2379", "register address")
)

func main() {
	flag.Parse()

	resolver := naming.NewResolver(*service)
	balancer := grpc.RoundRobin(resolver)

	conn, err := grpc.Dial(*addr, grpc.WithInsecure(), grpc.WithBalancer(balancer))
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(2 * time.Second)
	for t := range ticker.C {
		client := pb.NewHelloServiceClient(conn)
		response, err := client.SayHello(context.Background(), &pb.HelloRequest{Greeting: "world"})
		if err != nil {
			panic(err)
		}
		fmt.Printf("%v: reply is %s\n", t, response.Reply)
	}
}
