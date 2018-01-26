package opts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuotedStringSetWithQuotes(t *testing.T) ***REMOVED***
	value := ""
	qs := NewQuotedString(&value)
	assert.NoError(t, qs.Set(`"something"`))
	assert.Equal(t, "something", qs.String())
	assert.Equal(t, "something", value)
***REMOVED***

func TestQuotedStringSetWithMismatchedQuotes(t *testing.T) ***REMOVED***
	value := ""
	qs := NewQuotedString(&value)
	assert.NoError(t, qs.Set(`"something'`))
	assert.Equal(t, `"something'`, qs.String())
***REMOVED***

func TestQuotedStringSetWithNoQuotes(t *testing.T) ***REMOVED***
	value := ""
	qs := NewQuotedString(&value)
	assert.NoError(t, qs.Set("something"))
	assert.Equal(t, "something", qs.String())
***REMOVED***
