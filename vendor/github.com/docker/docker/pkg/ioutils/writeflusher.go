package ioutils

import (
	"io"
	"sync"
)

// WriteFlusher wraps the Write and Flush operation ensuring that every write
// is a flush. In addition, the Close method can be called to intercept
// Read/Write calls if the targets lifecycle has already ended.
type WriteFlusher struct ***REMOVED***
	w           io.Writer
	flusher     flusher
	flushed     chan struct***REMOVED******REMOVED***
	flushedOnce sync.Once
	closed      chan struct***REMOVED******REMOVED***
	closeLock   sync.Mutex
***REMOVED***

type flusher interface ***REMOVED***
	Flush()
***REMOVED***

var errWriteFlusherClosed = io.EOF

func (wf *WriteFlusher) Write(b []byte) (n int, err error) ***REMOVED***
	select ***REMOVED***
	case <-wf.closed:
		return 0, errWriteFlusherClosed
	default:
	***REMOVED***

	n, err = wf.w.Write(b)
	wf.Flush() // every write is a flush.
	return n, err
***REMOVED***

// Flush the stream immediately.
func (wf *WriteFlusher) Flush() ***REMOVED***
	select ***REMOVED***
	case <-wf.closed:
		return
	default:
	***REMOVED***

	wf.flushedOnce.Do(func() ***REMOVED***
		close(wf.flushed)
	***REMOVED***)
	wf.flusher.Flush()
***REMOVED***

// Flushed returns the state of flushed.
// If it's flushed, return true, or else it return false.
func (wf *WriteFlusher) Flushed() bool ***REMOVED***
	// BUG(stevvooe): Remove this method. Its use is inherently racy. Seems to
	// be used to detect whether or a response code has been issued or not.
	// Another hook should be used instead.
	var flushed bool
	select ***REMOVED***
	case <-wf.flushed:
		flushed = true
	default:
	***REMOVED***
	return flushed
***REMOVED***

// Close closes the write flusher, disallowing any further writes to the
// target. After the flusher is closed, all calls to write or flush will
// result in an error.
func (wf *WriteFlusher) Close() error ***REMOVED***
	wf.closeLock.Lock()
	defer wf.closeLock.Unlock()

	select ***REMOVED***
	case <-wf.closed:
		return errWriteFlusherClosed
	default:
		close(wf.closed)
	***REMOVED***
	return nil
***REMOVED***

// NewWriteFlusher returns a new WriteFlusher.
func NewWriteFlusher(w io.Writer) *WriteFlusher ***REMOVED***
	var fl flusher
	if f, ok := w.(flusher); ok ***REMOVED***
		fl = f
	***REMOVED*** else ***REMOVED***
		fl = &NopFlusher***REMOVED******REMOVED***
	***REMOVED***
	return &WriteFlusher***REMOVED***w: w, flusher: fl, closed: make(chan struct***REMOVED******REMOVED***), flushed: make(chan struct***REMOVED******REMOVED***)***REMOVED***
***REMOVED***
