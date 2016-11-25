package etcd

import (
	"fmt"
	"log"

	lib "gomicro/naming/lib"

	etcd "github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"google.golang.org/grpc/naming"
)

// EtcdWatcher is the implementaion of grpc.naming.Watcher
type EtcdWatcher struct {
	resolver *EtcdResolver // the etcd resolver
	client   *etcd.Client  // the etcd client
	addrs    []string      // the addresses of service
}

// Close do nothing
func (watcher *EtcdWatcher) Close() {}

// Next to return to the updates
func (watcher *EtcdWatcher) Next() ([]*naming.Update, error) {
	// key is the etcd key/value dir to watch
	key := fmt.Sprintf("/%s/%s", Prefix, watcher.resolver.ServiceName)
	// new a keys api
	keyAPI := etcd.NewKeysAPI(*watcher.client)

	// if the addresses is nil means it's initially called
	if watcher.addrs == nil {
		// query addresses from etcd
		s, err := keyAPI.Get(context.Background(), key, &etcd.GetOptions{Recursive: true})
		if err != nil {
			log.Printf("naming: get value of key '%s' from etcd error: %s\n", key, err.Error())
		}
		// extract the addresses from the etcd response
		addrs, empty := extractAddrs(s)
		// drop the empty directory in etcd
		dropEmptyDir(keyAPI, empty)

		// addrs is not empty, return the updates
		if len(addrs) != 0 {
			watcher.addrs = addrs
			return lib.GenUpdates([]string{}, addrs), nil
		}
		// addrs is empty, should to watch new data
	}

	// generate the etcd watcher
	w := keyAPI.Watcher(key, &etcd.WatcherOptions{Recursive: true})
	for {
		_, err := w.Next(context.Background())
		if err == nil {
			// query addresses from etcd
			s, err := keyAPI.Get(context.Background(), key, &etcd.GetOptions{Recursive: true})
			if err != nil {
				continue
			}
			// extract the addresses from the etcd response
			addrs, empty := extractAddrs(s)
			// drop the empty directory in etcd
			dropEmptyDir(keyAPI, empty)

			// 1. delete the addresses in cache, not in etcd
			// 2. add the addresses in etcd, not in cache
			updates := lib.GenUpdates(watcher.addrs, addrs)

			// update the watcher.addrs
			watcher.addrs = addrs
			// if addrs updated, return it
			if len(updates) != 0 {
				return updates, nil
			}
		}
	}
}

// helper function to extract addrs from etcd response
func extractAddrs(s *etcd.Response) (addrs, empty []string) {
	addrs, empty = []string{}, []string{}

	// check the response
	if s == nil || s.Node == nil || s.Node.Nodes == nil || len(s.Node.Nodes) == 0 {
		return addrs, empty
	}

	for _, node := range s.Node.Nodes {
		// node should contain host & port
		host, port := "", ""
		for _, v := range node.Nodes {
			// get the last 4 characters
			what := v.Key[len(v.Key)-4 : len(v.Key)]
			if what == "host" {
				host = v.Value
			}
			if what == "port" {
				port = v.Value
			}
		}

		// if one of host&port has no value, the addr is set partly, should not return
		if host != "" && port != "" {
			addrs = append(addrs, fmt.Sprintf("%s:%s", host, port))
		}
		if host == "" && port == "" {
			empty = append(empty, node.Key)
		}
	}

	return addrs, empty
}

func dropEmptyDir(keyAPI etcd.KeysAPI, empty []string) {
	if keyAPI == nil || len(empty) == 0 {
		return
	}

	for _, key := range empty {
		_, err := keyAPI.Delete(context.Background(), key, &etcd.DeleteOptions{Recursive: true})
		if err != nil {
			log.Printf("naming: delete empty service dir (%s) error: %s\n", key, err.Error())
		}
	}
}
