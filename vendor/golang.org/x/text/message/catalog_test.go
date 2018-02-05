// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package message

import (
	"strings"
	"testing"

	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

func TestMatchLanguage(t *testing.T) ***REMOVED***
	c := catalog.NewBuilder(catalog.Fallback(language.English))
	c.SetString(language.Bengali, "", "")
	c.SetString(language.English, "", "")
	c.SetString(language.German, "", "")

	testCases := []struct ***REMOVED***
		args string // '|'-separated list
		want string
	***REMOVED******REMOVED******REMOVED***
		args: "de-CH",
		want: "de",
	***REMOVED***, ***REMOVED***
		args: "bn-u-nu-latn|en-US,en;q=0.9,de;q=0.8,nl;q=0.7",
		want: "bn-u-nu-latn",
	***REMOVED***, ***REMOVED***
		args: "gr",
		want: "en",
	***REMOVED******REMOVED***
	for _, tc := range testCases ***REMOVED***
		DefaultCatalog = c
		t.Run(tc.args, func(t *testing.T) ***REMOVED***
			got := MatchLanguage(strings.Split(tc.args, "|")...)
			if got != language.Make(tc.want) ***REMOVED***
				t.Errorf("got %q; want %q", got, tc.want)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
