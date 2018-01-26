package signal

import (
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSignal(t *testing.T) ***REMOVED***
	_, checkAtoiError := ParseSignal("0")
	assert.EqualError(t, checkAtoiError, "Invalid signal: 0")

	_, error := ParseSignal("SIG")
	assert.EqualError(t, error, "Invalid signal: SIG")

	for sigStr := range SignalMap ***REMOVED***
		responseSignal, error := ParseSignal(sigStr)
		assert.NoError(t, error)
		signal := SignalMap[sigStr]
		assert.EqualValues(t, signal, responseSignal)
	***REMOVED***
***REMOVED***

func TestValidSignalForPlatform(t *testing.T) ***REMOVED***
	isValidSignal := ValidSignalForPlatform(syscall.Signal(0))
	assert.EqualValues(t, false, isValidSignal)

	for _, sigN := range SignalMap ***REMOVED***
		isValidSignal = ValidSignalForPlatform(syscall.Signal(sigN))
		assert.EqualValues(t, true, isValidSignal)
	***REMOVED***
***REMOVED***
