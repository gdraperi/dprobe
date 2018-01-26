package fs

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/sirupsen/logrus"
)

// ChangeKind is the type of modification that
// a change is making.
type ChangeKind int

const (
	// ChangeKindUnmodified represents an unmodified
	// file
	ChangeKindUnmodified = iota

	// ChangeKindAdd represents an addition of
	// a file
	ChangeKindAdd

	// ChangeKindModify represents a change to
	// an existing file
	ChangeKindModify

	// ChangeKindDelete represents a delete of
	// a file
	ChangeKindDelete
)

func (k ChangeKind) String() string ***REMOVED***
	switch k ***REMOVED***
	case ChangeKindUnmodified:
		return "unmodified"
	case ChangeKindAdd:
		return "add"
	case ChangeKindModify:
		return "modify"
	case ChangeKindDelete:
		return "delete"
	default:
		return ""
	***REMOVED***
***REMOVED***

// Change represents single change between a diff and its parent.
type Change struct ***REMOVED***
	Kind ChangeKind
	Path string
***REMOVED***

// ChangeFunc is the type of function called for each change
// computed during a directory changes calculation.
type ChangeFunc func(ChangeKind, string, os.FileInfo, error) error

// Changes computes changes between two directories calling the
// given change function for each computed change. The first
// directory is intended to the base directory and second
// directory the changed directory.
//
// The change callback is called by the order of path names and
// should be appliable in that order.
//  Due to this apply ordering, the following is true
//  - Removed directory trees only create a single change for the root
//    directory removed. Remaining changes are implied.
//  - A directory which is modified to become a file will not have
//    delete entries for sub-path items, their removal is implied
//    by the removal of the parent directory.
//
// Opaque directories will not be treated specially and each file
// removed from the base directory will show up as a removal.
//
// File content comparisons will be done on files which have timestamps
// which may have been truncated. If either of the files being compared
// has a zero value nanosecond value, each byte will be compared for
// differences. If 2 files have the same seconds value but different
// nanosecond values where one of those values is zero, the files will
// be considered unchanged if the content is the same. This behavior
// is to account for timestamp truncation during archiving.
func Changes(ctx context.Context, a, b string, changeFn ChangeFunc) error ***REMOVED***
	if a == "" ***REMOVED***
		logrus.Debugf("Using single walk diff for %s", b)
		return addDirChanges(ctx, changeFn, b)
	***REMOVED*** else if diffOptions := detectDirDiff(b, a); diffOptions != nil ***REMOVED***
		logrus.Debugf("Using single walk diff for %s from %s", diffOptions.diffDir, a)
		return diffDirChanges(ctx, changeFn, a, diffOptions)
	***REMOVED***

	logrus.Debugf("Using double walk diff for %s from %s", b, a)
	return doubleWalkDiff(ctx, changeFn, a, b)
***REMOVED***

