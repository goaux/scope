package scope_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/goaux/scope"
)

func ExampleContext() {
	var g sync.WaitGroup
	for ctx := range scope.Context(context.WithTimeout(context.TODO(), 100*time.Millisecond)) {
		fmt.Println("pass 0")
		g.Add(1)
		go func() {
			defer g.Done()
			<-ctx.Done()
			fmt.Println("pass 1", ctx.Err())
		}()
		// ctx will be canceled at the end of the loop body.
	}
	g.Wait()
	fmt.Println("ok")
	// Output:
	// pass 0
	// pass 1 context canceled
	// ok
}

func ExampleWithCancel() {
	var g sync.WaitGroup
	for ctx := range scope.WithCancel(context.TODO()) {
		fmt.Println("pass 0")
		g.Add(1)
		go func() {
			defer g.Done()
			<-ctx.Done()
			fmt.Println("pass 1")
		}()
		// ctx will be canceled at the end of the loop body.
	}
	g.Wait()
	fmt.Println("ok")
	// Output:
	// pass 0
	// pass 1
	// ok
}

func ExampleWithTimeout() {
	var g sync.WaitGroup
	for ctx := range scope.WithTimeout(context.TODO(), 100*time.Millisecond) {
		fmt.Println("pass 0")
		g.Add(1)
		go func() {
			defer g.Done()
			<-ctx.Done()
			fmt.Println("pass 1", ctx.Err())
		}()
		// ctx will be canceled at the end of the loop body.
	}
	g.Wait()
	fmt.Println("ok")
	// Output:
	// pass 0
	// pass 1 context canceled
	// ok
}

func ExampleWithTimeoutCause() {
	var g sync.WaitGroup
	for ctx := range scope.WithTimeoutCause(context.TODO(), 100*time.Millisecond, errors.New("hello")) {
		fmt.Println("pass 0")
		g.Add(1)
		go func() {
			defer g.Done()
			<-ctx.Done()
			fmt.Println("pass 1", ctx.Err())
		}()
		// ctx will be canceled at the end of the loop body.
	}
	g.Wait()
	fmt.Println("ok")
	// Output:
	// pass 0
	// pass 1 context canceled
	// ok
}

func ExampleWithDeadline() {
	var g sync.WaitGroup
	for ctx := range scope.WithDeadline(context.TODO(), time.Now().Add(100*time.Millisecond)) {
		fmt.Println("pass 0")
		g.Add(1)
		go func() {
			defer g.Done()
			<-ctx.Done()
			fmt.Println("pass 1", ctx.Err())
		}()
		// ctx will be canceled at the end of the loop body.
	}
	g.Wait()
	fmt.Println("ok")
	// Output:
	// pass 0
	// pass 1 context canceled
	// ok
}
