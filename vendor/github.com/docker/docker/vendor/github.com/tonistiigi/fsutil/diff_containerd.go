package fsutil

import (
	"os"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

// Everything below is copied from containerd/fs. TODO: remove duplication @dmcgowan

// Const redefined because containerd/fs doesn't build on !linux

// ChangeKind is the type of modification that
// a change is making.
type ChangeKind int

const (
	// ChangeKindAdd represents an addition of
	// a file
	ChangeKindAdd ChangeKind = iota

	// ChangeKindModify represents a change to
	// an existing file
	ChangeKindModify

	// ChangeKindDelete represents a delete of
	// a file
	ChangeKindDelete
)

// ChangeFunc is the type of function called for each change
// computed during a directory changes calculation.
type ChangeFunc func(ChangeKind, string, os.FileInfo, error) error

type currentPath struct ***REMOVED***
	path string
	f    os.FileInfo
	//	fullPath string
***REMOVED***

// doubleWalkDiff walks both directories to create a diff
func doubleWalkDiff(ctx context.Context, changeFn ChangeFunc, a, b walkerFn) (err error) ***REMOVED***
	g, ctx := errgroup.WithContext(ctx)

	var (
		c1 = make(chan *currentPath, 128)
		c2 = make(chan *currentPath, 128)

		f1, f2 *currentPath
		rmdir  string
	)
	g.Go(func() error ***REMOVED***
		defer close(c1)
		return a(ctx, c1)
	***REMOVED***)
	g.Go(func() error ***REMOVED***
		defer close(c2)
		return b(ctx, c2)
	***REMOVED***)
	g.Go(func() error ***REMOVED***
	loop0:
		for c1 != nil || c2 != nil ***REMOVED***
			if f1 == nil && c1 != nil ***REMOVED***
				f1, err = nextPath(ctx, c1)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				if f1 == nil ***REMOVED***
					c1 = nil
				***REMOVED***
			***REMOVED***

			if f2 == nil && c2 != nil ***REMOVED***
				f2, err = nextPath(ctx, c2)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				if f2 == nil ***REMOVED***
					c2 = nil
				***REMOVED***
			***REMOVED***
			if f1 == nil && f2 == nil ***REMOVED***
				continue
			***REMOVED***

			var f os.FileInfo
			k, p := pathChange(f1, f2)
			switch k ***REMOVED***
			case ChangeKindAdd:
				if rmdir != "" ***REMOVED***
					rmdir = ""
				***REMOVED***
				f = f2.f
				f2 = nil
			case ChangeKindDelete:
				// Check if this file is already removed by being
				// under of a removed directory
				if rmdir != "" && strings.HasPrefix(f1.path, rmdir) ***REMOVED***
					f1 = nil
					continue
				***REMOVED*** else if rmdir == "" && f1.f.IsDir() ***REMOVED***
					rmdir = f1.path + string(os.PathSeparator)
				***REMOVED*** else if rmdir != "" ***REMOVED***
					rmdir = ""
				***REMOVED***
				f1 = nil
			case ChangeKindModify:
				same, err := sameFile(f1, f2)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				if f1.f.IsDir() && !f2.f.IsDir() ***REMOVED***
					rmdir = f1.path + string(os.PathSeparator)
				***REMOVED*** else if rmdir != "" ***REMOVED***
					rmdir = ""
				***REMOVED***
				f = f2.f
				f1 = nil
				f2 = nil
				if same ***REMOVED***
					continue loop0
				***REMOVED***
			***REMOVED***
			if err := changeFn(k, p, f, nil); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***)

	return g.Wait()
***REMOVED***

func pathChange(lower, upper *currentPath) (ChangeKind, string) ***REMOVED***
	if lower == nil ***REMOVED***
		if upper == nil ***REMOVED***
			panic("cannot compare nil paths")
		***REMOVED***
		return ChangeKindAdd, upper.path
	***REMOVED***
	if upper == nil ***REMOVED***
		return ChangeKindDelete, lower.path
	***REMOVED***

	switch i := ComparePath(lower.path, upper.path); ***REMOVED***
	case i < 0:
		// File in lower that is not in upper
		return ChangeKindDelete, lower.path
	case i > 0:
		// File in upper that is not in lower
		return ChangeKindAdd, upper.path
	default:
		return ChangeKindModify, upper.path
	***REMOVED***
***REMOVED***

func sameFile(f1, f2 *currentPath) (same bool, retErr error) ***REMOVED***
	// If not a directory also check size, modtime, and content
	if !f1.f.IsDir() ***REMOVED***
		if f1.f.Size() != f2.f.Size() ***REMOVED***
			return false, nil
		***REMOVED***

		t1 := f1.f.ModTime()
		t2 := f2.f.ModTime()
		if t1.UnixNano() != t2.UnixNano() ***REMOVED***
			return false, nil
		***REMOVED***
	***REMOVED***

	ls1, ok := f1.f.Sys().(*Stat)
	if !ok ***REMOVED***
		return false, nil
	***REMOVED***
	ls2, ok := f1.f.Sys().(*Stat)
	if !ok ***REMOVED***
		return false, nil
	***REMOVED***

	return compareStat(ls1, ls2)
***REMOVED***

// compareStat returns whether the stats are equivalent,
// whether the files are considered the same file, and
// an error
func compareStat(ls1, ls2 *Stat) (bool, error) ***REMOVED***
	return ls1.Mode == ls2.Mode && ls1.Uid == ls2.Uid && ls1.Gid == ls2.Gid && ls1.Devmajor == ls2.Devmajor && ls1.Devminor == ls2.Devminor && ls1.Linkname == ls2.Linkname, nil
***REMOVED***

func nextPath(ctx context.Context, pathC <-chan *currentPath) (*currentPath, error) ***REMOVED***
	select ***REMOVED***
	case <-ctx.Done():
		return nil, ctx.Err()
	case p := <-pathC:
		return p, nil
	***REMOVED***
***REMOVED***
