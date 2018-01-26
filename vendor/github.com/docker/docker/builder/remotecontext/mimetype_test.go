package remotecontext

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectContentType(t *testing.T) ***REMOVED***
	input := []byte("That is just a plain text")

	contentType, _, err := detectContentType(input)
	require.NoError(t, err)
	assert.Equal(t, "text/plain", contentType)
***REMOVED***
