// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package catalog

import (
	"bytes"
	"path"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/text/internal/catmsg"
	"golang.org/x/text/language"
)

type entry struct ***REMOVED***
	tag, key string
	msg      interface***REMOVED******REMOVED***
***REMOVED***

func langs(s string) []language.Tag ***REMOVED***
	t, _, _ := language.ParseAcceptLanguage(s)
	return t
***REMOVED***

type testCase struct ***REMOVED***
	desc     string
	cat      []entry
	lookup   []entry
	fallback string
	match    []string
	tags     []language.Tag
***REMOVED***

var testCases = []testCase***REMOVED******REMOVED***
	desc: "empty catalog",
	lookup: []entry***REMOVED***
		***REMOVED***"en", "key", ""***REMOVED***,
		***REMOVED***"en", "", ""***REMOVED***,
		***REMOVED***"nl", "", ""***REMOVED***,
	***REMOVED***,
	match: []string***REMOVED***
		"gr -> und",
		"en-US -> und",
		"af -> und",
	***REMOVED***,
	tags: nil, // not an empty list.
***REMOVED***, ***REMOVED***
	desc: "one entry",
	cat: []entry***REMOVED***
		***REMOVED***"en", "hello", "Hello!"***REMOVED***,
	***REMOVED***,
	lookup: []entry***REMOVED***
		***REMOVED***"und", "hello", ""***REMOVED***,
		***REMOVED***"nl", "hello", ""***REMOVED***,
		***REMOVED***"en", "hello", "Hello!"***REMOVED***,
		***REMOVED***"en-US", "hello", "Hello!"***REMOVED***,
		***REMOVED***"en-GB", "hello", "Hello!"***REMOVED***,
		***REMOVED***"en-oxendict", "hello", "Hello!"***REMOVED***,
		***REMOVED***"en-oxendict-u-ms-metric", "hello", "Hello!"***REMOVED***,
	***REMOVED***,
	match: []string***REMOVED***
		"gr -> en",
		"en-US -> en",
	***REMOVED***,
	tags: langs("en"),
***REMOVED***, ***REMOVED***
	desc: "hierarchical languages",
	cat: []entry***REMOVED***
		***REMOVED***"en", "hello", "Hello!"***REMOVED***,
		***REMOVED***"en-GB", "hello", "Hellø!"***REMOVED***,
		***REMOVED***"en-US", "hello", "Howdy!"***REMOVED***,
		***REMOVED***"en", "greetings", "Greetings!"***REMOVED***,
		***REMOVED***"gsw", "hello", "Grüetzi!"***REMOVED***,
	***REMOVED***,
	lookup: []entry***REMOVED***
		***REMOVED***"und", "hello", ""***REMOVED***,
		***REMOVED***"nl", "hello", ""***REMOVED***,
		***REMOVED***"en", "hello", "Hello!"***REMOVED***,
		***REMOVED***"en-US", "hello", "Howdy!"***REMOVED***,
		***REMOVED***"en-GB", "hello", "Hellø!"***REMOVED***,
		***REMOVED***"en-oxendict", "hello", "Hello!"***REMOVED***,
		***REMOVED***"en-US-oxendict-u-ms-metric", "hello", "Howdy!"***REMOVED***,

		***REMOVED***"und", "greetings", ""***REMOVED***,
		***REMOVED***"nl", "greetings", ""***REMOVED***,
		***REMOVED***"en", "greetings", "Greetings!"***REMOVED***,
		***REMOVED***"en-US", "greetings", "Greetings!"***REMOVED***,
		***REMOVED***"en-GB", "greetings", "Greetings!"***REMOVED***,
		***REMOVED***"en-oxendict", "greetings", "Greetings!"***REMOVED***,
		***REMOVED***"en-US-oxendict-u-ms-metric", "greetings", "Greetings!"***REMOVED***,
	***REMOVED***,
	fallback: "gsw",
	match: []string***REMOVED***
		"gr -> gsw",
		"en-US -> en-US",
	***REMOVED***,
	tags: langs("gsw, en, en-GB, en-US"),
