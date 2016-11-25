package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	naming "gomicro/naming/etcd"
	"gomicro/naming/examples/pb"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	service = flag.String("service", "hello_service", "service name")
	port    = flag.Int("port", 1701, "listening port")
	addr    = flag.String("addr", "http://127.0.0.1:2379", "register address")
)

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		panic(err)
	}

	// register a service to etcd
	err = naming.Register(*service, "127.0.0.1", *port, *addr, time.Second*3, 5)
	if err != nil {
		panic(err)
	}

	log.Printf("starting hello service at %d", *port)
	s := grpc.NewServer()
	pb.RegisterHelloServiceServer(s, &HelloServer{})
	s.Serve(listener)
}

// HelloServer 创建对象，实现服务
type HelloServer struct{}

// SayHello 实现服务
func (HelloServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Println("getting request from client")
	reply := "Hello, " + r.Greeting
	return &pb.HelloResponse{Reply: reply}, nil
}
