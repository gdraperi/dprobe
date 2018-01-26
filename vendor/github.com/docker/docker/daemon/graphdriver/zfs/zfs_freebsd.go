package zfs

import (
	"fmt"
	"strings"

	"github.com/docker/docker/daemon/graphdriver"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

func checkRootdirFs(rootdir string) error ***REMOVED***
	var buf unix.Statfs_t
	if err := unix.Statfs(rootdir, &buf); err != nil ***REMOVED***
		return fmt.Errorf("Failed to access '%s': %s", rootdir, err)
	***REMOVED***

	// on FreeBSD buf.Fstypename contains ['z', 'f', 's', 0 ... ]
	if (buf.Fstypename[0] != 122) || (buf.Fstypename[1] != 102) || (buf.Fstypename[2] != 115) || (buf.Fstypename[3] != 0) ***REMOVED***
		logrus.Debugf("[zfs] no zfs dataset found for rootdir '%s'", rootdir)
		return graphdriver.ErrPrerequisites
	***REMOVED***

	return nil
***REMOVED***

func getMountpoint(id string) string ***REMOVED***
	maxlen := 12

	// we need to preserve filesystem suffix
	suffix := strings.SplitN(id, "-", 2)

	if len(suffix) > 1 ***REMOVED***
		return id[:maxlen] + "-" + suffix[1]
	***REMOVED***

	return id[:maxlen]
***REMOVED***
