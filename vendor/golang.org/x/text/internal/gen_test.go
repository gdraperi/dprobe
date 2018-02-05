// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"testing"

	"golang.org/x/text/language"
)

func TestParents(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		tag, parent string
	***REMOVED******REMOVED***
		***REMOVED***"af", "und"***REMOVED***,
		***REMOVED***"en", "und"***REMOVED***,
		***REMOVED***"en-001", "en"***REMOVED***,
		***REMOVED***"en-AU", "en-001"***REMOVED***,
		***REMOVED***"en-US", "en"***REMOVED***,
		***REMOVED***"en-US-u-va-posix", "en-US"***REMOVED***,
		***REMOVED***"ca-ES-valencia", "ca-ES"***REMOVED***,
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		tag, ok := language.CompactIndex(language.MustParse(tc.tag))
		if !ok ***REMOVED***
			t.Fatalf("Could not get index of flag %s", tc.tag)
		***REMOVED***
		want, ok := language.CompactIndex(language.MustParse(tc.parent))
		if !ok ***REMOVED***
			t.Fatalf("Could not get index of parent %s of tag %s", tc.parent, tc.tag)
		***REMOVED***
		if got := int(Parent[tag]); got != want ***REMOVED***
			t.Errorf("Parent[%s] = %d; want %d (%s)", tc.tag, got, want, tc.parent)
		***REMOVED***
	***REMOVED***
***REMOVED***
