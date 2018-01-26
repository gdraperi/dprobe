// +build !windows

package fs

import (
	"bytes"
	"os"
	"syscall"

	"github.com/containerd/continuity/sysx"
	"github.com/pkg/errors"
)

// detectDirDiff returns diff dir options if a directory could
// be found in the mount info for upper which is the direct
// diff with the provided lower directory
func detectDirDiff(upper, lower string) *diffDirOptions ***REMOVED***
	// TODO: get mount options for upper
	// TODO: detect AUFS
	// TODO: detect overlay
	return nil
***REMOVED***

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
	if err != nil && err != sysx.ENODATA ***REMOVED***
		return false, errors.Wrapf(err, "failed to get xattr for %s", p1)
	***REMOVED***
	c2, err := sysx.LGetxattr(p2, "security.capability")
	if err != nil && err != sysx.ENODATA ***REMOVED***
		return false, errors.Wrapf(err, "failed to get xattr for %s", p2)
	***REMOVED***
	return bytes.Equal(c1, c2), nil
***REMOVED***

func isLinked(f os.FileInfo) bool ***REMOVED***
	s, ok := f.Sys().(*syscall.Stat_t)
	if !ok ***REMOVED***
		return false
	***REMOVED***
	return !f.IsDir() && s.Nlink > 1
***REMOVED***
