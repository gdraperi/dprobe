package testutil

import (
	"io"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ErrorContains checks that the error is not nil, and contains the expected
// substring.
func ErrorContains(t require.TestingT, err error, expectedError string, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	require.Error(t, err, msgAndArgs...)
	assert.Contains(t, err.Error(), expectedError, msgAndArgs...)
***REMOVED***

// DevZero acts like /dev/zero but in an OS-independent fashion.
var DevZero io.Reader = devZero***REMOVED******REMOVED***

type devZero struct***REMOVED******REMOVED***

func (d devZero) Read(p []byte) (n int, err error) ***REMOVED***
	for i := range p ***REMOVED***
		p[i] = 0
	***REMOVED***
	return len(p), nil
***REMOVED***
