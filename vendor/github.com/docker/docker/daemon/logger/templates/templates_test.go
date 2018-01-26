package templates

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewParse(t *testing.T) ***REMOVED***
	tm, err := NewParse("foo", "this is a ***REMOVED******REMOVED*** . ***REMOVED******REMOVED***")
	assert.NoError(t, err)

	var b bytes.Buffer
	assert.NoError(t, tm.Execute(&b, "string"))
	want := "this is a string"
	assert.Equal(t, want, b.String())
***REMOVED***
