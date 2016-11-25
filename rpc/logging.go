package rpc

import (
	"time"

	"gomicro/log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Logging interceptor for grpc
func Logging(ctx context.Context, request interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (response interface{}, err error) {
	start := time.Now()

	log.CtxInfof(ctx, "calling %s, request=%s", info.FullMethod, marshal(request))
	response, err = handler(ctx, request)
	log.CtxInfof(ctx, "finished %s, cost=%v, response=%v, err=%v", info.FullMethod, time.Since(start), response, err)

	return response, err
}
