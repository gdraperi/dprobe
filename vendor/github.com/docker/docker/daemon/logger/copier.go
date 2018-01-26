package logger

import (
	"bytes"
	"io"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	// readSize is the maximum bytes read during a single read
	// operation.
	readSize = 2 * 1024

	// defaultBufSize provides a reasonable default for loggers that do
	// not have an external limit to impose on log line size.
	defaultBufSize = 16 * 1024
)

// Copier can copy logs from specified sources to Logger and attach Timestamp.
// Writes are concurrent, so you need implement some sync in your logger.
type Copier struct ***REMOVED***
	// srcs is map of name -> reader pairs, for example "stdout", "stderr"
	srcs      map[string]io.Reader
	dst       Logger
	copyJobs  sync.WaitGroup
	closeOnce sync.Once
	closed    chan struct***REMOVED******REMOVED***
***REMOVED***

// NewCopier creates a new Copier
func NewCopier(srcs map[string]io.Reader, dst Logger) *Copier ***REMOVED***
	return &Copier***REMOVED***
		srcs:   srcs,
		dst:    dst,
		closed: make(chan struct***REMOVED******REMOVED***),
	***REMOVED***
***REMOVED***

// Run starts logs copying
func (c *Copier) Run() ***REMOVED***
	for src, w := range c.srcs ***REMOVED***
		c.copyJobs.Add(1)
		go c.copySrc(src, w)
	***REMOVED***
***REMOVED***

func (c *Copier) copySrc(name string, src io.Reader) ***REMOVED***
	defer c.copyJobs.Done()

	bufSize := defaultBufSize
	if sizedLogger, ok := c.dst.(SizedLogger); ok ***REMOVED***
		bufSize = sizedLogger.BufSize()
	***REMOVED***
	buf := make([]byte, bufSize)

	n := 0
	eof := false

	for ***REMOVED***
		select ***REMOVED***
		case <-c.closed:
			return
		default:
			// Work out how much more data we are okay with reading this time.
			upto := n + readSize
			if upto > cap(buf) ***REMOVED***
				upto = cap(buf)
			***REMOVED***
			// Try to read that data.
			if upto > n ***REMOVED***
				read, err := src.Read(buf[n:upto])
				if err != nil ***REMOVED***
					if err != io.EOF ***REMOVED***
						logrus.Errorf("Error scanning log stream: %s", err)
						return
					***REMOVED***
					eof = true
				***REMOVED***
				n += read
			***REMOVED***
			// If we have no data to log, and there's no more coming, we're done.
			if n == 0 && eof ***REMOVED***
				return
			***REMOVED***
			// Break up the data that we've buffered up into lines, and log each in turn.
			p := 0
			for q := bytes.IndexByte(buf[p:n], '\n'); q >= 0; q = bytes.IndexByte(buf[p:n], '\n') ***REMOVED***
				select ***REMOVED***
				case <-c.closed:
					return
				default:
					msg := NewMessage()
					msg.Source = name
					msg.Timestamp = time.Now().UTC()
					msg.Line = append(msg.Line, buf[p:p+q]...)

					if logErr := c.dst.Log(msg); logErr != nil ***REMOVED***
						logrus.Errorf("Failed to log msg %q for logger %s: %s", msg.Line, c.dst.Name(), logErr)
					***REMOVED***
				***REMOVED***
				p += q + 1
			***REMOVED***
			// If there's no more coming, or the buffer is full but
			// has no newlines, log whatever we haven't logged yet,
			// noting that it's a partial log line.
			if eof || (p == 0 && n == len(buf)) ***REMOVED***
				if p < n ***REMOVED***
					msg := NewMessage()
					msg.Source = name
					msg.Timestamp = time.Now().UTC()
					msg.Line = append(msg.Line, buf[p:n]...)
					msg.Partial = true

					if logErr := c.dst.Log(msg); logErr != nil ***REMOVED***
						logrus.Errorf("Failed to log msg %q for logger %s: %s", msg.Line, c.dst.Name(), logErr)
					***REMOVED***
					p = 0
					n = 0
				***REMOVED***
				if eof ***REMOVED***
					return
				***REMOVED***
			***REMOVED***
			// Move any unlogged data to the front of the buffer in preparation for another read.
			if p > 0 ***REMOVED***
				copy(buf[0:], buf[p:n])
				n -= p
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// Wait waits until all copying is done
func (c *Copier) Wait() ***REMOVED***
	c.copyJobs.Wait()
***REMOVED***

// Close closes the copier
func (c *Copier) Close() ***REMOVED***
	c.closeOnce.Do(func() ***REMOVED***
		close(c.closed)
	***REMOVED***)
***REMOVED***
