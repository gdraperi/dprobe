package streamformatter

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStreamWriterStdout(t *testing.T) ***REMOVED***
	buffer := &bytes.Buffer***REMOVED******REMOVED***
	content := "content"
	sw := NewStdoutWriter(buffer)
	size, err := sw.Write([]byte(content))

	require.NoError(t, err)
	assert.Equal(t, len(content), size)

	expected := `***REMOVED***"stream":"content"***REMOVED***` + streamNewline
	assert.Equal(t, expected, buffer.String())
***REMOVED***

func TestStreamWriterStderr(t *testing.T) ***REMOVED***
	buffer := &bytes.Buffer***REMOVED******REMOVED***
	content := "content"
	sw := NewStderrWriter(buffer)
	size, err := sw.Write([]byte(content))

	require.NoError(t, err)
	assert.Equal(t, len(content), size)

	expected := `***REMOVED***"stream":"\u001b[91mcontent\u001b[0m"***REMOVED***` + streamNewline
	assert.Equal(t, expected, buffer.String())
***REMOVED***
