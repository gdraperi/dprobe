// +build windows

package client

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"syscall"
	"time"
	"unsafe"

	"github.com/sirupsen/logrus"
)

var (
	modkernel32   = syscall.NewLazyDLL("kernel32.dll")
	procCopyFileW = modkernel32.NewProc("CopyFileW")
)

// writeFileFromReader writes an output file from an io.Reader
func writeFileFromReader(path string, reader io.Reader, timeoutSeconds int, context string) (int64, error) ***REMOVED***
	outFile, err := os.Create(path)
	if err != nil ***REMOVED***
		return 0, fmt.Errorf("opengcs: writeFileFromReader: failed to create %s: %s", path, err)
	***REMOVED***
	defer outFile.Close()
	return copyWithTimeout(outFile, reader, 0, timeoutSeconds, context)
***REMOVED***

// copyWithTimeout is a wrapper for io.Copy using a timeout duration
func copyWithTimeout(dst io.Writer, src io.Reader, size int64, timeoutSeconds int, context string) (int64, error) ***REMOVED***
	logrus.Debugf("opengcs: copywithtimeout: size %d: timeout %d: (%s)", size, timeoutSeconds, context)

	type resultType struct ***REMOVED***
		err   error
		bytes int64
	***REMOVED***

	done := make(chan resultType, 1)
	go func() ***REMOVED***
		result := resultType***REMOVED******REMOVED***
		if logrus.GetLevel() < logrus.DebugLevel || logDataFromUVM == 0 ***REMOVED***
			result.bytes, result.err = io.Copy(dst, src)
		***REMOVED*** else ***REMOVED***
			// In advanced debug mode where we log (hexdump format) what is copied
			// up to the number of bytes defined by environment variable
			// OPENGCS_LOG_DATA_FROM_UVM
			var buf bytes.Buffer
			tee := io.TeeReader(src, &buf)
			result.bytes, result.err = io.Copy(dst, tee)
			if result.err == nil ***REMOVED***
				size := result.bytes
				if size > logDataFromUVM ***REMOVED***
					size = logDataFromUVM
				***REMOVED***
				if size > 0 ***REMOVED***
					bytes := make([]byte, size)
					if _, err := buf.Read(bytes); err == nil ***REMOVED***
						logrus.Debugf(fmt.Sprintf("opengcs: copyWithTimeout\n%s", hex.Dump(bytes)))
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		done <- result
	***REMOVED***()

	var result resultType
	timedout := time.After(time.Duration(timeoutSeconds) * time.Second)

	select ***REMOVED***
	case <-timedout:
		return 0, fmt.Errorf("opengcs: copyWithTimeout: timed out (%s)", context)
	case result = <-done:
		if result.err != nil && result.err != io.EOF ***REMOVED***
			// See https://github.com/golang/go/blob/f3f29d1dea525f48995c1693c609f5e67c046893/src/os/exec/exec_windows.go for a clue as to why we are doing this :)
			if se, ok := result.err.(syscall.Errno); ok ***REMOVED***
				const (
					errNoData     = syscall.Errno(232)
					errBrokenPipe = syscall.Errno(109)
				)
				if se == errNoData || se == errBrokenPipe ***REMOVED***
					logrus.Debugf("opengcs: copyWithTimeout: hit NoData or BrokenPipe: %d: %s", se, context)
					return result.bytes, nil
				***REMOVED***
			***REMOVED***
			return 0, fmt.Errorf("opengcs: copyWithTimeout: error reading: '%s' after %d bytes (%s)", result.err, result.bytes, context)
		***REMOVED***
	***REMOVED***
	logrus.Debugf("opengcs: copyWithTimeout: success - copied %d bytes (%s)", result.bytes, context)
	return result.bytes, nil
***REMOVED***

// CopyFile is a utility for copying a file - used for the sandbox cache.
// Uses CopyFileW win32 API for performance
func CopyFile(srcFile, destFile string, overwrite bool) error ***REMOVED***
	var bFailIfExists uint32 = 1
	if overwrite ***REMOVED***
		bFailIfExists = 0
	***REMOVED***

	lpExistingFileName, err := syscall.UTF16PtrFromString(srcFile)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	lpNewFileName, err := syscall.UTF16PtrFromString(destFile)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	r1, _, err := syscall.Syscall(
		procCopyFileW.Addr(),
		3,
		uintptr(unsafe.Pointer(lpExistingFileName)),
		uintptr(unsafe.Pointer(lpNewFileName)),
		uintptr(bFailIfExists))
	if r1 == 0 ***REMOVED***
		return fmt.Errorf("failed CopyFileW Win32 call from '%s' to '%s': %s", srcFile, destFile, err)
	***REMOVED***
	return nil
***REMOVED***
