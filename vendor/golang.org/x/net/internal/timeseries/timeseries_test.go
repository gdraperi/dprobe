// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package timeseries

import (
	"math"
	"testing"
	"time"
)

func isNear(x *Float, y float64, tolerance float64) bool ***REMOVED***
	return math.Abs(x.Value()-y) < tolerance
***REMOVED***

func isApproximate(x *Float, y float64) bool ***REMOVED***
	return isNear(x, y, 1e-2)
***REMOVED***

func checkApproximate(t *testing.T, o Observable, y float64) ***REMOVED***
	x := o.(*Float)
	if !isApproximate(x, y) ***REMOVED***
		t.Errorf("Wanted %g, got %g", y, x.Value())
	***REMOVED***
***REMOVED***

func checkNear(t *testing.T, o Observable, y, tolerance float64) ***REMOVED***
	x := o.(*Float)
	if !isNear(x, y, tolerance) ***REMOVED***
		t.Errorf("Wanted %g +- %g, got %g", y, tolerance, x.Value())
	***REMOVED***
***REMOVED***

var baseTime = time.Date(2013, 1, 1, 0, 0, 0, 0, time.UTC)

func tu(s int64) time.Time ***REMOVED***
	return baseTime.Add(time.Duration(s) * time.Second)
***REMOVED***

func tu2(s int64, ns int64) time.Time ***REMOVED***
	return baseTime.Add(time.Duration(s)*time.Second + time.Duration(ns)*time.Nanosecond)
***REMOVED***

func TestBasicTimeSeries(t *testing.T) ***REMOVED***
	ts := NewTimeSeries(NewFloat)
	fo := new(Float)
	*fo = Float(10)
	ts.AddWithTime(fo, tu(1))
	ts.AddWithTime(fo, tu(1))
	ts.AddWithTime(fo, tu(1))
	ts.AddWithTime(fo, tu(1))
	checkApproximate(t, ts.Range(tu(0), tu(1)), 40)
	checkApproximate(t, ts.Total(), 40)
	ts.AddWithTime(fo, tu(3))
	ts.AddWithTime(fo, tu(3))
	ts.AddWithTime(fo, tu(3))
	checkApproximate(t, ts.Range(tu(0), tu(2)), 40)
	checkApproximate(t, ts.Range(tu(2), tu(4)), 30)
	checkApproximate(t, ts.Total(), 70)
	ts.AddWithTime(fo, tu(1))
	ts.AddWithTime(fo, tu(1))
	checkApproximate(t, ts.Range(tu(0), tu(2)), 60)
	checkApproximate(t, ts.Range(tu(2), tu(4)), 30)
	checkApproximate(t, ts.Total(), 90)
	*fo = Float(100)
	ts.AddWithTime(fo, tu(100))
	checkApproximate(t, ts.Range(tu(99), tu(100)), 100)
	checkApproximate(t, ts.Range(tu(0), tu(4)), 36)
	checkApproximate(t, ts.Total(), 190)
	*fo = Float(10)
	ts.AddWithTime(fo, tu(1))
	ts.AddWithTime(fo, tu(1))
	checkApproximate(t, ts.Range(tu(0), tu(4)), 44)
	checkApproximate(t, ts.Range(tu(37), tu2(100, 100e6)), 100)
	checkApproximate(t, ts.Range(tu(50), tu2(100, 100e6)), 100)
	checkApproximate(t, ts.Range(tu(99), tu2(100, 100e6)), 100)
	checkApproximate(t, ts.Total(), 210)

	for i, l := range ts.ComputeRange(tu(36), tu(100), 64) ***REMOVED***
		if i == 63 ***REMOVED***
			checkApproximate(t, l, 100)
		***REMOVED*** else ***REMOVED***
			checkApproximate(t, l, 0)
		***REMOVED***
	***REMOVED***

	checkApproximate(t, ts.Range(tu(0), tu(100)), 210)
	checkApproximate(t, ts.Range(tu(10), tu(100)), 100)

	for i, l := range ts.ComputeRange(tu(0), tu(100), 100) ***REMOVED***
		if i < 10 ***REMOVED***
			checkApproximate(t, l, 11)
		***REMOVED*** else if i >= 90 ***REMOVED***
			checkApproximate(t, l, 10)
		***REMOVED*** else ***REMOVED***
			checkApproximate(t, l, 0)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestFloat(t *testing.T) ***REMOVED***
	f := Float(1)
	if g, w := f.String(), "1"; g != w ***REMOVED***
		t.Errorf("Float(1).String = %q; want %q", g, w)
	***REMOVED***
	f2 := Float(2)
	var o Observable = &f2
	f.Add(o)
	if g, w := f.Value(), 3.0; g != w ***REMOVED***
		t.Errorf("Float post-add = %v; want %v", g, w)
	***REMOVED***
	f.Multiply(2)
	if g, w := f.Value(), 6.0; g != w ***REMOVED***
		t.Errorf("Float post-multiply = %v; want %v", g, w)
	***REMOVED***
	f.Clear()
	if g, w := f.Value(), 0.0; g != w ***REMOVED***
		t.Errorf("Float post-clear = %v; want %v", g, w)
	***REMOVED***
	f.CopyFrom(&f2)
	if g, w := f.Value(), 2.0; g != w ***REMOVED***
		t.Errorf("Float post-CopyFrom = %v; want %v", g, w)
	***REMOVED***
***REMOVED***

type mockClock struct ***REMOVED***
	time time.Time
***REMOVED***

func (m *mockClock) Time() time.Time ***REMOVED*** return m.time ***REMOVED***
func (m *mockClock) Set(t time.Time) ***REMOVED*** m.time = t ***REMOVED***

const buckets = 6

var testResolutions = []time.Duration***REMOVED***
	10 * time.Second,  // level holds one minute of observations
	100 * time.Second, // level holds ten minutes of observations
	10 * time.Minute,  // level holds one hour of observations
***REMOVED***

// TestTimeSeries uses a small number of buckets to force a higher
// error rate on approximations from the timeseries.
type TestTimeSeries struct ***REMOVED***
	timeSeries
***REMOVED***

func TestExpectedErrorRate(t *testing.T) ***REMOVED***
	ts := new(TestTimeSeries)
	fake := new(mockClock)
	fake.Set(time.Now())
	ts.timeSeries.init(testResolutions, NewFloat, buckets, fake)
	for i := 1; i <= 61*61; i++ ***REMOVED***
		fake.Set(fake.Time().Add(1 * time.Second))
		ob := Float(1)
		ts.AddWithTime(&ob, fake.Time())

		// The results should be accurate within one missing bucket (1/6) of the observations recorded.
		checkNear(t, ts.Latest(0, buckets), min(float64(i), 60), 10)
		checkNear(t, ts.Latest(1, buckets), min(float64(i), 600), 100)
		checkNear(t, ts.Latest(2, buckets), min(float64(i), 3600), 600)
	***REMOVED***
***REMOVED***

func min(a, b float64) float64 ***REMOVED***
	if a < b ***REMOVED***
		return a
	***REMOVED***
	return b
***REMOVED***
