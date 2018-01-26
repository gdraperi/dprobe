// +build windows

package dockerfile

import "testing"

func TestNormalizeWorkdir(t *testing.T) ***REMOVED***
	tests := []struct***REMOVED*** platform, current, requested, expected, etext string ***REMOVED******REMOVED***
		***REMOVED***"windows", ``, ``, ``, `cannot normalize nothing`***REMOVED***,
		***REMOVED***"windows", ``, `C:`, ``, `C:. is not a directory. If you are specifying a drive letter, please add a trailing '\'`***REMOVED***,
		***REMOVED***"windows", ``, `C:.`, ``, `C:. is not a directory. If you are specifying a drive letter, please add a trailing '\'`***REMOVED***,
		***REMOVED***"windows", `c:`, `\a`, ``, `c:. is not a directory. If you are specifying a drive letter, please add a trailing '\'`***REMOVED***,
		***REMOVED***"windows", `c:.`, `\a`, ``, `c:. is not a directory. If you are specifying a drive letter, please add a trailing '\'`***REMOVED***,
		***REMOVED***"windows", ``, `a`, `C:\a`, ``***REMOVED***,
		***REMOVED***"windows", ``, `c:\foo`, `C:\foo`, ``***REMOVED***,
		***REMOVED***"windows", ``, `c:\\foo`, `C:\foo`, ``***REMOVED***,
		***REMOVED***"windows", ``, `\foo`, `C:\foo`, ``***REMOVED***,
		***REMOVED***"windows", ``, `\\foo`, `C:\foo`, ``***REMOVED***,
		***REMOVED***"windows", ``, `/foo`, `C:\foo`, ``***REMOVED***,
		***REMOVED***"windows", ``, `C:/foo`, `C:\foo`, ``***REMOVED***,
		***REMOVED***"windows", `C:\foo`, `bar`, `C:\foo\bar`, ``***REMOVED***,
		***REMOVED***"windows", `C:\foo`, `/bar`, `C:\bar`, ``***REMOVED***,
		***REMOVED***"windows", `C:\foo`, `\bar`, `C:\bar`, ``***REMOVED***,
		***REMOVED***"linux", ``, ``, ``, `cannot normalize nothing`***REMOVED***,
		***REMOVED***"linux", ``, `foo`, `/foo`, ``***REMOVED***,
		***REMOVED***"linux", ``, `/foo`, `/foo`, ``***REMOVED***,
		***REMOVED***"linux", `/foo`, `bar`, `/foo/bar`, ``***REMOVED***,
		***REMOVED***"linux", `/foo`, `/bar`, `/bar`, ``***REMOVED***,
		***REMOVED***"linux", `\a`, `b\c`, `/a/b/c`, ``***REMOVED***,
	***REMOVED***
	for _, i := range tests ***REMOVED***
		r, e := normalizeWorkdir(i.platform, i.current, i.requested)

		if i.etext != "" && e == nil ***REMOVED***
			t.Fatalf("TestNormalizeWorkingDir Expected error %s for '%s' '%s', got no error", i.etext, i.current, i.requested)
		***REMOVED***

		if i.etext != "" && e.Error() != i.etext ***REMOVED***
			t.Fatalf("TestNormalizeWorkingDir Expected error %s for '%s' '%s', got %s", i.etext, i.current, i.requested, e.Error())
		***REMOVED***

		if r != i.expected ***REMOVED***
			t.Fatalf("TestNormalizeWorkingDir Expected '%s' for '%s' '%s', got '%s'", i.expected, i.current, i.requested, r)
		***REMOVED***
	***REMOVED***
***REMOVED***
