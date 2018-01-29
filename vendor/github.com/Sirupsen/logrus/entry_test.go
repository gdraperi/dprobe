package logrus

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntryWithError(t *testing.T) ***REMOVED***

	assert := assert.New(t)

	defer func() ***REMOVED***
		ErrorKey = "error"
	***REMOVED***()

	err := fmt.Errorf("kaboom at layer %d", 4711)

	assert.Equal(err, WithError(err).Data["error"])

	logger := New()
	logger.Out = &bytes.Buffer***REMOVED******REMOVED***
	entry := NewEntry(logger)

	assert.Equal(err, entry.WithError(err).Data["error"])

	ErrorKey = "err"

	assert.Equal(err, entry.WithError(err).Data["err"])

***REMOVED***

func TestEntryPanicln(t *testing.T) ***REMOVED***
	errBoom := fmt.Errorf("boom time")

	defer func() ***REMOVED***
		p := recover()
		assert.NotNil(t, p)

		switch pVal := p.(type) ***REMOVED***
		case *Entry:
			assert.Equal(t, "kaboom", pVal.Message)
			assert.Equal(t, errBoom, pVal.Data["err"])
		default:
			t.Fatalf("want type *Entry, got %T: %#v", pVal, pVal)
		***REMOVED***
	***REMOVED***()

	logger := New()
	logger.Out = &bytes.Buffer***REMOVED******REMOVED***
	entry := NewEntry(logger)
	entry.WithField("err", errBoom).Panicln("kaboom")
***REMOVED***

func TestEntryPanicf(t *testing.T) ***REMOVED***
	errBoom := fmt.Errorf("boom again")

	defer func() ***REMOVED***
		p := recover()
		assert.NotNil(t, p)

		switch pVal := p.(type) ***REMOVED***
		case *Entry:
			assert.Equal(t, "kaboom true", pVal.Message)
			assert.Equal(t, errBoom, pVal.Data["err"])
		default:
			t.Fatalf("want type *Entry, got %T: %#v", pVal, pVal)
		***REMOVED***
	***REMOVED***()

	logger := New()
	logger.Out = &bytes.Buffer***REMOVED******REMOVED***
	entry := NewEntry(logger)
	entry.WithField("err", errBoom).Panicf("kaboom %v", true)
***REMOVED***
