package fs

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

var (
	errTooManyLinks = errors.New("too many links")
)

type currentPath struct ***REMOVED***
	path     string
	f        os.FileInfo
	fullPath string
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
	// TODO: compare by directory

	switch i := strings.Compare(lower.path, upper.path); ***REMOVED***
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

func sameFile(f1, f2 *currentPath) (bool, error) ***REMOVED***
	if os.SameFile(f1.f, f2.f) ***REMOVED***
		return true, nil
	***REMOVED***

	equalStat, err := compareSysStat(f1.f.Sys(), f2.f.Sys())
	if err != nil || !equalStat ***REMOVED***
		return equalStat, err
	***REMOVED***

	if eq, err := compareCapabilities(f1.fullPath, f2.fullPath); err != nil || !eq ***REMOVED***
		return eq, err
	***REMOVED***

	// If not a directory also check size, modtime, and content
	if !f1.f.IsDir() ***REMOVED***
		if f1.f.Size() != f2.f.Size() ***REMOVED***
			return false, nil
		***REMOVED***
		t1 := f1.f.ModTime()
		t2 := f2.f.ModTime()

		if t1.Unix() != t2.Unix() ***REMOVED***
			return false, nil
		***REMOVED***

		// If the timestamp may have been truncated in one of the
		// files, check content of file to determine difference
		if t1.Nanosecond() == 0 || t2.Nanosecond() == 0 ***REMOVED***
			var eq bool
			if (f1.f.Mode() & os.ModeSymlink) == os.ModeSymlink ***REMOVED***
				eq, err = compareSymlinkTarget(f1.fullPath, f2.fullPath)
			***REMOVED*** else if f1.f.Size() > 0 ***REMOVED***
				eq, err = compareFileContent(f1.fullPath, f2.fullPath)
			***REMOVED***
			if err != nil || !eq ***REMOVED***
				return eq, err
			***REMOVED***
		***REMOVED*** else if t1.Nanosecond() != t2.Nanosecond() ***REMOVED***
			return false, nil
		***REMOVED***
	***REMOVED***

	return true, nil
***REMOVED***

func compareSymlinkTarget(p1, p2 string) (bool, error) ***REMOVED***
	t1, err := os.Readlink(p1)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	t2, err := os.Readlink(p2)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	return t1 == t2, nil
***REMOVED***

const compareChuckSize = 32 * 1024

// compareFileContent compares the content of 2 same sized files
// by comparing each byte.
func compareFileContent(p1, p2 string) (bool, error) ***REMOVED***
	f1, err := os.Open(p1)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	defer f1.Close()
	f2, err := os.Open(p2)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	defer f2.Close()

	b1 := make([]byte, compareChuckSize)
	b2 := make([]byte, compareChuckSize)
	for ***REMOVED***
		n1, err1 := f1.Read(b1)
		if err1 != nil && err1 != io.EOF ***REMOVED***
			return false, err1
		***REMOVED***
		n2, err2 := f2.Read(b2)
		if err2 != nil && err2 != io.EOF ***REMOVED***
			return false, err2
		***REMOVED***
		if n1 != n2 || !bytes.Equal(b1[:n1], b2[:n2]) ***REMOVED***
			return false, nil
		***REMOVED***
		if err1 == io.EOF && err2 == io.EOF ***REMOVED***
			return true, nil
		***REMOVED***
	***REMOVED***
***REMOVED***

func pathWalk(ctx context.Context, root string, pathC chan<- *currentPath) error ***REMOVED***
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

		p := &currentPath***REMOVED***
			path:     path,
			f:        f,
			fullPath: filepath.Join(root, path),
		***REMOVED***

		select ***REMOVED***
		case <-ctx.Done():
			return ctx.Err()
		case pathC <- p:
			return nil
		***REMOVED***
	***REMOVED***)
***REMOVED***

func nextPath(ctx context.Context, pathC <-chan *currentPath) (*currentPath, error) ***REMOVED***
	select ***REMOVED***
	case <-ctx.Done():
		return nil, ctx.Err()
	case p := <-pathC:
		return p, nil
	***REMOVED***
***REMOVED***

// RootPath joins a path with a root, evaluating and bounding any
// symlink to the root directory.
func RootPath(root, path string) (string, error) ***REMOVED***
	if path == "" ***REMOVED***
		return root, nil
	***REMOVED***
	var linksWalked int // to protect against cycles
	for ***REMOVED***
		i := linksWalked
		newpath, err := walkLinks(root, path, &linksWalked)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		path = newpath
		if i == linksWalked ***REMOVED***
			newpath = filepath.Join("/", newpath)
			if path == newpath ***REMOVED***
				return filepath.Join(root, newpath), nil
			***REMOVED***
			path = newpath
		***REMOVED***
	***REMOVED***
***REMOVED***

func walkLink(root, path string, linksWalked *int) (newpath string, islink bool, err error) ***REMOVED***
	if *linksWalked > 255 ***REMOVED***
		return "", false, errTooManyLinks
	***REMOVED***

	path = filepath.Join("/", path)
	if path == "/" ***REMOVED***
		return path, false, nil
	***REMOVED***
	realPath := filepath.Join(root, path)

	fi, err := os.Lstat(realPath)
	if err != nil ***REMOVED***
		// If path does not yet exist, treat as non-symlink
		if os.IsNotExist(err) ***REMOVED***
			return path, false, nil
		***REMOVED***
		return "", false, err
	***REMOVED***
	if fi.Mode()&os.ModeSymlink == 0 ***REMOVED***
		return path, false, nil
	***REMOVED***
	newpath, err = os.Readlink(realPath)
	if err != nil ***REMOVED***
		return "", false, err
	***REMOVED***
	if filepath.IsAbs(newpath) && strings.HasPrefix(newpath, root) ***REMOVED***
		newpath = newpath[:len(root)]
		if !strings.HasPrefix(newpath, "/") ***REMOVED***
			newpath = "/" + newpath
		***REMOVED***
	***REMOVED***
	*linksWalked++
	return newpath, true, nil
***REMOVED***

func walkLinks(root, path string, linksWalked *int) (string, error) ***REMOVED***
	switch dir, file := filepath.Split(path); ***REMOVED***
	case dir == "":
		newpath, _, err := walkLink(root, file, linksWalked)
		return newpath, err
	case file == "":
		if os.IsPathSeparator(dir[len(dir)-1]) ***REMOVED***
			if dir == "/" ***REMOVED***
				return dir, nil
			***REMOVED***
			return walkLinks(root, dir[:len(dir)-1], linksWalked)
		***REMOVED***
		newpath, _, err := walkLink(root, dir, linksWalked)
		return newpath, err
	default:
		newdir, err := walkLinks(root, dir, linksWalked)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		newpath, islink, err := walkLink(root, filepath.Join(newdir, file), linksWalked)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		if !islink ***REMOVED***
			return newpath, nil
		***REMOVED***
		if filepath.IsAbs(newpath) ***REMOVED***
			return newpath, nil
		***REMOVED***
		return filepath.Join(newdir, newpath), nil
	***REMOVED***
***REMOVED***
