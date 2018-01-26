package reexec

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() ***REMOVED***
	Register("reexec", func() ***REMOVED***
		panic("Return Error")
	***REMOVED***)
	Init()
***REMOVED***

func TestRegister(t *testing.T) ***REMOVED***
	defer func() ***REMOVED***
		if r := recover(); r != nil ***REMOVED***
			require.Equal(t, `reexec func already registered under name "reexec"`, r)
		***REMOVED***
	***REMOVED***()
	Register("reexec", func() ***REMOVED******REMOVED***)
***REMOVED***

func TestCommand(t *testing.T) ***REMOVED***
	cmd := Command("reexec")
	w, err := cmd.StdinPipe()
	require.NoError(t, err, "Error on pipe creation: %v", err)
	defer w.Close()

	err = cmd.Start()
	require.NoError(t, err, "Error on re-exec cmd: %v", err)
	err = cmd.Wait()
	require.EqualError(t, err, "exit status 2")
***REMOVED***

func TestNaiveSelf(t *testing.T) ***REMOVED***
	if os.Getenv("TEST_CHECK") == "1" ***REMOVED***
		os.Exit(2)
	***REMOVED***
	cmd := exec.Command(naiveSelf(), "-test.run=TestNaiveSelf")
	cmd.Env = append(os.Environ(), "TEST_CHECK=1")
	err := cmd.Start()
	require.NoError(t, err, "Unable to start command")
	err = cmd.Wait()
	require.EqualError(t, err, "exit status 2")

	os.Args[0] = "mkdir"
	assert.NotEqual(t, naiveSelf(), os.Args[0])
***REMOVED***