***REMOVED***, ***REMOVED***
	desc: "variables",
	cat: []entry***REMOVED***
		***REMOVED***"en", "hello %s", []Message***REMOVED***
			Var("person", String("Jane")),
			String("Hello $***REMOVED***person***REMOVED***!"),
		***REMOVED******REMOVED***,
		***REMOVED***"en", "hello error", []Message***REMOVED***
			Var("person", String("Jane")),
			noMatchMessage***REMOVED******REMOVED***, // trigger sequence path.
			String("Hello $***REMOVED***person."),
		***REMOVED******REMOVED***,
		***REMOVED***"en", "fallback to var value", []Message***REMOVED***
			Var("you", noMatchMessage***REMOVED******REMOVED***, noMatchMessage***REMOVED******REMOVED***),
			String("Hello $***REMOVED***you***REMOVED***."),
		***REMOVED******REMOVED***,
		***REMOVED***"en", "scopes", []Message***REMOVED***
			Var("person1", String("Mark")),
			Var("person2", String("Jane")),
			Var("couple",
				Var("person1", String("Joe")),
				String("$***REMOVED***person1***REMOVED*** and $***REMOVED***person2***REMOVED***")),
			String("Hello $***REMOVED***couple***REMOVED***."),
		***REMOVED******REMOVED***,
		***REMOVED***"en", "missing var", String("Hello $***REMOVED***missing***REMOVED***.")***REMOVED***,
	***REMOVED***,
	lookup: []entry***REMOVED***
		***REMOVED***"en", "hello %s", "Hello Jane!"***REMOVED***,
		***REMOVED***"en", "hello error", "Hello $!(MISSINGBRACE)"***REMOVED***,
		***REMOVED***"en", "fallback to var value", "Hello you."***REMOVED***,
		***REMOVED***"en", "scopes", "Hello Joe and Jane."***REMOVED***,
		***REMOVED***"en", "missing var", "Hello missing."***REMOVED***,
	***REMOVED***,
	tags: langs("en"),
***REMOVED***, ***REMOVED***
	desc: "macros",
	cat: []entry***REMOVED***
		***REMOVED***"en", "macro1", String("Hello $***REMOVED***macro1(1)***REMOVED***.")***REMOVED***,
		***REMOVED***"en", "macro2", String("Hello $***REMOVED*** macro1(2) ***REMOVED***!")***REMOVED***,
		***REMOVED***"en", "macroWS", String("Hello $***REMOVED*** macro1( 2 ) ***REMOVED***!")***REMOVED***,
		***REMOVED***"en", "missing", String("Hello $***REMOVED*** missing(1 ***REMOVED***.")***REMOVED***,
		***REMOVED***"en", "badnum", String("Hello $***REMOVED*** badnum(1b) ***REMOVED***.")***REMOVED***,
		***REMOVED***"en", "undefined", String("Hello $***REMOVED*** undefined(1) ***REMOVED***.")***REMOVED***,
		***REMOVED***"en", "macroU", String("Hello $***REMOVED*** macroU(2) ***REMOVED***!")***REMOVED***,
	***REMOVED***,
	lookup: []entry***REMOVED***
		***REMOVED***"en", "macro1", "Hello Joe."***REMOVED***,
		***REMOVED***"en", "macro2", "Hello Joe!"***REMOVED***,
		***REMOVED***"en-US", "macroWS", "Hello Joe!"***REMOVED***,
		***REMOVED***"en-NL", "missing", "Hello $!(MISSINGPAREN)."***REMOVED***,
		***REMOVED***"en", "badnum", "Hello $!(BADNUM)."***REMOVED***,
		***REMOVED***"en", "undefined", "Hello undefined."***REMOVED***,
		***REMOVED***"en", "macroU", "Hello macroU!"***REMOVED***,
	***REMOVED***,
	tags: langs("en"),
***REMOVED******REMOVED***

func setMacros(b *Builder) ***REMOVED***
	b.SetMacro(language.English, "macro1", String("Joe"))
	b.SetMacro(language.Und, "macro2", String("$***REMOVED***macro1(1)***REMOVED***"))
	b.SetMacro(language.English, "macroU", noMatchMessage***REMOVED******REMOVED***)
***REMOVED***

type buildFunc func(t *testing.T, tc testCase) Catalog

func initBuilder(t *testing.T, tc testCase) Catalog ***REMOVED***
	options := []Option***REMOVED******REMOVED***
	if tc.fallback != "" ***REMOVED***
		options = append(options, Fallback(language.MustParse(tc.fallback)))
	***REMOVED***
	cat := NewBuilder(options...)
	for _, e := range tc.cat ***REMOVED***
		tag := language.MustParse(e.tag)
		switch msg := e.msg.(type) ***REMOVED***
		case string:

			cat.SetString(tag, e.key, msg)
		case Message:
			cat.Set(tag, e.key, msg)
		case []Message:
			cat.Set(tag, e.key, msg...)
		***REMOVED***
	***REMOVED***
	setMacros(cat)
	return cat
