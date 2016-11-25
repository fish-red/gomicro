package etcd

import (
	"errors"
	"fmt"
	"strings"

	etcd "github.com/coreos/etcd/client"
	"google.golang.org/grpc/naming"
)

// EtcdResolver is the implementaion of grpc.naming.Resolver
type EtcdResolver struct {
	ServiceName string // service name for resolving
}

// NewResolver return EtcdResolver with service name
func NewResolver(serviceName string) *EtcdResolver {
	return &EtcdResolver{serviceName}
}

// Resolve to resolve the service from etcd, target is the dial address of etcd
// target example: "http://127.0.0.1:2379;http://127.0.0.1:12379;http://127.0.0.1:22379"
func (resolver *EtcdResolver) Resolve(target string) (naming.Watcher, error) {
	if resolver.ServiceName == "" {
		return nil, errors.New("naming: no service name provided")
	}

	// new an etcd client
	endpoints := strings.Split(target, ",")
	conf := etcd.Config{Endpoints: endpoints}
	client, err := etcd.New(conf)
	if err != nil {
		return nil, fmt.Errorf("naming: create etcd error: %s", err.Error())
	}

	// return EtcdWatcher
	watcher := &EtcdWatcher{
		resolver: resolver,
		client:   &client,
	}

	return watcher, nil
}
