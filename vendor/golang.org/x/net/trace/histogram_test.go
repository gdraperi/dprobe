// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trace

import (
	"math"
	"testing"
)

type sumTest struct ***REMOVED***
	value        int64
	sum          int64
	sumOfSquares float64
	total        int64
***REMOVED***

var sumTests = []sumTest***REMOVED***
	***REMOVED***100, 100, 10000, 1***REMOVED***,
	***REMOVED***50, 150, 12500, 2***REMOVED***,
	***REMOVED***50, 200, 15000, 3***REMOVED***,
	***REMOVED***50, 250, 17500, 4***REMOVED***,
***REMOVED***

type bucketingTest struct ***REMOVED***
	in     int64
	log    int
	bucket int
***REMOVED***

var bucketingTests = []bucketingTest***REMOVED***
	***REMOVED***0, 0, 0***REMOVED***,
	***REMOVED***1, 1, 0***REMOVED***,
	***REMOVED***2, 2, 1***REMOVED***,
	***REMOVED***3, 2, 1***REMOVED***,
	***REMOVED***4, 3, 2***REMOVED***,
	***REMOVED***1000, 10, 9***REMOVED***,
	***REMOVED***1023, 10, 9***REMOVED***,
	***REMOVED***1024, 11, 10***REMOVED***,
	***REMOVED***1000000, 20, 19***REMOVED***,
***REMOVED***

type multiplyTest struct ***REMOVED***
	in                   int64
	ratio                float64
	expectedSum          int64
	expectedTotal        int64
	expectedSumOfSquares float64
***REMOVED***

var multiplyTests = []multiplyTest***REMOVED***
	***REMOVED***15, 2.5, 37, 2, 562.5***REMOVED***,
	***REMOVED***128, 4.6, 758, 13, 77953.9***REMOVED***,
***REMOVED***

type percentileTest struct ***REMOVED***
	fraction float64
	expected int64
***REMOVED***

var percentileTests = []percentileTest***REMOVED***
	***REMOVED***0.25, 48***REMOVED***,
	***REMOVED***0.5, 96***REMOVED***,
	***REMOVED***0.6, 109***REMOVED***,
	***REMOVED***0.75, 128***REMOVED***,
	***REMOVED***0.90, 205***REMOVED***,
	***REMOVED***0.95, 230***REMOVED***,
	***REMOVED***0.99, 256***REMOVED***,
***REMOVED***

func TestSum(t *testing.T) ***REMOVED***
	var h histogram

	for _, test := range sumTests ***REMOVED***
		h.addMeasurement(test.value)
		sum := h.sum
		if sum != test.sum ***REMOVED***
			t.Errorf("h.Sum = %v WANT: %v", sum, test.sum)
		***REMOVED***

		sumOfSquares := h.sumOfSquares
		if sumOfSquares != test.sumOfSquares ***REMOVED***
			t.Errorf("h.SumOfSquares = %v WANT: %v", sumOfSquares, test.sumOfSquares)
		***REMOVED***

		total := h.total()
		if total != test.total ***REMOVED***
			t.Errorf("h.Total = %v WANT: %v", total, test.total)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestMultiply(t *testing.T) ***REMOVED***
	var h histogram
	for i, test := range multiplyTests ***REMOVED***
		h.addMeasurement(test.in)
		h.Multiply(test.ratio)
		if h.sum != test.expectedSum ***REMOVED***
			t.Errorf("#%v: h.sum = %v WANT: %v", i, h.sum, test.expectedSum)
		***REMOVED***
		if h.total() != test.expectedTotal ***REMOVED***
			t.Errorf("#%v: h.total = %v WANT: %v", i, h.total(), test.expectedTotal)
		***REMOVED***
		if h.sumOfSquares != test.expectedSumOfSquares ***REMOVED***
			t.Errorf("#%v: h.SumOfSquares = %v WANT: %v", i, test.expectedSumOfSquares, h.sumOfSquares)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestBucketingFunctions(t *testing.T) ***REMOVED***
	for _, test := range bucketingTests ***REMOVED***
		log := log2(test.in)
		if log != test.log ***REMOVED***
			t.Errorf("log2 = %v WANT: %v", log, test.log)
		***REMOVED***

		bucket := getBucket(test.in)
		if bucket != test.bucket ***REMOVED***
			t.Errorf("getBucket = %v WANT: %v", bucket, test.bucket)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestAverage(t *testing.T) ***REMOVED***
	a := new(histogram)
	average := a.average()
	if average != 0 ***REMOVED***
		t.Errorf("Average of empty histogram was %v WANT: 0", average)
	***REMOVED***

	a.addMeasurement(1)
	a.addMeasurement(1)
	a.addMeasurement(3)
	const expected = float64(5) / float64(3)
	average = a.average()

	if !isApproximate(average, expected) ***REMOVED***
		t.Errorf("Average = %g WANT: %v", average, expected)
	***REMOVED***
***REMOVED***

func TestStandardDeviation(t *testing.T) ***REMOVED***
	a := new(histogram)
	add(a, 10, 1<<4)
	add(a, 10, 1<<5)
	add(a, 10, 1<<6)
	stdDev := a.standardDeviation()
	const expected = 19.95

	if !isApproximate(stdDev, expected) ***REMOVED***
		t.Errorf("StandardDeviation = %v WANT: %v", stdDev, expected)
	***REMOVED***

	// No values
	a = new(histogram)
	stdDev = a.standardDeviation()

	if !isApproximate(stdDev, 0) ***REMOVED***
		t.Errorf("StandardDeviation = %v WANT: 0", stdDev)
	***REMOVED***

	add(a, 1, 1<<4)
	if !isApproximate(stdDev, 0) ***REMOVED***
		t.Errorf("StandardDeviation = %v WANT: 0", stdDev)
	***REMOVED***

	add(a, 10, 1<<4)
	if !isApproximate(stdDev, 0) ***REMOVED***
		t.Errorf("StandardDeviation = %v WANT: 0", stdDev)
	***REMOVED***
