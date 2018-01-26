// +build windows

package lcow

import (
	"errors"
	"os"
	pathpkg "path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/containerd/continuity/pathdriver"
)

var _ pathdriver.PathDriver = &lcowfs***REMOVED******REMOVED***

// Continuity Path functions can be done locally
func (l *lcowfs) Join(path ...string) string ***REMOVED***
	return pathpkg.Join(path...)
***REMOVED***

func (l *lcowfs) IsAbs(path string) bool ***REMOVED***
	return pathpkg.IsAbs(path)
***REMOVED***

func sameWord(a, b string) bool ***REMOVED***
	return a == b
***REMOVED***

// Implementation taken from the Go standard library
func (l *lcowfs) Rel(basepath, targpath string) (string, error) ***REMOVED***
	baseVol := ""
	targVol := ""
	base := l.Clean(basepath)
	targ := l.Clean(targpath)
	if sameWord(targ, base) ***REMOVED***
		return ".", nil
	***REMOVED***
	base = base[len(baseVol):]
	targ = targ[len(targVol):]
	if base == "." ***REMOVED***
		base = ""
	***REMOVED***
	// Can't use IsAbs - `\a` and `a` are both relative in Windows.
	baseSlashed := len(base) > 0 && base[0] == l.Separator()
	targSlashed := len(targ) > 0 && targ[0] == l.Separator()
	if baseSlashed != targSlashed || !sameWord(baseVol, targVol) ***REMOVED***
		return "", errors.New("Rel: can't make " + targpath + " relative to " + basepath)
	***REMOVED***
	// Position base[b0:bi] and targ[t0:ti] at the first differing elements.
	bl := len(base)
	tl := len(targ)
	var b0, bi, t0, ti int
	for ***REMOVED***
		for bi < bl && base[bi] != l.Separator() ***REMOVED***
			bi++
		***REMOVED***
		for ti < tl && targ[ti] != l.Separator() ***REMOVED***
			ti++
		***REMOVED***
		if !sameWord(targ[t0:ti], base[b0:bi]) ***REMOVED***
			break
		***REMOVED***
		if bi < bl ***REMOVED***
			bi++
		***REMOVED***
		if ti < tl ***REMOVED***
			ti++
		***REMOVED***
		b0 = bi
		t0 = ti
	***REMOVED***
	if base[b0:bi] == ".." ***REMOVED***
		return "", errors.New("Rel: can't make " + targpath + " relative to " + basepath)
	***REMOVED***
	if b0 != bl ***REMOVED***
		// Base elements left. Must go up before going down.
		seps := strings.Count(base[b0:bl], string(l.Separator()))
		size := 2 + seps*3
		if tl != t0 ***REMOVED***
			size += 1 + tl - t0
		***REMOVED***
		buf := make([]byte, size)
		n := copy(buf, "..")
		for i := 0; i < seps; i++ ***REMOVED***
			buf[n] = l.Separator()
			copy(buf[n+1:], "..")
			n += 3
		***REMOVED***
		if t0 != tl ***REMOVED***
			buf[n] = l.Separator()
			copy(buf[n+1:], targ[t0:])
		***REMOVED***
		return string(buf), nil
	***REMOVED***
	return targ[t0:], nil
***REMOVED***

func (l *lcowfs) Base(path string) string ***REMOVED***
	return pathpkg.Base(path)
***REMOVED***

func (l *lcowfs) Dir(path string) string ***REMOVED***
	return pathpkg.Dir(path)
***REMOVED***

func (l *lcowfs) Clean(path string) string ***REMOVED***
	return pathpkg.Clean(path)
***REMOVED***

func (l *lcowfs) Split(path string) (dir, file string) ***REMOVED***
	return pathpkg.Split(path)
***REMOVED***

func (l *lcowfs) Separator() byte ***REMOVED***
	return '/'
***REMOVED***

func (l *lcowfs) Abs(path string) (string, error) ***REMOVED***
	// Abs is supposed to add the current working directory, which is meaningless in lcow.
	// So, return an error.
	return "", ErrNotSupported
***REMOVED***

// Implementation taken from the Go standard library
func (l *lcowfs) Walk(root string, walkFn filepath.WalkFunc) error ***REMOVED***
	info, err := l.Lstat(root)
	if err != nil ***REMOVED***
		err = walkFn(root, nil, err)
	***REMOVED*** else ***REMOVED***
		err = l.walk(root, info, walkFn)
	***REMOVED***
	if err == filepath.SkipDir ***REMOVED***
		return nil
	***REMOVED***
	return err
***REMOVED***

// walk recursively descends path, calling w.
func (l *lcowfs) walk(path string, info os.FileInfo, walkFn filepath.WalkFunc) error ***REMOVED***
	err := walkFn(path, info, nil)
	if err != nil ***REMOVED***
		if info.IsDir() && err == filepath.SkipDir ***REMOVED***
			return nil
		***REMOVED***
		return err
	***REMOVED***

	if !info.IsDir() ***REMOVED***
		return nil
	***REMOVED***

	names, err := l.readDirNames(path)
	if err != nil ***REMOVED***
		return walkFn(path, info, err)
	***REMOVED***

	for _, name := range names ***REMOVED***
		filename := l.Join(path, name)
		fileInfo, err := l.Lstat(filename)
		if err != nil ***REMOVED***
			if err := walkFn(filename, fileInfo, err); err != nil && err != filepath.SkipDir ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			err = l.walk(filename, fileInfo, walkFn)
			if err != nil ***REMOVED***
				if !fileInfo.IsDir() || err != filepath.SkipDir ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// readDirNames reads the directory named by dirname and returns
// a sorted list of directory entries.
func (l *lcowfs) readDirNames(dirname string) ([]string, error) ***REMOVED***
	f, err := l.Open(dirname)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	names := make([]string, len(files), len(files))
	for i := range files ***REMOVED***
		names[i] = files[i].Name()
	***REMOVED***

	sort.Strings(names)
	return names, nil
***REMOVED***

// Note that Go's filepath.FromSlash/ToSlash convert between OS paths and '/'. Since the path separator
// for LCOW (and Unix) is '/', they are no-ops.
func (l *lcowfs) FromSlash(path string) string ***REMOVED***
	return path
***REMOVED***

func (l *lcowfs) ToSlash(path string) string ***REMOVED***
	return path
***REMOVED***

func (l *lcowfs) Match(pattern, name string) (matched bool, err error) ***REMOVED***
	return pathpkg.Match(pattern, name)
***REMOVED***
