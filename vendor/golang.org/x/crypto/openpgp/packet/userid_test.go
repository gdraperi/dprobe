// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packet

import (
	"testing"
)

var userIdTests = []struct ***REMOVED***
	id                   string
	name, comment, email string
***REMOVED******REMOVED***
	***REMOVED***"", "", "", ""***REMOVED***,
	***REMOVED***"John Smith", "John Smith", "", ""***REMOVED***,
	***REMOVED***"John Smith ()", "John Smith", "", ""***REMOVED***,
	***REMOVED***"John Smith () <>", "John Smith", "", ""***REMOVED***,
	***REMOVED***"(comment", "", "comment", ""***REMOVED***,
	***REMOVED***"(comment)", "", "comment", ""***REMOVED***,
	***REMOVED***"<email", "", "", "email"***REMOVED***,
	***REMOVED***"<email>   sdfk", "", "", "email"***REMOVED***,
	***REMOVED***"  John Smith  (  Comment ) asdkflj < email > lksdfj", "John Smith", "Comment", "email"***REMOVED***,
	***REMOVED***"  John Smith  < email > lksdfj", "John Smith", "", "email"***REMOVED***,
	***REMOVED***"(<foo", "", "<foo", ""***REMOVED***,
	***REMOVED***"René Descartes (العربي)", "René Descartes", "العربي", ""***REMOVED***,
***REMOVED***

func TestParseUserId(t *testing.T) ***REMOVED***
	for i, test := range userIdTests ***REMOVED***
		name, comment, email := parseUserId(test.id)
		if name != test.name ***REMOVED***
			t.Errorf("%d: name mismatch got:%s want:%s", i, name, test.name)
		***REMOVED***
		if comment != test.comment ***REMOVED***
			t.Errorf("%d: comment mismatch got:%s want:%s", i, comment, test.comment)
		***REMOVED***
		if email != test.email ***REMOVED***
			t.Errorf("%d: email mismatch got:%s want:%s", i, email, test.email)
		***REMOVED***
	***REMOVED***
***REMOVED***

var newUserIdTests = []struct ***REMOVED***
	name, comment, email, id string
***REMOVED******REMOVED***
	***REMOVED***"foo", "", "", "foo"***REMOVED***,
	***REMOVED***"", "bar", "", "(bar)"***REMOVED***,
	***REMOVED***"", "", "baz", "<baz>"***REMOVED***,
	***REMOVED***"foo", "bar", "", "foo (bar)"***REMOVED***,
	***REMOVED***"foo", "", "baz", "foo <baz>"***REMOVED***,
	***REMOVED***"", "bar", "baz", "(bar) <baz>"***REMOVED***,
	***REMOVED***"foo", "bar", "baz", "foo (bar) <baz>"***REMOVED***,
***REMOVED***

func TestNewUserId(t *testing.T) ***REMOVED***
	for i, test := range newUserIdTests ***REMOVED***
		uid := NewUserId(test.name, test.comment, test.email)
		if uid == nil ***REMOVED***
			t.Errorf("#%d: returned nil", i)
			continue
		***REMOVED***
		if uid.Id != test.id ***REMOVED***
			t.Errorf("#%d: got '%s', want '%s'", i, uid.Id, test.id)
		***REMOVED***
	***REMOVED***
***REMOVED***

var invalidNewUserIdTests = []struct ***REMOVED***
	name, comment, email string
***REMOVED******REMOVED***
	***REMOVED***"foo(", "", ""***REMOVED***,
	***REMOVED***"foo<", "", ""***REMOVED***,
	***REMOVED***"", "bar)", ""***REMOVED***,
	***REMOVED***"", "bar<", ""***REMOVED***,
	***REMOVED***"", "", "baz>"***REMOVED***,
	***REMOVED***"", "", "baz)"***REMOVED***,
	***REMOVED***"", "", "baz\x00"***REMOVED***,
***REMOVED***

func TestNewUserIdWithInvalidInput(t *testing.T) ***REMOVED***
	for i, test := range invalidNewUserIdTests ***REMOVED***
		if uid := NewUserId(test.name, test.comment, test.email); uid != nil ***REMOVED***
			t.Errorf("#%d: returned non-nil value: %#v", i, uid)
		***REMOVED***
	***REMOVED***
***REMOVED***
