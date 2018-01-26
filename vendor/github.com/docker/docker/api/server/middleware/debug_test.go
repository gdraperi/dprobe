package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaskSecretKeys(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		path     string
		input    map[string]interface***REMOVED******REMOVED***
		expected map[string]interface***REMOVED******REMOVED***
	***REMOVED******REMOVED***
		***REMOVED***
			path:     "/v1.30/secrets/create",
			input:    map[string]interface***REMOVED******REMOVED******REMOVED***"Data": "foo", "Name": "name", "Labels": map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED******REMOVED***,
			expected: map[string]interface***REMOVED******REMOVED******REMOVED***"Data": "*****", "Name": "name", "Labels": map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED******REMOVED***,
		***REMOVED***,
		***REMOVED***
			path:     "/v1.30/secrets/create//",
			input:    map[string]interface***REMOVED******REMOVED******REMOVED***"Data": "foo", "Name": "name", "Labels": map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED******REMOVED***,
			expected: map[string]interface***REMOVED******REMOVED******REMOVED***"Data": "*****", "Name": "name", "Labels": map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED******REMOVED***,
		***REMOVED***,

		***REMOVED***
			path:     "/secrets/create?key=val",
			input:    map[string]interface***REMOVED******REMOVED******REMOVED***"Data": "foo", "Name": "name", "Labels": map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED******REMOVED***,
			expected: map[string]interface***REMOVED******REMOVED******REMOVED***"Data": "*****", "Name": "name", "Labels": map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED******REMOVED***,
		***REMOVED***,
		***REMOVED***
			path: "/v1.30/some/other/path",
			input: map[string]interface***REMOVED******REMOVED******REMOVED***
				"password": "pass",
				"other": map[string]interface***REMOVED******REMOVED******REMOVED***
					"secret":       "secret",
					"jointoken":    "jointoken",
					"unlockkey":    "unlockkey",
					"signingcakey": "signingcakey",
				***REMOVED***,
			***REMOVED***,
			expected: map[string]interface***REMOVED******REMOVED******REMOVED***
				"password": "*****",
				"other": map[string]interface***REMOVED******REMOVED******REMOVED***
					"secret":       "*****",
					"jointoken":    "*****",
					"unlockkey":    "*****",
					"signingcakey": "*****",
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for _, testcase := range tests ***REMOVED***
		maskSecretKeys(testcase.input, testcase.path)
		assert.Equal(t, testcase.expected, testcase.input)
	***REMOVED***
***REMOVED***
