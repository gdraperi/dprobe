package broadcaster

import (
	"io"
	"sync"
)

// Unbuffered accumulates multiple io.WriteCloser by stream.
type Unbuffered struct ***REMOVED***
	mu      sync.Mutex
	writers []io.WriteCloser
***REMOVED***

// Add adds new io.WriteCloser.
func (w *Unbuffered) Add(writer io.WriteCloser) ***REMOVED***
	w.mu.Lock()
	w.writers = append(w.writers, writer)
	w.mu.Unlock()
***REMOVED***

// Write writes bytes to all writers. Failed writers will be evicted during
// this call.
func (w *Unbuffered) Write(p []byte) (n int, err error) ***REMOVED***
	w.mu.Lock()
	var evict []int
	for i, sw := range w.writers ***REMOVED***
		if n, err := sw.Write(p); err != nil || n != len(p) ***REMOVED***
			// On error, evict the writer
			evict = append(evict, i)
		***REMOVED***
	***REMOVED***
	for n, i := range evict ***REMOVED***
		w.writers = append(w.writers[:i-n], w.writers[i-n+1:]...)
	***REMOVED***
	w.mu.Unlock()
	return len(p), nil
***REMOVED***

// Clean closes and removes all writers. Last non-eol-terminated part of data
// will be saved.
func (w *Unbuffered) Clean() error ***REMOVED***
	w.mu.Lock()
	for _, sw := range w.writers ***REMOVED***
		sw.Close()
	***REMOVED***
	w.writers = nil
	w.mu.Unlock()
	return nil
***REMOVED***
