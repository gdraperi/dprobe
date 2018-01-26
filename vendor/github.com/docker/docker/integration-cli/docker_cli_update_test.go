package main

import (
	"strings"
	"time"

	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/cli"
	"github.com/go-check/check"
	"github.com/gotestyourself/gotestyourself/icmd"
)

func (s *DockerSuite) TestUpdateRestartPolicy(c *check.C) ***REMOVED***
	out := cli.DockerCmd(c, "run", "-d", "--restart=on-failure:3", "busybox", "sh", "-c", "sleep 1 && false").Combined()
	timeout := 60 * time.Second
	if testEnv.OSType == "windows" ***REMOVED***
		timeout = 180 * time.Second
	***REMOVED***

	id := strings.TrimSpace(string(out))

	// update restart policy to on-failure:5
	cli.DockerCmd(c, "update", "--restart=on-failure:5", id)

	cli.WaitExited(c, id, timeout)

	count := inspectField(c, id, "RestartCount")
	c.Assert(count, checker.Equals, "5")

	maximumRetryCount := inspectField(c, id, "HostConfig.RestartPolicy.MaximumRetryCount")
	c.Assert(maximumRetryCount, checker.Equals, "5")
***REMOVED***

func (s *DockerSuite) TestUpdateRestartWithAutoRemoveFlag(c *check.C) ***REMOVED***
	out := runSleepingContainer(c, "--rm")
	id := strings.TrimSpace(out)

	// update restart policy for an AutoRemove container
	cli.Docker(cli.Args("update", "--restart=always", id)).Assert(c, icmd.Expected***REMOVED***
		ExitCode: 1,
		Err:      "Restart policy cannot be updated because AutoRemove is enabled for the container",
	***REMOVED***)
***REMOVED***
