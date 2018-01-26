// Copyright 2017, The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package cmp

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/google/go-cmp/cmp/internal/value"
)

type defaultReporter struct ***REMOVED***
	Option
	diffs  []string // List of differences, possibly truncated
	ndiffs int      // Total number of differences
	nbytes int      // Number of bytes in diffs
	nlines int      // Number of lines in diffs
***REMOVED***

var _ reporter = (*defaultReporter)(nil)

func (r *defaultReporter) Report(x, y reflect.Value, eq bool, p Path) ***REMOVED***
	if eq ***REMOVED***
		return // Ignore equal results
	***REMOVED***
	const maxBytes = 4096
	const maxLines = 256
	r.ndiffs++
	if r.nbytes < maxBytes && r.nlines < maxLines ***REMOVED***
		sx := value.Format(x, true)
		sy := value.Format(y, true)
		if sx == sy ***REMOVED***
			// Stringer is not helpful, so rely on more exact formatting.
			sx = value.Format(x, false)
			sy = value.Format(y, false)
		***REMOVED***
		s := fmt.Sprintf("%#v:\n\t-: %s\n\t+: %s\n", p, sx, sy)
		r.diffs = append(r.diffs, s)
		r.nbytes += len(s)
		r.nlines += strings.Count(s, "\n")
	***REMOVED***
***REMOVED***

func (r *defaultReporter) String() string ***REMOVED***
	s := strings.Join(r.diffs, "")
	if r.ndiffs == len(r.diffs) ***REMOVED***
		return s
	***REMOVED***
	return fmt.Sprintf("%s... %d more differences ...", s, len(r.diffs)-r.ndiffs)
***REMOVED***
