package rc

import (
	"gomicro/log"
	naming "gomicro/naming/etcd"

	"fmt"

	"reflect"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	serviceConns = newSafeMap()
)

// StartServiceConns start grpc connections with balancer
func StartServiceConns(address string, serviceList []string) {
	for _, serviceName := range serviceList {
		go func(name string) {
			// new a naming resolver
			resolver := naming.NewResolver(name)
			// new a grpc balancer
			balancer := grpc.RoundRobin(resolver)

			// new a grpc connection and buffer the connection
			conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBalancer(balancer))
			if err != nil {
				log.Printf("connect to '%s' service failed: %v", name, err)
			}
			serviceConns.Set(name, conn)
		}(serviceName)
	}
}

// CloseServiceConns close all established connections
func CloseServiceConns() {
	for _, conn := range serviceConns.List() {
		conn.Close()
	}
}

// DoRPC is helper func that make life easier
// ctx: context
// client: grpc client
// serviceName: name of service
// method: method name that you want to use
// request: grpc request
func DoRPC(ctx context.Context, client interface{}, serviceName string, method string, request interface{}) (response interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("do RPC '%s' error: %v", method, r)
		}
	}()

	// get a connection by the service name
	conn := serviceConns.Get(serviceName)
	if conn == nil {
		return nil, fmt.Errorf("service conn '%s' not found", serviceName)
	}

	// get NewServiceClient's reflect.Value
	vClient := reflect.ValueOf(client)
	var vParams []reflect.Value
	vParams = append(vParams, reflect.ValueOf(conn))

	// c[0] is serviceServer reflect.Value
	c := vClient.Call(vParams)

	// rpc parameter: context and request
	v := make([]reflect.Value, 2)
	v[0] = reflect.ValueOf(ctx)
	v[1] = reflect.ValueOf(request)

	// rpc method call
	f := c[0].MethodByName(method)
	result := f.Call(v)
	if !result[1].IsNil() {
		return nil, result[1].Interface().(error)
	}

	return result[0].Interface(), nil
}
