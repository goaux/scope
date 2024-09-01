package scope

import (
	"context"
	"iter"
	"time"
)

// Context returns an iterator that yields ctx as is.
// The loop body will be executed exactly once.
// The cancel will be called after the loop body executed.
func Context(
	ctx context.Context,
	cancel context.CancelFunc,
) iter.Seq[context.Context] {
	return func(yield func(context.Context) bool) {
		defer cancel()
		yield(ctx)
	}
}

// WithCancel returns Context(context.WithCancel(parent)).
//
// See [Context], [context.WithCancel].
func WithCancel(
	parent context.Context,
) iter.Seq[context.Context] {
	return Context(context.WithCancel(parent))
}

// WithTimeout returns Context(context.WithTimeout(parent, timeout)).
//
// See [Context], [context.WithTimeout].
func WithTimeout(
	parent context.Context,
	timeout time.Duration,
) iter.Seq[context.Context] {
	return Context(context.WithTimeout(parent, timeout))
}

// WithTimeoutCause returns Context(context.WithTimeoutCause(parent, timeout, cause)).
//
// See [Context], [context.WithTimeoutCause].
func WithTimeoutCause(
	parent context.Context,
	timeout time.Duration,
	cause error,
) iter.Seq[context.Context] {
	return Context(context.WithTimeoutCause(parent, timeout, cause))
}

// WithDeadline returns Context(context.WithDeadline(parent, d)).
//
// See [Context], [context.WithDeadline].
func WithDeadline(
	parent context.Context,
	d time.Time,
) iter.Seq[context.Context] {
	return Context(context.WithDeadline(parent, d))
}
