// +build linux

package overlay2

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"syscall"

	"github.com/docker/docker/pkg/system"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

// doesSupportNativeDiff checks whether the filesystem has a bug
// which copies up the opaque flag when copying up an opaque
// directory or the kernel enable CONFIG_OVERLAY_FS_REDIRECT_DIR.
// When these exist naive diff should be used.
func doesSupportNativeDiff(d string) error ***REMOVED***
	td, err := ioutil.TempDir(d, "opaque-bug-check")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer func() ***REMOVED***
		if err := os.RemoveAll(td); err != nil ***REMOVED***
			logrus.Warnf("Failed to remove check directory %v: %v", td, err)
		***REMOVED***
	***REMOVED***()

	// Make directories l1/d, l1/d1, l2/d, l3, work, merged
	if err := os.MkdirAll(filepath.Join(td, "l1", "d"), 0755); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := os.MkdirAll(filepath.Join(td, "l1", "d1"), 0755); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := os.MkdirAll(filepath.Join(td, "l2", "d"), 0755); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := os.Mkdir(filepath.Join(td, "l3"), 0755); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := os.Mkdir(filepath.Join(td, "work"), 0755); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := os.Mkdir(filepath.Join(td, "merged"), 0755); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Mark l2/d as opaque
	if err := system.Lsetxattr(filepath.Join(td, "l2", "d"), "trusted.overlay.opaque", []byte("y"), 0); err != nil ***REMOVED***
		return errors.Wrap(err, "failed to set opaque flag on middle layer")
	***REMOVED***

	opts := fmt.Sprintf("lowerdir=%s:%s,upperdir=%s,workdir=%s", path.Join(td, "l2"), path.Join(td, "l1"), path.Join(td, "l3"), path.Join(td, "work"))
	if err := unix.Mount("overlay", filepath.Join(td, "merged"), "overlay", 0, opts); err != nil ***REMOVED***
		return errors.Wrap(err, "failed to mount overlay")
	***REMOVED***
	defer func() ***REMOVED***
		if err := unix.Unmount(filepath.Join(td, "merged"), 0); err != nil ***REMOVED***
			logrus.Warnf("Failed to unmount check directory %v: %v", filepath.Join(td, "merged"), err)
		***REMOVED***
	***REMOVED***()

	// Touch file in d to force copy up of opaque directory "d" from "l2" to "l3"
	if err := ioutil.WriteFile(filepath.Join(td, "merged", "d", "f"), []byte***REMOVED******REMOVED***, 0644); err != nil ***REMOVED***
		return errors.Wrap(err, "failed to write to merged directory")
	***REMOVED***

	// Check l3/d does not have opaque flag
	xattrOpaque, err := system.Lgetxattr(filepath.Join(td, "l3", "d"), "trusted.overlay.opaque")
	if err != nil ***REMOVED***
		return errors.Wrap(err, "failed to read opaque flag on upper layer")
	***REMOVED***
	if string(xattrOpaque) == "y" ***REMOVED***
		return errors.New("opaque flag erroneously copied up, consider update to kernel 4.8 or later to fix")
	***REMOVED***

	// rename "d1" to "d2"
	if err := os.Rename(filepath.Join(td, "merged", "d1"), filepath.Join(td, "merged", "d2")); err != nil ***REMOVED***
		// if rename failed with syscall.EXDEV, the kernel doesn't have CONFIG_OVERLAY_FS_REDIRECT_DIR enabled
		if err.(*os.LinkError).Err == syscall.EXDEV ***REMOVED***
			return nil
		***REMOVED***
		return errors.Wrap(err, "failed to rename dir in merged directory")
	***REMOVED***
	// get the xattr of "d2"
	xattrRedirect, err := system.Lgetxattr(filepath.Join(td, "l3", "d2"), "trusted.overlay.redirect")
	if err != nil ***REMOVED***
		return errors.Wrap(err, "failed to read redirect flag on upper layer")
	***REMOVED***

	if string(xattrRedirect) == "d1" ***REMOVED***
		return errors.New("kernel has CONFIG_OVERLAY_FS_REDIRECT_DIR enabled")
	***REMOVED***

	return nil
***REMOVED***

// supportsMultipleLowerDir checks if the system supports multiple lowerdirs,
// which is required for the overlay2 driver. On 4.x kernels, multiple lowerdirs
// are always available (so this check isn't needed), and backported to RHEL and
// CentOS 3.x kernels (3.10.0-693.el7.x86_64 and up). This function is to detect
// support on those kernels, without doing a kernel version compare.
func supportsMultipleLowerDir(d string) error ***REMOVED***
	td, err := ioutil.TempDir(d, "multiple-lowerdir-check")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer func() ***REMOVED***
		if err := os.RemoveAll(td); err != nil ***REMOVED***
			logrus.Warnf("Failed to remove check directory %v: %v", td, err)
		***REMOVED***
	***REMOVED***()

	for _, dir := range []string***REMOVED***"lower1", "lower2", "upper", "work", "merged"***REMOVED*** ***REMOVED***
		if err := os.Mkdir(filepath.Join(td, dir), 0755); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	opts := fmt.Sprintf("lowerdir=%s:%s,upperdir=%s,workdir=%s", path.Join(td, "lower2"), path.Join(td, "lower1"), path.Join(td, "upper"), path.Join(td, "work"))
	if err := unix.Mount("overlay", filepath.Join(td, "merged"), "overlay", 0, opts); err != nil ***REMOVED***
		return errors.Wrap(err, "failed to mount overlay")
	***REMOVED***
	if err := unix.Unmount(filepath.Join(td, "merged"), 0); err != nil ***REMOVED***
		logrus.Warnf("Failed to unmount check directory %v: %v", filepath.Join(td, "merged"), err)
	***REMOVED***
	return nil
***REMOVED***
