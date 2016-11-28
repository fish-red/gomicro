package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"

	"gomicro/log"
	"gomicro/rpc"
	"gomicro/rpc/examples/pb"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
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
	grpc_prometheus.Register(s)
	http.Handle("/metrics", prometheus.Handler())
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
}

// ErrorHello 实现服务
func (HelloServer) ErrorHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.CtxInfof(ctx, "error hello")
	return nil, grpc.Errorf(codes.Canceled, "just try to error")
}
