package stdcopy

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync"
)

// StdType is the type of standard stream
// a writer can multiplex to.
type StdType byte

const (
	// Stdin represents standard input stream type.
	Stdin StdType = iota
	// Stdout represents standard output stream type.
	Stdout
	// Stderr represents standard error steam type.
	Stderr
	// Systemerr represents errors originating from the system that make it
	// into the the multiplexed stream.
	Systemerr

	stdWriterPrefixLen = 8
	stdWriterFdIndex   = 0
	stdWriterSizeIndex = 4

	startingBufLen = 32*1024 + stdWriterPrefixLen + 1
)

var bufPool = &sync.Pool***REMOVED***New: func() interface***REMOVED******REMOVED*** ***REMOVED*** return bytes.NewBuffer(nil) ***REMOVED******REMOVED***

// stdWriter is wrapper of io.Writer with extra customized info.
type stdWriter struct ***REMOVED***
	io.Writer
	prefix byte
***REMOVED***

// Write sends the buffer to the underneath writer.
// It inserts the prefix header before the buffer,
// so stdcopy.StdCopy knows where to multiplex the output.
// It makes stdWriter to implement io.Writer.
func (w *stdWriter) Write(p []byte) (n int, err error) ***REMOVED***
	if w == nil || w.Writer == nil ***REMOVED***
		return 0, errors.New("Writer not instantiated")
	***REMOVED***
	if p == nil ***REMOVED***
		return 0, nil
	***REMOVED***

	header := [stdWriterPrefixLen]byte***REMOVED***stdWriterFdIndex: w.prefix***REMOVED***
	binary.BigEndian.PutUint32(header[stdWriterSizeIndex:], uint32(len(p)))
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Write(header[:])
	buf.Write(p)

	n, err = w.Writer.Write(buf.Bytes())
	n -= stdWriterPrefixLen
	if n < 0 ***REMOVED***
		n = 0
	***REMOVED***

	buf.Reset()
	bufPool.Put(buf)
	return
***REMOVED***

// NewStdWriter instantiates a new Writer.
// Everything written to it will be encapsulated using a custom format,
// and written to the underlying `w` stream.
// This allows multiple write streams (e.g. stdout and stderr) to be muxed into a single connection.
// `t` indicates the id of the stream to encapsulate.
// It can be stdcopy.Stdin, stdcopy.Stdout, stdcopy.Stderr.
func NewStdWriter(w io.Writer, t StdType) io.Writer ***REMOVED***
	return &stdWriter***REMOVED***
		Writer: w,
		prefix: byte(t),
	***REMOVED***
***REMOVED***

// StdCopy is a modified version of io.Copy.
//
// StdCopy will demultiplex `src`, assuming that it contains two streams,
// previously multiplexed together using a StdWriter instance.
// As it reads from `src`, StdCopy will write to `dstout` and `dsterr`.
//
// StdCopy will read until it hits EOF on `src`. It will then return a nil error.
// In other words: if `err` is non nil, it indicates a real underlying error.
//
// `written` will hold the total number of bytes written to `dstout` and `dsterr`.
func StdCopy(dstout, dsterr io.Writer, src io.Reader) (written int64, err error) ***REMOVED***
	var (
		buf       = make([]byte, startingBufLen)
		bufLen    = len(buf)
		nr, nw    int
		er, ew    error
		out       io.Writer
		frameSize int
	)

	for ***REMOVED***
		// Make sure we have at least a full header
		for nr < stdWriterPrefixLen ***REMOVED***
			var nr2 int
			nr2, er = src.Read(buf[nr:])
			nr += nr2
			if er == io.EOF ***REMOVED***
				if nr < stdWriterPrefixLen ***REMOVED***
					return written, nil
				***REMOVED***
				break
			***REMOVED***
			if er != nil ***REMOVED***
				return 0, er
			***REMOVED***
		***REMOVED***

		stream := StdType(buf[stdWriterFdIndex])
		// Check the first byte to know where to write
		switch stream ***REMOVED***
		case Stdin:
			fallthrough
		case Stdout:
			// Write on stdout
			out = dstout
		case Stderr:
			// Write on stderr
			out = dsterr
		case Systemerr:
			// If we're on Systemerr, we won't write anywhere.
			// NB: if this code changes later, make sure you don't try to write
			// to outstream if Systemerr is the stream
			out = nil
		default:
			return 0, fmt.Errorf("Unrecognized input header: %d", buf[stdWriterFdIndex])
		***REMOVED***

		// Retrieve the size of the frame
		frameSize = int(binary.BigEndian.Uint32(buf[stdWriterSizeIndex : stdWriterSizeIndex+4]))

		// Check if the buffer is big enough to read the frame.
		// Extend it if necessary.
		if frameSize+stdWriterPrefixLen > bufLen ***REMOVED***
			buf = append(buf, make([]byte, frameSize+stdWriterPrefixLen-bufLen+1)...)
			bufLen = len(buf)
		***REMOVED***

		// While the amount of bytes read is less than the size of the frame + header, we keep reading
		for nr < frameSize+stdWriterPrefixLen ***REMOVED***
			var nr2 int
			nr2, er = src.Read(buf[nr:])
			nr += nr2
			if er == io.EOF ***REMOVED***
				if nr < frameSize+stdWriterPrefixLen ***REMOVED***
					return written, nil
				***REMOVED***
				break
			***REMOVED***
			if er != nil ***REMOVED***
				return 0, er
			***REMOVED***
		***REMOVED***

		// we might have an error from the source mixed up in our multiplexed
		// stream. if we do, return it.
		if stream == Systemerr ***REMOVED***
			return written, fmt.Errorf("error from daemon in stream: %s", string(buf[stdWriterPrefixLen:frameSize+stdWriterPrefixLen]))
		***REMOVED***

		// Write the retrieved frame (without header)
		nw, ew = out.Write(buf[stdWriterPrefixLen : frameSize+stdWriterPrefixLen])
		if ew != nil ***REMOVED***
			return 0, ew
		***REMOVED***

		// If the frame has not been fully written: error
		if nw != frameSize ***REMOVED***
			return 0, io.ErrShortWrite
		***REMOVED***
		written += int64(nw)

		// Move the rest of the buffer to the beginning
		copy(buf, buf[frameSize+stdWriterPrefixLen:])
		// Move the index
		nr -= frameSize + stdWriterPrefixLen
	***REMOVED***
***REMOVED***
