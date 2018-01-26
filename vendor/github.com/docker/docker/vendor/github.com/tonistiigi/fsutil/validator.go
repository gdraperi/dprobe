package fsutil

import (
	"os"
	"path"
	"runtime"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

type parent struct ***REMOVED***
	dir  string
	last string
***REMOVED***

type Validator struct ***REMOVED***
	parentDirs []parent
***REMOVED***

func (v *Validator) HandleChange(kind ChangeKind, p string, fi os.FileInfo, err error) (retErr error) ***REMOVED***
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	// test that all paths are in order and all parent dirs were present
	if v.parentDirs == nil ***REMOVED***
		v.parentDirs = make([]parent, 1, 10)
	***REMOVED***
	if runtime.GOOS == "windows" ***REMOVED***
		p = strings.Replace(p, "\\", "", -1)
	***REMOVED***
	if p != path.Clean(p) ***REMOVED***
		return errors.Errorf("invalid unclean path %s", p)
	***REMOVED***
	if path.IsAbs(p) ***REMOVED***
		return errors.Errorf("abolute path %s not allowed", p)
	***REMOVED***
	dir := path.Dir(p)
	base := path.Base(p)
	if dir == "." ***REMOVED***
		dir = ""
	***REMOVED***
	if dir == ".." || strings.HasPrefix(p, "../") ***REMOVED***
		return errors.Errorf("invalid path: %s", p)
	***REMOVED***

	// find a parent dir from saved records
	i := sort.Search(len(v.parentDirs), func(i int) bool ***REMOVED***
		return ComparePath(v.parentDirs[len(v.parentDirs)-1-i].dir, dir) <= 0
	***REMOVED***)
	i = len(v.parentDirs) - 1 - i
	if i != len(v.parentDirs)-1 ***REMOVED*** // skipping back to grandparent
		v.parentDirs = v.parentDirs[:i+1]
	***REMOVED***

	if dir != v.parentDirs[len(v.parentDirs)-1].dir || v.parentDirs[i].last >= base ***REMOVED***
		return errors.Errorf("changes out of order: %q %q", p, path.Join(v.parentDirs[i].dir, v.parentDirs[i].last))
	***REMOVED***
	v.parentDirs[i].last = base
	if kind != ChangeKindDelete && fi.IsDir() ***REMOVED***
		v.parentDirs = append(v.parentDirs, parent***REMOVED***
			dir:  path.Join(dir, base),
			last: "",
		***REMOVED***)
	***REMOVED***
	// todo: validate invalid mode combinations
	return err
***REMOVED***

func ComparePath(p1, p2 string) int ***REMOVED***
	// byte-by-byte comparison to be compatible with str<>str
	min := min(len(p1), len(p2))
	for i := 0; i < min; i++ ***REMOVED***
		switch ***REMOVED***
		case p1[i] == p2[i]:
			continue
		case p2[i] != '/' && p1[i] < p2[i] || p1[i] == '/':
			return -1
		default:
			return 1
		***REMOVED***
	***REMOVED***
	return len(p1) - len(p2)
***REMOVED***

func min(x, y int) int ***REMOVED***
	if x < y ***REMOVED***
		return x
	***REMOVED***
	return y
***REMOVED***
