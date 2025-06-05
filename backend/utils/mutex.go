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
	mu, _ := km.mutexes.LoadOrStore(key, &sync.Mutex{})
	mutex := mu.(*sync.Mutex)

	// Try to get the lock with a timeout
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if mutex.TryLock() {
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}
	return false
}

// Unlock unlocks the mutex for the given key
func (km *KeyMutex) Unlock(key string) {
	if mu, ok := km.mutexes.Load(key); ok {
		mu.(*sync.Mutex).Unlock()
	}
}
