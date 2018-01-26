package opts

import (
	"fmt"
	"os"
	"runtime"
	"testing"
)

func TestValidateEnv(t *testing.T) ***REMOVED***
	testcase := []struct ***REMOVED***
		value    string
		expected string
		err      error
	***REMOVED******REMOVED***
		***REMOVED***
			value:    "a",
			expected: "a",
		***REMOVED***,
		***REMOVED***
			value:    "something",
			expected: "something",
		***REMOVED***,
		***REMOVED***
			value:    "_=a",
			expected: "_=a",
		***REMOVED***,
		***REMOVED***
			value:    "env1=value1",
			expected: "env1=value1",
		***REMOVED***,
		***REMOVED***
			value:    "_env1=value1",
			expected: "_env1=value1",
		***REMOVED***,
		***REMOVED***
			value:    "env2=value2=value3",
			expected: "env2=value2=value3",
		***REMOVED***,
		***REMOVED***
			value:    "env3=abc!qwe",
			expected: "env3=abc!qwe",
		***REMOVED***,
		***REMOVED***
			value:    "env_4=value 4",
			expected: "env_4=value 4",
		***REMOVED***,
		***REMOVED***
			value:    "PATH",
			expected: fmt.Sprintf("PATH=%v", os.Getenv("PATH")),
		***REMOVED***,
		***REMOVED***
			value: "=a",
			err:   fmt.Errorf(fmt.Sprintf("invalid environment variable: %s", "=a")),
		***REMOVED***,
		***REMOVED***
			value:    "PATH=something",
			expected: "PATH=something",
		***REMOVED***,
		***REMOVED***
			value:    "asd!qwe",
			expected: "asd!qwe",
		***REMOVED***,
		***REMOVED***
			value:    "1asd",
			expected: "1asd",
		***REMOVED***,
		***REMOVED***
			value:    "123",
			expected: "123",
		***REMOVED***,
		***REMOVED***
			value:    "some space",
			expected: "some space",
		***REMOVED***,
		***REMOVED***
			value:    "  some space before",
			expected: "  some space before",
		***REMOVED***,
		***REMOVED***
			value:    "some space after  ",
			expected: "some space after  ",
		***REMOVED***,
		***REMOVED***
			value: "=",
			err:   fmt.Errorf(fmt.Sprintf("invalid environment variable: %s", "=")),
		***REMOVED***,
	***REMOVED***

	// Environment variables are case in-sensitive on Windows
	if runtime.GOOS == "windows" ***REMOVED***
		tmp := struct ***REMOVED***
			value    string
			expected string
			err      error
		***REMOVED******REMOVED***
			value:    "PaTh",
			expected: fmt.Sprintf("PaTh=%v", os.Getenv("PATH")),
		***REMOVED***
		testcase = append(testcase, tmp)

	***REMOVED***

	for _, r := range testcase ***REMOVED***
		actual, err := ValidateEnv(r.value)

		if err != nil ***REMOVED***
			if r.err == nil ***REMOVED***
				t.Fatalf("Expected err is nil, got err[%v]", err)
			***REMOVED***
			if err.Error() != r.err.Error() ***REMOVED***
				t.Fatalf("Expected err[%v], got err[%v]", r.err, err)
			***REMOVED***
		***REMOVED***

		if err == nil && r.err != nil ***REMOVED***
			t.Fatalf("Expected err[%v], but err is nil", r.err)
		***REMOVED***

		if actual != r.expected ***REMOVED***
			t.Fatalf("Expected [%v], got [%v]", r.expected, actual)
		***REMOVED***
	***REMOVED***
***REMOVED***
