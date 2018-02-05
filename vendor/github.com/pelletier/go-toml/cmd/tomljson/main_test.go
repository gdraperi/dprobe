package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func expectBufferEquality(t *testing.T, name string, buffer *bytes.Buffer, expected string) ***REMOVED***
	output := buffer.String()
	if output != expected ***REMOVED***
		t.Errorf("incorrect %s:\n%s\n\nexpected %s:\n%s", name, output, name, expected)
		t.Log([]rune(output))
		t.Log([]rune(expected))
	***REMOVED***
***REMOVED***

func expectProcessMainResults(t *testing.T, input string, args []string, exitCode int, expectedOutput string, expectedError string) ***REMOVED***
	inputReader := strings.NewReader(input)
	outputBuffer := new(bytes.Buffer)
	errorBuffer := new(bytes.Buffer)

	returnCode := processMain(args, inputReader, outputBuffer, errorBuffer)

	expectBufferEquality(t, "output", outputBuffer, expectedOutput)
	expectBufferEquality(t, "error", errorBuffer, expectedError)

	if returnCode != exitCode ***REMOVED***
		t.Error("incorrect return code:", returnCode, "expected", exitCode)
	***REMOVED***
***REMOVED***

func TestProcessMainReadFromStdin(t *testing.T) ***REMOVED***
	input := `
		[mytoml]
		a = 42`
	expectedOutput := `***REMOVED***
  "mytoml": ***REMOVED***
    "a": 42
  ***REMOVED***
***REMOVED***
`
	expectedError := ``
	expectedExitCode := 0

	expectProcessMainResults(t, input, []string***REMOVED******REMOVED***, expectedExitCode, expectedOutput, expectedError)
***REMOVED***

func TestProcessMainReadFromFile(t *testing.T) ***REMOVED***
	input := `
		[mytoml]
		a = 42`

	tmpfile, err := ioutil.TempFile("", "example.toml")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if _, err := tmpfile.Write([]byte(input)); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	defer os.Remove(tmpfile.Name())

	expectedOutput := `***REMOVED***
  "mytoml": ***REMOVED***
    "a": 42
  ***REMOVED***
***REMOVED***
`
	expectedError := ``
	expectedExitCode := 0

	expectProcessMainResults(t, ``, []string***REMOVED***tmpfile.Name()***REMOVED***, expectedExitCode, expectedOutput, expectedError)
***REMOVED***

func TestProcessMainReadFromMissingFile(t *testing.T) ***REMOVED***
	expectedError := `open /this/file/does/not/exist: no such file or directory
`
	expectProcessMainResults(t, ``, []string***REMOVED***"/this/file/does/not/exist"***REMOVED***, -1, ``, expectedError)
***REMOVED***
