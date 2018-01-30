// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webdav

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseIfHeader(t *testing.T) ***REMOVED***
	// The "section x.y.z" test cases come from section x.y.z of the spec at
	// http://www.webdav.org/specs/rfc4918.html
	testCases := []struct ***REMOVED***
		desc  string
		input string
		want  ifHeader
	***REMOVED******REMOVED******REMOVED***
		"bad: empty",
		``,
		ifHeader***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"bad: no parens",
		`foobar`,
		ifHeader***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"bad: empty list #1",
		`()`,
		ifHeader***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"bad: empty list #2",
		`(a) (b c) () (d)`,
		ifHeader***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"bad: no list after resource #1",
		`<foo>`,
		ifHeader***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"bad: no list after resource #2",
		`<foo> <bar> (a)`,
		ifHeader***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"bad: no list after resource #3",
		`<foo> (a) (b) <bar>`,
		ifHeader***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"bad: no-tag-list followed by tagged-list",
		`(a) (b) <foo> (c)`,
		ifHeader***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"bad: unfinished list",
		`(a`,
		ifHeader***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"bad: unfinished ETag",
		`([b`,
		ifHeader***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"bad: unfinished Notted list",
		`(Not a`,
		ifHeader***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"bad: double Not",
		`(Not Not a)`,
		ifHeader***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"good: one list with a Token",
		`(a)`,
		ifHeader***REMOVED***
			lists: []ifList***REMOVED******REMOVED***
				conditions: []Condition***REMOVED******REMOVED***
					Token: `a`,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		"good: one list with an ETag",
		`([a])`,
		ifHeader***REMOVED***
			lists: []ifList***REMOVED******REMOVED***
				conditions: []Condition***REMOVED******REMOVED***
					ETag: `a`,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		"good: one list with three Nots",
		`(Not a Not b Not [d])`,
		ifHeader***REMOVED***
			lists: []ifList***REMOVED******REMOVED***
				conditions: []Condition***REMOVED******REMOVED***
					Not:   true,
					Token: `a`,
				***REMOVED***, ***REMOVED***
					Not:   true,
					Token: `b`,
				***REMOVED***, ***REMOVED***
					Not:  true,
					ETag: `d`,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		"good: two lists",
		`(a) (b)`,
		ifHeader***REMOVED***
			lists: []ifList***REMOVED******REMOVED***
				conditions: []Condition***REMOVED******REMOVED***
					Token: `a`,
				***REMOVED******REMOVED***,
			***REMOVED***, ***REMOVED***
				conditions: []Condition***REMOVED******REMOVED***
					Token: `b`,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		"good: two Notted lists",
		`(Not a) (Not b)`,
		ifHeader***REMOVED***
			lists: []ifList***REMOVED******REMOVED***
				conditions: []Condition***REMOVED******REMOVED***
					Not:   true,
					Token: `a`,
				***REMOVED******REMOVED***,
			***REMOVED***, ***REMOVED***
				conditions: []Condition***REMOVED******REMOVED***
					Not:   true,
					Token: `b`,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		"section 7.5.1",
		`<http://www.example.com/users/f/fielding/index.html> 
			(<urn:uuid:f81d4fae-7dec-11d0-a765-00a0c91e6bf6>)`,
		ifHeader***REMOVED***
			lists: []ifList***REMOVED******REMOVED***
				resourceTag: `http://www.example.com/users/f/fielding/index.html`,
				conditions: []Condition***REMOVED******REMOVED***
					Token: `urn:uuid:f81d4fae-7dec-11d0-a765-00a0c91e6bf6`,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		"section 7.5.2 #1",
		`(<urn:uuid:150852e2-3847-42d5-8cbe-0f4f296f26cf>)`,
		ifHeader***REMOVED***
			lists: []ifList***REMOVED******REMOVED***
				conditions: []Condition***REMOVED******REMOVED***
					Token: `urn:uuid:150852e2-3847-42d5-8cbe-0f4f296f26cf`,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		"section 7.5.2 #2",
		`<http://example.com/locked/>
			(<urn:uuid:150852e2-3847-42d5-8cbe-0f4f296f26cf>)`,
		ifHeader***REMOVED***
			lists: []ifList***REMOVED******REMOVED***
				resourceTag: `http://example.com/locked/`,
				conditions: []Condition***REMOVED******REMOVED***
					Token: `urn:uuid:150852e2-3847-42d5-8cbe-0f4f296f26cf`,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		"section 7.5.2 #3",
		`<http://example.com/locked/member>
			(<urn:uuid:150852e2-3847-42d5-8cbe-0f4f296f26cf>)`,
		ifHeader***REMOVED***
			lists: []ifList***REMOVED******REMOVED***
				resourceTag: `http://example.com/locked/member`,
				conditions: []Condition***REMOVED******REMOVED***
					Token: `urn:uuid:150852e2-3847-42d5-8cbe-0f4f296f26cf`,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		"section 9.9.6",
		`(<urn:uuid:fe184f2e-6eec-41d0-c765-01adc56e6bb4>) 
			(<urn:uuid:e454f3f3-acdc-452a-56c7-00a5c91e4b77>)`,
		ifHeader***REMOVED***
			lists: []ifList***REMOVED******REMOVED***
				conditions: []Condition***REMOVED******REMOVED***
					Token: `urn:uuid:fe184f2e-6eec-41d0-c765-01adc56e6bb4`,
				***REMOVED******REMOVED***,
			***REMOVED***, ***REMOVED***
				conditions: []Condition***REMOVED******REMOVED***
					Token: `urn:uuid:e454f3f3-acdc-452a-56c7-00a5c91e4b77`,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		"section 9.10.8",
		`(<urn:uuid:e71d4fae-5dec-22d6-fea5-00a0c91e6be4>)`,
		ifHeader***REMOVED***
			lists: []ifList***REMOVED******REMOVED***
				conditions: []Condition***REMOVED******REMOVED***
					Token: `urn:uuid:e71d4fae-5dec-22d6-fea5-00a0c91e6be4`,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		"section 10.4.6",
		`(<urn:uuid:181d4fae-7d8c-11d0-a765-00a0c91e6bf2> 
			["I am an ETag"])
			(["I am another ETag"])`,
		ifHeader***REMOVED***
			lists: []ifList***REMOVED******REMOVED***
				conditions: []Condition***REMOVED******REMOVED***
					Token: `urn:uuid:181d4fae-7d8c-11d0-a765-00a0c91e6bf2`,
				***REMOVED***, ***REMOVED***
					ETag: `"I am an ETag"`,
				***REMOVED******REMOVED***,
			***REMOVED***, ***REMOVED***
				conditions: []Condition***REMOVED******REMOVED***
					ETag: `"I am another ETag"`,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		"section 10.4.7",
		`(Not <urn:uuid:181d4fae-7d8c-11d0-a765-00a0c91e6bf2> 
			<urn:uuid:58f202ac-22cf-11d1-b12d-002035b29092>)`,
		ifHeader***REMOVED***
			lists: []ifList***REMOVED******REMOVED***
				conditions: []Condition***REMOVED******REMOVED***
					Not:   true,
					Token: `urn:uuid:181d4fae-7d8c-11d0-a765-00a0c91e6bf2`,
				***REMOVED***, ***REMOVED***
					Token: `urn:uuid:58f202ac-22cf-11d1-b12d-002035b29092`,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		"section 10.4.8",
		`(<urn:uuid:181d4fae-7d8c-11d0-a765-00a0c91e6bf2>) 
			(Not <DAV:no-lock>)`,
		ifHeader***REMOVED***
			lists: []ifList***REMOVED******REMOVED***
				conditions: []Condition***REMOVED******REMOVED***
					Token: `urn:uuid:181d4fae-7d8c-11d0-a765-00a0c91e6bf2`,
				***REMOVED******REMOVED***,
			***REMOVED***, ***REMOVED***
				conditions: []Condition***REMOVED******REMOVED***
					Not:   true,
					Token: `DAV:no-lock`,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		"section 10.4.9",
		`</resource1> 
			(<urn:uuid:181d4fae-7d8c-11d0-a765-00a0c91e6bf2> 
			[W/"A weak ETag"]) (["strong ETag"])`,
		ifHeader***REMOVED***
			lists: []ifList***REMOVED******REMOVED***
				resourceTag: `/resource1`,
				conditions: []Condition***REMOVED******REMOVED***
					Token: `urn:uuid:181d4fae-7d8c-11d0-a765-00a0c91e6bf2`,
				***REMOVED***, ***REMOVED***
					ETag: `W/"A weak ETag"`,
				***REMOVED******REMOVED***,
			***REMOVED***, ***REMOVED***
				resourceTag: `/resource1`,
				conditions: []Condition***REMOVED******REMOVED***
					ETag: `"strong ETag"`,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		"section 10.4.10",
		`<http://www.example.com/specs/> 
			(<urn:uuid:181d4fae-7d8c-11d0-a765-00a0c91e6bf2>)`,
		ifHeader***REMOVED***
			lists: []ifList***REMOVED******REMOVED***
				resourceTag: `http://www.example.com/specs/`,
				conditions: []Condition***REMOVED******REMOVED***
					Token: `urn:uuid:181d4fae-7d8c-11d0-a765-00a0c91e6bf2`,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		"section 10.4.11 #1",
		`</specs/rfc2518.doc> (["4217"])`,
		ifHeader***REMOVED***
			lists: []ifList***REMOVED******REMOVED***
				resourceTag: `/specs/rfc2518.doc`,
				conditions: []Condition***REMOVED******REMOVED***
					ETag: `"4217"`,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		"section 10.4.11 #2",
		`</specs/rfc2518.doc> (Not ["4217"])`,
		ifHeader***REMOVED***
			lists: []ifList***REMOVED******REMOVED***
				resourceTag: `/specs/rfc2518.doc`,
				conditions: []Condition***REMOVED******REMOVED***
					Not:  true,
					ETag: `"4217"`,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED******REMOVED***

	for _, tc := range testCases ***REMOVED***
		got, ok := parseIfHeader(strings.Replace(tc.input, "\n", "", -1))
		if gotEmpty := reflect.DeepEqual(got, ifHeader***REMOVED******REMOVED***); gotEmpty == ok ***REMOVED***
			t.Errorf("%s: should be different: empty header == %t, ok == %t", tc.desc, gotEmpty, ok)
			continue
		***REMOVED***
		if !reflect.DeepEqual(got, tc.want) ***REMOVED***
			t.Errorf("%s:\ngot  %v\nwant %v", tc.desc, got, tc.want)
			continue
		***REMOVED***
	***REMOVED***
***REMOVED***
