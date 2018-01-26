package fsutil

import (
	"os"

	"github.com/pkg/errors"
)

// Hardlinks validates that all targets for links were part of the changes

type Hardlinks struct ***REMOVED***
	seenFiles map[string]struct***REMOVED******REMOVED***
***REMOVED***

func (v *Hardlinks) HandleChange(kind ChangeKind, p string, fi os.FileInfo, err error) error ***REMOVED***
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if v.seenFiles == nil ***REMOVED***
		v.seenFiles = make(map[string]struct***REMOVED******REMOVED***)
	***REMOVED***

	if kind == ChangeKindDelete ***REMOVED***
		return nil
	***REMOVED***

	stat, ok := fi.Sys().(*Stat)
	if !ok ***REMOVED***
		return errors.Errorf("invalid change without stat info: %s", p)
	***REMOVED***

	if fi.IsDir() || fi.Mode()&os.ModeSymlink != 0 ***REMOVED***
		return nil
	***REMOVED***

	if len(stat.Linkname) > 0 ***REMOVED***
		if _, ok := v.seenFiles[stat.Linkname]; !ok ***REMOVED***
			return errors.Errorf("invalid link %s to unknown path: %q", p, stat.Linkname)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		v.seenFiles[p] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***

	return nil
***REMOVED***
