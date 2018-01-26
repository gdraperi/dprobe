// +build windows

package fs

import (
	"context"
	"os"
	"path/filepath"
)

func diskUsage(roots ...string) (Usage, error) ***REMOVED***
	var (
		size int64
	)

	// TODO(stevvooe): Support inodes (or equivalent) for windows.

	for _, root := range roots ***REMOVED***
		if err := filepath.Walk(root, func(path string, fi os.FileInfo, err error) error ***REMOVED***
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			size += fi.Size()
			return nil
		***REMOVED***); err != nil ***REMOVED***
			return Usage***REMOVED******REMOVED***, err
		***REMOVED***
	***REMOVED***

	return Usage***REMOVED***
		Size: size,
	***REMOVED***, nil
***REMOVED***

func diffUsage(ctx context.Context, a, b string) (Usage, error) ***REMOVED***
	var (
		size int64
	)

	if err := Changes(ctx, a, b, func(kind ChangeKind, _ string, fi os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if kind == ChangeKindAdd || kind == ChangeKindModify ***REMOVED***
			size += fi.Size()

			return nil

		***REMOVED***
		return nil
	***REMOVED***); err != nil ***REMOVED***
		return Usage***REMOVED******REMOVED***, err
	***REMOVED***

	return Usage***REMOVED***
		Size: size,
	***REMOVED***, nil
***REMOVED***
