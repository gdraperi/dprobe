package cli

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/integration-cli/daemon"
	"github.com/docker/docker/integration-cli/environment"
	"github.com/gotestyourself/gotestyourself/icmd"
	"github.com/pkg/errors"
)

var testEnv *environment.Execution

// SetTestEnvironment sets a static test environment
// TODO: decouple this package from environment
func SetTestEnvironment(env *environment.Execution) ***REMOVED***
	testEnv = env
***REMOVED***

// CmdOperator defines functions that can modify a command
type CmdOperator func(*icmd.Cmd) func()

type testingT interface ***REMOVED***
	Fatal(args ...interface***REMOVED******REMOVED***)
	Fatalf(string, ...interface***REMOVED******REMOVED***)
***REMOVED***

// DockerCmd executes the specified docker command and expect a success
func DockerCmd(t testingT, args ...string) *icmd.Result ***REMOVED***
	return Docker(Args(args...)).Assert(t, icmd.Success)
***REMOVED***

// BuildCmd executes the specified docker build command and expect a success
func BuildCmd(t testingT, name string, cmdOperators ...CmdOperator) *icmd.Result ***REMOVED***
	return Docker(Build(name), cmdOperators...).Assert(t, icmd.Success)
***REMOVED***

// InspectCmd executes the specified docker inspect command and expect a success
func InspectCmd(t testingT, name string, cmdOperators ...CmdOperator) *icmd.Result ***REMOVED***
	return Docker(Inspect(name), cmdOperators...).Assert(t, icmd.Success)
***REMOVED***

// WaitRun will wait for the specified container to be running, maximum 5 seconds.
func WaitRun(t testingT, name string, cmdOperators ...CmdOperator) ***REMOVED***
	WaitForInspectResult(t, name, "***REMOVED******REMOVED***.State.Running***REMOVED******REMOVED***", "true", 5*time.Second, cmdOperators...)
***REMOVED***

// WaitExited will wait for the specified container to state exit, subject
// to a maximum time limit in seconds supplied by the caller
func WaitExited(t testingT, name string, timeout time.Duration, cmdOperators ...CmdOperator) ***REMOVED***
	WaitForInspectResult(t, name, "***REMOVED******REMOVED***.State.Status***REMOVED******REMOVED***", "exited", timeout, cmdOperators...)
***REMOVED***

// WaitRestart will wait for the specified container to restart once
func WaitRestart(t testingT, name string, timeout time.Duration, cmdOperators ...CmdOperator) ***REMOVED***
	WaitForInspectResult(t, name, "***REMOVED******REMOVED***.RestartCount***REMOVED******REMOVED***", "1", timeout, cmdOperators...)
***REMOVED***

// WaitForInspectResult waits for the specified expression to be equals to the specified expected string in the given time.
func WaitForInspectResult(t testingT, name, expr, expected string, timeout time.Duration, cmdOperators ...CmdOperator) ***REMOVED***
	after := time.After(timeout)

	args := []string***REMOVED***"inspect", "-f", expr, name***REMOVED***
	for ***REMOVED***
		result := Docker(Args(args...), cmdOperators...)
		if result.Error != nil ***REMOVED***
			if !strings.Contains(strings.ToLower(result.Stderr()), "no such") ***REMOVED***
				t.Fatalf("error executing docker inspect: %v\n%s",
					result.Stderr(), result.Stdout())
			***REMOVED***
			select ***REMOVED***
			case <-after:
				t.Fatal(result.Error)
			default:
				time.Sleep(10 * time.Millisecond)
				continue
			***REMOVED***
		***REMOVED***

		out := strings.TrimSpace(result.Stdout())
		if out == expected ***REMOVED***
			break
		***REMOVED***

		select ***REMOVED***
		case <-after:
			t.Fatalf("condition \"%q == %q\" not true in time (%v)", out, expected, timeout)
		default:
		***REMOVED***

		time.Sleep(100 * time.Millisecond)
	***REMOVED***
***REMOVED***

// Docker executes the specified docker command
func Docker(cmd icmd.Cmd, cmdOperators ...CmdOperator) *icmd.Result ***REMOVED***
	for _, op := range cmdOperators ***REMOVED***
		deferFn := op(&cmd)
		if deferFn != nil ***REMOVED***
			defer deferFn()
		***REMOVED***
	***REMOVED***
	appendDocker(&cmd)
	if err := validateArgs(cmd.Command...); err != nil ***REMOVED***
		return &icmd.Result***REMOVED***
			Error: err,
		***REMOVED***
	***REMOVED***
	return icmd.RunCmd(cmd)
