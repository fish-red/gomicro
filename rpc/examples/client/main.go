package main

import (
	"fmt"

	"gomicro/rpc/examples/pb"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	// conn, err := grpc.Dial("127.0.0.1:1701", grpc.WithInsecure())
	conn, err := grpc.Dial(
		"127.0.0.1:1701", grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
		grpc.WithStreamInterceptor(grpc_prometheus.StreamClientInterceptor),
	)
	if err != nil {
		panic(err)
	}

	client := pb.NewHelloServiceClient(conn)

	{
		ctx := metadata.NewContext(context.Background(), metadata.Pairs("guid", uuid.New()))
		response, err := client.NormalHello(ctx, &pb.HelloRequest{Greeting: "world"})
		fmt.Printf("normal hello: reponse=%#v, error=%v\n", response, err)
	}

	// {
	// 	ctx := metadata.NewContext(context.Background(), metadata.Pairs("tid", "normal-panic-request"))
	// 	response, err := client.PanicHello(ctx, &pb.HelloRequest{Greeting: "world"})
	// 	fmt.Printf("panic hello: response=%#v, error=%v\n", response, err)
	// }
}
