package symlink

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/pkg/longpath"
	"golang.org/x/sys/windows"
)

func toShort(path string) (string, error) ***REMOVED***
	p, err := windows.UTF16FromString(path)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	b := p // GetShortPathName says we can reuse buffer
	n, err := windows.GetShortPathName(&p[0], &b[0], uint32(len(b)))
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if n > uint32(len(b)) ***REMOVED***
		b = make([]uint16, n)
		if _, err = windows.GetShortPathName(&p[0], &b[0], uint32(len(b))); err != nil ***REMOVED***
			return "", err
		***REMOVED***
	***REMOVED***
	return windows.UTF16ToString(b), nil
***REMOVED***

func toLong(path string) (string, error) ***REMOVED***
	p, err := windows.UTF16FromString(path)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	b := p // GetLongPathName says we can reuse buffer
	n, err := windows.GetLongPathName(&p[0], &b[0], uint32(len(b)))
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if n > uint32(len(b)) ***REMOVED***
		b = make([]uint16, n)
		n, err = windows.GetLongPathName(&p[0], &b[0], uint32(len(b)))
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
	***REMOVED***
	b = b[:n]
	return windows.UTF16ToString(b), nil
***REMOVED***

func evalSymlinks(path string) (string, error) ***REMOVED***
	path, err := walkSymlinks(path)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	p, err := toShort(path)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	p, err = toLong(p)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	// windows.GetLongPathName does not change the case of the drive letter,
	// but the result of EvalSymlinks must be unique, so we have
	// EvalSymlinks(`c:\a`) == EvalSymlinks(`C:\a`).
	// Make drive letter upper case.
	if len(p) >= 2 && p[1] == ':' && 'a' <= p[0] && p[0] <= 'z' ***REMOVED***
		p = string(p[0]+'A'-'a') + p[1:]
	***REMOVED*** else if len(p) >= 6 && p[5] == ':' && 'a' <= p[4] && p[4] <= 'z' ***REMOVED***
		p = p[:3] + string(p[4]+'A'-'a') + p[5:]
	***REMOVED***
	return filepath.Clean(p), nil
***REMOVED***

const utf8RuneSelf = 0x80

func walkSymlinks(path string) (string, error) ***REMOVED***
	const maxIter = 255
	originalPath := path
	// consume path by taking each frontmost path element,
	// expanding it if it's a symlink, and appending it to b
	var b bytes.Buffer
	for n := 0; path != ""; n++ ***REMOVED***
		if n > maxIter ***REMOVED***
			return "", errors.New("EvalSymlinks: too many links in " + originalPath)
		***REMOVED***

		// A path beginning with `\\?\` represents the root, so automatically
		// skip that part and begin processing the next segment.
		if strings.HasPrefix(path, longpath.Prefix) ***REMOVED***
			b.WriteString(longpath.Prefix)
			path = path[4:]
			continue
		***REMOVED***

		// find next path component, p
		var i = -1
		for j, c := range path ***REMOVED***
			if c < utf8RuneSelf && os.IsPathSeparator(uint8(c)) ***REMOVED***
				i = j
				break
			***REMOVED***
		***REMOVED***
		var p string
		if i == -1 ***REMOVED***
			p, path = path, ""
		***REMOVED*** else ***REMOVED***
			p, path = path[:i], path[i+1:]
		***REMOVED***

		if p == "" ***REMOVED***
			if b.Len() == 0 ***REMOVED***
				// must be absolute path
				b.WriteRune(filepath.Separator)
			***REMOVED***
			continue
		***REMOVED***

		// If this is the first segment after the long path prefix, accept the
		// current segment as a volume root or UNC share and move on to the next.
		if b.String() == longpath.Prefix ***REMOVED***
			b.WriteString(p)
			b.WriteRune(filepath.Separator)
			continue
		***REMOVED***

		fi, err := os.Lstat(b.String() + p)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		if fi.Mode()&os.ModeSymlink == 0 ***REMOVED***
			b.WriteString(p)
			if path != "" || (b.Len() == 2 && len(p) == 2 && p[1] == ':') ***REMOVED***
				b.WriteRune(filepath.Separator)
			***REMOVED***
			continue
		***REMOVED***

		// it's a symlink, put it at the front of path
		dest, err := os.Readlink(b.String() + p)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		if filepath.IsAbs(dest) || os.IsPathSeparator(dest[0]) ***REMOVED***
			b.Reset()
		***REMOVED***
		path = dest + string(filepath.Separator) + path
	***REMOVED***
	return filepath.Clean(b.String()), nil
***REMOVED***

func isDriveOrRoot(p string) bool ***REMOVED***
	if p == string(filepath.Separator) ***REMOVED***
		return true
	***REMOVED***

	length := len(p)
	if length >= 2 ***REMOVED***
		if p[length-1] == ':' && (('a' <= p[length-2] && p[length-2] <= 'z') || ('A' <= p[length-2] && p[length-2] <= 'Z')) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
