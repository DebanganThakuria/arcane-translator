package utils

import (
	"sync"
	"time"
)

var Mutex KeyMutex

// KeyMutex provides mutex locks for specific keys
type KeyMutex struct {
	mutexes sync.Map // maps key -> *sync.Mutex
}

// TryLock locks the mutex for the given key
func (km *KeyMutex) TryLock(key string, timeout time.Duration) bool {
	mu, loaded := km.mutexes.LoadOrStore(key, &sync.Mutex{})
	if !loaded {
		return false
	}

	// Try to get the lock with a timeout
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if mu.(*sync.Mutex).TryLock() {
			return true
		}
		time.Sleep(10 * time.Millisecond) // Small delay between attempts
	}
	return false
}

// Unlock unlocks the mutex for the given key
func (km *KeyMutex) Unlock(key string) {
	mu, _ := km.mutexes.Load(key)
	mu.(*sync.Mutex).Unlock()
}
