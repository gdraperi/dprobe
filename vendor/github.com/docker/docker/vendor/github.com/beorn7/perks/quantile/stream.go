// Package quantile computes approximate quantiles over an unbounded data
// stream within low memory and CPU bounds.
//
// A small amount of accuracy is traded to achieve the above properties.
//
// Multiple streams can be merged before calling Query to generate a single set
// of results. This is meaningful when the streams represent the same type of
// data. See Merge and Samples.
//
// For more detailed information about the algorithm used, see:
//
// Effective Computation of Biased Quantiles over Data Streams
//
// http://www.cs.rutgers.edu/~muthu/bquant.pdf
package quantile

import (
	"math"
	"sort"
)

// Sample holds an observed value and meta information for compression. JSON
// tags have been added for convenience.
type Sample struct ***REMOVED***
	Value float64 `json:",string"`
	Width float64 `json:",string"`
	Delta float64 `json:",string"`
***REMOVED***

// Samples represents a slice of samples. It implements sort.Interface.
type Samples []Sample

func (a Samples) Len() int           ***REMOVED*** return len(a) ***REMOVED***
func (a Samples) Less(i, j int) bool ***REMOVED*** return a[i].Value < a[j].Value ***REMOVED***
func (a Samples) Swap(i, j int)      ***REMOVED*** a[i], a[j] = a[j], a[i] ***REMOVED***

type invariant func(s *stream, r float64) float64

// NewLowBiased returns an initialized Stream for low-biased quantiles
// (e.g. 0.01, 0.1, 0.5) where the needed quantiles are not known a priori, but
// error guarantees can still be given even for the lower ranks of the data
// distribution.
//
// The provided epsilon is a relative error, i.e. the true quantile of a value
// returned by a query is guaranteed to be within (1±Epsilon)*Quantile.
//
// See http://www.cs.rutgers.edu/~muthu/bquant.pdf for time, space, and error
// properties.
func NewLowBiased(epsilon float64) *Stream ***REMOVED***
	ƒ := func(s *stream, r float64) float64 ***REMOVED***
		return 2 * epsilon * r
	***REMOVED***
	return newStream(ƒ)
***REMOVED***

// NewHighBiased returns an initialized Stream for high-biased quantiles
// (e.g. 0.01, 0.1, 0.5) where the needed quantiles are not known a priori, but
// error guarantees can still be given even for the higher ranks of the data
// distribution.
//
// The provided epsilon is a relative error, i.e. the true quantile of a value
// returned by a query is guaranteed to be within 1-(1±Epsilon)*(1-Quantile).
//
// See http://www.cs.rutgers.edu/~muthu/bquant.pdf for time, space, and error
// properties.
func NewHighBiased(epsilon float64) *Stream ***REMOVED***
	ƒ := func(s *stream, r float64) float64 ***REMOVED***
		return 2 * epsilon * (s.n - r)
	***REMOVED***
	return newStream(ƒ)
***REMOVED***

// NewTargeted returns an initialized Stream concerned with a particular set of
// quantile values that are supplied a priori. Knowing these a priori reduces
// space and computation time. The targets map maps the desired quantiles to
// their absolute errors, i.e. the true quantile of a value returned by a query
// is guaranteed to be within (Quantile±Epsilon).
//
// See http://www.cs.rutgers.edu/~muthu/bquant.pdf for time, space, and error properties.
func NewTargeted(targets map[float64]float64) *Stream ***REMOVED***
	ƒ := func(s *stream, r float64) float64 ***REMOVED***
		var m = math.MaxFloat64
		var f float64
		for quantile, epsilon := range targets ***REMOVED***
			if quantile*s.n <= r ***REMOVED***
				f = (2 * epsilon * r) / quantile
			***REMOVED*** else ***REMOVED***
				f = (2 * epsilon * (s.n - r)) / (1 - quantile)
			***REMOVED***
			if f < m ***REMOVED***
				m = f
			***REMOVED***
		***REMOVED***
		return m
	***REMOVED***
	return newStream(ƒ)
***REMOVED***

// Stream computes quantiles for a stream of float64s. It is not thread-safe by
// design. Take care when using across multiple goroutines.
type Stream struct ***REMOVED***
	*stream
	b      Samples
	sorted bool
***REMOVED***

func newStream(ƒ invariant) *Stream ***REMOVED***
	x := &stream***REMOVED***ƒ: ƒ***REMOVED***
	return &Stream***REMOVED***x, make(Samples, 0, 500), true***REMOVED***
***REMOVED***

// Insert inserts v into the stream.
func (s *Stream) Insert(v float64) ***REMOVED***
	s.insert(Sample***REMOVED***Value: v, Width: 1***REMOVED***)
***REMOVED***

func (s *Stream) insert(sample Sample) ***REMOVED***
	s.b = append(s.b, sample)
	s.sorted = false
	if len(s.b) == cap(s.b) ***REMOVED***
		s.flush()
	***REMOVED***
***REMOVED***

