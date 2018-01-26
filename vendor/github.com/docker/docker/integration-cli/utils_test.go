package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/docker/docker/internal/testutil"
	"github.com/go-check/check"
	"github.com/gotestyourself/gotestyourself/icmd"
	"github.com/pkg/errors"
)

func getPrefixAndSlashFromDaemonPlatform() (prefix, slash string) ***REMOVED***
	if testEnv.OSType == "windows" ***REMOVED***
		return "c:", `\`
	***REMOVED***
	return "", "/"
***REMOVED***

// TODO: update code to call cmd.RunCmd directly, and remove this function
// Deprecated: use gotestyourself/gotestyourself/icmd
func runCommandWithOutput(execCmd *exec.Cmd) (string, int, error) ***REMOVED***
	result := icmd.RunCmd(transformCmd(execCmd))
	return result.Combined(), result.ExitCode, result.Error
***REMOVED***

// Temporary shim for migrating commands to the new function
func transformCmd(execCmd *exec.Cmd) icmd.Cmd ***REMOVED***
	return icmd.Cmd***REMOVED***
		Command: execCmd.Args,
		Env:     execCmd.Env,
		Dir:     execCmd.Dir,
		Stdin:   execCmd.Stdin,
		Stdout:  execCmd.Stdout,
	***REMOVED***
***REMOVED***

// ParseCgroupPaths parses 'procCgroupData', which is output of '/proc/<pid>/cgroup', and returns
// a map which cgroup name as key and path as value.
func ParseCgroupPaths(procCgroupData string) map[string]string ***REMOVED***
	cgroupPaths := map[string]string***REMOVED******REMOVED***
	for _, line := range strings.Split(procCgroupData, "\n") ***REMOVED***
		parts := strings.Split(line, ":")
		if len(parts) != 3 ***REMOVED***
			continue
		***REMOVED***
		cgroupPaths[parts[1]] = parts[2]
	***REMOVED***
	return cgroupPaths
***REMOVED***

// RandomTmpDirPath provides a temporary path with rand string appended.
// does not create or checks if it exists.
func RandomTmpDirPath(s string, platform string) string ***REMOVED***
	// TODO: why doesn't this use os.TempDir() ?
	tmp := "/tmp"
	if platform == "windows" ***REMOVED***
		tmp = os.Getenv("TEMP")
	***REMOVED***
	path := filepath.Join(tmp, fmt.Sprintf("%s.%s", s, testutil.GenerateRandomAlphaOnlyString(10)))
	if platform == "windows" ***REMOVED***
		return filepath.FromSlash(path) // Using \
	***REMOVED***
	return filepath.ToSlash(path) // Using /
***REMOVED***

// RunCommandPipelineWithOutput runs the array of commands with the output
// of each pipelined with the following (like cmd1 | cmd2 | cmd3 would do).
// It returns the final output, the exitCode different from 0 and the error
// if something bad happened.
// Deprecated: use icmd instead
func RunCommandPipelineWithOutput(cmds ...*exec.Cmd) (output string, err error) ***REMOVED***
	if len(cmds) < 2 ***REMOVED***
		return "", errors.New("pipeline does not have multiple cmds")
	***REMOVED***

	// connect stdin of each cmd to stdout pipe of previous cmd
	for i, cmd := range cmds ***REMOVED***
		if i > 0 ***REMOVED***
			prevCmd := cmds[i-1]
			cmd.Stdin, err = prevCmd.StdoutPipe()

			if err != nil ***REMOVED***
				return "", fmt.Errorf("cannot set stdout pipe for %s: %v", cmd.Path, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// start all cmds except the last
	for _, cmd := range cmds[:len(cmds)-1] ***REMOVED***
		if err = cmd.Start(); err != nil ***REMOVED***
			return "", fmt.Errorf("starting %s failed with error: %v", cmd.Path, err)
		***REMOVED***
	***REMOVED***

	defer func() ***REMOVED***
		var pipeErrMsgs []string
		// wait all cmds except the last to release their resources
		for _, cmd := range cmds[:len(cmds)-1] ***REMOVED***
			if pipeErr := cmd.Wait(); pipeErr != nil ***REMOVED***
				pipeErrMsgs = append(pipeErrMsgs, fmt.Sprintf("command %s failed with error: %v", cmd.Path, pipeErr))
			***REMOVED***
		***REMOVED***
		if len(pipeErrMsgs) > 0 && err == nil ***REMOVED***
			err = fmt.Errorf("pipelineError from Wait: %v", strings.Join(pipeErrMsgs, ", "))
		***REMOVED***
	***REMOVED***()

	// wait on last cmd
	out, err := cmds[len(cmds)-1].CombinedOutput()
	return string(out), err
***REMOVED***

type elementListOptions struct ***REMOVED***
	element, format string
***REMOVED***

func existingElements(c *check.C, opts elementListOptions) []string ***REMOVED***
	args := []string***REMOVED******REMOVED***
	switch opts.element ***REMOVED***
	case "container":
		args = append(args, "ps", "-a")
	case "image":
		args = append(args, "images", "-a")
	case "network":
		args = append(args, "network", "ls")
	case "plugin":
		args = append(args, "plugin", "ls")
	case "volume":
		args = append(args, "volume", "ls")
	***REMOVED***
	if opts.format != "" ***REMOVED***
		args = append(args, "--format", opts.format)
	***REMOVED***
	out, _ := dockerCmd(c, args...)
	lines := []string***REMOVED******REMOVED***
	for _, l := range strings.Split(out, "\n") ***REMOVED***
		if l != "" ***REMOVED***
			lines = append(lines, l)
		***REMOVED***
	***REMOVED***
	return lines
***REMOVED***

// ExistingContainerIDs returns a list of currently existing container IDs.
func ExistingContainerIDs(c *check.C) []string ***REMOVED***
	return existingElements(c, elementListOptions***REMOVED***element: "container", format: "***REMOVED******REMOVED***.ID***REMOVED******REMOVED***"***REMOVED***)
***REMOVED***

// ExistingContainerNames returns a list of existing container names.
func ExistingContainerNames(c *check.C) []string ***REMOVED***
	return existingElements(c, elementListOptions***REMOVED***element: "container", format: "***REMOVED******REMOVED***.Names***REMOVED******REMOVED***"***REMOVED***)
***REMOVED***

// RemoveLinesForExistingElements removes existing elements from the output of a
// docker command.
// This function takes an output []string and returns a []string.
func RemoveLinesForExistingElements(output, existing []string) []string ***REMOVED***
	for _, e := range existing ***REMOVED***
		index := -1
		for i, line := range output ***REMOVED***
			if strings.Contains(line, e) ***REMOVED***
				index = i
				break
			***REMOVED***
		***REMOVED***
		if index != -1 ***REMOVED***
			output = append(output[:index], output[index+1:]...)
		***REMOVED***
	***REMOVED***
	return output
***REMOVED***

// RemoveOutputForExistingElements removes existing elements from the output of
// a docker command.
// This function takes an output string and returns a string.
func RemoveOutputForExistingElements(output string, existing []string) string ***REMOVED***
	res := RemoveLinesForExistingElements(strings.Split(output, "\n"), existing)
	return strings.Join(res, "\n")
***REMOVED***
