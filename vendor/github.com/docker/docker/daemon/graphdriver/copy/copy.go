// +build linux

package copy

/*
#include <linux/fs.h>

#ifndef FICLONE
#define FICLONE		_IOW(0x94, 9, int)
#endif
*/
import "C"
import (
	"container/list"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/docker/docker/pkg/pools"
	"github.com/docker/docker/pkg/system"
	rsystem "github.com/opencontainers/runc/libcontainer/system"
	"golang.org/x/sys/unix"
)

// Mode indicates whether to use hardlink or copy content
type Mode int

const (
	// Content creates a new file, and copies the content of the file
	Content Mode = iota
	// Hardlink creates a new hardlink to the existing file
	Hardlink
)

func copyRegular(srcPath, dstPath string, fileinfo os.FileInfo, copyWithFileRange, copyWithFileClone *bool) error ***REMOVED***
	srcFile, err := os.Open(srcPath)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer srcFile.Close()

	// If the destination file already exists, we shouldn't blow it away
	dstFile, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, fileinfo.Mode())
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer dstFile.Close()

	if *copyWithFileClone ***REMOVED***
		_, _, err = unix.Syscall(unix.SYS_IOCTL, dstFile.Fd(), C.FICLONE, srcFile.Fd())
		if err == nil ***REMOVED***
			return nil
		***REMOVED***

		*copyWithFileClone = false
		if err == unix.EXDEV ***REMOVED***
			*copyWithFileRange = false
		***REMOVED***
	***REMOVED***
	if *copyWithFileRange ***REMOVED***
		err = doCopyWithFileRange(srcFile, dstFile, fileinfo)
		// Trying the file_clone may not have caught the exdev case
		// as the ioctl may not have been available (therefore EINVAL)
		if err == unix.EXDEV || err == unix.ENOSYS ***REMOVED***
			*copyWithFileRange = false
		***REMOVED*** else ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return legacyCopy(srcFile, dstFile)
***REMOVED***

func doCopyWithFileRange(srcFile, dstFile *os.File, fileinfo os.FileInfo) error ***REMOVED***
	amountLeftToCopy := fileinfo.Size()

	for amountLeftToCopy > 0 ***REMOVED***
		n, err := unix.CopyFileRange(int(srcFile.Fd()), nil, int(dstFile.Fd()), nil, int(amountLeftToCopy), 0)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		amountLeftToCopy = amountLeftToCopy - int64(n)
	***REMOVED***

	return nil
***REMOVED***

func legacyCopy(srcFile io.Reader, dstFile io.Writer) error ***REMOVED***
	_, err := pools.Copy(dstFile, srcFile)

	return err
***REMOVED***

