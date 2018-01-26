package instructions

import "testing"

type testCase struct ***REMOVED***
	name       string
	args       []string
	attributes map[string]bool
	expected   []string
***REMOVED***

func initTestCases() []testCase ***REMOVED***
	testCases := []testCase***REMOVED******REMOVED***

	testCases = append(testCases, testCase***REMOVED***
		name:       "empty args",
		args:       []string***REMOVED******REMOVED***,
		attributes: make(map[string]bool),
		expected:   []string***REMOVED******REMOVED***,
	***REMOVED***)

	jsonAttributes := make(map[string]bool)
	jsonAttributes["json"] = true

	testCases = append(testCases, testCase***REMOVED***
		name:       "json attribute with one element",
		args:       []string***REMOVED***"foo"***REMOVED***,
		attributes: jsonAttributes,
		expected:   []string***REMOVED***"foo"***REMOVED***,
	***REMOVED***)

	testCases = append(testCases, testCase***REMOVED***
		name:       "json attribute with two elements",
		args:       []string***REMOVED***"foo", "bar"***REMOVED***,
		attributes: jsonAttributes,
		expected:   []string***REMOVED***"foo", "bar"***REMOVED***,
	***REMOVED***)

	testCases = append(testCases, testCase***REMOVED***
		name:       "no attributes",
		args:       []string***REMOVED***"foo", "bar"***REMOVED***,
		attributes: nil,
		expected:   []string***REMOVED***"foo bar"***REMOVED***,
	***REMOVED***)

	return testCases
***REMOVED***

func TestHandleJSONArgs(t *testing.T) ***REMOVED***
	testCases := initTestCases()

	for _, test := range testCases ***REMOVED***
		arguments := handleJSONArgs(test.args, test.attributes)

		if len(arguments) != len(test.expected) ***REMOVED***
			t.Fatalf("In test \"%s\": length of returned slice is incorrect. Expected: %d, got: %d", test.name, len(test.expected), len(arguments))
		***REMOVED***

		for i := range test.expected ***REMOVED***
			if arguments[i] != test.expected[i] ***REMOVED***
				t.Fatalf("In test \"%s\": element as position %d is incorrect. Expected: %s, got: %s", test.name, i, test.expected[i], arguments[i])
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
