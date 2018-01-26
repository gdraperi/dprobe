// +build !windows

package main

import (
	"strings"

	"github.com/go-check/check"
)

// Test case for #30166 (target was not validated)
func (s *DockerSuite) TestCreateTmpfsMountsTarget(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	type testCase struct ***REMOVED***
		target        string
		expectedError string
	***REMOVED***
	cases := []testCase***REMOVED***
		***REMOVED***
			target:        ".",
			expectedError: "mount path must be absolute",
		***REMOVED***,
		***REMOVED***
			target:        "foo",
			expectedError: "mount path must be absolute",
		***REMOVED***,
		***REMOVED***
			target:        "/",
			expectedError: "destination can't be '/'",
		***REMOVED***,
		***REMOVED***
			target:        "//",
			expectedError: "destination can't be '/'",
		***REMOVED***,
	***REMOVED***
	for _, x := range cases ***REMOVED***
		out, _, _ := dockerCmdWithError("create", "--tmpfs", x.target, "busybox", "sh")
		if x.expectedError != "" && !strings.Contains(out, x.expectedError) ***REMOVED***
			c.Fatalf("mounting tmpfs over %q should fail with %q, but got %q",
				x.target, x.expectedError, out)
		***REMOVED***
	***REMOVED***
***REMOVED***
