package main

import "sync"

// VDR is a simplified DID Document registry
type VDR struct {
	// keys maps nodeIDs to Keys
	keys  map[string]string
	mutex sync.Mutex
}

func NewVDR() VDR {
	return VDR {
		keys: map[string]string{},
		mutex: sync.Mutex{},
	}
}

func (vdr *VDR) Matches(nodeID string, key string) bool {
	vdr.mutex.Lock()
	defer vdr.mutex.Unlock()

	key, ok := vdr.keys[nodeID]
	if !ok {
		vdr.keys[nodeID] = key
	}
	return vdr.keys[nodeID] == key
}

func (vdr *VDR) Put(nodeID string, key string) {
	vdr.mutex.Lock()
	defer vdr.mutex.Unlock()

	vdr.keys[nodeID] = key
}