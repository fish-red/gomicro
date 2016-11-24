package rpc

import (
	"runtime"

	"gomicro/log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	// MaxStackSize runtime输出缓冲
	MaxStackSize = 4096
)

// Recovery interceptor to handle grpc panic
func Recovery(ctx context.Context, request interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (response interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			// log stack
			stack := make([]byte, MaxStackSize)
			stack = stack[:runtime.Stack(stack, false)]
			log.CtxErrorf(ctx, "panic grpc invoke: %s, err=%v, stack:\n%s", info.FullMethod, r, string(stack))

			// if panic, set custom error to 'err', in order that client and sense it.
			err = grpc.Errorf(codes.Internal, "panic error: %v", r)
		}
	}()

	return handler(ctx, request)
}
