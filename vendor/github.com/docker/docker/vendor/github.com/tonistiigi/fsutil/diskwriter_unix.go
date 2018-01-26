// +build !windows

package fsutil

import (
	"os"
	"syscall"

	"github.com/containerd/continuity/sysx"
	"github.com/pkg/errors"
)

func rewriteMetadata(p string, stat *Stat) error ***REMOVED***
	for key, value := range stat.Xattrs ***REMOVED***
		sysx.Setxattr(p, key, value, 0)
	***REMOVED***

	if err := os.Lchown(p, int(stat.Uid), int(stat.Gid)); err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to lchown %s", p)
	***REMOVED***

	if os.FileMode(stat.Mode)&os.ModeSymlink == 0 ***REMOVED***
		if err := os.Chmod(p, os.FileMode(stat.Mode)); err != nil ***REMOVED***
			return errors.Wrapf(err, "failed to chown %s", p)
		***REMOVED***
	***REMOVED***

	if err := chtimes(p, stat.ModTime); err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to chtimes %s", p)
	***REMOVED***

	return nil
***REMOVED***

// handleTarTypeBlockCharFifo is an OS-specific helper function used by
// createTarFile to handle the following types of header: Block; Char; Fifo
func handleTarTypeBlockCharFifo(path string, stat *Stat) error ***REMOVED***
	mode := uint32(stat.Mode & 07777)
	if os.FileMode(stat.Mode)&os.ModeCharDevice != 0 ***REMOVED***
		mode |= syscall.S_IFCHR
	***REMOVED*** else if os.FileMode(stat.Mode)&os.ModeNamedPipe != 0 ***REMOVED***
		mode |= syscall.S_IFIFO
	***REMOVED*** else ***REMOVED***
		mode |= syscall.S_IFBLK
	***REMOVED***

	if err := syscall.Mknod(path, mode, int(mkdev(stat.Devmajor, stat.Devminor))); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***
