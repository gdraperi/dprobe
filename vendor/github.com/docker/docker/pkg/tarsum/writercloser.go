package tarsum

import (
	"io"
)

type writeCloseFlusher interface ***REMOVED***
	io.WriteCloser
	Flush() error
***REMOVED***

type nopCloseFlusher struct ***REMOVED***
	io.Writer
***REMOVED***

func (n *nopCloseFlusher) Close() error ***REMOVED***
	return nil
***REMOVED***

func (n *nopCloseFlusher) Flush() error ***REMOVED***
	return nil
***REMOVED***
