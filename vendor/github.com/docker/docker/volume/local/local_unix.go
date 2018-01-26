// +build linux freebsd

// Package local provides the default implementation for volumes. It
// is used to mount data volume containers and directories local to
// the host server.
package local

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"

	"github.com/docker/docker/pkg/mount"
)

var (
	oldVfsDir = filepath.Join("vfs", "dir")

	validOpts = map[string]bool***REMOVED***
		"type":   true, // specify the filesystem type for mount, e.g. nfs
		"o":      true, // generic mount options
		"device": true, // device to mount from
	***REMOVED***
)

type optsConfig struct ***REMOVED***
	MountType   string
	MountOpts   string
	MountDevice string
***REMOVED***

func (o *optsConfig) String() string ***REMOVED***
	return fmt.Sprintf("type='%s' device='%s' o='%s'", o.MountType, o.MountDevice, o.MountOpts)
***REMOVED***

// scopedPath verifies that the path where the volume is located
// is under Docker's root and the valid local paths.
func (r *Root) scopedPath(realPath string) bool ***REMOVED***
	// Volumes path for Docker version >= 1.7
	if strings.HasPrefix(realPath, filepath.Join(r.scope, volumesPathName)) && realPath != filepath.Join(r.scope, volumesPathName) ***REMOVED***
		return true
	***REMOVED***

	// Volumes path for Docker version < 1.7
	if strings.HasPrefix(realPath, filepath.Join(r.scope, oldVfsDir)) ***REMOVED***
		return true
	***REMOVED***

	return false
***REMOVED***

func setOpts(v *localVolume, opts map[string]string) error ***REMOVED***
	if len(opts) == 0 ***REMOVED***
		return nil
	***REMOVED***
	if err := validateOpts(opts); err != nil ***REMOVED***
		return err
	***REMOVED***

	v.opts = &optsConfig***REMOVED***
		MountType:   opts["type"],
		MountOpts:   opts["o"],
		MountDevice: opts["device"],
	***REMOVED***
	return nil
***REMOVED***

func (v *localVolume) mount() error ***REMOVED***
	if v.opts.MountDevice == "" ***REMOVED***
		return fmt.Errorf("missing device in volume options")
	***REMOVED***
	mountOpts := v.opts.MountOpts
	if v.opts.MountType == "nfs" ***REMOVED***
		if addrValue := getAddress(v.opts.MountOpts); addrValue != "" && net.ParseIP(addrValue).To4() == nil ***REMOVED***
			ipAddr, err := net.ResolveIPAddr("ip", addrValue)
			if err != nil ***REMOVED***
				return errors.Wrapf(err, "error resolving passed in nfs address")
			***REMOVED***
			mountOpts = strings.Replace(mountOpts, "addr="+addrValue, "addr="+ipAddr.String(), 1)
		***REMOVED***
	***REMOVED***
	err := mount.Mount(v.opts.MountDevice, v.path, v.opts.MountType, mountOpts)
	return errors.Wrapf(err, "error while mounting volume with options: %s", v.opts)
***REMOVED***

func (v *localVolume) CreatedAt() (time.Time, error) ***REMOVED***
	fileInfo, err := os.Stat(v.path)
	if err != nil ***REMOVED***
		return time.Time***REMOVED******REMOVED***, err
	***REMOVED***
	sec, nsec := fileInfo.Sys().(*syscall.Stat_t).Ctim.Unix()
	return time.Unix(sec, nsec), nil
***REMOVED***
