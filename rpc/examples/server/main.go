package main

import (
	"flag"
	"fmt"
	"net"

	"gomicro/log"
	"gomicro/rpc"
	"gomicro/rpc/examples/pb"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var (
	port = flag.Int("port", 1701, "listening port")
)

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		panic(err)
	}

	log.CtxPrintf(nil, "starting hello service at %d", *port)
	s := rpc.NewServer()
	pb.RegisterHelloServiceServer(s, &HelloServer{})
	s.Serve(listener)
}

// HelloServer 创建对象，实现服务
type HelloServer struct{}

// NormalHello 实现服务
func (HelloServer) NormalHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.CtxInfof(ctx, "normal hello")
	reply := "Hello, " + r.Greeting
	return &pb.HelloResponse{Reply: reply}, nil
}

// PanicHello 实现服务
func (HelloServer) PanicHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.CtxInfof(ctx, "panic hello")
	panic(fmt.Errorf("nothing"))

	reply := "Hello, " + r.Greeting
	return &pb.HelloResponse{Reply: reply}, nil
}

// ErrorHello 实现服务
func (HelloServer) ErrorHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.CtxInfof(ctx, "error hello")
	return nil, grpc.Errorf(codes.Canceled, "just try to error")
}
