package metadata

import (
	"context"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

type transactionKey struct***REMOVED******REMOVED***

// WithTransactionContext returns a new context holding the provided
// bolt transaction. Functions which require a bolt transaction will
// first check to see if a transaction is already created on the
// context before creating their own.
func WithTransactionContext(ctx context.Context, tx *bolt.Tx) context.Context ***REMOVED***
	return context.WithValue(ctx, transactionKey***REMOVED******REMOVED***, tx)
***REMOVED***

type transactor interface ***REMOVED***
	View(fn func(*bolt.Tx) error) error
	Update(fn func(*bolt.Tx) error) error
***REMOVED***

// view gets a bolt db transaction either from the context
// or starts a new one with the provided bolt database.
func view(ctx context.Context, db transactor, fn func(*bolt.Tx) error) error ***REMOVED***
	tx, ok := ctx.Value(transactionKey***REMOVED******REMOVED***).(*bolt.Tx)
	if !ok ***REMOVED***
		return db.View(fn)
	***REMOVED***
	return fn(tx)
***REMOVED***

// update gets a writable bolt db transaction either from the context
// or starts a new one with the provided bolt database.
func update(ctx context.Context, db transactor, fn func(*bolt.Tx) error) error ***REMOVED***
	tx, ok := ctx.Value(transactionKey***REMOVED******REMOVED***).(*bolt.Tx)
	if !ok ***REMOVED***
		return db.Update(fn)
	***REMOVED*** else if !tx.Writable() ***REMOVED***
		return errors.Wrap(bolt.ErrTxNotWritable, "unable to use transaction from context")
	***REMOVED***
	return fn(tx)
***REMOVED***
