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

	timespec := []unix.Timespec***REMOVED***unix.Timespec(sys.StatAtime(st)), unix.Timespec(sys.StatMtime(st))***REMOVED***
	if err := unix.UtimesNanoAt(unix.AT_FDCWD, name, timespec, unix.AT_SYMLINK_NOFOLLOW); err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to utime %s", name)
	***REMOVED***

	return nil
***REMOVED***

func copyFileContent(dst, src *os.File) error ***REMOVED***
	st, err := src.Stat()
	if err != nil ***REMOVED***
		return errors.Wrap(err, "unable to stat source")
	***REMOVED***

	n, err := unix.CopyFileRange(int(src.Fd()), nil, int(dst.Fd()), nil, int(st.Size()), 0)
	if err != nil ***REMOVED***
		if err != unix.ENOSYS && err != unix.EXDEV ***REMOVED***
			return errors.Wrap(err, "copy file range failed")
		***REMOVED***

		buf := bufferPool.Get().(*[]byte)
		_, err = io.CopyBuffer(dst, src, *buf)
		bufferPool.Put(buf)
		return err
	***REMOVED***

	if int64(n) != st.Size() ***REMOVED***
		return errors.Wrapf(err, "short copy: %d of %d", int64(n), st.Size())
	***REMOVED***

	return nil
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
