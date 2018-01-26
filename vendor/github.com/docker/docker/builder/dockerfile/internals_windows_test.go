// +build windows

package dockerfile

import (
	"fmt"
	"testing"

	"github.com/docker/docker/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeDest(t *testing.T) ***REMOVED***
	tests := []struct***REMOVED*** current, requested, expected, etext string ***REMOVED******REMOVED***
		***REMOVED***``, `D:\`, ``, `Windows does not support destinations not on the system drive (C:)`***REMOVED***,
		***REMOVED***``, `e:/`, ``, `Windows does not support destinations not on the system drive (C:)`***REMOVED***,
		***REMOVED***`invalid`, `./c1`, ``, `Current WorkingDir invalid is not platform consistent`***REMOVED***,
		***REMOVED***`C:`, ``, ``, `Current WorkingDir C: is not platform consistent`***REMOVED***,
		***REMOVED***`C`, ``, ``, `Current WorkingDir C is not platform consistent`***REMOVED***,
		***REMOVED***`D:\`, `.`, ``, "Windows does not support relative paths when WORKDIR is not the system drive"***REMOVED***,
		***REMOVED***``, `D`, `D`, ``***REMOVED***,
		***REMOVED***``, `./a1`, `.\a1`, ``***REMOVED***,
		***REMOVED***``, `.\b1`, `.\b1`, ``***REMOVED***,
		***REMOVED***``, `/`, `\`, ``***REMOVED***,
		***REMOVED***``, `\`, `\`, ``***REMOVED***,
		***REMOVED***``, `c:/`, `\`, ``***REMOVED***,
		***REMOVED***``, `c:\`, `\`, ``***REMOVED***,
		***REMOVED***``, `.`, `.`, ``***REMOVED***,
		***REMOVED***`C:\wdd`, `./a1`, `\wdd\a1`, ``***REMOVED***,
		***REMOVED***`C:\wde`, `.\b1`, `\wde\b1`, ``***REMOVED***,
		***REMOVED***`C:\wdf`, `/`, `\`, ``***REMOVED***,
		***REMOVED***`C:\wdg`, `\`, `\`, ``***REMOVED***,
		***REMOVED***`C:\wdh`, `c:/`, `\`, ``***REMOVED***,
		***REMOVED***`C:\wdi`, `c:\`, `\`, ``***REMOVED***,
		***REMOVED***`C:\wdj`, `.`, `\wdj`, ``***REMOVED***,
		***REMOVED***`C:\wdk`, `foo/bar`, `\wdk\foo\bar`, ``***REMOVED***,
		***REMOVED***`C:\wdl`, `foo\bar`, `\wdl\foo\bar`, ``***REMOVED***,
		***REMOVED***`C:\wdm`, `foo/bar/`, `\wdm\foo\bar\`, ``***REMOVED***,
		***REMOVED***`C:\wdn`, `foo\bar/`, `\wdn\foo\bar\`, ``***REMOVED***,
	***REMOVED***
	for _, testcase := range tests ***REMOVED***
		msg := fmt.Sprintf("Input: %s, %s", testcase.current, testcase.requested)
		actual, err := normalizeDest(testcase.current, testcase.requested, "windows")
		if testcase.etext == "" ***REMOVED***
			if !assert.NoError(t, err, msg) ***REMOVED***
				continue
			***REMOVED***
			assert.Equal(t, testcase.expected, actual, msg)
		***REMOVED*** else ***REMOVED***
			testutil.ErrorContains(t, err, testcase.etext)
		***REMOVED***
	***REMOVED***
***REMOVED***
