package rc

import (
	"sync"

	"google.golang.org/grpc"
)

type safeMap struct {
	lock    *sync.RWMutex
	clients map[string]*grpc.ClientConn
}

// NewSafeMap return a new thread safe map
func newSafeMap() *safeMap {
	return &safeMap{
		lock:    new(sync.RWMutex),
		clients: make(map[string]*grpc.ClientConn),
	}
}

// Get from maps return the key's value
func (m *safeMap) Get(k string) *grpc.ClientConn {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if val, ok := m.clients[k]; ok {
		return val
	}
	return nil
}

// Set Maps the given key and value. Returns false if the key is already in the map and changes nothing.
func (m *safeMap) Set(key string, v *grpc.ClientConn) bool {
	m.lock.Lock()
	defer m.lock.Unlock()

	if val, ok := m.clients[key]; !ok {
		m.clients[key] = v
	} else if val != v {
		m.clients[key] = v
	} else {
		return false
	}
	return true
}

// Check returns true if key is exist in the map.
func (m *safeMap) Check(key string) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if _, ok := m.clients[key]; ok {
		return true
	}
	return false
}

// Delete remove the key's value from the map
func (m *safeMap) Delete(key string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.clients[key]; ok {
		delete(m.clients, key)
	}
}

// List returns all the clients
func (m *safeMap) List() map[string]*grpc.ClientConn {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.clients
}
