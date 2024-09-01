package scope

import "iter"

// Rollbacker is a base type that [Tx] and [Tx2] accept.
type Rollbacker interface {
	Rollback() error
}

// Context returns an iterator that yields tx as is.
// The loop body will be executed exactly once.
// tx.Rollback will be called after the loop body executed.
func Tx[Tx Rollbacker](tx Tx) iter.Seq[Tx] {
	return func(yield func(Tx) bool) {
		defer tx.Rollback()
		yield(tx)
	}
}

// Context returns an iterator that yields tx and err as is.
// The loop body will be executed exactly once.
// If err is nil, tx.Rollback will be called after the loop body executed.
func Tx2[Tx Rollbacker](tx Tx, err error) iter.Seq2[Tx, error] {
	return func(yield func(Tx, error) bool) {
		if err == nil {
			defer tx.Rollback()
		}
		yield(tx, err)
	}
}
