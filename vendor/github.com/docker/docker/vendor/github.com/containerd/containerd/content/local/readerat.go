package local

import (
	"os"
)

// readerat implements io.ReaderAt in a completely stateless manner by opening
// the referenced file for each call to ReadAt.
type sizeReaderAt struct ***REMOVED***
	size int64
	fp   *os.File
***REMOVED***

func (ra sizeReaderAt) ReadAt(p []byte, offset int64) (int, error) ***REMOVED***
	return ra.fp.ReadAt(p, offset)
***REMOVED***

func (ra sizeReaderAt) Size() int64 ***REMOVED***
	return ra.size
***REMOVED***

func (ra sizeReaderAt) Close() error ***REMOVED***
	return ra.fp.Close()
***REMOVED***
