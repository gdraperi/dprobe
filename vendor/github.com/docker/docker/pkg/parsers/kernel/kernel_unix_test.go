// +build !windows

package kernel

import (
	"fmt"
	"testing"
)

func assertParseRelease(t *testing.T, release string, b *VersionInfo, result int) ***REMOVED***
	var (
		a *VersionInfo
	)
	a, _ = ParseRelease(release)

	if r := CompareKernelVersion(*a, *b); r != result ***REMOVED***
		t.Fatalf("Unexpected kernel version comparison result for (%v,%v). Found %d, expected %d", release, b, r, result)
	***REMOVED***
	if a.Flavor != b.Flavor ***REMOVED***
		t.Fatalf("Unexpected parsed kernel flavor.  Found %s, expected %s", a.Flavor, b.Flavor)
	***REMOVED***
***REMOVED***

// TestParseRelease tests the ParseRelease() function
func TestParseRelease(t *testing.T) ***REMOVED***
	assertParseRelease(t, "3.8.0", &VersionInfo***REMOVED***Kernel: 3, Major: 8, Minor: 0***REMOVED***, 0)
	assertParseRelease(t, "3.4.54.longterm-1", &VersionInfo***REMOVED***Kernel: 3, Major: 4, Minor: 54, Flavor: ".longterm-1"***REMOVED***, 0)
	assertParseRelease(t, "3.4.54.longterm-1", &VersionInfo***REMOVED***Kernel: 3, Major: 4, Minor: 54, Flavor: ".longterm-1"***REMOVED***, 0)
	assertParseRelease(t, "3.8.0-19-generic", &VersionInfo***REMOVED***Kernel: 3, Major: 8, Minor: 0, Flavor: "-19-generic"***REMOVED***, 0)
	assertParseRelease(t, "3.12.8tag", &VersionInfo***REMOVED***Kernel: 3, Major: 12, Minor: 8, Flavor: "tag"***REMOVED***, 0)
	assertParseRelease(t, "3.12-1-amd64", &VersionInfo***REMOVED***Kernel: 3, Major: 12, Minor: 0, Flavor: "-1-amd64"***REMOVED***, 0)
	assertParseRelease(t, "3.8.0", &VersionInfo***REMOVED***Kernel: 4, Major: 8, Minor: 0***REMOVED***, -1)
	// Errors
	invalids := []string***REMOVED***
		"3",
		"a",
		"a.a",
		"a.a.a-a",
	***REMOVED***
	for _, invalid := range invalids ***REMOVED***
		expectedMessage := fmt.Sprintf("Can't parse kernel version %v", invalid)
		if _, err := ParseRelease(invalid); err == nil || err.Error() != expectedMessage ***REMOVED***

		***REMOVED***
	***REMOVED***
***REMOVED***

func assertKernelVersion(t *testing.T, a, b VersionInfo, result int) ***REMOVED***
	if r := CompareKernelVersion(a, b); r != result ***REMOVED***
		t.Fatalf("Unexpected kernel version comparison result. Found %d, expected %d", r, result)
	***REMOVED***
***REMOVED***

// TestCompareKernelVersion tests the CompareKernelVersion() function
func TestCompareKernelVersion(t *testing.T) ***REMOVED***
	assertKernelVersion(t,
		VersionInfo***REMOVED***Kernel: 3, Major: 8, Minor: 0***REMOVED***,
		VersionInfo***REMOVED***Kernel: 3, Major: 8, Minor: 0***REMOVED***,
		0)
	assertKernelVersion(t,
		VersionInfo***REMOVED***Kernel: 2, Major: 6, Minor: 0***REMOVED***,
		VersionInfo***REMOVED***Kernel: 3, Major: 8, Minor: 0***REMOVED***,
		-1)
	assertKernelVersion(t,
		VersionInfo***REMOVED***Kernel: 3, Major: 8, Minor: 0***REMOVED***,
		VersionInfo***REMOVED***Kernel: 2, Major: 6, Minor: 0***REMOVED***,
		1)
	assertKernelVersion(t,
		VersionInfo***REMOVED***Kernel: 3, Major: 8, Minor: 0***REMOVED***,
		VersionInfo***REMOVED***Kernel: 3, Major: 8, Minor: 0***REMOVED***,
		0)
	assertKernelVersion(t,
		VersionInfo***REMOVED***Kernel: 3, Major: 8, Minor: 5***REMOVED***,
		VersionInfo***REMOVED***Kernel: 3, Major: 8, Minor: 0***REMOVED***,
		1)
	assertKernelVersion(t,
		VersionInfo***REMOVED***Kernel: 3, Major: 0, Minor: 20***REMOVED***,
		VersionInfo***REMOVED***Kernel: 3, Major: 8, Minor: 0***REMOVED***,
		-1)
	assertKernelVersion(t,
		VersionInfo***REMOVED***Kernel: 3, Major: 7, Minor: 20***REMOVED***,
		VersionInfo***REMOVED***Kernel: 3, Major: 8, Minor: 0***REMOVED***,
		-1)
	assertKernelVersion(t,
		VersionInfo***REMOVED***Kernel: 3, Major: 8, Minor: 20***REMOVED***,
		VersionInfo***REMOVED***Kernel: 3, Major: 7, Minor: 0***REMOVED***,
		1)
	assertKernelVersion(t,
		VersionInfo***REMOVED***Kernel: 3, Major: 8, Minor: 20***REMOVED***,
		VersionInfo***REMOVED***Kernel: 3, Major: 8, Minor: 0***REMOVED***,
		1)
	assertKernelVersion(t,
		VersionInfo***REMOVED***Kernel: 3, Major: 8, Minor: 0***REMOVED***,
		VersionInfo***REMOVED***Kernel: 3, Major: 8, Minor: 20***REMOVED***,
		-1)
***REMOVED***
