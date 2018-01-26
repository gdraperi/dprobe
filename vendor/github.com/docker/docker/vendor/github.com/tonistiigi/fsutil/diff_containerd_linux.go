package fsutil

import (
	"bytes"
	"syscall"

	"github.com/containerd/continuity/sysx"
	"github.com/pkg/errors"
)

// compareSysStat returns whether the stats are equivalent,
// whether the files are considered the same file, and
// an error
func compareSysStat(s1, s2 interface***REMOVED******REMOVED***) (bool, error) ***REMOVED***
	ls1, ok := s1.(*syscall.Stat_t)
	if !ok ***REMOVED***
		return false, nil
	***REMOVED***
	ls2, ok := s2.(*syscall.Stat_t)
	if !ok ***REMOVED***
		return false, nil
	***REMOVED***

	return ls1.Mode == ls2.Mode && ls1.Uid == ls2.Uid && ls1.Gid == ls2.Gid && ls1.Rdev == ls2.Rdev, nil
***REMOVED***

func compareCapabilities(p1, p2 string) (bool, error) ***REMOVED***
	c1, err := sysx.LGetxattr(p1, "security.capability")
	if err != nil && err != syscall.ENODATA ***REMOVED***
		return false, errors.Wrapf(err, "failed to get xattr for %s", p1)
	***REMOVED***
	c2, err := sysx.LGetxattr(p2, "security.capability")
	if err != nil && err != syscall.ENODATA ***REMOVED***
		return false, errors.Wrapf(err, "failed to get xattr for %s", p2)
	***REMOVED***
	return bytes.Equal(c1, c2), nil
***REMOVED***
