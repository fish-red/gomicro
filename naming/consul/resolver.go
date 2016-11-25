package consul

import (
	"errors"
	"fmt"

	consul "github.com/hashicorp/consul/api"
	"google.golang.org/grpc/naming"
)

// ConsulResolver is the implementaion of grpc.naming.Resolver
type ConsulResolver struct {
	ServiceName string
}

// NewResolver return ConsulResolver with service name
func NewResolver(serviceName string) *ConsulResolver {
	return &ConsulResolver{ServiceName: serviceName}
}

// Resolve to resolve the service from consul, target is the dial address of consul
func (resolver *ConsulResolver) Resolve(target string) (naming.Watcher, error) {
	if resolver.ServiceName == "" {
		return nil, errors.New("naming: no service name provided")
	}

	// generate consul client
	conf := &consul.Config{
		Scheme:  "http",
		Address: target,
	}
	client, err := consul.NewClient(conf)
	if err != nil {
		return nil, fmt.Errorf("naming: creat consul error: %v", err)
	}

	// return ConsulWatcher
	watcher := &ConsulWatcher{
		resolver: resolver,
		client:   client,
	}
	return watcher, nil
}
