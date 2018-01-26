package mount

import (
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

// Mount to the provided target path
func (m *Mount) Mount(target string) error ***REMOVED***
	flags, data := parseMountOptions(m.Options)

	// propagation types.
	const ptypes = unix.MS_SHARED | unix.MS_PRIVATE | unix.MS_SLAVE | unix.MS_UNBINDABLE

	// Ensure propagation type change flags aren't included in other calls.
	oflags := flags &^ ptypes

	// In the case of remounting with changed data (data != ""), need to call mount (moby/moby#34077).
	if flags&unix.MS_REMOUNT == 0 || data != "" ***REMOVED***
		// Initial call applying all non-propagation flags for mount
		// or remount with changed data
		if err := unix.Mount(m.Source, target, m.Type, uintptr(oflags), data); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if flags&ptypes != 0 ***REMOVED***
		// Change the propagation type.
		const pflags = ptypes | unix.MS_REC | unix.MS_SILENT
		if err := unix.Mount("", target, "", uintptr(flags&pflags), ""); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	const broflags = unix.MS_BIND | unix.MS_RDONLY
	if oflags&broflags == broflags ***REMOVED***
		// Remount the bind to apply read only.
		return unix.Mount("", target, "", uintptr(oflags|unix.MS_REMOUNT), "")
	***REMOVED***
	return nil
***REMOVED***

// Unmount the provided mount path with the flags
func Unmount(target string, flags int) error ***REMOVED***
	if err := unmount(target, flags); err != nil && err != unix.EINVAL ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func unmount(target string, flags int) error ***REMOVED***
	for i := 0; i < 50; i++ ***REMOVED***
		if err := unix.Unmount(target, flags); err != nil ***REMOVED***
			switch err ***REMOVED***
			case unix.EBUSY:
				time.Sleep(50 * time.Millisecond)
				continue
			default:
				return err
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***
	return errors.Wrapf(unix.EBUSY, "failed to unmount target %s", target)
***REMOVED***

// UnmountAll repeatedly unmounts the given mount point until there
// are no mounts remaining (EINVAL is returned by mount), which is
// useful for undoing a stack of mounts on the same mount point.
func UnmountAll(mount string, flags int) error ***REMOVED***
	for ***REMOVED***
		if err := unmount(mount, flags); err != nil ***REMOVED***
			// EINVAL is returned if the target is not a
			// mount point, indicating that we are
			// done. It can also indicate a few other
			// things (such as invalid flags) which we
			// unfortunately end up squelching here too.
			if err == unix.EINVAL ***REMOVED***
				return nil
			***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
***REMOVED***

// parseMountOptions takes fstab style mount options and parses them for
// use with a standard mount() syscall
func parseMountOptions(options []string) (int, string) ***REMOVED***
	var (
		flag int
		data []string
	)
	flags := map[string]struct ***REMOVED***
		clear bool
		flag  int
	***REMOVED******REMOVED***
		"async":         ***REMOVED***true, unix.MS_SYNCHRONOUS***REMOVED***,
		"atime":         ***REMOVED***true, unix.MS_NOATIME***REMOVED***,
		"bind":          ***REMOVED***false, unix.MS_BIND***REMOVED***,
		"defaults":      ***REMOVED***false, 0***REMOVED***,
		"dev":           ***REMOVED***true, unix.MS_NODEV***REMOVED***,
		"diratime":      ***REMOVED***true, unix.MS_NODIRATIME***REMOVED***,
		"dirsync":       ***REMOVED***false, unix.MS_DIRSYNC***REMOVED***,
		"exec":          ***REMOVED***true, unix.MS_NOEXEC***REMOVED***,
		"mand":          ***REMOVED***false, unix.MS_MANDLOCK***REMOVED***,
		"noatime":       ***REMOVED***false, unix.MS_NOATIME***REMOVED***,
		"nodev":         ***REMOVED***false, unix.MS_NODEV***REMOVED***,
		"nodiratime":    ***REMOVED***false, unix.MS_NODIRATIME***REMOVED***,
		"noexec":        ***REMOVED***false, unix.MS_NOEXEC***REMOVED***,
		"nomand":        ***REMOVED***true, unix.MS_MANDLOCK***REMOVED***,
		"norelatime":    ***REMOVED***true, unix.MS_RELATIME***REMOVED***,
		"nostrictatime": ***REMOVED***true, unix.MS_STRICTATIME***REMOVED***,
		"nosuid":        ***REMOVED***false, unix.MS_NOSUID***REMOVED***,
		"rbind":         ***REMOVED***false, unix.MS_BIND | unix.MS_REC***REMOVED***,
		"relatime":      ***REMOVED***false, unix.MS_RELATIME***REMOVED***,
		"remount":       ***REMOVED***false, unix.MS_REMOUNT***REMOVED***,
		"ro":            ***REMOVED***false, unix.MS_RDONLY***REMOVED***,
		"rw":            ***REMOVED***true, unix.MS_RDONLY***REMOVED***,
		"strictatime":   ***REMOVED***false, unix.MS_STRICTATIME***REMOVED***,
		"suid":          ***REMOVED***true, unix.MS_NOSUID***REMOVED***,
		"sync":          ***REMOVED***false, unix.MS_SYNCHRONOUS***REMOVED***,
	***REMOVED***
	for _, o := range options ***REMOVED***
		// If the option does not exist in the flags table or the flag
		// is not supported on the platform,
		// then it is a data value for a specific fs type
		if f, exists := flags[o]; exists && f.flag != 0 ***REMOVED***
			if f.clear ***REMOVED***
				flag &^= f.flag
			***REMOVED*** else ***REMOVED***
				flag |= f.flag
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			data = append(data, o)
		***REMOVED***
	***REMOVED***
	return flag, strings.Join(data, ",")
***REMOVED***
