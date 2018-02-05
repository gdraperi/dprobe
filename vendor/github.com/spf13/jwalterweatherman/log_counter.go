// Copyright © 2016 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package jwalterweatherman

import (
	"sync/atomic"
)

type logCounter struct ***REMOVED***
	counter uint64
***REMOVED***

func (c *logCounter) incr() ***REMOVED***
	atomic.AddUint64(&c.counter, 1)
***REMOVED***

func (c *logCounter) resetCounter() ***REMOVED***
	atomic.StoreUint64(&c.counter, 0)
***REMOVED***

func (c *logCounter) getCount() uint64 ***REMOVED***
	return atomic.LoadUint64(&c.counter)
***REMOVED***

func (c *logCounter) Write(p []byte) (n int, err error) ***REMOVED***
	c.incr()
	return len(p), nil
***REMOVED***

// LogCountForLevel returns the number of log invocations for a given threshold.
func (n *Notepad) LogCountForLevel(l Threshold) uint64 ***REMOVED***
	return n.logCounters[l].getCount()
***REMOVED***

// LogCountForLevelsGreaterThanorEqualTo returns the number of log invocations
// greater than or equal to a given threshold.
func (n *Notepad) LogCountForLevelsGreaterThanorEqualTo(threshold Threshold) uint64 ***REMOVED***
	var cnt uint64

	for i := int(threshold); i < len(n.logCounters); i++ ***REMOVED***
		cnt += n.LogCountForLevel(Threshold(i))
	***REMOVED***

	return cnt
***REMOVED***

// ResetLogCounters resets the invocation counters for all levels.
func (n *Notepad) ResetLogCounters() ***REMOVED***
	for _, np := range n.logCounters ***REMOVED***
		np.resetCounter()
	***REMOVED***
***REMOVED***
