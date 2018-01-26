package devices

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/opencontainers/runc/libcontainer/configs"

	"golang.org/x/sys/unix"
)

var (
	ErrNotADevice = errors.New("not a device node")
)

// Testing dependencies
var (
	unixLstat     = unix.Lstat
	ioutilReadDir = ioutil.ReadDir
)

// Given the path to a device and its cgroup_permissions(which cannot be easily queried) look up the information about a linux device and return that information as a Device struct.
func DeviceFromPath(path, permissions string) (*configs.Device, error) ***REMOVED***
	var stat unix.Stat_t
	err := unixLstat(path, &stat)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var (
		devNumber = stat.Rdev
		major     = unix.Major(devNumber)
	)
	if major == 0 ***REMOVED***
		return nil, ErrNotADevice
	***REMOVED***

	var (
		devType rune
		mode    = stat.Mode
	)
	switch ***REMOVED***
	case mode&unix.S_IFBLK == unix.S_IFBLK:
		devType = 'b'
	case mode&unix.S_IFCHR == unix.S_IFCHR:
		devType = 'c'
	***REMOVED***
	return &configs.Device***REMOVED***
		Type:        devType,
		Path:        path,
		Major:       int64(major),
		Minor:       int64(unix.Minor(devNumber)),
		Permissions: permissions,
		FileMode:    os.FileMode(mode),
		Uid:         stat.Uid,
		Gid:         stat.Gid,
	***REMOVED***, nil
***REMOVED***

func HostDevices() ([]*configs.Device, error) ***REMOVED***
	return getDevices("/dev")
***REMOVED***

func getDevices(path string) ([]*configs.Device, error) ***REMOVED***
	files, err := ioutilReadDir(path)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	out := []*configs.Device***REMOVED******REMOVED***
	for _, f := range files ***REMOVED***
		switch ***REMOVED***
		case f.IsDir():
			switch f.Name() ***REMOVED***
			// ".lxc" & ".lxd-mounts" added to address https://github.com/lxc/lxd/issues/2825
			case "pts", "shm", "fd", "mqueue", ".lxc", ".lxd-mounts":
				continue
			default:
				sub, err := getDevices(filepath.Join(path, f.Name()))
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***

				out = append(out, sub...)
				continue
			***REMOVED***
		case f.Name() == "console":
			continue
		***REMOVED***
		device, err := DeviceFromPath(filepath.Join(path, f.Name()), "rwm")
		if err != nil ***REMOVED***
			if err == ErrNotADevice ***REMOVED***
				continue
			***REMOVED***
			if os.IsNotExist(err) ***REMOVED***
				continue
			***REMOVED***
			return nil, err
		***REMOVED***
		out = append(out, device)
	***REMOVED***
	return out, nil
***REMOVED***
