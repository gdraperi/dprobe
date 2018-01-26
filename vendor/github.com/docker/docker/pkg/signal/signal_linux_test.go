// +build darwin linux

package signal

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCatchAll(t *testing.T) ***REMOVED***
	sigs := make(chan os.Signal, 1)
	CatchAll(sigs)
	defer StopCatch(sigs)

	listOfSignals := map[string]string***REMOVED***
		"CONT": syscall.SIGCONT.String(),
		"HUP":  syscall.SIGHUP.String(),
		"CHLD": syscall.SIGCHLD.String(),
		"ILL":  syscall.SIGILL.String(),
		"FPE":  syscall.SIGFPE.String(),
		"CLD":  syscall.SIGCLD.String(),
	***REMOVED***

	for sigStr := range listOfSignals ***REMOVED***
		signal, ok := SignalMap[sigStr]
		if ok ***REMOVED***
			go func() ***REMOVED***
				time.Sleep(1 * time.Millisecond)
				syscall.Kill(syscall.Getpid(), signal)
			***REMOVED***()

			s := <-sigs
			assert.EqualValues(t, s.String(), signal.String())
		***REMOVED***

	***REMOVED***
***REMOVED***

func TestStopCatch(t *testing.T) ***REMOVED***
	signal := SignalMap["HUP"]
	channel := make(chan os.Signal, 1)
	CatchAll(channel)
	go func() ***REMOVED***

		time.Sleep(1 * time.Millisecond)
		syscall.Kill(syscall.Getpid(), signal)
	***REMOVED***()
	signalString := <-channel
	assert.EqualValues(t, signalString.String(), signal.String())

	StopCatch(channel)
	_, ok := <-channel
	assert.EqualValues(t, ok, false)
***REMOVED***
