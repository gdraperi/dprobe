// Copyright 2013 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	// MinimumTick is the minimum supported time resolution. This has to be
	// at least time.Second in order for the code below to work.
	minimumTick = time.Millisecond
	// second is the Time duration equivalent to one second.
	second = int64(time.Second / minimumTick)
	// The number of nanoseconds per minimum tick.
	nanosPerTick = int64(minimumTick / time.Nanosecond)

	// Earliest is the earliest Time representable. Handy for
	// initializing a high watermark.
	Earliest = Time(math.MinInt64)
	// Latest is the latest Time representable. Handy for initializing
	// a low watermark.
	Latest = Time(math.MaxInt64)
)

// Time is the number of milliseconds since the epoch
// (1970-01-01 00:00 UTC) excluding leap seconds.
type Time int64

// Interval describes and interval between two timestamps.
type Interval struct ***REMOVED***
	Start, End Time
***REMOVED***

// Now returns the current time as a Time.
func Now() Time ***REMOVED***
	return TimeFromUnixNano(time.Now().UnixNano())
***REMOVED***

// TimeFromUnix returns the Time equivalent to the Unix Time t
// provided in seconds.
func TimeFromUnix(t int64) Time ***REMOVED***
	return Time(t * second)
***REMOVED***

// TimeFromUnixNano returns the Time equivalent to the Unix Time
// t provided in nanoseconds.
func TimeFromUnixNano(t int64) Time ***REMOVED***
	return Time(t / nanosPerTick)
***REMOVED***

// Equal reports whether two Times represent the same instant.
func (t Time) Equal(o Time) bool ***REMOVED***
	return t == o
***REMOVED***

// Before reports whether the Time t is before o.
func (t Time) Before(o Time) bool ***REMOVED***
	return t < o
***REMOVED***

// After reports whether the Time t is after o.
func (t Time) After(o Time) bool ***REMOVED***
	return t > o
***REMOVED***

// Add returns the Time t + d.
func (t Time) Add(d time.Duration) Time ***REMOVED***
	return t + Time(d/minimumTick)
***REMOVED***

// Sub returns the Duration t - o.
func (t Time) Sub(o Time) time.Duration ***REMOVED***
	return time.Duration(t-o) * minimumTick
***REMOVED***

// Time returns the time.Time representation of t.
func (t Time) Time() time.Time ***REMOVED***
	return time.Unix(int64(t)/second, (int64(t)%second)*nanosPerTick)
***REMOVED***

// Unix returns t as a Unix time, the number of seconds elapsed
// since January 1, 1970 UTC.
func (t Time) Unix() int64 ***REMOVED***
	return int64(t) / second
***REMOVED***

// UnixNano returns t as a Unix time, the number of nanoseconds elapsed
// since January 1, 1970 UTC.
func (t Time) UnixNano() int64 ***REMOVED***
	return int64(t) * nanosPerTick
***REMOVED***

// The number of digits after the dot.
var dotPrecision = int(math.Log10(float64(second)))

// String returns a string representation of the Time.
func (t Time) String() string ***REMOVED***
	return strconv.FormatFloat(float64(t)/float64(second), 'f', -1, 64)
***REMOVED***

// MarshalJSON implements the json.Marshaler interface.
func (t Time) MarshalJSON() ([]byte, error) ***REMOVED***
	return []byte(t.String()), nil
***REMOVED***

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Time) UnmarshalJSON(b []byte) error ***REMOVED***
	p := strings.Split(string(b), ".")
	switch len(p) ***REMOVED***
	case 1:
		v, err := strconv.ParseInt(string(p[0]), 10, 64)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*t = Time(v * second)

	case 2:
		v, err := strconv.ParseInt(string(p[0]), 10, 64)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		v *= second

		prec := dotPrecision - len(p[1])
		if prec < 0 ***REMOVED***
			p[1] = p[1][:dotPrecision]
		***REMOVED*** else if prec > 0 ***REMOVED***
			p[1] = p[1] + strings.Repeat("0", prec)
		***REMOVED***

		va, err := strconv.ParseInt(p[1], 10, 32)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		*t = Time(v + va)

	default:
		return fmt.Errorf("invalid time %q", string(b))
	***REMOVED***
	return nil
***REMOVED***

// Duration wraps time.Duration. It is used to parse the custom duration format
// from YAML.
// This type should not propagate beyond the scope of input/output processing.
type Duration time.Duration

var durationRE = regexp.MustCompile("^([0-9]+)(y|w|d|h|m|s|ms)$")

// StringToDuration parses a string into a time.Duration, assuming that a year
// always has 365d, a week always has 7d, and a day always has 24h.
func ParseDuration(durationStr string) (Duration, error) ***REMOVED***
	matches := durationRE.FindStringSubmatch(durationStr)
	if len(matches) != 3 ***REMOVED***
		return 0, fmt.Errorf("not a valid duration string: %q", durationStr)
	***REMOVED***
	var (
		n, _ = strconv.Atoi(matches[1])
		dur  = time.Duration(n) * time.Millisecond
	)
	switch unit := matches[2]; unit ***REMOVED***
	case "y":
		dur *= 1000 * 60 * 60 * 24 * 365
	case "w":
		dur *= 1000 * 60 * 60 * 24 * 7
	case "d":
		dur *= 1000 * 60 * 60 * 24
	case "h":
		dur *= 1000 * 60 * 60
	case "m":
		dur *= 1000 * 60
	case "s":
		dur *= 1000
	case "ms":
		// Value already correct
	default:
		return 0, fmt.Errorf("invalid time unit in duration string: %q", unit)
	***REMOVED***
	return Duration(dur), nil
***REMOVED***

func (d Duration) String() string ***REMOVED***
	var (
		ms   = int64(time.Duration(d) / time.Millisecond)
		unit = "ms"
	)
	factors := map[string]int64***REMOVED***
		"y":  1000 * 60 * 60 * 24 * 365,
		"w":  1000 * 60 * 60 * 24 * 7,
		"d":  1000 * 60 * 60 * 24,
		"h":  1000 * 60 * 60,
		"m":  1000 * 60,
		"s":  1000,
		"ms": 1,
	***REMOVED***

	switch int64(0) ***REMOVED***
	case ms % factors["y"]:
		unit = "y"
	case ms % factors["w"]:
		unit = "w"
	case ms % factors["d"]:
		unit = "d"
	case ms % factors["h"]:
		unit = "h"
	case ms % factors["m"]:
		unit = "m"
	case ms % factors["s"]:
		unit = "s"
	***REMOVED***
	return fmt.Sprintf("%v%v", ms/factors[unit], unit)
***REMOVED***

// MarshalYAML implements the yaml.Marshaler interface.
func (d Duration) MarshalYAML() (interface***REMOVED******REMOVED***, error) ***REMOVED***
	return d.String(), nil
***REMOVED***

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (d *Duration) UnmarshalYAML(unmarshal func(interface***REMOVED******REMOVED***) error) error ***REMOVED***
	var s string
	if err := unmarshal(&s); err != nil ***REMOVED***
		return err
	***REMOVED***
	dur, err := ParseDuration(s)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*d = dur
	return nil
***REMOVED***
