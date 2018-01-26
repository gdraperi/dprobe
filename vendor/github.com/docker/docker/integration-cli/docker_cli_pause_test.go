package main

import (
	"strings"
	"time"

	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/cli"
	"github.com/go-check/check"
)

func (s *DockerSuite) TestPause(c *check.C) ***REMOVED***
	testRequires(c, IsPausable)

	name := "testeventpause"
	runSleepingContainer(c, "-d", "--name", name)

	cli.DockerCmd(c, "pause", name)
	pausedContainers := strings.Fields(
		cli.DockerCmd(c, "ps", "-f", "status=paused", "-q", "-a").Combined(),
	)
	c.Assert(len(pausedContainers), checker.Equals, 1)

	cli.DockerCmd(c, "unpause", name)

	out := cli.DockerCmd(c, "events", "--since=0", "--until", daemonUnixTime(c)).Combined()
	events := strings.Split(strings.TrimSpace(out), "\n")
	actions := eventActionsByIDAndType(c, events, name, "container")

	c.Assert(actions[len(actions)-2], checker.Equals, "pause")
	c.Assert(actions[len(actions)-1], checker.Equals, "unpause")
***REMOVED***

func (s *DockerSuite) TestPauseMultipleContainers(c *check.C) ***REMOVED***
	testRequires(c, IsPausable)

	containers := []string***REMOVED***
		"testpausewithmorecontainers1",
		"testpausewithmorecontainers2",
	***REMOVED***
	for _, name := range containers ***REMOVED***
		runSleepingContainer(c, "-d", "--name", name)
	***REMOVED***
	cli.DockerCmd(c, append([]string***REMOVED***"pause"***REMOVED***, containers...)...)
	pausedContainers := strings.Fields(
		cli.DockerCmd(c, "ps", "-f", "status=paused", "-q", "-a").Combined(),
	)
	c.Assert(len(pausedContainers), checker.Equals, len(containers))

	cli.DockerCmd(c, append([]string***REMOVED***"unpause"***REMOVED***, containers...)...)

	out := cli.DockerCmd(c, "events", "--since=0", "--until", daemonUnixTime(c)).Combined()
	events := strings.Split(strings.TrimSpace(out), "\n")

	for _, name := range containers ***REMOVED***
		actions := eventActionsByIDAndType(c, events, name, "container")

		c.Assert(actions[len(actions)-2], checker.Equals, "pause")
		c.Assert(actions[len(actions)-1], checker.Equals, "unpause")
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestPauseFailsOnWindowsServerContainers(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsWindows, NotPausable)
	runSleepingContainer(c, "-d", "--name=test")
	out, _, _ := dockerCmdWithError("pause", "test")
	c.Assert(out, checker.Contains, "cannot pause Windows Server Containers")
***REMOVED***

func (s *DockerSuite) TestStopPausedContainer(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)

	id := runSleepingContainer(c)
	cli.WaitRun(c, id)
	cli.DockerCmd(c, "pause", id)
	cli.DockerCmd(c, "stop", id)
	cli.WaitForInspectResult(c, id, "***REMOVED******REMOVED***.State.Running***REMOVED******REMOVED***", "false", 30*time.Second)
***REMOVED***
