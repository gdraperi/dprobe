// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux darwin dragonfly freebsd openbsd netbsd solaris

package tar

import (
	"io"
	"os"
	"runtime"
	"syscall"
)

func init() ***REMOVED***
	sysSparseDetect = sparseDetectUnix
***REMOVED***

func sparseDetectUnix(f *os.File) (sph sparseHoles, err error) ***REMOVED***
	// SEEK_DATA and SEEK_HOLE originated from Solaris and support for it
	// has been added to most of the other major Unix systems.
	var seekData, seekHole = 3, 4 // SEEK_DATA/SEEK_HOLE from unistd.h

	if runtime.GOOS == "darwin" ***REMOVED***
		// Darwin has the constants swapped, compared to all other UNIX.
		seekData, seekHole = 4, 3
	***REMOVED***

	// Check for seekData/seekHole support.
	// Different OS and FS may differ in the exact errno that is returned when
	// there is no support. Rather than special-casing every possible errno
	// representing "not supported", just assume that a non-nil error means
	// that seekData/seekHole is not supported.
	if _, err := f.Seek(0, seekHole); err != nil ***REMOVED***
		return nil, nil
	***REMOVED***

	// Populate the SparseHoles.
	var last, pos int64 = -1, 0
	for ***REMOVED***
		// Get the location of the next hole section.
		if pos, err = fseek(f, pos, seekHole); pos == last || err != nil ***REMOVED***
			return sph, err
		***REMOVED***
		offset := pos
		last = pos

		// Get the location of the next data section.
		if pos, err = fseek(f, pos, seekData); pos == last || err != nil ***REMOVED***
			return sph, err
		***REMOVED***
		length := pos - offset
		last = pos

		if length > 0 ***REMOVED***
			sph = append(sph, SparseEntry***REMOVED***offset, length***REMOVED***)
		***REMOVED***
	***REMOVED***
***REMOVED***

func fseek(f *os.File, pos int64, whence int) (int64, error) ***REMOVED***
	pos, err := f.Seek(pos, whence)
	if errno(err) == syscall.ENXIO ***REMOVED***
		// SEEK_DATA returns ENXIO when past the last data fragment,
		// which makes determining the size of the last hole difficult.
		pos, err = f.Seek(0, io.SeekEnd)
	***REMOVED***
	return pos, err
***REMOVED***

func errno(err error) error ***REMOVED***
	if perr, ok := err.(*os.PathError); ok ***REMOVED***
		return perr.Err
	***REMOVED***
	return err
***REMOVED***
