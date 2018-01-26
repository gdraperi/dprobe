// Package tailfile provides helper functions to read the nth lines of any
// ReadSeeker.
package tailfile

import (
	"bytes"
	"errors"
	"io"
	"os"
)

const blockSize = 1024

var eol = []byte("\n")

// ErrNonPositiveLinesNumber is an error returned if the lines number was negative.
var ErrNonPositiveLinesNumber = errors.New("The number of lines to extract from the file must be positive")

//TailFile returns last n lines of reader f (could be a fil).
func TailFile(f io.ReadSeeker, n int) ([][]byte, error) ***REMOVED***
	if n <= 0 ***REMOVED***
		return nil, ErrNonPositiveLinesNumber
	***REMOVED***
	size, err := f.Seek(0, os.SEEK_END)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	block := -1
	var data []byte
	var cnt int
	for ***REMOVED***
		var b []byte
		step := int64(block * blockSize)
		left := size + step // how many bytes to beginning
		if left < 0 ***REMOVED***
			if _, err := f.Seek(0, os.SEEK_SET); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			b = make([]byte, blockSize+left)
			if _, err := f.Read(b); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			data = append(b, data...)
			break
		***REMOVED*** else ***REMOVED***
			b = make([]byte, blockSize)
			if _, err := f.Seek(left, os.SEEK_SET); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			if _, err := f.Read(b); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			data = append(b, data...)
		***REMOVED***
		cnt += bytes.Count(b, eol)
		if cnt > n ***REMOVED***
			break
		***REMOVED***
		block--
	***REMOVED***
	lines := bytes.Split(data, eol)
	if n < len(lines) ***REMOVED***
		return lines[len(lines)-n-1 : len(lines)-1], nil
	***REMOVED***
	return lines[:len(lines)-1], nil
***REMOVED***
