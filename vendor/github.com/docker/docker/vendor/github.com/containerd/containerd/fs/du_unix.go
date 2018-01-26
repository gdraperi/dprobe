// +build !windows

package fs

import (
	"context"
	"os"
	"path/filepath"
	"syscall"
)

type inode struct ***REMOVED***
	// TODO(stevvooe): Can probably reduce memory usage by not tracking
	// device, but we can leave this right for now.
	dev, ino uint64
***REMOVED***

func diskUsage(roots ...string) (Usage, error) ***REMOVED***

	var (
		size   int64
		inodes = map[inode]struct***REMOVED******REMOVED******REMOVED******REMOVED*** // expensive!
	)

	for _, root := range roots ***REMOVED***
		if err := filepath.Walk(root, func(path string, fi os.FileInfo, err error) error ***REMOVED***
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			stat := fi.Sys().(*syscall.Stat_t)

			inoKey := inode***REMOVED***dev: uint64(stat.Dev), ino: uint64(stat.Ino)***REMOVED***
			if _, ok := inodes[inoKey]; !ok ***REMOVED***
				inodes[inoKey] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
				size += fi.Size()
			***REMOVED***

			return nil
		***REMOVED***); err != nil ***REMOVED***
			return Usage***REMOVED******REMOVED***, err
		***REMOVED***
	***REMOVED***

	return Usage***REMOVED***
		Inodes: int64(len(inodes)),
		Size:   size,
	***REMOVED***, nil
***REMOVED***

func diffUsage(ctx context.Context, a, b string) (Usage, error) ***REMOVED***
	var (
		size   int64
		inodes = map[inode]struct***REMOVED******REMOVED******REMOVED******REMOVED*** // expensive!
	)

	if err := Changes(ctx, a, b, func(kind ChangeKind, _ string, fi os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if kind == ChangeKindAdd || kind == ChangeKindModify ***REMOVED***
			stat := fi.Sys().(*syscall.Stat_t)

			inoKey := inode***REMOVED***dev: uint64(stat.Dev), ino: uint64(stat.Ino)***REMOVED***
			if _, ok := inodes[inoKey]; !ok ***REMOVED***
				inodes[inoKey] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
				size += fi.Size()
			***REMOVED***

			return nil

		***REMOVED***
		return nil
	***REMOVED***); err != nil ***REMOVED***
		return Usage***REMOVED******REMOVED***, err
	***REMOVED***

	return Usage***REMOVED***
		Inodes: int64(len(inodes)),
		Size:   size,
	***REMOVED***, nil
***REMOVED***
