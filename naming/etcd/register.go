package etcd

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	etcd "github.com/coreos/etcd/client"
)

var (
	// Prefix for the service key in etcd
	Prefix     = "naming"
	keyAPI     etcd.KeysAPI
	serviceKey string
)

// Register is the helper function to self-register service into Etcd/Consul server
// should call Unregister when pocess stop
// name - service name
// host - service host
// port - service port
// target - etcd dial address, for example: "http://127.0.0.1:2379;http://127.0.0.1:12379"
// interval - interval of self-register to etcd
// ttl - ttl of the register information
func Register(name string, host string, port int, target string, interval time.Duration, ttl int) error {
	// get endpoints from register dial address
	endpoints := strings.Split(target, ",")
	conf := etcd.Config{Endpoints: endpoints}

	// new a etcd client
	client, err := etcd.New(conf)
	if err != nil {
		return fmt.Errorf("naming: create etcd client error: %v", err)
	}
	// new a keys api
	keyAPI = etcd.NewKeysAPI(client)

	serviceID := fmt.Sprintf("%s-%s-%d", name, host, port)
	serviceKey = fmt.Sprintf("/%s/%s/%s", Prefix, name, serviceID)
	hostKey := fmt.Sprintf("/%s/%s/%s/host", Prefix, name, serviceID)
	portKey := fmt.Sprintf("/%s/%s/%s/port", Prefix, name, serviceID)

	go func() {
		// invoke self-register with ticker
		ticker := time.NewTicker(interval)

		for {
			// wait for next tick point
			<-ticker.C

			_, err := keyAPI.Get(context.Background(), serviceKey, &etcd.GetOptions{Recursive: true})
			if err != nil {
				if _, err := keyAPI.Set(context.Background(), hostKey, host, nil); err != nil {
					log.Printf("naming: re-register service '%s' host to etcd error: %s\n", name, err.Error())
				}
				if _, err := keyAPI.Set(context.Background(), portKey, fmt.Sprintf("%d", port), nil); err != nil {
					log.Printf("naming: re-register service '%s' port to etcd error: %s\n", name, err.Error())
				}

				option := &etcd.SetOptions{TTL: time.Duration(ttl) * time.Second, PrevExist: etcd.PrevExist, Dir: true}
				if _, err := keyAPI.Set(context.Background(), serviceKey, "", option); err != nil {
					log.Printf("naming: set service '%s' ttl to etcd error: %s\n", name, err.Error())
				}
			} else {
				// refresh set to true for not notifying the watcher
				option := &etcd.SetOptions{
					TTL:       time.Duration(ttl) * time.Second,
					PrevExist: etcd.PrevExist,
					Dir:       true,
					Refresh:   true,
				}
				if _, err := keyAPI.Set(context.Background(), serviceKey, "", option); err != nil {
					log.Printf("naming: set service '%s' ttl to etcd error: %s\n", name, err.Error())
				}
			}

		}
	}()

	// initial register
	if _, err := keyAPI.Set(context.Background(), hostKey, host, nil); err != nil {
		log.Printf("naming: initial service '%s' host to etcd error: %s\n", name, err.Error())
	}
	if _, err := keyAPI.Set(context.Background(), portKey, fmt.Sprintf("%d", port), nil); err != nil {
		log.Printf("naming: initial service '%s' port to etcd error: %s\n", name, err.Error())
	}

	option := &etcd.SetOptions{TTL: time.Duration(ttl) * time.Second, PrevExist: etcd.PrevExist, Dir: true}
	if _, err := keyAPI.Set(context.Background(), serviceKey, "", option); err != nil {
		log.Printf("naming: set service '%s' ttl to etcd error: %s\n", name, err.Error())
	}

	return nil
}

// UnRegister delete service from etcd
func UnRegister() error {
	_, err := keyAPI.Delete(context.Background(), serviceKey, &etcd.DeleteOptions{Recursive: true})
	if err != nil {
		log.Println("naming: unregister service error: ", err.Error())
	} else {
		log.Println("naming: unregistered service from etcd server.")
	}

	return err
}
