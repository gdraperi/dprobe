package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testLengthHelper(generator func(int) string, t *testing.T) ***REMOVED***
	expectedLength := 20
	s := generator(expectedLength)
	assert.Equal(t, expectedLength, len(s))
***REMOVED***

func testUniquenessHelper(generator func(int) string, t *testing.T) ***REMOVED***
	repeats := 25
	set := make(map[string]struct***REMOVED******REMOVED***, repeats)
	for i := 0; i < repeats; i = i + 1 ***REMOVED***
		str := generator(64)
		assert.Equal(t, 64, len(str))
		_, ok := set[str]
		assert.False(t, ok, "Random number is repeated")
		set[str] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
***REMOVED***

func TestGenerateRandomAlphaOnlyStringLength(t *testing.T) ***REMOVED***
	testLengthHelper(GenerateRandomAlphaOnlyString, t)
***REMOVED***

func TestGenerateRandomAlphaOnlyStringUniqueness(t *testing.T) ***REMOVED***
	testUniquenessHelper(GenerateRandomAlphaOnlyString, t)
***REMOVED***
