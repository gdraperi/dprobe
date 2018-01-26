package parser

import (
	"testing"
)

var invalidJSONArraysOfStrings = []string***REMOVED***
	`["a",42,"b"]`,
	`["a",123.456,"b"]`,
	`["a",***REMOVED******REMOVED***,"b"]`,
	`["a",***REMOVED***"c": "d"***REMOVED***,"b"]`,
	`["a",["c"],"b"]`,
	`["a",true,"b"]`,
	`["a",false,"b"]`,
	`["a",null,"b"]`,
***REMOVED***

var validJSONArraysOfStrings = map[string][]string***REMOVED***
	`[]`:           ***REMOVED******REMOVED***,
	`[""]`:         ***REMOVED***""***REMOVED***,
	`["a"]`:        ***REMOVED***"a"***REMOVED***,
	`["a","b"]`:    ***REMOVED***"a", "b"***REMOVED***,
	`[ "a", "b" ]`: ***REMOVED***"a", "b"***REMOVED***,
	`[	"a",	"b"	]`: ***REMOVED***"a", "b"***REMOVED***,
	`	[	"a",	"b"	]	`: ***REMOVED***"a", "b"***REMOVED***,
	`["abc 123", "♥", "☃", "\" \\ \/ \b \f \n \r \t \u0000"]`: ***REMOVED***"abc 123", "♥", "☃", "\" \\ / \b \f \n \r \t \u0000"***REMOVED***,
***REMOVED***

func TestJSONArraysOfStrings(t *testing.T) ***REMOVED***
	for json, expected := range validJSONArraysOfStrings ***REMOVED***
		d := NewDefaultDirective()

		if node, _, err := parseJSON(json, d); err != nil ***REMOVED***
			t.Fatalf("%q should be a valid JSON array of strings, but wasn't! (err: %q)", json, err)
		***REMOVED*** else ***REMOVED***
			i := 0
			for node != nil ***REMOVED***
				if i >= len(expected) ***REMOVED***
					t.Fatalf("expected result is shorter than parsed result (%d vs %d+) in %q", len(expected), i+1, json)
				***REMOVED***
				if node.Value != expected[i] ***REMOVED***
					t.Fatalf("expected %q (not %q) in %q at pos %d", expected[i], node.Value, json, i)
				***REMOVED***
				node = node.Next
				i++
			***REMOVED***
			if i != len(expected) ***REMOVED***
				t.Fatalf("expected result is longer than parsed result (%d vs %d) in %q", len(expected), i+1, json)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	for _, json := range invalidJSONArraysOfStrings ***REMOVED***
		d := NewDefaultDirective()

		if _, _, err := parseJSON(json, d); err != errDockerfileNotStringArray ***REMOVED***
			t.Fatalf("%q should be an invalid JSON array of strings, but wasn't!", json)
		***REMOVED***
	***REMOVED***
***REMOVED***
