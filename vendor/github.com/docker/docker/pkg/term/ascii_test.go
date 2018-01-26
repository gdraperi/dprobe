package term

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToBytes(t *testing.T) ***REMOVED***
	codes, err := ToBytes("ctrl-a,a")
	require.NoError(t, err)
	assert.Equal(t, []byte***REMOVED***1, 97***REMOVED***, codes)

	_, err = ToBytes("shift-z")
	assert.Error(t, err)

	codes, err = ToBytes("ctrl-@,ctrl-[,~,ctrl-o")
	require.NoError(t, err)
	assert.Equal(t, []byte***REMOVED***0, 27, 126, 15***REMOVED***, codes)

	codes, err = ToBytes("DEL,+")
	require.NoError(t, err)
	assert.Equal(t, []byte***REMOVED***127, 43***REMOVED***, codes)
***REMOVED***
