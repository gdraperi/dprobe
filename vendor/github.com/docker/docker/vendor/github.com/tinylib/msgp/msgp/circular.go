package msgp

type timer interface ***REMOVED***
	StartTimer()
	StopTimer()
***REMOVED***

// EndlessReader is an io.Reader
// that loops over the same data
// endlessly. It is used for benchmarking.
type EndlessReader struct ***REMOVED***
	tb     timer
	data   []byte
	offset int
***REMOVED***

// NewEndlessReader returns a new endless reader
func NewEndlessReader(b []byte, tb timer) *EndlessReader ***REMOVED***
	return &EndlessReader***REMOVED***tb: tb, data: b, offset: 0***REMOVED***
***REMOVED***

// Read implements io.Reader. In practice, it
// always returns (len(p), nil), although it
// fills the supplied slice while the benchmark
// timer is stopped.
func (c *EndlessReader) Read(p []byte) (int, error) ***REMOVED***
	c.tb.StopTimer()
	var n int
	l := len(p)
	m := len(c.data)
	for n < l ***REMOVED***
		nn := copy(p[n:], c.data[c.offset:])
		n += nn
		c.offset += nn
		c.offset %= m
	***REMOVED***
	c.tb.StartTimer()
	return n, nil
***REMOVED***
