// +build solaris darwin freebsd

package fs

import (
	"io"
	"os"
	"syscall"

	"github.com/containerd/containerd/sys"
	"github.com/containerd/continuity/sysx"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

func copyFileInfo(fi os.FileInfo, name string) error ***REMOVED***
	st := fi.Sys().(*syscall.Stat_t)
	if err := os.Lchown(name, int(st.Uid), int(st.Gid)); err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to chown %s", name)
	***REMOVED***

	if (fi.Mode() & os.ModeSymlink) != os.ModeSymlink ***REMOVED***
		if err := os.Chmod(name, fi.Mode()); err != nil ***REMOVED***
			return errors.Wrapf(err, "failed to chmod %s", name)
		***REMOVED***
	***REMOVED***

	timespec := []syscall.Timespec***REMOVED***sys.StatAtime(st), sys.StatMtime(st)***REMOVED***
	if err := syscall.UtimesNano(name, timespec); err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to utime %s", name)
	***REMOVED***

	return nil
***REMOVED***

func copyFileContent(dst, src *os.File) error ***REMOVED***
	buf := bufferPool.Get().(*[]byte)
	_, err := io.CopyBuffer(dst, src, *buf)
	bufferPool.Put(buf)

	return err
***REMOVED***

func copyXAttrs(dst, src string) error ***REMOVED***
	xattrKeys, err := sysx.LListxattr(src)
	if err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to list xattrs on %s", src)
	***REMOVED***
	for _, xattr := range xattrKeys ***REMOVED***
		data, err := sysx.LGetxattr(src, xattr)
		if err != nil ***REMOVED***
			return errors.Wrapf(err, "failed to get xattr %q on %s", xattr, src)
		***REMOVED***
		if err := sysx.LSetxattr(dst, xattr, data, 0); err != nil ***REMOVED***
			return errors.Wrapf(err, "failed to set xattr %q on %s", xattr, dst)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func copyDevice(dst string, fi os.FileInfo) error ***REMOVED***
	st, ok := fi.Sys().(*syscall.Stat_t)
	if !ok ***REMOVED***
		return errors.New("unsupported stat type")
	***REMOVED***
	return unix.Mknod(dst, uint32(fi.Mode()), int(st.Rdev))
***REMOVED***
