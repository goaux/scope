package scope_test

import (
	"fmt"
	"sync"

	"github.com/goaux/scope"
)

func ExampleLock() {
	var mu sync.Mutex
	for range scope.Lock(&mu) {
		// mu is locked.
		fmt.Println("pass 0")
	}
	fmt.Println("pass 1")
	// Output:
	// pass 0
	// pass 1
}

func ExampleLock_rWMutex() {
	var mu sync.RWMutex
	for range scope.Lock(&mu) {
		// mu is locked for writing.
		fmt.Println("pass 0")
		// mu will be unlocked at the end of the loop body.
	}
	fmt.Println("pass 1")
	for range scope.Lock(mu.RLocker()) {
		// mu is locked for reading.
		fmt.Println("pass 2")
	}
	fmt.Println("pass 3")
	// Output:
	// pass 0
	// pass 1
	// pass 2
	// pass 3
}

func ExampleLock_inspect() {
	var mu TestLocker
	for range scope.Lock(&mu) {
		// mu is locked for writing.
		fmt.Println("pass 0")
		// mu will be unlocked at the end of the loop body.
	}
	fmt.Println("pass 1")
	for range scope.Lock(mu.RLocker()) {
		// mu is locked for reading.
		fmt.Println("pass 2")
		// mu will be unlocked at the end of the loop body.
	}
	fmt.Println("pass 3")
	// Output:
	// Lock
	// pass 0
	// Unlock
	// pass 1
	// RLock
	// pass 2
	// RUnlock
	// pass 3
}

type TestLocker struct{}

func (*TestLocker) Lock()    { fmt.Println("Lock") }
func (*TestLocker) Unlock()  { fmt.Println("Unlock") }
func (*TestLocker) RLock()   { fmt.Println("RLock") }
func (*TestLocker) RUnlock() { fmt.Println("RUnlock") }

func (*TestLocker) RLocker() sync.Locker { return (*rlocker)(nil) }

type rlocker TestLocker

func (r *rlocker) Lock()   { (*TestLocker)(r).RLock() }
func (r *rlocker) Unlock() { (*TestLocker)(r).RUnlock() }
