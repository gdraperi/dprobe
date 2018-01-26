package main

import (
	"bytes"
	"os/exec"
	"strings"
	"time"

	"github.com/docker/docker/integration-cli/checker"
	"github.com/go-check/check"
	"github.com/gotestyourself/gotestyourself/icmd"
)

// non-blocking wait with 0 exit code
func (s *DockerSuite) TestWaitNonBlockedExitZero(c *check.C) ***REMOVED***
	out, _ := dockerCmd(c, "run", "-d", "busybox", "sh", "-c", "true")
	containerID := strings.TrimSpace(out)

	err := waitInspect(containerID, "***REMOVED******REMOVED***.State.Running***REMOVED******REMOVED***", "false", 30*time.Second)
	c.Assert(err, checker.IsNil) //Container should have stopped by now

	out, _ = dockerCmd(c, "wait", containerID)
	c.Assert(strings.TrimSpace(out), checker.Equals, "0", check.Commentf("failed to set up container, %v", out))

***REMOVED***

// blocking wait with 0 exit code
func (s *DockerSuite) TestWaitBlockedExitZero(c *check.C) ***REMOVED***
	// Windows busybox does not support trap in this way, not sleep with sub-second
	// granularity. It will always exit 0x40010004.
	testRequires(c, DaemonIsLinux)
	out, _ := dockerCmd(c, "run", "-d", "busybox", "/bin/sh", "-c", "trap 'exit 0' TERM; while true; do usleep 10; done")
	containerID := strings.TrimSpace(out)

	c.Assert(waitRun(containerID), checker.IsNil)

	chWait := make(chan string)
	go func() ***REMOVED***
		chWait <- ""
		out := icmd.RunCommand(dockerBinary, "wait", containerID).Combined()
		chWait <- out
	***REMOVED***()

	<-chWait // make sure the goroutine is started
	time.Sleep(100 * time.Millisecond)
	dockerCmd(c, "stop", containerID)

	select ***REMOVED***
	case status := <-chWait:
		c.Assert(strings.TrimSpace(status), checker.Equals, "0", check.Commentf("expected exit 0, got %s", status))
	case <-time.After(2 * time.Second):
		c.Fatal("timeout waiting for `docker wait` to exit")
	***REMOVED***

***REMOVED***

// non-blocking wait with random exit code
func (s *DockerSuite) TestWaitNonBlockedExitRandom(c *check.C) ***REMOVED***
	out, _ := dockerCmd(c, "run", "-d", "busybox", "sh", "-c", "exit 99")
	containerID := strings.TrimSpace(out)

	err := waitInspect(containerID, "***REMOVED******REMOVED***.State.Running***REMOVED******REMOVED***", "false", 30*time.Second)
	c.Assert(err, checker.IsNil) //Container should have stopped by now
	out, _ = dockerCmd(c, "wait", containerID)
	c.Assert(strings.TrimSpace(out), checker.Equals, "99", check.Commentf("failed to set up container, %v", out))

***REMOVED***

// blocking wait with random exit code
func (s *DockerSuite) TestWaitBlockedExitRandom(c *check.C) ***REMOVED***
	// Cannot run on Windows as trap in Windows busybox does not support trap in this way.
	testRequires(c, DaemonIsLinux)
	out, _ := dockerCmd(c, "run", "-d", "busybox", "/bin/sh", "-c", "trap 'exit 99' TERM; while true; do usleep 10; done")
	containerID := strings.TrimSpace(out)
	c.Assert(waitRun(containerID), checker.IsNil)

	chWait := make(chan error)
	waitCmd := exec.Command(dockerBinary, "wait", containerID)
	waitCmdOut := bytes.NewBuffer(nil)
	waitCmd.Stdout = waitCmdOut
	c.Assert(waitCmd.Start(), checker.IsNil)
	go func() ***REMOVED***
		chWait <- waitCmd.Wait()
	***REMOVED***()

	dockerCmd(c, "stop", containerID)

	select ***REMOVED***
	case err := <-chWait:
		c.Assert(err, checker.IsNil, check.Commentf(waitCmdOut.String()))
		status, err := waitCmdOut.ReadString('\n')
		c.Assert(err, checker.IsNil)
		c.Assert(strings.TrimSpace(status), checker.Equals, "99", check.Commentf("expected exit 99, got %s", status))
	case <-time.After(2 * time.Second):
		waitCmd.Process.Kill()
		c.Fatal("timeout waiting for `docker wait` to exit")
	***REMOVED***
***REMOVED***
