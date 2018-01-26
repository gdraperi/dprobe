package oci

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opencontainers/runc/libcontainer/configs"
	"github.com/opencontainers/runc/libcontainer/devices"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

// Device transforms a libcontainer configs.Device to a specs.LinuxDevice object.
func Device(d *configs.Device) specs.LinuxDevice ***REMOVED***
	return specs.LinuxDevice***REMOVED***
		Type:     string(d.Type),
		Path:     d.Path,
		Major:    d.Major,
		Minor:    d.Minor,
		FileMode: fmPtr(int64(d.FileMode)),
		UID:      u32Ptr(int64(d.Uid)),
		GID:      u32Ptr(int64(d.Gid)),
	***REMOVED***
***REMOVED***

func deviceCgroup(d *configs.Device) specs.LinuxDeviceCgroup ***REMOVED***
	t := string(d.Type)
	return specs.LinuxDeviceCgroup***REMOVED***
		Allow:  true,
		Type:   t,
		Major:  &d.Major,
		Minor:  &d.Minor,
		Access: d.Permissions,
	***REMOVED***
***REMOVED***

// DevicesFromPath computes a list of devices and device permissions from paths (pathOnHost and pathInContainer) and cgroup permissions.
func DevicesFromPath(pathOnHost, pathInContainer, cgroupPermissions string) (devs []specs.LinuxDevice, devPermissions []specs.LinuxDeviceCgroup, err error) ***REMOVED***
	resolvedPathOnHost := pathOnHost

	// check if it is a symbolic link
	if src, e := os.Lstat(pathOnHost); e == nil && src.Mode()&os.ModeSymlink == os.ModeSymlink ***REMOVED***
		if linkedPathOnHost, e := filepath.EvalSymlinks(pathOnHost); e == nil ***REMOVED***
			resolvedPathOnHost = linkedPathOnHost
		***REMOVED***
	***REMOVED***

	device, err := devices.DeviceFromPath(resolvedPathOnHost, cgroupPermissions)
	// if there was no error, return the device
	if err == nil ***REMOVED***
		device.Path = pathInContainer
		return append(devs, Device(device)), append(devPermissions, deviceCgroup(device)), nil
	***REMOVED***

	// if the device is not a device node
	// try to see if it's a directory holding many devices
	if err == devices.ErrNotADevice ***REMOVED***

		// check if it is a directory
		if src, e := os.Stat(resolvedPathOnHost); e == nil && src.IsDir() ***REMOVED***

			// mount the internal devices recursively
			filepath.Walk(resolvedPathOnHost, func(dpath string, f os.FileInfo, e error) error ***REMOVED***
				childDevice, e := devices.DeviceFromPath(dpath, cgroupPermissions)
				if e != nil ***REMOVED***
					// ignore the device
					return nil
				***REMOVED***

				// add the device to userSpecified devices
				childDevice.Path = strings.Replace(dpath, resolvedPathOnHost, pathInContainer, 1)
				devs = append(devs, Device(childDevice))
				devPermissions = append(devPermissions, deviceCgroup(childDevice))

				return nil
			***REMOVED***)
		***REMOVED***
	***REMOVED***

	if len(devs) > 0 ***REMOVED***
		return devs, devPermissions, nil
	***REMOVED***

	return devs, devPermissions, fmt.Errorf("error gathering device information while adding custom device %q: %s", pathOnHost, err)
***REMOVED***
