package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/integration-cli/cli"
	"github.com/go-check/check"
	"github.com/gotestyourself/gotestyourself/icmd"
)

const attachWait = 5 * time.Second

func (s *DockerSuite) TestAttachMultipleAndRestart(c *check.C) ***REMOVED***
	endGroup := &sync.WaitGroup***REMOVED******REMOVED***
	startGroup := &sync.WaitGroup***REMOVED******REMOVED***
	endGroup.Add(3)
	startGroup.Add(3)

	cli.DockerCmd(c, "run", "--name", "attacher", "-d", "busybox", "/bin/sh", "-c", "while true; do sleep 1; echo hello; done")
	cli.WaitRun(c, "attacher")

	startDone := make(chan struct***REMOVED******REMOVED***)
	endDone := make(chan struct***REMOVED******REMOVED***)

	go func() ***REMOVED***
		endGroup.Wait()
		close(endDone)
	***REMOVED***()

	go func() ***REMOVED***
		startGroup.Wait()
		close(startDone)
	***REMOVED***()

	for i := 0; i < 3; i++ ***REMOVED***
		go func() ***REMOVED***
			cmd := exec.Command(dockerBinary, "attach", "attacher")

			defer func() ***REMOVED***
				cmd.Wait()
				endGroup.Done()
			***REMOVED***()

			out, err := cmd.StdoutPipe()
			if err != nil ***REMOVED***
				c.Fatal(err)
			***REMOVED***
			defer out.Close()

			if err := cmd.Start(); err != nil ***REMOVED***
				c.Fatal(err)
			***REMOVED***

			buf := make([]byte, 1024)

			if _, err := out.Read(buf); err != nil && err != io.EOF ***REMOVED***
				c.Fatal(err)
			***REMOVED***

			startGroup.Done()

			if !strings.Contains(string(buf), "hello") ***REMOVED***
				c.Fatalf("unexpected output %s expected hello\n", string(buf))
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	select ***REMOVED***
	case <-startDone:
	case <-time.After(attachWait):
		c.Fatalf("Attaches did not initialize properly")
	***REMOVED***

	cli.DockerCmd(c, "kill", "attacher")

	select ***REMOVED***
	case <-endDone:
	case <-time.After(attachWait):
		c.Fatalf("Attaches did not finish properly")
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestAttachTTYWithoutStdin(c *check.C) ***REMOVED***
	// TODO @jhowardmsft. Figure out how to get this running again reliable on Windows.
	// It works by accident at the moment. Sometimes. I've gone back to v1.13.0 and see the same.
	// On Windows, docker run -d -ti busybox causes the container to exit immediately.
	// Obviously a year back when I updated the test, that was not the case. However,
	// with this, and the test racing with the tear-down which panic's, sometimes CI
	// will just fail and `MISS` all the other tests. For now, disabling it. Will
	// open an issue to track re-enabling this and root-causing the problem.
	testRequires(c, DaemonIsLinux)
	out, _ := dockerCmd(c, "run", "-d", "-ti", "busybox")

	id := strings.TrimSpace(out)
	c.Assert(waitRun(id), check.IsNil)

	done := make(chan error)
	go func() ***REMOVED***
		defer close(done)

		cmd := exec.Command(dockerBinary, "attach", id)
		if _, err := cmd.StdinPipe(); err != nil ***REMOVED***
			done <- err
			return
		***REMOVED***

		expected := "the input device is not a TTY"
		if runtime.GOOS == "windows" ***REMOVED***
			expected += ".  If you are using mintty, try prefixing the command with 'winpty'"
		***REMOVED***
		if out, _, err := runCommandWithOutput(cmd); err == nil ***REMOVED***
			done <- fmt.Errorf("attach should have failed")
			return
		***REMOVED*** else if !strings.Contains(out, expected) ***REMOVED***
			done <- fmt.Errorf("attach failed with error %q: expected %q", out, expected)
			return
		***REMOVED***
	***REMOVED***()

	select ***REMOVED***
	case err := <-done:
		c.Assert(err, check.IsNil)
	case <-time.After(attachWait):
		c.Fatal("attach is running but should have failed")
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestAttachDisconnect(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	out, _ := dockerCmd(c, "run", "-di", "busybox", "/bin/cat")
	id := strings.TrimSpace(out)

	cmd := exec.Command(dockerBinary, "attach", id)
	stdin, err := cmd.StdinPipe()
	if err != nil ***REMOVED***
		c.Fatal(err)
	***REMOVED***
	defer stdin.Close()
	stdout, err := cmd.StdoutPipe()
	c.Assert(err, check.IsNil)
	defer stdout.Close()
	c.Assert(cmd.Start(), check.IsNil)
	defer func() ***REMOVED***
		cmd.Process.Kill()
		cmd.Wait()
	***REMOVED***()

	_, err = stdin.Write([]byte("hello\n"))
	c.Assert(err, check.IsNil)
	out, err = bufio.NewReader(stdout).ReadString('\n')
	c.Assert(err, check.IsNil)
	c.Assert(strings.TrimSpace(out), check.Equals, "hello")

	c.Assert(stdin.Close(), check.IsNil)

	// Expect container to still be running after stdin is closed
	running := inspectField(c, id, "State.Running")
	c.Assert(running, check.Equals, "true")
***REMOVED***

func (s *DockerSuite) TestAttachPausedContainer(c *check.C) ***REMOVED***
	testRequires(c, IsPausable)
	runSleepingContainer(c, "-d", "--name=test")
	dockerCmd(c, "pause", "test")

	result := dockerCmdWithResult("attach", "test")
	result.Assert(c, icmd.Expected***REMOVED***
		Error:    "exit status 1",
		ExitCode: 1,
		Err:      "You cannot attach to a paused container, unpause it first",
	***REMOVED***)
***REMOVED***
