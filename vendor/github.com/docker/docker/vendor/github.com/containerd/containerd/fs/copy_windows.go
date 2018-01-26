package fs

import (
	"io"
	"os"

	"github.com/pkg/errors"
)

func copyFileInfo(fi os.FileInfo, name string) error ***REMOVED***
	if err := os.Chmod(name, fi.Mode()); err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to chmod %s", name)
	***REMOVED***

	// TODO: copy windows specific metadata

	return nil
***REMOVED***

func copyFileContent(dst, src *os.File) error ***REMOVED***
	buf := bufferPool.Get().(*[]byte)
	_, err := io.CopyBuffer(dst, src, *buf)
	bufferPool.Put(buf)
	return err
***REMOVED***

func copyXAttrs(dst, src string) error ***REMOVED***
	return nil
***REMOVED***

func copyDevice(dst string, fi os.FileInfo) error ***REMOVED***
	return errors.New("device copy not supported")
***REMOVED***
