package useragent

import "testing"

func TestVersionInfo(t *testing.T) ***REMOVED***
	vi := VersionInfo***REMOVED***"foo", "bar"***REMOVED***
	if !vi.isValid() ***REMOVED***
		t.Fatalf("VersionInfo should be valid")
	***REMOVED***
	vi = VersionInfo***REMOVED***"", "bar"***REMOVED***
	if vi.isValid() ***REMOVED***
		t.Fatalf("Expected VersionInfo to be invalid")
	***REMOVED***
	vi = VersionInfo***REMOVED***"foo", ""***REMOVED***
	if vi.isValid() ***REMOVED***
		t.Fatalf("Expected VersionInfo to be invalid")
	***REMOVED***
***REMOVED***

func TestAppendVersions(t *testing.T) ***REMOVED***
	vis := []VersionInfo***REMOVED***
		***REMOVED***"foo", "1.0"***REMOVED***,
		***REMOVED***"bar", "0.1"***REMOVED***,
		***REMOVED***"pi", "3.1.4"***REMOVED***,
	***REMOVED***
	v := AppendVersions("base", vis...)
	expect := "base foo/1.0 bar/0.1 pi/3.1.4"
	if v != expect ***REMOVED***
		t.Fatalf("expected %q, got %q", expect, v)
	***REMOVED***
***REMOVED***
