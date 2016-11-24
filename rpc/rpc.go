// grpc interceptor chain builder & middlewares. Now, we are only using unary rpc,
// this package only support unary interceptor.

package rpc

import (
	"google.golang.org/grpc"
)

// NewServer 创建grpc服务
func NewServer() *grpc.Server {
	return grpc.NewServer(grpc.UnaryInterceptor(UnaryInterceptorChain(Recovery, Logging)))
}
