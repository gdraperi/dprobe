// Package signal provides helper functions for dealing with signals across
// various operating systems.
package signal

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

// CatchAll catches all signals and relays them to the specified channel.
func CatchAll(sigc chan os.Signal) ***REMOVED***
	handledSigs := []os.Signal***REMOVED******REMOVED***
	for _, s := range SignalMap ***REMOVED***
		handledSigs = append(handledSigs, s)
	***REMOVED***
	signal.Notify(sigc, handledSigs...)
***REMOVED***

// StopCatch stops catching the signals and closes the specified channel.
func StopCatch(sigc chan os.Signal) ***REMOVED***
	signal.Stop(sigc)
	close(sigc)
***REMOVED***

// ParseSignal translates a string to a valid syscall signal.
// It returns an error if the signal map doesn't include the given signal.
func ParseSignal(rawSignal string) (syscall.Signal, error) ***REMOVED***
	s, err := strconv.Atoi(rawSignal)
	if err == nil ***REMOVED***
		if s == 0 ***REMOVED***
			return -1, fmt.Errorf("Invalid signal: %s", rawSignal)
		***REMOVED***
		return syscall.Signal(s), nil
	***REMOVED***
	signal, ok := SignalMap[strings.TrimPrefix(strings.ToUpper(rawSignal), "SIG")]
	if !ok ***REMOVED***
		return -1, fmt.Errorf("Invalid signal: %s", rawSignal)
	***REMOVED***
	return signal, nil
***REMOVED***

// ValidSignalForPlatform returns true if a signal is valid on the platform
func ValidSignalForPlatform(sig syscall.Signal) bool ***REMOVED***
	for _, v := range SignalMap ***REMOVED***
		if v == sig ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
