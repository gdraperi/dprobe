package main

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/cli/build"
	"github.com/go-check/check"
)

func waitForHealthStatus(c *check.C, name string, prev string, expected string) ***REMOVED***
	prev = prev + "\n"
	expected = expected + "\n"
	for ***REMOVED***
		out, _ := dockerCmd(c, "inspect", "--format=***REMOVED******REMOVED***.State.Health.Status***REMOVED******REMOVED***", name)
		if out == expected ***REMOVED***
			return
		***REMOVED***
		c.Check(out, checker.Equals, prev)
		if out != prev ***REMOVED***
			return
		***REMOVED***
		time.Sleep(100 * time.Millisecond)
	***REMOVED***
***REMOVED***

func getHealth(c *check.C, name string) *types.Health ***REMOVED***
	out, _ := dockerCmd(c, "inspect", "--format=***REMOVED******REMOVED***json .State.Health***REMOVED******REMOVED***", name)
	var health types.Health
	err := json.Unmarshal([]byte(out), &health)
	c.Check(err, checker.Equals, nil)
	return &health
***REMOVED***

func (s *DockerSuite) TestHealth(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux) // busybox doesn't work on Windows

	existingContainers := ExistingContainerIDs(c)

	imageName := "testhealth"
	buildImageSuccessfully(c, imageName, build.WithDockerfile(`FROM busybox
		RUN echo OK > /status
		CMD ["/bin/sleep", "120"]
		STOPSIGNAL SIGKILL
		HEALTHCHECK --interval=1s --timeout=30s \
		  CMD cat /status`))

	// No health status before starting
	name := "test_health"
	cid, _ := dockerCmd(c, "create", "--name", name, imageName)
	out, _ := dockerCmd(c, "ps", "-a", "--format=***REMOVED******REMOVED***.ID***REMOVED******REMOVED*** ***REMOVED******REMOVED***.Status***REMOVED******REMOVED***")
	out = RemoveOutputForExistingElements(out, existingContainers)
	c.Check(out, checker.Equals, cid[:12]+" Created\n")

	// Inspect the options
	out, _ = dockerCmd(c, "inspect",
		"--format=timeout=***REMOVED******REMOVED***.Config.Healthcheck.Timeout***REMOVED******REMOVED*** interval=***REMOVED******REMOVED***.Config.Healthcheck.Interval***REMOVED******REMOVED*** retries=***REMOVED******REMOVED***.Config.Healthcheck.Retries***REMOVED******REMOVED*** test=***REMOVED******REMOVED***.Config.Healthcheck.Test***REMOVED******REMOVED***", name)
	c.Check(out, checker.Equals, "timeout=30s interval=1s retries=0 test=[CMD-SHELL cat /status]\n")

	// Start
	dockerCmd(c, "start", name)
	waitForHealthStatus(c, name, "starting", "healthy")

	// Make it fail
	dockerCmd(c, "exec", name, "rm", "/status")
	waitForHealthStatus(c, name, "healthy", "unhealthy")

	// Inspect the status
	out, _ = dockerCmd(c, "inspect", "--format=***REMOVED******REMOVED***.State.Health.Status***REMOVED******REMOVED***", name)
	c.Check(out, checker.Equals, "unhealthy\n")

	// Make it healthy again
	dockerCmd(c, "exec", name, "touch", "/status")
	waitForHealthStatus(c, name, "unhealthy", "healthy")

	// Remove container
	dockerCmd(c, "rm", "-f", name)

	// Disable the check from the CLI
	out, _ = dockerCmd(c, "create", "--name=noh", "--no-healthcheck", imageName)
	out, _ = dockerCmd(c, "inspect", "--format=***REMOVED******REMOVED***.Config.Healthcheck.Test***REMOVED******REMOVED***", "noh")
	c.Check(out, checker.Equals, "[NONE]\n")
	dockerCmd(c, "rm", "noh")

	// Disable the check with a new build
	buildImageSuccessfully(c, "no_healthcheck", build.WithDockerfile(`FROM testhealth
		HEALTHCHECK NONE`))

	out, _ = dockerCmd(c, "inspect", "--format=***REMOVED******REMOVED***.ContainerConfig.Healthcheck.Test***REMOVED******REMOVED***", "no_healthcheck")
	c.Check(out, checker.Equals, "[NONE]\n")

	// Enable the checks from the CLI
	_, _ = dockerCmd(c, "run", "-d", "--name=fatal_healthcheck",
		"--health-interval=1s",
		"--health-retries=3",
		"--health-cmd=cat /status",
		"no_healthcheck")
	waitForHealthStatus(c, "fatal_healthcheck", "starting", "healthy")
	health := getHealth(c, "fatal_healthcheck")
	c.Check(health.Status, checker.Equals, "healthy")
	c.Check(health.FailingStreak, checker.Equals, 0)
	last := health.Log[len(health.Log)-1]
	c.Check(last.ExitCode, checker.Equals, 0)
	c.Check(last.Output, checker.Equals, "OK\n")

	// Fail the check
	dockerCmd(c, "exec", "fatal_healthcheck", "rm", "/status")
	waitForHealthStatus(c, "fatal_healthcheck", "healthy", "unhealthy")

	failsStr, _ := dockerCmd(c, "inspect", "--format=***REMOVED******REMOVED***.State.Health.FailingStreak***REMOVED******REMOVED***", "fatal_healthcheck")
	fails, err := strconv.Atoi(strings.TrimSpace(failsStr))
	c.Check(err, check.IsNil)
	c.Check(fails >= 3, checker.Equals, true)
	dockerCmd(c, "rm", "-f", "fatal_healthcheck")

	// Check timeout
	// Note: if the interval is too small, it seems that Docker spends all its time running health
	// checks and never gets around to killing it.
	_, _ = dockerCmd(c, "run", "-d", "--name=test",
		"--health-interval=1s", "--health-cmd=sleep 5m", "--health-timeout=1s", imageName)
	waitForHealthStatus(c, "test", "starting", "unhealthy")
	health = getHealth(c, "test")
	last = health.Log[len(health.Log)-1]
	c.Check(health.Status, checker.Equals, "unhealthy")
	c.Check(last.ExitCode, checker.Equals, -1)
	c.Check(last.Output, checker.Equals, "Health check exceeded timeout (1s)")
	dockerCmd(c, "rm", "-f", "test")

	// Check JSON-format
	buildImageSuccessfully(c, imageName, build.WithDockerfile(`FROM busybox
		RUN echo OK > /status
		CMD ["/bin/sleep", "120"]
		STOPSIGNAL SIGKILL
		HEALTHCHECK --interval=1s --timeout=30s \
		  CMD ["cat", "/my status"]`))
	out, _ = dockerCmd(c, "inspect",
		"--format=***REMOVED******REMOVED***.Config.Healthcheck.Test***REMOVED******REMOVED***", imageName)
	c.Check(out, checker.Equals, "[CMD cat /my status]\n")

***REMOVED***

// GitHub #33021
func (s *DockerSuite) TestUnsetEnvVarHealthCheck(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux) // busybox doesn't work on Windows

	imageName := "testhealth"
	buildImageSuccessfully(c, imageName, build.WithDockerfile(`FROM busybox
HEALTHCHECK --interval=1s --timeout=5s --retries=5 CMD /bin/sh -c "sleep 1"
ENTRYPOINT /bin/sh -c "sleep 600"`))

	name := "env_test_health"
	// No health status before starting
	dockerCmd(c, "run", "-d", "--name", name, "-e", "FOO", imageName)
	defer func() ***REMOVED***
		dockerCmd(c, "rm", "-f", name)
		dockerCmd(c, "rmi", imageName)
	***REMOVED***()

	// Start
	dockerCmd(c, "start", name)
	waitForHealthStatus(c, name, "starting", "healthy")

***REMOVED***
