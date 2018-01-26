package fsutil

import (
	"hash"
	"os"

	"golang.org/x/net/context"
)

type walkerFn func(ctx context.Context, pathC chan<- *currentPath) error

func Changes(ctx context.Context, a, b walkerFn, changeFn ChangeFunc) error ***REMOVED***
	return nil
***REMOVED***

type HandleChangeFn func(ChangeKind, string, os.FileInfo, error) error

type ContentHasher func(*Stat) (hash.Hash, error)

func GetWalkerFn(root string) walkerFn ***REMOVED***
	return func(ctx context.Context, pathC chan<- *currentPath) error ***REMOVED***
		return Walk(ctx, root, nil, func(path string, f os.FileInfo, err error) error ***REMOVED***
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			p := &currentPath***REMOVED***
				path: path,
				f:    f,
			***REMOVED***

			select ***REMOVED***
			case <-ctx.Done():
				return ctx.Err()
			case pathC <- p:
				return nil
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func emptyWalker(ctx context.Context, pathC chan<- *currentPath) error ***REMOVED***
	return nil
***REMOVED***
