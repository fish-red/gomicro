package consul

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	consul "github.com/hashicorp/consul/api"
)

// Register is the helper function to self-register service into Etcd/Consul server
// name - service name
// host - service host
// port - service port
// target - consul dial address, for example: "127.0.0.1:8500"
// interval - interval of self-register to etcd
// ttl - ttl of the register information
func Register(name string, host string, port int, target string, interval time.Duration, ttl int) error {
	config := &consul.Config{Scheme: "http", Address: target}
	client, err := consul.NewClient(config)
	if err != nil {
		return fmt.Errorf("naming: create consul client error: %v", err)
	}

	serviceID := fmt.Sprintf("%s-%s-%d", name, host, port)
	// unregister if meet signhup
	go func() {
		signalChannel := make(chan os.Signal, 1)
		// register the signals
		signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
		// receive a signal
		receivedSignal := <-signalChannel
		log.Println("naming: receive signal: ", receivedSignal)

		// unregister the service
		err := client.Agent().ServiceDeregister(serviceID)
		if err != nil {
			log.Println("naming: unregister service error: ", err.Error())
		} else {
			log.Println("naming: unregistered service from consul server.")
		}

		// check if the service is unregistered
		err = client.Agent().CheckDeregister(serviceID)
		if err != nil {
			log.Println("naming: unregister check error: ", err.Error())
		}

		// exit the process
		os.Exit(-1)
	}()

	// goroutine to update ttl
	go func() {
		ticker := time.NewTicker(interval)
		for {
			// capture a tick point
			<-ticker.C
			err = client.Agent().UpdateTTL(serviceID, "", "passing")
			if err != nil {
				log.Println("naming: update ttl of service error: ", err.Error())
			}
		}
	}()

	// inital register service
	serviceRegister := &consul.AgentServiceRegistration{
		ID:      serviceID,
		Name:    name,
		Address: host,
		Port:    port,
	}
	err = client.Agent().ServiceRegister(serviceRegister)
	if err != nil {
		return fmt.Errorf("naming: initial register service '%s' host to consul error: %s", name, err.Error())
	}

	// inital register service check
	check := consul.AgentServiceCheck{TTL: fmt.Sprintf("%ds", ttl), Status: "passing"}
	err = client.Agent().CheckRegister(&consul.AgentCheckRegistration{
		ID:                serviceID,
		Name:              name,
		ServiceID:         serviceID,
		AgentServiceCheck: check,
	})
	if err != nil {
		return fmt.Errorf("naming: initial register service check to consul error: %s", err.Error())
	}

	return nil
}
