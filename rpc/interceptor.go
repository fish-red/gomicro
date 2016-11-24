package rpc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// UnaryInterceptorChain build the multi interceptors into one interceptor chain
func UnaryInterceptorChain(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, request interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (response interface{}, err error) {
		chain := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			chain = build(interceptors[i], chain, info)
		}

		return chain(ctx, request)
	}
}

func build(interceptor grpc.UnaryServerInterceptor, handler grpc.UnaryHandler, info *grpc.UnaryServerInfo) grpc.UnaryHandler {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return interceptor(ctx, request, info, handler)
	}
}
