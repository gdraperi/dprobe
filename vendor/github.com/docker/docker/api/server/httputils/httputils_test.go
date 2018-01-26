package httputils

import "testing"

// matchesContentType
func TestJsonContentType(t *testing.T) ***REMOVED***
	if !matchesContentType("application/json", "application/json") ***REMOVED***
		t.Fail()
	***REMOVED***

	if !matchesContentType("application/json; charset=utf-8", "application/json") ***REMOVED***
		t.Fail()
	***REMOVED***

	if matchesContentType("dockerapplication/json", "application/json") ***REMOVED***
		t.Fail()
	***REMOVED***
***REMOVED***