***REMOVED***

func TestPercentileBoundary(t *testing.T) ***REMOVED***
	a := new(histogram)
	add(a, 5, 1<<4)
	add(a, 10, 1<<6)
	add(a, 5, 1<<7)

	for _, test := range percentileTests ***REMOVED***
		percentile := a.percentileBoundary(test.fraction)
		if percentile != test.expected ***REMOVED***
			t.Errorf("h.PercentileBoundary (fraction=%v) = %v WANT: %v", test.fraction, percentile, test.expected)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestCopyFrom(t *testing.T) ***REMOVED***
	a := histogram***REMOVED***5, 25, []int64***REMOVED***1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18,
		19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38***REMOVED***, 4, -1***REMOVED***
	b := histogram***REMOVED***6, 36, []int64***REMOVED***2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
		20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39***REMOVED***, 5, -1***REMOVED***

	a.CopyFrom(&b)

	if a.String() != b.String() ***REMOVED***
		t.Errorf("a.String = %s WANT: %s", a.String(), b.String())
	***REMOVED***
***REMOVED***

func TestClear(t *testing.T) ***REMOVED***
	a := histogram***REMOVED***5, 25, []int64***REMOVED***1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18,
		19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38***REMOVED***, 4, -1***REMOVED***

	a.Clear()

	expected := "0, 0.000000, 0, 0, []"
	if a.String() != expected ***REMOVED***
		t.Errorf("a.String = %s WANT %s", a.String(), expected)
	***REMOVED***
***REMOVED***

func TestNew(t *testing.T) ***REMOVED***
	a := histogram***REMOVED***5, 25, []int64***REMOVED***1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18,
		19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38***REMOVED***, 4, -1***REMOVED***
	b := a.New()

	expected := "0, 0.000000, 0, 0, []"
	if b.(*histogram).String() != expected ***REMOVED***
		t.Errorf("b.(*histogram).String = %s WANT: %s", b.(*histogram).String(), expected)
	***REMOVED***
***REMOVED***

func TestAdd(t *testing.T) ***REMOVED***
	// The tests here depend on the associativity of addMeasurement and Add.
	// Add empty observation
	a := histogram***REMOVED***5, 25, []int64***REMOVED***1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18,
		19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38***REMOVED***, 4, -1***REMOVED***
	b := a.New()

	expected := a.String()
	a.Add(b)
	if a.String() != expected ***REMOVED***
		t.Errorf("a.String = %s WANT: %s", a.String(), expected)
	***REMOVED***

	// Add same bucketed value, no new buckets
	c := new(histogram)
	d := new(histogram)
	e := new(histogram)
	c.addMeasurement(12)
	d.addMeasurement(11)
	e.addMeasurement(12)
	e.addMeasurement(11)
	c.Add(d)
	if c.String() != e.String() ***REMOVED***
		t.Errorf("c.String = %s WANT: %s", c.String(), e.String())
	***REMOVED***

	// Add bucketed values
	f := new(histogram)
	g := new(histogram)
	h := new(histogram)
	f.addMeasurement(4)
	f.addMeasurement(12)
	f.addMeasurement(100)
	g.addMeasurement(18)
	g.addMeasurement(36)
	g.addMeasurement(255)
	h.addMeasurement(4)
	h.addMeasurement(12)
	h.addMeasurement(100)
	h.addMeasurement(18)
	h.addMeasurement(36)
	h.addMeasurement(255)
	f.Add(g)
	if f.String() != h.String() ***REMOVED***
		t.Errorf("f.String = %q WANT: %q", f.String(), h.String())
	***REMOVED***

	// add buckets to no buckets
	i := new(histogram)
	j := new(histogram)
	k := new(histogram)
	j.addMeasurement(18)
	j.addMeasurement(36)
	j.addMeasurement(255)
	k.addMeasurement(18)
	k.addMeasurement(36)
	k.addMeasurement(255)
	i.Add(j)
	if i.String() != k.String() ***REMOVED***
		t.Errorf("i.String = %q WANT: %q", i.String(), k.String())
	***REMOVED***

	// add buckets to single value (no overlap)
	l := new(histogram)
	m := new(histogram)
	n := new(histogram)
	l.addMeasurement(0)
	m.addMeasurement(18)
	m.addMeasurement(36)
	m.addMeasurement(255)
	n.addMeasurement(0)
	n.addMeasurement(18)
	n.addMeasurement(36)
	n.addMeasurement(255)
	l.Add(m)
	if l.String() != n.String() ***REMOVED***
		t.Errorf("l.String = %q WANT: %q", l.String(), n.String())
	***REMOVED***

	// mixed order
	o := new(histogram)
	p := new(histogram)
	o.addMeasurement(0)
	o.addMeasurement(2)
	o.addMeasurement(0)
	p.addMeasurement(0)
	p.addMeasurement(0)
	p.addMeasurement(2)
	if o.String() != p.String() ***REMOVED***
		t.Errorf("o.String = %q WANT: %q", o.String(), p.String())
	***REMOVED***
***REMOVED***

func add(h *histogram, times int, val int64) ***REMOVED***
	for i := 0; i < times; i++ ***REMOVED***
		h.addMeasurement(val)
	***REMOVED***
***REMOVED***

func isApproximate(x, y float64) bool ***REMOVED***
	return math.Abs(x-y) < 1e-2
***REMOVED***
