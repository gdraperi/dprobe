package progress

import (
	"io"
	"time"

	"golang.org/x/time/rate"
)

// Reader is a Reader with progress bar.
type Reader struct ***REMOVED***
	in          io.ReadCloser // Stream to read from
	out         Output        // Where to send progress bar to
	size        int64
	current     int64
	lastUpdate  int64
	id          string
	action      string
	rateLimiter *rate.Limiter
***REMOVED***

// NewProgressReader creates a new ProgressReader.
func NewProgressReader(in io.ReadCloser, out Output, size int64, id, action string) *Reader ***REMOVED***
	return &Reader***REMOVED***
		in:          in,
		out:         out,
		size:        size,
		id:          id,
		action:      action,
		rateLimiter: rate.NewLimiter(rate.Every(100*time.Millisecond), 1),
	***REMOVED***
***REMOVED***

func (p *Reader) Read(buf []byte) (n int, err error) ***REMOVED***
	read, err := p.in.Read(buf)
	p.current += int64(read)
	updateEvery := int64(1024 * 512) //512kB
	if p.size > 0 ***REMOVED***
		// Update progress for every 1% read if 1% < 512kB
		if increment := int64(0.01 * float64(p.size)); increment < updateEvery ***REMOVED***
			updateEvery = increment
		***REMOVED***
	***REMOVED***
	if p.current-p.lastUpdate > updateEvery || err != nil ***REMOVED***
		p.updateProgress(err != nil && read == 0)
		p.lastUpdate = p.current
	***REMOVED***

	return read, err
***REMOVED***

// Close closes the progress reader and its underlying reader.
func (p *Reader) Close() error ***REMOVED***
	if p.current < p.size ***REMOVED***
		// print a full progress bar when closing prematurely
		p.current = p.size
		p.updateProgress(false)
	***REMOVED***
	return p.in.Close()
***REMOVED***

func (p *Reader) updateProgress(last bool) ***REMOVED***
	if last || p.current == p.size || p.rateLimiter.Allow() ***REMOVED***
		p.out.WriteProgress(Progress***REMOVED***ID: p.id, Action: p.action, Current: p.current, Total: p.size, LastUpdate: last***REMOVED***)
	***REMOVED***
***REMOVED***