func copyXattr(srcPath, dstPath, attr string) error ***REMOVED***
	data, err := system.Lgetxattr(srcPath, attr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if data != nil ***REMOVED***
		if err := system.Lsetxattr(dstPath, attr, data, 0); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

type fileID struct ***REMOVED***
	dev uint64
	ino uint64
***REMOVED***

type dirMtimeInfo struct ***REMOVED***
	dstPath *string
	stat    *syscall.Stat_t
***REMOVED***

// DirCopy copies or hardlinks the contents of one directory to another,
// properly handling xattrs, and soft links
//
// Copying xattrs can be opted out of by passing false for copyXattrs.
func DirCopy(srcDir, dstDir string, copyMode Mode, copyXattrs bool) error ***REMOVED***
	copyWithFileRange := true
	copyWithFileClone := true

	// This is a map of source file inodes to dst file paths
	copiedFiles := make(map[fileID]string)

	dirsToSetMtimes := list.New()
	err := filepath.Walk(srcDir, func(srcPath string, f os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Rebase path
		relPath, err := filepath.Rel(srcDir, srcPath)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		dstPath := filepath.Join(dstDir, relPath)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		stat, ok := f.Sys().(*syscall.Stat_t)
		if !ok ***REMOVED***
			return fmt.Errorf("Unable to get raw syscall.Stat_t data for %s", srcPath)
		***REMOVED***

		isHardlink := false

		switch f.Mode() & os.ModeType ***REMOVED***
		case 0: // Regular file
			id := fileID***REMOVED***dev: stat.Dev, ino: stat.Ino***REMOVED***
			if copyMode == Hardlink ***REMOVED***
				isHardlink = true
				if err2 := os.Link(srcPath, dstPath); err2 != nil ***REMOVED***
					return err2
				***REMOVED***
			***REMOVED*** else if hardLinkDstPath, ok := copiedFiles[id]; ok ***REMOVED***
				if err2 := os.Link(hardLinkDstPath, dstPath); err2 != nil ***REMOVED***
					return err2
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				if err2 := copyRegular(srcPath, dstPath, f, &copyWithFileRange, &copyWithFileClone); err2 != nil ***REMOVED***
					return err2
				***REMOVED***
				copiedFiles[id] = dstPath
			***REMOVED***

		case os.ModeDir:
			if err := os.Mkdir(dstPath, f.Mode()); err != nil && !os.IsExist(err) ***REMOVED***
				return err
			***REMOVED***

		case os.ModeSymlink:
			link, err := os.Readlink(srcPath)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			if err := os.Symlink(link, dstPath); err != nil ***REMOVED***
				return err
			***REMOVED***

		case os.ModeNamedPipe:
			fallthrough
		case os.ModeSocket:
			if rsystem.RunningInUserNS() ***REMOVED***
				// cannot create a device if running in user namespace
				return nil
			***REMOVED***
			if err := unix.Mkfifo(dstPath, stat.Mode); err != nil ***REMOVED***
				return err
			***REMOVED***

		case os.ModeDevice:
			if err := unix.Mknod(dstPath, stat.Mode, int(stat.Rdev)); err != nil ***REMOVED***
				return err
			***REMOVED***

		default:
			return fmt.Errorf("unknown file type for %s", srcPath)
		***REMOVED***

		// Everything below is copying metadata from src to dst. All this metadata
		// already shares an inode for hardlinks.
		if isHardlink ***REMOVED***
			return nil
		***REMOVED***

		if err := os.Lchown(dstPath, int(stat.Uid), int(stat.Gid)); err != nil ***REMOVED***
			return err
		***REMOVED***

		if copyXattrs ***REMOVED***
			if err := doCopyXattrs(srcPath, dstPath); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		isSymlink := f.Mode()&os.ModeSymlink != 0

		// There is no LChmod, so ignore mode for symlink. Also, this
		// must happen after chown, as that can modify the file mode
		if !isSymlink ***REMOVED***
			if err := os.Chmod(dstPath, f.Mode()); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		// system.Chtimes doesn't support a NOFOLLOW flag atm
		// nolint: unconvert
		if f.IsDir() ***REMOVED***
			dirsToSetMtimes.PushFront(&dirMtimeInfo***REMOVED***dstPath: &dstPath, stat: stat***REMOVED***)
		***REMOVED*** else if !isSymlink ***REMOVED***
			aTime := time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec))
			mTime := time.Unix(int64(stat.Mtim.Sec), int64(stat.Mtim.Nsec))
			if err := system.Chtimes(dstPath, aTime, mTime); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			ts := []syscall.Timespec***REMOVED***stat.Atim, stat.Mtim***REMOVED***
			if err := system.LUtimesNano(dstPath, ts); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for e := dirsToSetMtimes.Front(); e != nil; e = e.Next() ***REMOVED***
		mtimeInfo := e.Value.(*dirMtimeInfo)
		ts := []syscall.Timespec***REMOVED***mtimeInfo.stat.Atim, mtimeInfo.stat.Mtim***REMOVED***
		if err := system.LUtimesNano(*mtimeInfo.dstPath, ts); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func doCopyXattrs(srcPath, dstPath string) error ***REMOVED***
	if err := copyXattr(srcPath, dstPath, "security.capability"); err != nil ***REMOVED***
		return err
	***REMOVED***

	// We need to copy this attribute if it appears in an overlay upper layer, as
	// this function is used to copy those. It is set by overlay if a directory
	// is removed and then re-created and should not inherit anything from the
	// same dir in the lower dir.
	return copyXattr(srcPath, dstPath, "trusted.overlay.opaque")
***REMOVED***
