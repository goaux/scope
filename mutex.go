package scope

import (
	"iter"
	"sync"
)

// Lock returns an iterator that yields a locked mu.
// The loop body executes exactly once.
// mu.Lock is called before the loop body executes.
// mu.Unlock is called after the loop body executes.
func Lock[Mu sync.Locker](mu Mu) iter.Seq[Mu] {
	return func(yield func(Mu) bool) {
		mu.Lock()
		defer mu.Unlock()
		yield(mu)
	}
}
