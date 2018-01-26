// Copyright (c) 2012 The Go Authors. All rights reserved.
// 
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
// 
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
// 
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package check

import (
	"fmt"
	"runtime"
	"time"
)

var memStats runtime.MemStats

// testingB is a type passed to Benchmark functions to manage benchmark
// timing and to specify the number of iterations to run.
type timer struct ***REMOVED***
	start     time.Time // Time test or benchmark started
	duration  time.Duration
	N         int
	bytes     int64
	timerOn   bool
	benchTime time.Duration
	// The initial states of memStats.Mallocs and memStats.TotalAlloc.
	startAllocs uint64
	startBytes  uint64
	// The net total of this test after being run.
	netAllocs uint64
	netBytes  uint64
***REMOVED***

// StartTimer starts timing a test. This function is called automatically
// before a benchmark starts, but it can also used to resume timing after
// a call to StopTimer.
func (c *C) StartTimer() ***REMOVED***
	if !c.timerOn ***REMOVED***
		c.start = time.Now()
		c.timerOn = true

		runtime.ReadMemStats(&memStats)
		c.startAllocs = memStats.Mallocs
		c.startBytes = memStats.TotalAlloc
	***REMOVED***
***REMOVED***

// StopTimer stops timing a test. This can be used to pause the timer
// while performing complex initialization that you don't
// want to measure.
func (c *C) StopTimer() ***REMOVED***
	if c.timerOn ***REMOVED***
		c.duration += time.Now().Sub(c.start)
		c.timerOn = false
		runtime.ReadMemStats(&memStats)
		c.netAllocs += memStats.Mallocs - c.startAllocs
		c.netBytes += memStats.TotalAlloc - c.startBytes
	***REMOVED***
***REMOVED***

// ResetTimer sets the elapsed benchmark time to zero.
// It does not affect whether the timer is running.
func (c *C) ResetTimer() ***REMOVED***
	if c.timerOn ***REMOVED***
		c.start = time.Now()
		runtime.ReadMemStats(&memStats)
		c.startAllocs = memStats.Mallocs
		c.startBytes = memStats.TotalAlloc
	***REMOVED***
	c.duration = 0
	c.netAllocs = 0
	c.netBytes = 0
***REMOVED***

// SetBytes informs the number of bytes that the benchmark processes
// on each iteration. If this is called in a benchmark it will also
// report MB/s.
func (c *C) SetBytes(n int64) ***REMOVED***
	c.bytes = n
***REMOVED***

func (c *C) nsPerOp() int64 ***REMOVED***
	if c.N <= 0 ***REMOVED***
		return 0
	***REMOVED***
	return c.duration.Nanoseconds() / int64(c.N)
***REMOVED***

func (c *C) mbPerSec() float64 ***REMOVED***
	if c.bytes <= 0 || c.duration <= 0 || c.N <= 0 ***REMOVED***
		return 0
	***REMOVED***
	return (float64(c.bytes) * float64(c.N) / 1e6) / c.duration.Seconds()
***REMOVED***

func (c *C) timerString() string ***REMOVED***
	if c.N <= 0 ***REMOVED***
		return fmt.Sprintf("%3.3fs", float64(c.duration.Nanoseconds())/1e9)
	***REMOVED***
	mbs := c.mbPerSec()
	mb := ""
	if mbs != 0 ***REMOVED***
		mb = fmt.Sprintf("\t%7.2f MB/s", mbs)
	***REMOVED***
	nsop := c.nsPerOp()
	ns := fmt.Sprintf("%10d ns/op", nsop)
	if c.N > 0 && nsop < 100 ***REMOVED***
		// The format specifiers here make sure that
		// the ones digits line up for all three possible formats.
		if nsop < 10 ***REMOVED***
			ns = fmt.Sprintf("%13.2f ns/op", float64(c.duration.Nanoseconds())/float64(c.N))
		***REMOVED*** else ***REMOVED***
			ns = fmt.Sprintf("%12.1f ns/op", float64(c.duration.Nanoseconds())/float64(c.N))
		***REMOVED***
	***REMOVED***
	memStats := ""
	if c.benchMem ***REMOVED***
		allocedBytes := fmt.Sprintf("%8d B/op", int64(c.netBytes)/int64(c.N))
		allocs := fmt.Sprintf("%8d allocs/op", int64(c.netAllocs)/int64(c.N))
		memStats = fmt.Sprintf("\t%s\t%s", allocedBytes, allocs)
	***REMOVED***
	return fmt.Sprintf("%8d\t%s%s%s", c.N, ns, mb, memStats)
***REMOVED***

func min(x, y int) int ***REMOVED***
	if x > y ***REMOVED***
		return y
	***REMOVED***
	return x
***REMOVED***

func max(x, y int) int ***REMOVED***
	if x < y ***REMOVED***
		return y
	***REMOVED***
	return x
***REMOVED***

// roundDown10 rounds a number down to the nearest power of 10.
func roundDown10(n int) int ***REMOVED***
	var tens = 0
	// tens = floor(log_10(n))
	for n > 10 ***REMOVED***
		n = n / 10
		tens++
	***REMOVED***
	// result = 10^tens
	result := 1
	for i := 0; i < tens; i++ ***REMOVED***
		result *= 10
	***REMOVED***
	return result
***REMOVED***

// roundUp rounds x up to a number of the form [1eX, 2eX, 5eX].
func roundUp(n int) int ***REMOVED***
	base := roundDown10(n)
	if n < (2 * base) ***REMOVED***
		return 2 * base
	***REMOVED***
	if n < (5 * base) ***REMOVED***
		return 5 * base
	***REMOVED***
	return 10 * base
***REMOVED***
