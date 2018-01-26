// +build linux

package devmapper

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

// FIXME: this is copy-pasted from the aufs driver.
// It should be moved into the core.

// Mounted returns true if a mount point exists.
func Mounted(mountpoint string) (bool, error) ***REMOVED***
	var mntpointSt unix.Stat_t
	if err := unix.Stat(mountpoint, &mntpointSt); err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			return false, nil
		***REMOVED***
		return false, err
	***REMOVED***
	var parentSt unix.Stat_t
	if err := unix.Stat(filepath.Join(mountpoint, ".."), &parentSt); err != nil ***REMOVED***
		return false, err
	***REMOVED***
	return mntpointSt.Dev != parentSt.Dev, nil
***REMOVED***

type probeData struct ***REMOVED***
	fsName string
	magic  string
	offset uint64
***REMOVED***

// ProbeFsType returns the filesystem name for the given device id.
func ProbeFsType(device string) (string, error) ***REMOVED***
	probes := []probeData***REMOVED***
		***REMOVED***"btrfs", "_BHRfS_M", 0x10040***REMOVED***,
		***REMOVED***"ext4", "\123\357", 0x438***REMOVED***,
		***REMOVED***"xfs", "XFSB", 0***REMOVED***,
	***REMOVED***

	maxLen := uint64(0)
	for _, p := range probes ***REMOVED***
		l := p.offset + uint64(len(p.magic))
		if l > maxLen ***REMOVED***
			maxLen = l
		***REMOVED***
	***REMOVED***

	file, err := os.Open(device)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	defer file.Close()

	buffer := make([]byte, maxLen)
	l, err := file.Read(buffer)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	if uint64(l) != maxLen ***REMOVED***
		return "", fmt.Errorf("devmapper: unable to detect filesystem type of %s, short read", device)
	***REMOVED***

	for _, p := range probes ***REMOVED***
		if bytes.Equal([]byte(p.magic), buffer[p.offset:p.offset+uint64(len(p.magic))]) ***REMOVED***
			return p.fsName, nil
		***REMOVED***
	***REMOVED***

	return "", fmt.Errorf("devmapper: Unknown filesystem type on %s", device)
***REMOVED***

func joinMountOptions(a, b string) string ***REMOVED***
	if a == "" ***REMOVED***
		return b
	***REMOVED***
	if b == "" ***REMOVED***
		return a
	***REMOVED***
	return a + "," + b
***REMOVED***
