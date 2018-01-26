// +build linux

package signal

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildTestBinary(t *testing.T, tmpdir string, prefix string) (string, string) ***REMOVED***
	tmpDir, err := ioutil.TempDir(tmpdir, prefix)
	require.NoError(t, err)
	exePath := tmpDir + "/" + prefix
	wd, _ := os.Getwd()
	testHelperCode := wd + "/testfiles/main.go"
	cmd := exec.Command("go", "build", "-o", exePath, testHelperCode)
	err = cmd.Run()
	require.NoError(t, err)
	return exePath, tmpDir
***REMOVED***

func TestTrap(t *testing.T) ***REMOVED***
	var sigmap = []struct ***REMOVED***
		name     string
		signal   os.Signal
		multiple bool
	***REMOVED******REMOVED***
		***REMOVED***"TERM", syscall.SIGTERM, false***REMOVED***,
		***REMOVED***"QUIT", syscall.SIGQUIT, true***REMOVED***,
		***REMOVED***"INT", os.Interrupt, false***REMOVED***,
		***REMOVED***"TERM", syscall.SIGTERM, true***REMOVED***,
		***REMOVED***"INT", os.Interrupt, true***REMOVED***,
	***REMOVED***
	exePath, tmpDir := buildTestBinary(t, "", "main")
	defer os.RemoveAll(tmpDir)

	for _, v := range sigmap ***REMOVED***
		cmd := exec.Command(exePath)
		cmd.Env = append(os.Environ(), fmt.Sprintf("SIGNAL_TYPE=%s", v.name))
		if v.multiple ***REMOVED***
			cmd.Env = append(cmd.Env, "IF_MULTIPLE=1")
		***REMOVED***
		err := cmd.Start()
		require.NoError(t, err)
		err = cmd.Wait()
		if e, ok := err.(*exec.ExitError); ok ***REMOVED***
			code := e.Sys().(syscall.WaitStatus).ExitStatus()
			if v.multiple ***REMOVED***
				assert.Equal(t, 128+int(v.signal.(syscall.Signal)), code)
			***REMOVED*** else ***REMOVED***
				assert.Equal(t, 99, code)
			***REMOVED***
			continue
		***REMOVED***
		t.Fatal("process didn't end with any error")
	***REMOVED***

***REMOVED***

func TestDumpStacks(t *testing.T) ***REMOVED***
	directory, err := ioutil.TempDir("", "test-dump-tasks")
	assert.NoError(t, err)
	defer os.RemoveAll(directory)
	dumpPath, err := DumpStacks(directory)
	assert.NoError(t, err)
	readFile, _ := ioutil.ReadFile(dumpPath)
	fileData := string(readFile)
	assert.Contains(t, fileData, "goroutine")
***REMOVED***

func TestDumpStacksWithEmptyInput(t *testing.T) ***REMOVED***
	path, err := DumpStacks("")
	assert.NoError(t, err)
	assert.Equal(t, os.Stderr.Name(), path)
***REMOVED***
