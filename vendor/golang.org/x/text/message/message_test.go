// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package message

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"golang.org/x/text/internal"
	"golang.org/x/text/internal/format"
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

type formatFunc func(s fmt.State, v rune)

func (f formatFunc) Format(s fmt.State, v rune) ***REMOVED*** f(s, v) ***REMOVED***

func TestBinding(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		tag   string
		value interface***REMOVED******REMOVED***
		want  string
	***REMOVED******REMOVED***
		***REMOVED***"en", 1, "1"***REMOVED***,
		***REMOVED***"en", "2", "2"***REMOVED***,
		***REMOVED*** // Language is passed.
			"en",
			formatFunc(func(fs fmt.State, v rune) ***REMOVED***
				s := fs.(format.State)
				io.WriteString(s, s.Language().String())
			***REMOVED***),
			"en",
		***REMOVED***,
	***REMOVED***
	for i, tc := range testCases ***REMOVED***
		p := NewPrinter(language.MustParse(tc.tag))
		if got := p.Sprint(tc.value); got != tc.want ***REMOVED***
			t.Errorf("%d:%s:Sprint(%v) = %q; want %q", i, tc.tag, tc.value, got, tc.want)
		***REMOVED***
		var buf bytes.Buffer
		p.Fprint(&buf, tc.value)
		if got := buf.String(); got != tc.want ***REMOVED***
			t.Errorf("%d:%s:Fprint(%v) = %q; want %q", i, tc.tag, tc.value, got, tc.want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestLocalization(t *testing.T) ***REMOVED***
	type test struct ***REMOVED***
		tag  string
		key  Reference
		args []interface***REMOVED******REMOVED***
		want string
	***REMOVED***
	args := func(x ...interface***REMOVED******REMOVED***) []interface***REMOVED******REMOVED*** ***REMOVED*** return x ***REMOVED***
	empty := []interface***REMOVED******REMOVED******REMOVED******REMOVED***
	joe := []interface***REMOVED******REMOVED******REMOVED***"Joe"***REMOVED***
	joeAndMary := []interface***REMOVED******REMOVED******REMOVED***"Joe", "Mary"***REMOVED***

	testCases := []struct ***REMOVED***
		desc string
		cat  []entry
		test []test
	***REMOVED******REMOVED******REMOVED***
		desc: "empty",
		test: []test***REMOVED***
			***REMOVED***"en", "key", empty, "key"***REMOVED***,
			***REMOVED***"en", "", empty, ""***REMOVED***,
			***REMOVED***"nl", "", empty, ""***REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		desc: "hierarchical languages",
		cat: []entry***REMOVED***
			***REMOVED***"en", "hello %s", "Hello %s!"***REMOVED***,
			***REMOVED***"en-GB", "hello %s", "Hellø %s!"***REMOVED***,
			***REMOVED***"en-US", "hello %s", "Howdy %s!"***REMOVED***,
			***REMOVED***"en", "greetings %s and %s", "Greetings %s and %s!"***REMOVED***,
		***REMOVED***,
		test: []test***REMOVED***
			***REMOVED***"und", "hello %s", joe, "hello Joe"***REMOVED***,
			***REMOVED***"nl", "hello %s", joe, "hello Joe"***REMOVED***,
			***REMOVED***"en", "hello %s", joe, "Hello Joe!"***REMOVED***,
			***REMOVED***"en-US", "hello %s", joe, "Howdy Joe!"***REMOVED***,
			***REMOVED***"en-GB", "hello %s", joe, "Hellø Joe!"***REMOVED***,
			***REMOVED***"en-oxendict", "hello %s", joe, "Hello Joe!"***REMOVED***,
			***REMOVED***"en-US-oxendict-u-ms-metric", "hello %s", joe, "Howdy Joe!"***REMOVED***,

			***REMOVED***"und", "greetings %s and %s", joeAndMary, "greetings Joe and Mary"***REMOVED***,
			***REMOVED***"nl", "greetings %s and %s", joeAndMary, "greetings Joe and Mary"***REMOVED***,
			***REMOVED***"en", "greetings %s and %s", joeAndMary, "Greetings Joe and Mary!"***REMOVED***,
			***REMOVED***"en-US", "greetings %s and %s", joeAndMary, "Greetings Joe and Mary!"***REMOVED***,
			***REMOVED***"en-GB", "greetings %s and %s", joeAndMary, "Greetings Joe and Mary!"***REMOVED***,
			***REMOVED***"en-oxendict", "greetings %s and %s", joeAndMary, "Greetings Joe and Mary!"***REMOVED***,
			***REMOVED***"en-US-oxendict-u-ms-metric", "greetings %s and %s", joeAndMary, "Greetings Joe and Mary!"***REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		desc: "references",
		cat: []entry***REMOVED***
			***REMOVED***"en", "hello", "Hello!"***REMOVED***,
		***REMOVED***,
		test: []test***REMOVED***
			***REMOVED***"en", "hello", empty, "Hello!"***REMOVED***,
			***REMOVED***"en", Key("hello", "fallback"), empty, "Hello!"***REMOVED***,
			***REMOVED***"en", Key("xxx", "fallback"), empty, "fallback"***REMOVED***,
			***REMOVED***"und", Key("hello", "fallback"), empty, "fallback"***REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		desc: "zero substitution", // work around limitation of fmt
		cat: []entry***REMOVED***
			***REMOVED***"en", "hello %s", "Hello!"***REMOVED***,
			***REMOVED***"en", "hi %s and %s", "Hello %[2]s!"***REMOVED***,
		***REMOVED***,
		test: []test***REMOVED***
			***REMOVED***"en", "hello %s", joe, "Hello!"***REMOVED***,
			***REMOVED***"en", "hello %s", joeAndMary, "Hello!"***REMOVED***,
			***REMOVED***"en", "hi %s and %s", joeAndMary, "Hello Mary!"***REMOVED***,
			// The following tests resolve to the fallback string.
			***REMOVED***"und", "hello", joeAndMary, "hello"***REMOVED***,
			***REMOVED***"und", "hello %%%%", joeAndMary, "hello %%"***REMOVED***,
			***REMOVED***"und", "hello %#%%4.2%  ", joeAndMary, "hello %%  "***REMOVED***,
			***REMOVED***"und", "hello %s", joeAndMary, "hello Joe%!(EXTRA string=Mary)"***REMOVED***,
			***REMOVED***"und", "hello %+%%s", joeAndMary, "hello %Joe%!(EXTRA string=Mary)"***REMOVED***,
			***REMOVED***"und", "hello %-42%%s ", joeAndMary, "hello %Joe %!(EXTRA string=Mary)"***REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		desc: "number formatting", // work around limitation of fmt
		cat: []entry***REMOVED***
			***REMOVED***"und", "files", "%d files left"***REMOVED***,
			***REMOVED***"und", "meters", "%.2f meters"***REMOVED***,
			***REMOVED***"de", "files", "%d Dateien übrig"***REMOVED***,
		***REMOVED***,
		test: []test***REMOVED***
			***REMOVED***"en", "meters", args(3000.2), "3,000.20 meters"***REMOVED***,
			***REMOVED***"en-u-nu-gujr", "files", args(123456), "૧૨૩,૪૫૬ files left"***REMOVED***,
			***REMOVED***"de", "files", args(1234), "1.234 Dateien übrig"***REMOVED***,
			***REMOVED***"de-CH", "files", args(1234), "1’234 Dateien übrig"***REMOVED***,
			***REMOVED***"de-CH-u-nu-mong", "files", args(1234), "᠑’᠒᠓᠔ Dateien übrig"***REMOVED***,
		***REMOVED***,
	***REMOVED******REMOVED***

	for _, tc := range testCases ***REMOVED***
		cat, _ := initCat(tc.cat)

		for i, pt := range tc.test ***REMOVED***
			t.Run(fmt.Sprintf("%s:%d", tc.desc, i), func(t *testing.T) ***REMOVED***
				p := NewPrinter(language.MustParse(pt.tag), Catalog(cat))

				if got := p.Sprintf(pt.key, pt.args...); got != pt.want ***REMOVED***
					t.Errorf("Sprintf(%q, %v) = %s; want %s",
						pt.key, pt.args, got, pt.want)
					return // Next error will likely be the same.
				***REMOVED***

				w := &bytes.Buffer***REMOVED******REMOVED***
				p.Fprintf(w, pt.key, pt.args...)
				if got := w.String(); got != pt.want ***REMOVED***
					t.Errorf("Fprintf(%q, %v) = %s; want %s",
						pt.key, pt.args, got, pt.want)
				***REMOVED***
			***REMOVED***)
		***REMOVED***
	***REMOVED***
***REMOVED***

type entry struct***REMOVED*** tag, key, msg string ***REMOVED***

func initCat(entries []entry) (*catalog.Builder, []language.Tag) ***REMOVED***
	tags := []language.Tag***REMOVED******REMOVED***
	cat := catalog.NewBuilder()
	for _, e := range entries ***REMOVED***
		tag := language.MustParse(e.tag)
		tags = append(tags, tag)
		cat.SetString(tag, e.key, e.msg)
	***REMOVED***
	return cat, internal.UniqueTags(tags)
***REMOVED***
