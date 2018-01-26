package mount

import (
	"sort"
	"strings"

	"syscall"

	"github.com/sirupsen/logrus"
)

// GetMounts retrieves a list of mounts for the current running process.
func GetMounts() ([]*Info, error) ***REMOVED***
	return parseMountTable()
***REMOVED***

// Mounted determines if a specified mountpoint has been mounted.
// On Linux it looks at /proc/self/mountinfo.
func Mounted(mountpoint string) (bool, error) ***REMOVED***
	entries, err := parseMountTable()
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	// Search the table for the mountpoint
	for _, e := range entries ***REMOVED***
		if e.Mountpoint == mountpoint ***REMOVED***
			return true, nil
		***REMOVED***
	***REMOVED***
	return false, nil
***REMOVED***

// Mount will mount filesystem according to the specified configuration, on the
// condition that the target path is *not* already mounted. Options must be
// specified like the mount or fstab unix commands: "opt1=val1,opt2=val2". See
// flags.go for supported option flags.
func Mount(device, target, mType, options string) error ***REMOVED***
	flag, _ := parseOptions(options)
	if flag&REMOUNT != REMOUNT ***REMOVED***
		if mounted, err := Mounted(target); err != nil || mounted ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return ForceMount(device, target, mType, options)
***REMOVED***

// ForceMount will mount a filesystem according to the specified configuration,
// *regardless* if the target path is not already mounted. Options must be
// specified like the mount or fstab unix commands: "opt1=val1,opt2=val2". See
// flags.go for supported option flags.
func ForceMount(device, target, mType, options string) error ***REMOVED***
	flag, data := parseOptions(options)
	return mount(device, target, mType, uintptr(flag), data)
***REMOVED***

// Unmount lazily unmounts a filesystem on supported platforms, otherwise
// does a normal unmount.
func Unmount(target string) error ***REMOVED***
	if mounted, err := Mounted(target); err != nil || !mounted ***REMOVED***
		return err
	***REMOVED***
	return unmount(target, mntDetach)
***REMOVED***

// RecursiveUnmount unmounts the target and all mounts underneath, starting with
// the deepsest mount first.
func RecursiveUnmount(target string) error ***REMOVED***
	mounts, err := GetMounts()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Make the deepest mount be first
	sort.Sort(sort.Reverse(byMountpoint(mounts)))

	for i, m := range mounts ***REMOVED***
		if !strings.HasPrefix(m.Mountpoint, target) ***REMOVED***
			continue
		***REMOVED***
		logrus.Debugf("Trying to unmount %s", m.Mountpoint)
		err = unmount(m.Mountpoint, mntDetach)
		if err != nil ***REMOVED***
			// If the error is EINVAL either this whole package is wrong (invalid flags passed to unmount(2)) or this is
			// not a mountpoint (which is ok in this case).
			// Meanwhile calling `Mounted()` is very expensive.
			//
			// We've purposefully used `syscall.EINVAL` here instead of `unix.EINVAL` to avoid platform branching
			// Since `EINVAL` is defined for both Windows and Linux in the `syscall` package (and other platforms),
			//   this is nicer than defining a custom value that we can refer to in each platform file.
			if err == syscall.EINVAL ***REMOVED***
				continue
			***REMOVED***
			if i == len(mounts)-1 ***REMOVED***
				if mounted, e := Mounted(m.Mountpoint); e != nil || mounted ***REMOVED***
					return err
				***REMOVED***
				continue
			***REMOVED***
			// This is some submount, we can ignore this error for now, the final unmount will fail if this is a real problem
			logrus.WithError(err).Warnf("Failed to unmount submount %s", m.Mountpoint)
			continue
		***REMOVED***

		logrus.Debugf("Unmounted %s", m.Mountpoint)
	***REMOVED***
	return nil
***REMOVED***
