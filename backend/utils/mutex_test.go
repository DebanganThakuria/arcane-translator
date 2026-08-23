package utils

import (
	"sync"
	"testing"
	"time"
)

func TestKeyMutexExcludesSameKey(t *testing.T) {
	var km KeyMutex

	if !km.TryLock("novel-1", 10*time.Millisecond) {
		t.Fatal("first lock should be granted")
	}

	if km.TryLock("novel-1", 10*time.Millisecond) {
		t.Fatal("second lock on the same key should be refused")
	}

	km.Unlock("novel-1")

	if !km.TryLock("novel-1", 10*time.Millisecond) {
		t.Fatal("lock should be available again after unlock")
	}
	km.Unlock("novel-1")
}

func TestKeyMutexIndependentKeys(t *testing.T) {
	var km KeyMutex

	if !km.TryLock("novel-1", 10*time.Millisecond) {
		t.Fatal("lock on novel-1 refused")
	}
	if !km.TryLock("novel-2", 10*time.Millisecond) {
		t.Fatal("a different key must not be blocked by novel-1")
	}

	km.Unlock("novel-1")
	km.Unlock("novel-2")
}

// The previous implementation kept one mutex per key forever, so a long-running
// server accumulated an entry for every novel URL it had ever seen.
func TestKeyMutexReleasesEntries(t *testing.T) {
	var km KeyMutex

	for i := range 100 {
		key := "novel-" + string(rune('a'+i%26)) + string(rune('0'+i/26))
		if !km.TryLock(key, 10*time.Millisecond) {
			t.Fatalf("lock %d refused", i)
		}
		km.Unlock(key)
	}

	km.mu.Lock()
	remaining := len(km.locks)
	km.mu.Unlock()

	if remaining != 0 {
		t.Errorf("%d entries retained after every lock was released, want 0", remaining)
	}
}

// A refused acquisition must not leave its reference behind either.
func TestKeyMutexReleasesAfterTimeout(t *testing.T) {
	var km KeyMutex

	if !km.TryLock("busy", 10*time.Millisecond) {
		t.Fatal("first lock refused")
	}
	if km.TryLock("busy", 5*time.Millisecond) {
		t.Fatal("second lock should have timed out")
	}
	km.Unlock("busy")

	km.mu.Lock()
	remaining := len(km.locks)
	km.mu.Unlock()

	if remaining != 0 {
		t.Errorf("%d entries retained after a timed-out acquisition, want 0", remaining)
	}
}

// Unlocking something never locked used to panic the whole server.
func TestKeyMutexUnlockUnheldIsSafe(t *testing.T) {
	var km KeyMutex
	km.Unlock("never-locked")

	if !km.TryLock("never-locked", 10*time.Millisecond) {
		t.Fatal("key should still be usable")
	}
	km.Unlock("never-locked")
}

// Run under -race: exactly one goroutine may hold a given key at a time.
func TestKeyMutexConcurrent(t *testing.T) {
	var (
		km      KeyMutex
		wg      sync.WaitGroup
		mu      sync.Mutex
		holding int
		granted int
	)

	for range 50 {
		wg.Go(func() {
			if !km.TryLock("shared", 2*time.Second) {
				return
			}

			mu.Lock()
			holding++
			granted++
			if holding > 1 {
				t.Errorf("%d goroutines held the same key at once", holding)
			}
			mu.Unlock()

			time.Sleep(time.Millisecond)

			mu.Lock()
			holding--
			mu.Unlock()

			km.Unlock("shared")
		})
	}

	wg.Wait()

	if granted == 0 {
		t.Fatal("no goroutine ever acquired the lock")
	}
}
