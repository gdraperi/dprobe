package term

import (
	"io"
)

// EscapeError is special error which returned by a TTY proxy reader's Read()
// method in case its detach escape sequence is read.
type EscapeError struct***REMOVED******REMOVED***

func (EscapeError) Error() string ***REMOVED***
	return "read escape sequence"
***REMOVED***

// escapeProxy is used only for attaches with a TTY. It is used to proxy
// stdin keypresses from the underlying reader and look for the passed in
// escape key sequence to signal a detach.
type escapeProxy struct ***REMOVED***
	escapeKeys   []byte
	escapeKeyPos int
	r            io.Reader
***REMOVED***

// NewEscapeProxy returns a new TTY proxy reader which wraps the given reader
// and detects when the specified escape keys are read, in which case the Read
// method will return an error of type EscapeError.
func NewEscapeProxy(r io.Reader, escapeKeys []byte) io.Reader ***REMOVED***
	return &escapeProxy***REMOVED***
		escapeKeys: escapeKeys,
		r:          r,
	***REMOVED***
***REMOVED***

func (r *escapeProxy) Read(buf []byte) (int, error) ***REMOVED***
	nr, err := r.r.Read(buf)

	preserve := func() ***REMOVED***
		// this preserves the original key presses in the passed in buffer
		nr += r.escapeKeyPos
		preserve := make([]byte, 0, r.escapeKeyPos+len(buf))
		preserve = append(preserve, r.escapeKeys[:r.escapeKeyPos]...)
		preserve = append(preserve, buf...)
		r.escapeKeyPos = 0
		copy(buf[0:nr], preserve)
	***REMOVED***

	if nr != 1 || err != nil ***REMOVED***
		if r.escapeKeyPos > 0 ***REMOVED***
			preserve()
		***REMOVED***
		return nr, err
	***REMOVED***

	if buf[0] != r.escapeKeys[r.escapeKeyPos] ***REMOVED***
		if r.escapeKeyPos > 0 ***REMOVED***
			preserve()
		***REMOVED***
		return nr, nil
	***REMOVED***

	if r.escapeKeyPos == len(r.escapeKeys)-1 ***REMOVED***
		return 0, EscapeError***REMOVED******REMOVED***
	***REMOVED***

	// Looks like we've got an escape key, but we need to match again on the next
	// read.
	// Store the current escape key we found so we can look for the next one on
	// the next read.
	// Since this is an escape key, make sure we don't let the caller read it
	// If later on we find that this is not the escape sequence, we'll add the
	// keys back
	r.escapeKeyPos++
	return nr - r.escapeKeyPos, nil
***REMOVED***
