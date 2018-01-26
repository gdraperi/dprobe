/*Package icmd executes binaries and provides convenient assertions for testing the results.
 */
package icmd

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type testingT interface ***REMOVED***
	Fatalf(string, ...interface***REMOVED******REMOVED***)
***REMOVED***

type helperT interface ***REMOVED***
	Helper()
***REMOVED***

// None is a token to inform Result.Assert that the output should be empty
const None = "[NOTHING]"

type lockedBuffer struct ***REMOVED***
	m   sync.RWMutex
	buf bytes.Buffer
***REMOVED***

func (buf *lockedBuffer) Write(b []byte) (int, error) ***REMOVED***
	buf.m.Lock()
	defer buf.m.Unlock()
	return buf.buf.Write(b)
***REMOVED***

func (buf *lockedBuffer) String() string ***REMOVED***
	buf.m.RLock()
	defer buf.m.RUnlock()
	return buf.buf.String()
***REMOVED***

// Result stores the result of running a command
type Result struct ***REMOVED***
	Cmd      *exec.Cmd
	ExitCode int
	Error    error
	// Timeout is true if the command was killed because it ran for too long
	Timeout   bool
	outBuffer *lockedBuffer
	errBuffer *lockedBuffer
***REMOVED***

// Assert compares the Result against the Expected struct, and fails the test if
// any of the expectations are not met.
// TODO: deprecate and replace with assert.CompareFunc
func (r *Result) Assert(t testingT, exp Expected) *Result ***REMOVED***
	if ht, ok := t.(helperT); ok ***REMOVED***
		ht.Helper()
	***REMOVED***
	err := r.Compare(exp)
	if err == nil ***REMOVED***
		return r
	***REMOVED***
	t.Fatalf(err.Error() + "\n")
	return nil
***REMOVED***

