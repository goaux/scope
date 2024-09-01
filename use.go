// Package scope provides iterators that automatically closing, unlocking or
// canceling resources such as files when they are no longer needed.
//
// The iterator returned by this package will execute the loop body exactly once.
//
// The defer idiom is commonly used in Go to allocate and release resources.
//
//	file, err := os.Open(name)
//	if err != nil {
//		panic(err)
//	}
//	defer file.Close()
//	_ = file // use file here
//
// Deferred actions are executed when the function returns, so the release of
// resources is delayed until the function returns.
//
// Sometimes that's what you want, and sometimes it's not.
//
// If you need to release resources immediately, you can use a function literal
// to specify the scope:
//
//	func() {
//		file, err := os.Open(name)
//		if err != nil {
//			return
//		}
//		defer file.Close()
//		_ = file // use file here
//	}()
//
// This package provides another solution for such situations.
// The following code snippet achieves the same result as the code above.
//
//	for file, err := range scope.Use2(os.Open(name)) {
//		if err != nil {
//			break
//		}
//		_ = file // use file here
//	}
package scope

import (
	"io"
	"iter"
)

// Use returns an iterator that yields resource as is.
// The loop body will be executed exactly once.
// resource.Close will be called after the loop body executed.
func Use[T io.Closer](resource T) iter.Seq[T] {
	return func(yield func(T) bool) {
		defer resource.Close()
		yield(resource)
	}
}

// Use2 returns an iterator that yields resource and err as is.
// The loop body will be executed exactly once.
// If err is nil, resource.Close will be called after the loop body executed.
func Use2[T io.Closer](resource T, err error) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		if err == nil {
			defer resource.Close()
		}
		yield(resource, err)
	}
}
