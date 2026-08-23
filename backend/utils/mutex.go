package utils

import (
	"sync"
	"time"
)

// Mutex guards long-running per-resource work, so two requests cannot translate
// the same chapter or refresh the same novel at once.
var Mutex KeyMutex

// KeyMutex hands out one lock per key.
//
// Entries are reference counted and removed once the last holder releases them.
// A plain map of key to mutex never shrinks, which on this workload means one
// live mutex per novel URL ever seen.
type KeyMutex struct {
	mu    sync.Mutex
	locks map[string]*keyedLock
}

type keyedLock struct {
	// Buffered to one slot: sending acquires, receiving releases. A channel
	// rather than sync.Mutex because it supports waiting with a timeout.
	ch chan struct{}
	// Number of goroutines holding or waiting for this lock.
	refs int
}

// acquire returns the lock for key, creating it if needed, and records the
// caller as a referrer.
func (km *KeyMutex) acquire(key string) *keyedLock {
	km.mu.Lock()
	defer km.mu.Unlock()

	if km.locks == nil {
		km.locks = make(map[string]*keyedLock)
	}

	lock, ok := km.locks[key]
	if !ok {
		lock = &keyedLock{ch: make(chan struct{}, 1)}
		km.locks[key] = lock
	}
	lock.refs++

	return lock
}

// release drops one reference and deletes the entry when it reaches zero.
func (km *KeyMutex) release(key string) {
	km.mu.Lock()
	defer km.mu.Unlock()

	lock, ok := km.locks[key]
	if !ok {
		return
	}

	lock.refs--
	if lock.refs <= 0 {
		delete(km.locks, key)
	}
}

// TryLock takes the lock for key, waiting up to timeout. It reports whether the
// lock was taken; only call Unlock when it returns true.
func (km *KeyMutex) TryLock(key string, timeout time.Duration) bool {
	lock := km.acquire(key)

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case lock.ch <- struct{}{}:
		return true
	case <-timer.C:
		// Give up the reference we took, so an abandoned key is not retained.
		km.release(key)
		return false
	}
}

// Unlock releases a lock previously taken by TryLock.
func (km *KeyMutex) Unlock(key string) {
	km.mu.Lock()
	lock, ok := km.locks[key]
	km.mu.Unlock()

	if !ok {
		return
	}

	select {
	case <-lock.ch:
	default:
		// Not currently held; releasing again would be a bug in the caller
		// rather than something to panic the whole server over.
		return
	}

	km.release(key)
}