// Compare returns a formatted error with the command, stdout, stderr, exit
// code, and any failed expectations
// nolint: gocyclo
func (r *Result) Compare(exp Expected) error ***REMOVED***
	errors := []string***REMOVED******REMOVED***
	add := func(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
		errors = append(errors, fmt.Sprintf(format, args...))
	***REMOVED***

	if exp.ExitCode != r.ExitCode ***REMOVED***
		add("ExitCode was %d expected %d", r.ExitCode, exp.ExitCode)
	***REMOVED***
	if exp.Timeout != r.Timeout ***REMOVED***
		if exp.Timeout ***REMOVED***
			add("Expected command to timeout")
		***REMOVED*** else ***REMOVED***
			add("Expected command to finish, but it hit the timeout")
		***REMOVED***
	***REMOVED***
	if !matchOutput(exp.Out, r.Stdout()) ***REMOVED***
		add("Expected stdout to contain %q", exp.Out)
	***REMOVED***
	if !matchOutput(exp.Err, r.Stderr()) ***REMOVED***
		add("Expected stderr to contain %q", exp.Err)
	***REMOVED***
	switch ***REMOVED***
	// If a non-zero exit code is expected there is going to be an error.
	// Don't require an error message as well as an exit code because the
	// error message is going to be "exit status <code> which is not useful
	case exp.Error == "" && exp.ExitCode != 0:
	case exp.Error == "" && r.Error != nil:
		add("Expected no error")
	case exp.Error != "" && r.Error == nil:
		add("Expected error to contain %q, but there was no error", exp.Error)
	case exp.Error != "" && !strings.Contains(r.Error.Error(), exp.Error):
		add("Expected error to contain %q", exp.Error)
	***REMOVED***

	if len(errors) == 0 ***REMOVED***
		return nil
	***REMOVED***
	return fmt.Errorf("%s\nFailures:\n%s", r, strings.Join(errors, "\n"))
***REMOVED***

func matchOutput(expected string, actual string) bool ***REMOVED***
	switch expected ***REMOVED***
	case None:
		return actual == ""
	default:
		return strings.Contains(actual, expected)
	***REMOVED***
***REMOVED***

func (r *Result) String() string ***REMOVED***
	var timeout string
	if r.Timeout ***REMOVED***
		timeout = " (timeout)"
	***REMOVED***

	return fmt.Sprintf(`
Command:  %s
ExitCode: %d%s
Error:    %v
Stdout:   %v
Stderr:   %v
`,
		strings.Join(r.Cmd.Args, " "),
		r.ExitCode,
		timeout,
		r.Error,
		r.Stdout(),
		r.Stderr())
***REMOVED***

// Expected is the expected output from a Command. This struct is compared to a
// Result struct by Result.Assert().
type Expected struct ***REMOVED***
	ExitCode int
	Timeout  bool
	Error    string
	Out      string
	Err      string
***REMOVED***

// Success is the default expected result. A Success result is one with a 0
// ExitCode.
var Success = Expected***REMOVED******REMOVED***

// Stdout returns the stdout of the process as a string
func (r *Result) Stdout() string ***REMOVED***
	return r.outBuffer.String()
***REMOVED***

// Stderr returns the stderr of the process as a string
func (r *Result) Stderr() string ***REMOVED***
	return r.errBuffer.String()
***REMOVED***

// Combined returns the stdout and stderr combined into a single string
func (r *Result) Combined() string ***REMOVED***
	return r.outBuffer.String() + r.errBuffer.String()
***REMOVED***

func (r *Result) setExitError(err error) ***REMOVED***
	if err == nil ***REMOVED***
		return
	***REMOVED***
	r.Error = err
	r.ExitCode = processExitCode(err)
***REMOVED***

// Cmd contains the arguments and options for a process to run as part of a test
// suite.
type Cmd struct ***REMOVED***
	Command []string
	Timeout time.Duration
	Stdin   io.Reader
	Stdout  io.Writer
	Dir     string
	Env     []string
***REMOVED***

// Command create a simple Cmd with the specified command and arguments
func Command(command string, args ...string) Cmd ***REMOVED***
	return Cmd***REMOVED***Command: append([]string***REMOVED***command***REMOVED***, args...)***REMOVED***
***REMOVED***

// RunCmd runs a command and returns a Result
func RunCmd(cmd Cmd, cmdOperators ...CmdOp) *Result ***REMOVED***
	for _, op := range cmdOperators ***REMOVED***
		op(&cmd)
	***REMOVED***
	result := StartCmd(cmd)
	if result.Error != nil ***REMOVED***
		return result
	***REMOVED***
	return WaitOnCmd(cmd.Timeout, result)
***REMOVED***

// RunCommand runs a command with default options, and returns a result
func RunCommand(command string, args ...string) *Result ***REMOVED***
	return RunCmd(Command(command, args...))
***REMOVED***

// StartCmd starts a command, but doesn't wait for it to finish
func StartCmd(cmd Cmd) *Result ***REMOVED***
	result := buildCmd(cmd)
	if result.Error != nil ***REMOVED***
		return result
	***REMOVED***
	result.setExitError(result.Cmd.Start())
	return result
***REMOVED***

func buildCmd(cmd Cmd) *Result ***REMOVED***
	var execCmd *exec.Cmd
	switch len(cmd.Command) ***REMOVED***
	case 1:
		execCmd = exec.Command(cmd.Command[0])
	default:
		execCmd = exec.Command(cmd.Command[0], cmd.Command[1:]...)
	***REMOVED***
	outBuffer := new(lockedBuffer)
	errBuffer := new(lockedBuffer)

	execCmd.Stdin = cmd.Stdin
	execCmd.Dir = cmd.Dir
	execCmd.Env = cmd.Env
	if cmd.Stdout != nil ***REMOVED***
		execCmd.Stdout = io.MultiWriter(outBuffer, cmd.Stdout)
	***REMOVED*** else ***REMOVED***
		execCmd.Stdout = outBuffer
	***REMOVED***
	execCmd.Stderr = errBuffer
	return &Result***REMOVED***
		Cmd:       execCmd,
		outBuffer: outBuffer,
		errBuffer: errBuffer,
	***REMOVED***
***REMOVED***

// WaitOnCmd waits for a command to complete. If timeout is non-nil then
// only wait until the timeout.
func WaitOnCmd(timeout time.Duration, result *Result) *Result ***REMOVED***
	if timeout == time.Duration(0) ***REMOVED***
		result.setExitError(result.Cmd.Wait())
		return result
	***REMOVED***

	done := make(chan error, 1)
	// Wait for command to exit in a goroutine
	go func() ***REMOVED***
		done <- result.Cmd.Wait()
	***REMOVED***()

	select ***REMOVED***
	case <-time.After(timeout):
		killErr := result.Cmd.Process.Kill()
		if killErr != nil ***REMOVED***
			fmt.Printf("failed to kill (pid=%d): %v\n", result.Cmd.Process.Pid, killErr)
		***REMOVED***
		result.Timeout = true
	case err := <-done:
		result.setExitError(err)
	***REMOVED***
	return result
***REMOVED***