func addDirChanges(ctx context.Context, changeFn ChangeFunc, root string) error ***REMOVED***
	return filepath.Walk(root, func(path string, f os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Rebase path
		path, err = filepath.Rel(root, path)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		path = filepath.Join(string(os.PathSeparator), path)

		// Skip root
		if path == string(os.PathSeparator) ***REMOVED***
			return nil
		***REMOVED***

		return changeFn(ChangeKindAdd, path, f, nil)
	***REMOVED***)
***REMOVED***

// diffDirOptions is used when the diff can be directly calculated from
// a diff directory to its base, without walking both trees.
type diffDirOptions struct ***REMOVED***
	diffDir      string
	skipChange   func(string) (bool, error)
	deleteChange func(string, string, os.FileInfo) (string, error)
***REMOVED***

// diffDirChanges walks the diff directory and compares changes against the base.
func diffDirChanges(ctx context.Context, changeFn ChangeFunc, base string, o *diffDirOptions) error ***REMOVED***
	changedDirs := make(map[string]struct***REMOVED******REMOVED***)
	return filepath.Walk(o.diffDir, func(path string, f os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Rebase path
		path, err = filepath.Rel(o.diffDir, path)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		path = filepath.Join(string(os.PathSeparator), path)

		// Skip root
		if path == string(os.PathSeparator) ***REMOVED***
			return nil
		***REMOVED***

		// TODO: handle opaqueness, start new double walker at this
		// location to get deletes, and skip tree in single walker

		if o.skipChange != nil ***REMOVED***
			if skip, err := o.skipChange(path); skip ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		var kind ChangeKind

		deletedFile, err := o.deleteChange(o.diffDir, path, f)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Find out what kind of modification happened
		if deletedFile != "" ***REMOVED***
			path = deletedFile
			kind = ChangeKindDelete
			f = nil
		***REMOVED*** else ***REMOVED***
			// Otherwise, the file was added
			kind = ChangeKindAdd

			// ...Unless it already existed in a base, in which case, it's a modification
			stat, err := os.Stat(filepath.Join(base, path))
			if err != nil && !os.IsNotExist(err) ***REMOVED***
				return err
			***REMOVED***
			if err == nil ***REMOVED***
				// The file existed in the base, so that's a modification

				// However, if it's a directory, maybe it wasn't actually modified.
				// If you modify /foo/bar/baz, then /foo will be part of the changed files only because it's the parent of bar
				if stat.IsDir() && f.IsDir() ***REMOVED***
					if f.Size() == stat.Size() && f.Mode() == stat.Mode() && sameFsTime(f.ModTime(), stat.ModTime()) ***REMOVED***
						// Both directories are the same, don't record the change
						return nil
					***REMOVED***
				***REMOVED***
				kind = ChangeKindModify
			***REMOVED***
		***REMOVED***

		// If /foo/bar/file.txt is modified, then /foo/bar must be part of the changed files.
		// This block is here to ensure the change is recorded even if the
		// modify time, mode and size of the parent directory in the rw and ro layers are all equal.
		// Check https://github.com/docker/docker/pull/13590 for details.
		if f.IsDir() ***REMOVED***
			changedDirs[path] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
		if kind == ChangeKindAdd || kind == ChangeKindDelete ***REMOVED***
			parent := filepath.Dir(path)
			if _, ok := changedDirs[parent]; !ok && parent != "/" ***REMOVED***
				pi, err := os.Stat(filepath.Join(o.diffDir, parent))
				if err := changeFn(ChangeKindModify, parent, pi, err); err != nil ***REMOVED***
					return err
				***REMOVED***
				changedDirs[parent] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			***REMOVED***
		***REMOVED***

		return changeFn(kind, path, f, nil)
	***REMOVED***)
***REMOVED***

// doubleWalkDiff walks both directories to create a diff
func doubleWalkDiff(ctx context.Context, changeFn ChangeFunc, a, b string) (err error) ***REMOVED***
	g, ctx := errgroup.WithContext(ctx)

	var (
		c1 = make(chan *currentPath)
		c2 = make(chan *currentPath)

		f1, f2         *currentPath
		rmdir          string
		lastEmittedDir = string(filepath.Separator)
		parents        []os.FileInfo
	)
	g.Go(func() error ***REMOVED***
		defer close(c1)
		return pathWalk(ctx, a, c1)
	***REMOVED***)
	g.Go(func() error ***REMOVED***
		defer close(c2)
		return pathWalk(ctx, b, c2)
	***REMOVED***)
	g.Go(func() error ***REMOVED***
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

			var (
				f    os.FileInfo
				emit = true
			)
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
					if !isLinked(f) ***REMOVED***
						emit = false
					***REMOVED***
					k = ChangeKindUnmodified
				***REMOVED***
			***REMOVED***
			if emit ***REMOVED***
				emittedDir, emitParents := commonParents(lastEmittedDir, p, parents)
				for _, pf := range emitParents ***REMOVED***
					p := filepath.Join(emittedDir, pf.Name())
					if err := changeFn(ChangeKindUnmodified, p, pf, nil); err != nil ***REMOVED***
						return err
					***REMOVED***
					emittedDir = p
				***REMOVED***

				if err := changeFn(k, p, f, nil); err != nil ***REMOVED***
					return err
				***REMOVED***

				if f != nil && f.IsDir() ***REMOVED***
					lastEmittedDir = p
				***REMOVED*** else ***REMOVED***
					lastEmittedDir = emittedDir
				***REMOVED***

				parents = parents[:0]
			***REMOVED*** else if f.IsDir() ***REMOVED***
				lastEmittedDir, parents = commonParents(lastEmittedDir, p, parents)
				parents = append(parents, f)
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***)

	return g.Wait()
***REMOVED***

func commonParents(base, updated string, dirs []os.FileInfo) (string, []os.FileInfo) ***REMOVED***
	if basePrefix := makePrefix(base); strings.HasPrefix(updated, basePrefix) ***REMOVED***
		var (
			parents []os.FileInfo
			last    = base
		)
		for _, d := range dirs ***REMOVED***
			next := filepath.Join(last, d.Name())
			if strings.HasPrefix(updated, makePrefix(last)) ***REMOVED***
				parents = append(parents, d)
				last = next
			***REMOVED*** else ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
		return base, parents
	***REMOVED***

	baseS := strings.Split(base, string(filepath.Separator))
	updatedS := strings.Split(updated, string(filepath.Separator))
	commonS := []string***REMOVED***string(filepath.Separator)***REMOVED***

	min := len(baseS)
	if len(updatedS) < min ***REMOVED***
		min = len(updatedS)
	***REMOVED***
	for i := 0; i < min; i++ ***REMOVED***
		if baseS[i] == updatedS[i] ***REMOVED***
			commonS = append(commonS, baseS[i])
		***REMOVED*** else ***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	return filepath.Join(commonS...), []os.FileInfo***REMOVED******REMOVED***
***REMOVED***

func makePrefix(d string) string ***REMOVED***
	if d == "" || d[len(d)-1] != filepath.Separator ***REMOVED***
		return d + string(filepath.Separator)
	***REMOVED***
	return d
***REMOVED***
