// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"fmt"
	"strings"
	"testing"

	"golang.org/x/text/language"
)

func TestUnique(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		in, want string
	***REMOVED******REMOVED***
		***REMOVED***"", "[]"***REMOVED***,
		***REMOVED***"en", "[en]"***REMOVED***,
		***REMOVED***"en en", "[en]"***REMOVED***,
		***REMOVED***"en en en", "[en]"***REMOVED***,
		***REMOVED***"en-u-cu-eur en", "[en en-u-cu-eur]"***REMOVED***,
		***REMOVED***"nl en", "[en nl]"***REMOVED***,
		***REMOVED***"pt-Pt pt", "[pt pt-PT]"***REMOVED***,
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		tags := []language.Tag***REMOVED******REMOVED***
		for _, s := range strings.Split(tc.in, " ") ***REMOVED***
			if s != "" ***REMOVED***
				tags = append(tags, language.MustParse(s))
			***REMOVED***
		***REMOVED***
		if got := fmt.Sprint(UniqueTags(tags)); got != tc.want ***REMOVED***
			t.Errorf("Unique(%s) = %s; want %s", tc.in, got, tc.want)
		***REMOVED***
	***REMOVED***
***REMOVED***
