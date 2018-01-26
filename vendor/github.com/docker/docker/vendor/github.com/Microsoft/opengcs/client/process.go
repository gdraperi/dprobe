// +build windows

package client

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Microsoft/hcsshim"
	"github.com/sirupsen/logrus"
)

// Process is the structure pertaining to a process running in a utility VM.
type process struct ***REMOVED***
	Process hcsshim.Process
	Stdin   io.WriteCloser
	Stdout  io.ReadCloser
	Stderr  io.ReadCloser
***REMOVED***

// createUtilsProcess is a convenient wrapper for hcsshim.createUtilsProcess to use when
// communicating with a utility VM.
func (config *Config) createUtilsProcess(commandLine string) (process, error) ***REMOVED***
	logrus.Debugf("opengcs: createUtilsProcess")

	if config.Uvm == nil ***REMOVED***
		return process***REMOVED******REMOVED***, fmt.Errorf("cannot create utils process as no utility VM is in configuration")
	***REMOVED***

	var (
		err  error
		proc process
	)

	env := make(map[string]string)
	env["PATH"] = "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:"
	processConfig := &hcsshim.ProcessConfig***REMOVED***
		EmulateConsole:    false,
		CreateStdInPipe:   true,
		CreateStdOutPipe:  true,
		CreateStdErrPipe:  true,
		CreateInUtilityVm: true,
		WorkingDirectory:  "/bin",
		Environment:       env,
		CommandLine:       commandLine,
	***REMOVED***
	proc.Process, err = config.Uvm.CreateProcess(processConfig)
	if err != nil ***REMOVED***
		return process***REMOVED******REMOVED***, fmt.Errorf("failed to create process (%+v) in utility VM: %s", config, err)
	***REMOVED***

	if proc.Stdin, proc.Stdout, proc.Stderr, err = proc.Process.Stdio(); err != nil ***REMOVED***
		proc.Process.Kill() // Should this have a timeout?
		proc.Process.Close()
		return process***REMOVED******REMOVED***, fmt.Errorf("failed to get stdio pipes for process %+v: %s", config, err)
	***REMOVED***

	logrus.Debugf("opengcs: createUtilsProcess success: pid %d", proc.Process.Pid())
	return proc, nil
***REMOVED***

// RunProcess runs the given command line program in the utilityVM. It takes in
// an input to the reader to feed into stdin and returns stdout to output.
// IMPORTANT: It is the responsibility of the caller to call Close() on the returned process.
func (config *Config) RunProcess(commandLine string, stdin io.Reader, stdout io.Writer, stderr io.Writer) (hcsshim.Process, error) ***REMOVED***
	logrus.Debugf("opengcs: RunProcess: %s", commandLine)
	process, err := config.createUtilsProcess(commandLine)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Send the data into the process's stdin
	if stdin != nil ***REMOVED***
		if _, err = copyWithTimeout(process.Stdin,
			stdin,
			0,
			config.UvmTimeoutSeconds,
			fmt.Sprintf("send to stdin of %s", commandLine)); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		// Don't need stdin now we've sent everything. This signals GCS that we are finished sending data.
		if err := process.Process.CloseStdin(); err != nil && !hcsshim.IsNotExist(err) && !hcsshim.IsAlreadyClosed(err) ***REMOVED***
			// This error will occur if the compute system is currently shutting down
			if perr, ok := err.(*hcsshim.ProcessError); ok && perr.Err != hcsshim.ErrVmcomputeOperationInvalidState ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if stdout != nil ***REMOVED***
		// Copy the data over to the writer.
		if _, err := copyWithTimeout(stdout,
			process.Stdout,
			0,
			config.UvmTimeoutSeconds,
			fmt.Sprintf("RunProcess: copy back from %s", commandLine)); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	if stderr != nil ***REMOVED***
		// Copy the data over to the writer.
		if _, err := copyWithTimeout(stderr,
			process.Stderr,
			0,
			config.UvmTimeoutSeconds,
			fmt.Sprintf("RunProcess: copy back from %s", commandLine)); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	logrus.Debugf("opengcs: runProcess success: %s", commandLine)
	return process.Process, nil
***REMOVED***

func debugCommand(s string) string ***REMOVED***
	return fmt.Sprintf(`echo -e 'DEBUG COMMAND: %s\\n--------------\\n';%s;echo -e '\\n\\n';`, s, s)
***REMOVED***

// DebugGCS extracts logs from the GCS. It's a useful hack for debugging,
// but not necessarily optimal, but all that is available to us in RS3.
func (config *Config) DebugGCS() ***REMOVED***
	if logrus.GetLevel() < logrus.DebugLevel || len(os.Getenv("OPENGCS_DEBUG_ENABLE")) == 0 ***REMOVED***
		return
	***REMOVED***

	var out bytes.Buffer
	cmd := os.Getenv("OPENGCS_DEBUG_COMMAND")
	if cmd == "" ***REMOVED***
		cmd = `sh -c "`
		cmd += debugCommand("kill -10 `pidof gcs`") // SIGUSR1 for stackdump
		cmd += debugCommand("ls -l /tmp")
		cmd += debugCommand("cat /tmp/gcs.log")
		cmd += debugCommand("cat /tmp/gcs/gcs-stacks*")
		cmd += debugCommand("cat /tmp/gcs/paniclog*")
		cmd += debugCommand("ls -l /tmp/gcs")
		cmd += debugCommand("ls -l /tmp/gcs/*")
		cmd += debugCommand("cat /tmp/gcs/*/config.json")
		cmd += debugCommand("ls -lR /var/run/gcsrunc")
		cmd += debugCommand("cat /tmp/gcs/global-runc.log")
		cmd += debugCommand("cat /tmp/gcs/*/runc.log")
		cmd += debugCommand("ps -ef")
		cmd += `"`
	***REMOVED***
	proc, err := config.RunProcess(cmd, nil, &out, nil)
	defer func() ***REMOVED***
		if proc != nil ***REMOVED***
			proc.Kill()
			proc.Close()
		***REMOVED***
	***REMOVED***()
	if err != nil ***REMOVED***
		logrus.Debugln("benign failure getting gcs logs: ", err)
	***REMOVED***
	if proc != nil ***REMOVED***
		proc.WaitTimeout(time.Duration(int(time.Second) * 30))
	***REMOVED***
	logrus.Debugf("GCS Debugging:\n%s\n\nEnd GCS Debugging", strings.TrimSpace(out.String()))
***REMOVED***
