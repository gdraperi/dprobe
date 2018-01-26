package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
)

var bufferPool = &sync.Pool***REMOVED***
	New: func() interface***REMOVED******REMOVED*** ***REMOVED***
		buffer := make([]byte, 32*1024)
		return &buffer
	***REMOVED***,
***REMOVED***

// CopyDir copies the directory from src to dst.
// Most efficient copy of files is attempted.
func CopyDir(dst, src string) error ***REMOVED***
	inodes := map[uint64]string***REMOVED******REMOVED***
	return copyDirectory(dst, src, inodes)
***REMOVED***

func copyDirectory(dst, src string, inodes map[uint64]string) error ***REMOVED***
	stat, err := os.Stat(src)
	if err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to stat %s", src)
	***REMOVED***
	if !stat.IsDir() ***REMOVED***
		return errors.Errorf("source is not directory")
	***REMOVED***

	if st, err := os.Stat(dst); err != nil ***REMOVED***
		if err := os.Mkdir(dst, stat.Mode()); err != nil ***REMOVED***
			return errors.Wrapf(err, "failed to mkdir %s", dst)
		***REMOVED***
	***REMOVED*** else if !st.IsDir() ***REMOVED***
		return errors.Errorf("cannot copy to non-directory: %s", dst)
	***REMOVED*** else ***REMOVED***
		if err := os.Chmod(dst, stat.Mode()); err != nil ***REMOVED***
			return errors.Wrapf(err, "failed to chmod on %s", dst)
		***REMOVED***
	***REMOVED***

	fis, err := ioutil.ReadDir(src)
	if err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to read %s", src)
	***REMOVED***

	if err := copyFileInfo(stat, dst); err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to copy file info for %s", dst)
	***REMOVED***

	for _, fi := range fis ***REMOVED***
		source := filepath.Join(src, fi.Name())
		target := filepath.Join(dst, fi.Name())

		switch ***REMOVED***
		case fi.IsDir():
			if err := copyDirectory(target, source, inodes); err != nil ***REMOVED***
				return err
			***REMOVED***
			continue
		case (fi.Mode() & os.ModeType) == 0:
			link, err := getLinkSource(target, fi, inodes)
			if err != nil ***REMOVED***
				return errors.Wrap(err, "failed to get hardlink")
			***REMOVED***
			if link != "" ***REMOVED***
				if err := os.Link(link, target); err != nil ***REMOVED***
					return errors.Wrap(err, "failed to create hard link")
				***REMOVED***
			***REMOVED*** else if err := copyFile(source, target); err != nil ***REMOVED***
				return errors.Wrap(err, "failed to copy files")
			***REMOVED***
		case (fi.Mode() & os.ModeSymlink) == os.ModeSymlink:
			link, err := os.Readlink(source)
			if err != nil ***REMOVED***
				return errors.Wrapf(err, "failed to read link: %s", source)
			***REMOVED***
			if err := os.Symlink(link, target); err != nil ***REMOVED***
				return errors.Wrapf(err, "failed to create symlink: %s", target)
			***REMOVED***
		case (fi.Mode() & os.ModeDevice) == os.ModeDevice:
			if err := copyDevice(target, fi); err != nil ***REMOVED***
				return errors.Wrapf(err, "failed to create device")
			***REMOVED***
		default:
			// TODO: Support pipes and sockets
			return errors.Wrapf(err, "unsupported mode %s", fi.Mode())
		***REMOVED***
		if err := copyFileInfo(fi, target); err != nil ***REMOVED***
			return errors.Wrap(err, "failed to copy file info")
		***REMOVED***

		if err := copyXAttrs(target, source); err != nil ***REMOVED***
			return errors.Wrap(err, "failed to copy xattrs")
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func copyFile(source, target string) error ***REMOVED***
	src, err := os.Open(source)
	if err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to open source %s", source)
	***REMOVED***
	defer src.Close()
	tgt, err := os.Create(target)
	if err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to open target %s", target)
	***REMOVED***
	defer tgt.Close()

	return copyFileContent(tgt, src)
***REMOVED***
