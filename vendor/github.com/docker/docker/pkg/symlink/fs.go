// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.BSD file.

// This code is a modified version of path/filepath/symlink.go from the Go standard library.

package symlink

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/pkg/system"
)

// FollowSymlinkInScope is a wrapper around evalSymlinksInScope that returns an
// absolute path. This function handles paths in a platform-agnostic manner.
func FollowSymlinkInScope(path, root string) (string, error) ***REMOVED***
	path, err := filepath.Abs(filepath.FromSlash(path))
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	root, err = filepath.Abs(filepath.FromSlash(root))
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return evalSymlinksInScope(path, root)
***REMOVED***

// evalSymlinksInScope will evaluate symlinks in `path` within a scope `root` and return
// a result guaranteed to be contained within the scope `root`, at the time of the call.
// Symlinks in `root` are not evaluated and left as-is.
// Errors encountered while attempting to evaluate symlinks in path will be returned.
// Non-existing paths are valid and do not constitute an error.
// `path` has to contain `root` as a prefix, or else an error will be returned.
// Trying to break out from `root` does not constitute an error.
//
// Example:
//   If /foo/bar -> /outside,
//   FollowSymlinkInScope("/foo/bar", "/foo") == "/foo/outside" instead of "/outside"
//
// IMPORTANT: it is the caller's responsibility to call evalSymlinksInScope *after* relevant symlinks
// are created and not to create subsequently, additional symlinks that could potentially make a
// previously-safe path, unsafe. Example: if /foo/bar does not exist, evalSymlinksInScope("/foo/bar", "/foo")
// would return "/foo/bar". If one makes /foo/bar a symlink to /baz subsequently, then "/foo/bar" should
// no longer be considered safely contained in "/foo".
func evalSymlinksInScope(path, root string) (string, error) ***REMOVED***
	root = filepath.Clean(root)
	if path == root ***REMOVED***
		return path, nil
	***REMOVED***
	if !strings.HasPrefix(path, root) ***REMOVED***
		return "", errors.New("evalSymlinksInScope: " + path + " is not in " + root)
	***REMOVED***
	const maxIter = 255
	originalPath := path
	// given root of "/a" and path of "/a/b/../../c" we want path to be "/b/../../c"
	path = path[len(root):]
	if root == string(filepath.Separator) ***REMOVED***
		path = string(filepath.Separator) + path
	***REMOVED***
	if !strings.HasPrefix(path, string(filepath.Separator)) ***REMOVED***
		return "", errors.New("evalSymlinksInScope: " + path + " is not in " + root)
	***REMOVED***
	path = filepath.Clean(path)
	// consume path by taking each frontmost path element,
	// expanding it if it's a symlink, and appending it to b
	var b bytes.Buffer
	// b here will always be considered to be the "current absolute path inside
	// root" when we append paths to it, we also append a slash and use
	// filepath.Clean after the loop to trim the trailing slash
	for n := 0; path != ""; n++ ***REMOVED***
		if n > maxIter ***REMOVED***
			return "", errors.New("evalSymlinksInScope: too many links in " + originalPath)
		***REMOVED***

		// find next path component, p
		i := strings.IndexRune(path, filepath.Separator)
		var p string
		if i == -1 ***REMOVED***
			p, path = path, ""
		***REMOVED*** else ***REMOVED***
			p, path = path[:i], path[i+1:]
		***REMOVED***

		if p == "" ***REMOVED***
			continue
		***REMOVED***

		// this takes a b.String() like "b/../" and a p like "c" and turns it
		// into "/b/../c" which then gets filepath.Cleaned into "/c" and then
		// root gets prepended and we Clean again (to remove any trailing slash
		// if the first Clean gave us just "/")
		cleanP := filepath.Clean(string(filepath.Separator) + b.String() + p)
		if isDriveOrRoot(cleanP) ***REMOVED***
			// never Lstat "/" itself, or drive letters on Windows
			b.Reset()
			continue
		***REMOVED***
		fullP := filepath.Clean(root + cleanP)

		fi, err := os.Lstat(fullP)
		if os.IsNotExist(err) ***REMOVED***
			// if p does not exist, accept it
			b.WriteString(p)
			b.WriteRune(filepath.Separator)
			continue
		***REMOVED***
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		if fi.Mode()&os.ModeSymlink == 0 ***REMOVED***
			b.WriteString(p)
			b.WriteRune(filepath.Separator)
			continue
		***REMOVED***

		// it's a symlink, put it at the front of path
		dest, err := os.Readlink(fullP)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		if system.IsAbs(dest) ***REMOVED***
			b.Reset()
		***REMOVED***
		path = dest + string(filepath.Separator) + path
	***REMOVED***

	// see note above on "fullP := ..." for why this is double-cleaned and
	// what's happening here
	return filepath.Clean(root + filepath.Clean(string(filepath.Separator)+b.String())), nil
***REMOVED***

// EvalSymlinks returns the path name after the evaluation of any symbolic
// links.
// If path is relative the result will be relative to the current directory,
// unless one of the components is an absolute symbolic link.
// This version has been updated to support long paths prepended with `\\?\`.
func EvalSymlinks(path string) (string, error) ***REMOVED***
	return evalSymlinks(path)
***REMOVED***
