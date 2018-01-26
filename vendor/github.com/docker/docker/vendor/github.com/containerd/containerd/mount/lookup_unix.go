// +build !windows

package mount

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	"github.com/pkg/errors"
)

// Lookup returns the mount info corresponds to the path.
func Lookup(dir string) (Info, error) ***REMOVED***
	var dirStat syscall.Stat_t
	dir = filepath.Clean(dir)
	if err := syscall.Stat(dir, &dirStat); err != nil ***REMOVED***
		return Info***REMOVED******REMOVED***, errors.Wrapf(err, "failed to access %q", dir)
	***REMOVED***

	mounts, err := Self()
	if err != nil ***REMOVED***
		return Info***REMOVED******REMOVED***, err
	***REMOVED***

	// Sort descending order by Info.Mountpoint
	sort.Slice(mounts, func(i, j int) bool ***REMOVED***
		return mounts[j].Mountpoint < mounts[i].Mountpoint
	***REMOVED***)
	for _, m := range mounts ***REMOVED***
		// Note that m.***REMOVED***Major, Minor***REMOVED*** are generally unreliable for our purpose here
		// https://www.spinics.net/lists/linux-btrfs/msg58908.html
		var st syscall.Stat_t
		if err := syscall.Stat(m.Mountpoint, &st); err != nil ***REMOVED***
			// may fail; ignore err
			continue
		***REMOVED***
		if st.Dev == dirStat.Dev && strings.HasPrefix(dir, m.Mountpoint) ***REMOVED***
			return m, nil
		***REMOVED***
	***REMOVED***

	return Info***REMOVED******REMOVED***, fmt.Errorf("failed to find the mount info for %q", dir)
***REMOVED***