***REMOVED***

type dictionary map[string]string

func (d dictionary) Lookup(key string) (data string, ok bool) ***REMOVED***
	data, ok = d[key]
	return data, ok
***REMOVED***

func initCatalog(t *testing.T, tc testCase) Catalog ***REMOVED***
	m := map[string]Dictionary***REMOVED******REMOVED***
	for _, e := range tc.cat ***REMOVED***
		m[e.tag] = dictionary***REMOVED******REMOVED***
	***REMOVED***
	for _, e := range tc.cat ***REMOVED***
		var msg Message
		switch x := e.msg.(type) ***REMOVED***
		case string:
			msg = String(x)
		case Message:
			msg = x
		case []Message:
			msg = firstInSequence(x)
		***REMOVED***
		data, _ := catmsg.Compile(language.MustParse(e.tag), nil, msg)
		m[e.tag].(dictionary)[e.key] = data
	***REMOVED***
	options := []Option***REMOVED******REMOVED***
	if tc.fallback != "" ***REMOVED***
		options = append(options, Fallback(language.MustParse(tc.fallback)))
	***REMOVED***
	c, err := NewFromMap(m, options...)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	// TODO: implement macros for fixed catalogs.
	b := NewBuilder()
	setMacros(b)
	c.(*catalog).macros.index = b.macros.index
	return c
***REMOVED***

func TestMatcher(t *testing.T) ***REMOVED***
	test := func(t *testing.T, init buildFunc) ***REMOVED***
		for _, tc := range testCases ***REMOVED***
			for _, s := range tc.match ***REMOVED***
				a := strings.Split(s, "->")
				t.Run(path.Join(tc.desc, a[0]), func(t *testing.T) ***REMOVED***
					cat := init(t, tc)
					got, _ := language.MatchStrings(cat.Matcher(), a[0])
					want := language.MustParse(strings.TrimSpace(a[1]))
					if got != want ***REMOVED***
						t.Errorf("got %q; want %q", got, want)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	t.Run("Builder", func(t *testing.T) ***REMOVED*** test(t, initBuilder) ***REMOVED***)
	t.Run("Catalog", func(t *testing.T) ***REMOVED*** test(t, initCatalog) ***REMOVED***)
***REMOVED***

func TestCatalog(t *testing.T) ***REMOVED***
	test := func(t *testing.T, init buildFunc) ***REMOVED***
		for _, tc := range testCases ***REMOVED***
			cat := init(t, tc)
			wantTags := tc.tags
			if got := cat.Languages(); !reflect.DeepEqual(got, wantTags) ***REMOVED***
				t.Errorf("%s:Languages: got %v; want %v", tc.desc, got, wantTags)
			***REMOVED***

			for _, e := range tc.lookup ***REMOVED***
				t.Run(path.Join(tc.desc, e.tag, e.key), func(t *testing.T) ***REMOVED***
					tag := language.MustParse(e.tag)
					buf := testRenderer***REMOVED******REMOVED***
					ctx := cat.Context(tag, &buf)
					want := e.msg.(string)
					err := ctx.Execute(e.key)
					gotFound := err != ErrNotFound
					wantFound := want != ""
					if gotFound != wantFound ***REMOVED***
						t.Fatalf("err: got %v (%v); want %v", gotFound, err, wantFound)
					***REMOVED***
					if got := buf.buf.String(); got != want ***REMOVED***
						t.Errorf("Lookup:\ngot  %q\nwant %q", got, want)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	t.Run("Builder", func(t *testing.T) ***REMOVED*** test(t, initBuilder) ***REMOVED***)
	t.Run("Catalog", func(t *testing.T) ***REMOVED*** test(t, initCatalog) ***REMOVED***)
***REMOVED***

type testRenderer struct ***REMOVED***
	buf bytes.Buffer
***REMOVED***

func (f *testRenderer) Arg(i int) interface***REMOVED******REMOVED*** ***REMOVED*** return nil ***REMOVED***
func (f *testRenderer) Render(s string)       ***REMOVED*** f.buf.WriteString(s) ***REMOVED***

var msgNoMatch = catmsg.Register("no match", func(d *catmsg.Decoder) bool ***REMOVED***
	return false // no match
***REMOVED***)

type noMatchMessage struct***REMOVED******REMOVED***

func (noMatchMessage) Compile(e *catmsg.Encoder) error ***REMOVED***
	e.EncodeMessageType(msgNoMatch)
	return catmsg.ErrIncomplete
***REMOVED***