// Query returns the computed qth percentiles value. If s was created with
// NewTargeted, and q is not in the set of quantiles provided a priori, Query
// will return an unspecified result.
func (s *Stream) Query(q float64) float64 ***REMOVED***
	if !s.flushed() ***REMOVED***
		// Fast path when there hasn't been enough data for a flush;
		// this also yields better accuracy for small sets of data.
		l := len(s.b)
		if l == 0 ***REMOVED***
			return 0
		***REMOVED***
		i := int(math.Ceil(float64(l) * q))
		if i > 0 ***REMOVED***
			i -= 1
		***REMOVED***
		s.maybeSort()
		return s.b[i].Value
	***REMOVED***
	s.flush()
	return s.stream.query(q)
***REMOVED***

// Merge merges samples into the underlying streams samples. This is handy when
// merging multiple streams from separate threads, database shards, etc.
//
// ATTENTION: This method is broken and does not yield correct results. The
// underlying algorithm is not capable of merging streams correctly.
func (s *Stream) Merge(samples Samples) ***REMOVED***
	sort.Sort(samples)
	s.stream.merge(samples)
***REMOVED***

// Reset reinitializes and clears the list reusing the samples buffer memory.
func (s *Stream) Reset() ***REMOVED***
	s.stream.reset()
	s.b = s.b[:0]
***REMOVED***

// Samples returns stream samples held by s.
func (s *Stream) Samples() Samples ***REMOVED***
	if !s.flushed() ***REMOVED***
		return s.b
	***REMOVED***
	s.flush()
	return s.stream.samples()
***REMOVED***

// Count returns the total number of samples observed in the stream
// since initialization.
func (s *Stream) Count() int ***REMOVED***
	return len(s.b) + s.stream.count()
***REMOVED***

func (s *Stream) flush() ***REMOVED***
	s.maybeSort()
	s.stream.merge(s.b)
	s.b = s.b[:0]
***REMOVED***

func (s *Stream) maybeSort() ***REMOVED***
	if !s.sorted ***REMOVED***
		s.sorted = true
		sort.Sort(s.b)
	***REMOVED***
***REMOVED***

func (s *Stream) flushed() bool ***REMOVED***
	return len(s.stream.l) > 0
***REMOVED***

type stream struct ***REMOVED***
	n float64
	l []Sample
	ƒ invariant
***REMOVED***

func (s *stream) reset() ***REMOVED***
	s.l = s.l[:0]
	s.n = 0
***REMOVED***

func (s *stream) insert(v float64) ***REMOVED***
	s.merge(Samples***REMOVED******REMOVED***v, 1, 0***REMOVED******REMOVED***)
***REMOVED***

func (s *stream) merge(samples Samples) ***REMOVED***
	// TODO(beorn7): This tries to merge not only individual samples, but
	// whole summaries. The paper doesn't mention merging summaries at
	// all. Unittests show that the merging is inaccurate. Find out how to
	// do merges properly.
	var r float64
	i := 0
	for _, sample := range samples ***REMOVED***
		for ; i < len(s.l); i++ ***REMOVED***
			c := s.l[i]
			if c.Value > sample.Value ***REMOVED***
				// Insert at position i.
				s.l = append(s.l, Sample***REMOVED******REMOVED***)
				copy(s.l[i+1:], s.l[i:])
				s.l[i] = Sample***REMOVED***
					sample.Value,
					sample.Width,
					math.Max(sample.Delta, math.Floor(s.ƒ(s, r))-1),
					// TODO(beorn7): How to calculate delta correctly?
				***REMOVED***
				i++
				goto inserted
			***REMOVED***
			r += c.Width
		***REMOVED***
		s.l = append(s.l, Sample***REMOVED***sample.Value, sample.Width, 0***REMOVED***)
		i++
	inserted:
		s.n += sample.Width
		r += sample.Width
	***REMOVED***
	s.compress()
***REMOVED***

func (s *stream) count() int ***REMOVED***
	return int(s.n)
***REMOVED***

func (s *stream) query(q float64) float64 ***REMOVED***
	t := math.Ceil(q * s.n)
	t += math.Ceil(s.ƒ(s, t) / 2)
	p := s.l[0]
	var r float64
	for _, c := range s.l[1:] ***REMOVED***
		r += p.Width
		if r+c.Width+c.Delta > t ***REMOVED***
			return p.Value
		***REMOVED***
		p = c
	***REMOVED***
	return p.Value
***REMOVED***

func (s *stream) compress() ***REMOVED***
	if len(s.l) < 2 ***REMOVED***
		return
	***REMOVED***
	x := s.l[len(s.l)-1]
	xi := len(s.l) - 1
	r := s.n - 1 - x.Width

	for i := len(s.l) - 2; i >= 0; i-- ***REMOVED***
		c := s.l[i]
		if c.Width+x.Width+x.Delta <= s.ƒ(s, r) ***REMOVED***
			x.Width += c.Width
			s.l[xi] = x
			// Remove element at i.
			copy(s.l[i:], s.l[i+1:])
			s.l = s.l[:len(s.l)-1]
			xi -= 1
		***REMOVED*** else ***REMOVED***
			x = c
			xi = i
		***REMOVED***
		r -= c.Width
	***REMOVED***
***REMOVED***

func (s *stream) samples() Samples ***REMOVED***
	samples := make(Samples, len(s.l))
	copy(samples, s.l)
	return samples
***REMOVED***
