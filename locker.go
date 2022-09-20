package state

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type locker uint32

const maxBackoff = 16

func (l *locker) Lock() {
	backoff := 1
	for !atomic.CompareAndSwapUint32((*uint32)(l), 0, 1) {
		// Leverage the exponential backoff algorithm, see https://en.wikipedia.org/wiki/Exponential_backoff.
		for i := 0; i < backoff; i++ {
			runtime.Gosched()
		}
		if backoff < maxBackoff {
			backoff <<= 1
		}
	}
}

func (l *locker) Unlock() {
	atomic.StoreUint32((*uint32)(l), 0)
}

// NewStateLock instantiates a spin-lock.
func NewStateLock() sync.Locker {
	return new(locker)
}
