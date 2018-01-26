package main

import (
	"strings"

	"github.com/docker/docker/integration-cli/checker"
	"github.com/go-check/check"
	"github.com/gotestyourself/gotestyourself/icmd"
)

func (s *DockerSuite) TestTopMultipleArgs(c *check.C) ***REMOVED***
	out := runSleepingContainer(c, "-d")
	cleanedContainerID := strings.TrimSpace(out)

	var expected icmd.Expected
	switch testEnv.OSType ***REMOVED***
	case "windows":
		expected = icmd.Expected***REMOVED***ExitCode: 1, Err: "Windows does not support arguments to top"***REMOVED***
	default:
		expected = icmd.Expected***REMOVED***Out: "PID"***REMOVED***
	***REMOVED***
	result := dockerCmdWithResult("top", cleanedContainerID, "-o", "pid")
	result.Assert(c, expected)
***REMOVED***

func (s *DockerSuite) TestTopNonPrivileged(c *check.C) ***REMOVED***
	out := runSleepingContainer(c, "-d")
	cleanedContainerID := strings.TrimSpace(out)

	out1, _ := dockerCmd(c, "top", cleanedContainerID)
	out2, _ := dockerCmd(c, "top", cleanedContainerID)
	dockerCmd(c, "kill", cleanedContainerID)

	// Windows will list the name of the launched executable which in this case is busybox.exe, without the parameters.
	// Linux will display the command executed in the container
	var lookingFor string
	if testEnv.OSType == "windows" ***REMOVED***
		lookingFor = "busybox.exe"
	***REMOVED*** else ***REMOVED***
		lookingFor = "top"
	***REMOVED***

	c.Assert(out1, checker.Contains, lookingFor, check.Commentf("top should've listed `%s` in the process list, but failed the first time", lookingFor))
	c.Assert(out2, checker.Contains, lookingFor, check.Commentf("top should've listed `%s` in the process list, but failed the second time", lookingFor))
***REMOVED***

// TestTopWindowsCoreProcesses validates that there are lines for the critical
// processes which are found in a Windows container. Note Windows is architecturally
// very different to Linux in this regard.
func (s *DockerSuite) TestTopWindowsCoreProcesses(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsWindows)
	out := runSleepingContainer(c, "-d")
	cleanedContainerID := strings.TrimSpace(out)
	out1, _ := dockerCmd(c, "top", cleanedContainerID)
	lookingFor := []string***REMOVED***"smss.exe", "csrss.exe", "wininit.exe", "services.exe", "lsass.exe", "CExecSvc.exe"***REMOVED***
	for i, s := range lookingFor ***REMOVED***
		c.Assert(out1, checker.Contains, s, check.Commentf("top should've listed `%s` in the process list, but failed. Test case %d", s, i))
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestTopPrivileged(c *check.C) ***REMOVED***
	// Windows does not support --privileged
	testRequires(c, DaemonIsLinux, NotUserNamespace)
	out, _ := dockerCmd(c, "run", "--privileged", "-i", "-d", "busybox", "top")
	cleanedContainerID := strings.TrimSpace(out)

	out1, _ := dockerCmd(c, "top", cleanedContainerID)
	out2, _ := dockerCmd(c, "top", cleanedContainerID)
	dockerCmd(c, "kill", cleanedContainerID)

	c.Assert(out1, checker.Contains, "top", check.Commentf("top should've listed `top` in the process list, but failed the first time"))
	c.Assert(out2, checker.Contains, "top", check.Commentf("top should've listed `top` in the process list, but failed the second time"))
***REMOVED***