***REMOVED***

// validateArgs is a checker to ensure tests are not running commands which are
// not supported on platforms. Specifically on Windows this is 'busybox top'.
func validateArgs(args ...string) error ***REMOVED***
	if testEnv.OSType != "windows" ***REMOVED***
		return nil
	***REMOVED***
	foundBusybox := -1
	for key, value := range args ***REMOVED***
		if strings.ToLower(value) == "busybox" ***REMOVED***
			foundBusybox = key
		***REMOVED***
		if (foundBusybox != -1) && (key == foundBusybox+1) && (strings.ToLower(value) == "top") ***REMOVED***
			return errors.New("cannot use 'busybox top' in tests on Windows. Use runSleepingContainer()")
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Build executes the specified docker build command
func Build(name string) icmd.Cmd ***REMOVED***
	return icmd.Command("build", "-t", name)
***REMOVED***

// Inspect executes the specified docker inspect command
func Inspect(name string) icmd.Cmd ***REMOVED***
	return icmd.Command("inspect", name)
***REMOVED***

// Format sets the specified format with --format flag
func Format(format string) func(*icmd.Cmd) func() ***REMOVED***
	return func(cmd *icmd.Cmd) func() ***REMOVED***
		cmd.Command = append(
			[]string***REMOVED***cmd.Command[0]***REMOVED***,
			append([]string***REMOVED***"--format", fmt.Sprintf("***REMOVED******REMOVED***%s***REMOVED******REMOVED***", format)***REMOVED***, cmd.Command[1:]...)...,
		)
		return nil
	***REMOVED***
***REMOVED***

func appendDocker(cmd *icmd.Cmd) ***REMOVED***
	cmd.Command = append([]string***REMOVED***testEnv.DockerBinary()***REMOVED***, cmd.Command...)
***REMOVED***

// Args build an icmd.Cmd struct from the specified arguments
func Args(args ...string) icmd.Cmd ***REMOVED***
	switch len(args) ***REMOVED***
	case 0:
		return icmd.Cmd***REMOVED******REMOVED***
	case 1:
		return icmd.Command(args[0])
	default:
		return icmd.Command(args[0], args[1:]...)
	***REMOVED***
***REMOVED***

// Daemon points to the specified daemon
func Daemon(d *daemon.Daemon) func(*icmd.Cmd) func() ***REMOVED***
	return func(cmd *icmd.Cmd) func() ***REMOVED***
		cmd.Command = append([]string***REMOVED***"--host", d.Sock()***REMOVED***, cmd.Command...)
		return nil
	***REMOVED***
***REMOVED***

// WithTimeout sets the timeout for the command to run
func WithTimeout(timeout time.Duration) func(cmd *icmd.Cmd) func() ***REMOVED***
	return func(cmd *icmd.Cmd) func() ***REMOVED***
		cmd.Timeout = timeout
		return nil
	***REMOVED***
***REMOVED***

// WithEnvironmentVariables sets the specified environment variables for the command to run
func WithEnvironmentVariables(envs ...string) func(cmd *icmd.Cmd) func() ***REMOVED***
	return func(cmd *icmd.Cmd) func() ***REMOVED***
		cmd.Env = envs
		return nil
	***REMOVED***
***REMOVED***

// WithFlags sets the specified flags for the command to run
func WithFlags(flags ...string) func(*icmd.Cmd) func() ***REMOVED***
	return func(cmd *icmd.Cmd) func() ***REMOVED***
		cmd.Command = append(cmd.Command, flags...)
		return nil
	***REMOVED***
***REMOVED***

// InDir sets the folder in which the command should be executed
func InDir(path string) func(*icmd.Cmd) func() ***REMOVED***
	return func(cmd *icmd.Cmd) func() ***REMOVED***
		cmd.Dir = path
		return nil
	***REMOVED***
***REMOVED***

// WithStdout sets the standard output writer of the command
func WithStdout(writer io.Writer) func(*icmd.Cmd) func() ***REMOVED***
	return func(cmd *icmd.Cmd) func() ***REMOVED***
		cmd.Stdout = writer
		return nil
	***REMOVED***
***REMOVED***

// WithStdin sets the standard input reader for the command
func WithStdin(stdin io.Reader) func(*icmd.Cmd) func() ***REMOVED***
	return func(cmd *icmd.Cmd) func() ***REMOVED***
		cmd.Stdin = stdin
		return nil
	***REMOVED***
***REMOVED***
