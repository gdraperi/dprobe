package main

import (
	"strings"

	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/pkg/stringid"
	"github.com/go-check/check"
	"github.com/gotestyourself/gotestyourself/icmd"
)

func (s *DockerSuite) TestRenameStoppedContainer(c *check.C) ***REMOVED***
	out, _ := dockerCmd(c, "run", "--name", "first_name", "-d", "busybox", "sh")

	cleanedContainerID := strings.TrimSpace(out)
	dockerCmd(c, "wait", cleanedContainerID)

	name := inspectField(c, cleanedContainerID, "Name")
	newName := "new_name" + stringid.GenerateNonCryptoID()
	dockerCmd(c, "rename", "first_name", newName)

	name = inspectField(c, cleanedContainerID, "Name")
	c.Assert(name, checker.Equals, "/"+newName, check.Commentf("Failed to rename container %s", name))

***REMOVED***

func (s *DockerSuite) TestRenameRunningContainer(c *check.C) ***REMOVED***
	out, _ := dockerCmd(c, "run", "--name", "first_name", "-d", "busybox", "sh")

	newName := "new_name" + stringid.GenerateNonCryptoID()
	cleanedContainerID := strings.TrimSpace(out)
	dockerCmd(c, "rename", "first_name", newName)

	name := inspectField(c, cleanedContainerID, "Name")
	c.Assert(name, checker.Equals, "/"+newName, check.Commentf("Failed to rename container %s", name))
***REMOVED***

func (s *DockerSuite) TestRenameRunningContainerAndReuse(c *check.C) ***REMOVED***
	out := runSleepingContainer(c, "--name", "first_name")
	c.Assert(waitRun("first_name"), check.IsNil)

	newName := "new_name"
	ContainerID := strings.TrimSpace(out)
	dockerCmd(c, "rename", "first_name", newName)

	name := inspectField(c, ContainerID, "Name")
	c.Assert(name, checker.Equals, "/"+newName, check.Commentf("Failed to rename container"))

	out = runSleepingContainer(c, "--name", "first_name")
	c.Assert(waitRun("first_name"), check.IsNil)
	newContainerID := strings.TrimSpace(out)
	name = inspectField(c, newContainerID, "Name")
	c.Assert(name, checker.Equals, "/first_name", check.Commentf("Failed to reuse container name"))
***REMOVED***

func (s *DockerSuite) TestRenameCheckNames(c *check.C) ***REMOVED***
	dockerCmd(c, "run", "--name", "first_name", "-d", "busybox", "sh")

	newName := "new_name" + stringid.GenerateNonCryptoID()
	dockerCmd(c, "rename", "first_name", newName)

	name := inspectField(c, newName, "Name")
	c.Assert(name, checker.Equals, "/"+newName, check.Commentf("Failed to rename container %s", name))

	result := dockerCmdWithResult("inspect", "-f=***REMOVED******REMOVED***.Name***REMOVED******REMOVED***", "--type=container", "first_name")
	result.Assert(c, icmd.Expected***REMOVED***
		ExitCode: 1,
		Err:      "No such container: first_name",
	***REMOVED***)
***REMOVED***

// TODO: move to unit test
func (s *DockerSuite) TestRenameInvalidName(c *check.C) ***REMOVED***
	runSleepingContainer(c, "--name", "myname")

	out, _, err := dockerCmdWithError("rename", "myname", "new:invalid")
	c.Assert(err, checker.NotNil, check.Commentf("Renaming container to invalid name should have failed: %s", out))
	c.Assert(out, checker.Contains, "Invalid container name", check.Commentf("%v", err))

	out, _ = dockerCmd(c, "ps", "-a")
	c.Assert(out, checker.Contains, "myname", check.Commentf("Output of docker ps should have included 'myname': %s", out))
***REMOVED***

func (s *DockerSuite) TestRenameAnonymousContainer(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)

	dockerCmd(c, "network", "create", "network1")
	out, _ := dockerCmd(c, "create", "-it", "--net", "network1", "busybox", "top")

	anonymousContainerID := strings.TrimSpace(out)

	dockerCmd(c, "rename", anonymousContainerID, "container1")
	dockerCmd(c, "start", "container1")

	count := "-c"
	if testEnv.OSType == "windows" ***REMOVED***
		count = "-n"
	***REMOVED***

	_, _, err := dockerCmdWithError("run", "--net", "network1", "busybox", "ping", count, "1", "container1")
	c.Assert(err, check.IsNil, check.Commentf("Embedded DNS lookup fails after renaming anonymous container: %v", err))
***REMOVED***

func (s *DockerSuite) TestRenameContainerWithSameName(c *check.C) ***REMOVED***
	out := runSleepingContainer(c, "--name", "old")
	ContainerID := strings.TrimSpace(out)

	out, _, err := dockerCmdWithError("rename", "old", "old")
	c.Assert(err, checker.NotNil, check.Commentf("Renaming a container with the same name should have failed"))
	c.Assert(out, checker.Contains, "Renaming a container with the same name", check.Commentf("%v", err))

	out, _, err = dockerCmdWithError("rename", ContainerID, "old")
	c.Assert(err, checker.NotNil, check.Commentf("Renaming a container with the same name should have failed"))
	c.Assert(out, checker.Contains, "Renaming a container with the same name", check.Commentf("%v", err))
***REMOVED***

// Test case for #23973
func (s *DockerSuite) TestRenameContainerWithLinkedContainer(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)

	db1, _ := dockerCmd(c, "run", "--name", "db1", "-d", "busybox", "top")
	dockerCmd(c, "run", "--name", "app1", "-d", "--link", "db1:/mysql", "busybox", "top")
	dockerCmd(c, "rename", "app1", "app2")
	out, _, err := dockerCmdWithError("inspect", "--format=***REMOVED******REMOVED*** .Id ***REMOVED******REMOVED***", "app2/mysql")
	c.Assert(err, checker.IsNil)
	c.Assert(strings.TrimSpace(out), checker.Equals, strings.TrimSpace(db1))
***REMOVED***
