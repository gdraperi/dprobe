// Testing support for go-toml

package toml

import (
	"testing"
)

func TestPositionString(t *testing.T) ***REMOVED***
	p := Position***REMOVED***123, 456***REMOVED***
	expected := "(123, 456)"
	value := p.String()

	if value != expected ***REMOVED***
		t.Errorf("Expected %v, got %v instead", expected, value)
	***REMOVED***
***REMOVED***

func TestInvalid(t *testing.T) ***REMOVED***
	for i, v := range []Position***REMOVED***
		***REMOVED***0, 1234***REMOVED***,
		***REMOVED***1234, 0***REMOVED***,
		***REMOVED***0, 0***REMOVED***,
	***REMOVED*** ***REMOVED***
		if !v.Invalid() ***REMOVED***
			t.Errorf("Position at %v is valid: %v", i, v)
		***REMOVED***
	***REMOVED***
***REMOVED***
