package consul

import (
	"fmt"
	"log"
	"time"

	"gomicro/naming/lib"

	consul "github.com/hashicorp/consul/api"
	"google.golang.org/grpc/naming"
)

// ConsulWatcher is the implementation of grpc.naming.Watcher
type ConsulWatcher struct {
	resolver  *ConsulResolver // the consul resolver
	client    *consul.Client  // the consul client
	lastIndex uint64          // LastIndex to watch consul

	// addrs is the service address cache
	// before check: every value shoud be 1
	// after check: 1 - deleted  2 - nothing  3 - new added
	addrs []string
}

// Close do nonthing
func (watcher *ConsulWatcher) Close() {}

// Next to return the updates
func (watcher *ConsulWatcher) Next() ([]*naming.Update, error) {
	// watcher.addrs is nil means it is initial called
	if watcher.addrs == nil {
		// return addrs to balancer, use ticker to query consul till data gotten
		addrs, lastIndex, _ := watcher.queryConsul(nil)
		if len(addrs) != 0 {
			watcher.addrs = addrs
			watcher.lastIndex = lastIndex
			return lib.GenUpdates([]string{}, addrs), nil
		}
	}

	for {
		// query the consul to get the addresses for service
		addrs, lastIndex, err := watcher.queryConsul(&consul.QueryOptions{WaitIndex: watcher.lastIndex})
		if err != nil {
			log.Printf("naming: get addresses of '%s' from consul error: %s\n", watcher.resolver.ServiceName, err.Error())
			time.Sleep(1 * time.Second)
			continue
		}

		// 1. delete the addresses in cache not in consul
		// 2. add the addresses in consul not in cache
		updates := lib.GenUpdates(watcher.addrs, addrs)

		// update the addresses and last index in cache
		watcher.addrs = addrs
		watcher.lastIndex = lastIndex
		if len(updates) != 0 {
			return updates, nil
		}
	}
}

// queryConsul is helper function to query consul
func (watcher *ConsulWatcher) queryConsul(option *consul.QueryOptions) ([]string, uint64, error) {
	// query the addresses from consul by service name
	services, meta, err := watcher.client.Health().Service(watcher.resolver.ServiceName, "", true, option)
	if err != nil {
		return nil, 0, err
	}

	var addrs []string
	for _, s := range services {
		// addr should like: 127.0.0.1:8001
		addrs = append(addrs, fmt.Sprintf("%s:%d", s.Service.Address, s.Service.Port))
	}

	return addrs, meta.LastIndex, nil
}
