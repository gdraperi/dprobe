package main

import (
	"encoding/json"
	"strings"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/versions/v1p20"
	"github.com/docker/docker/client"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/go-check/check"
	"github.com/stretchr/testify/assert"
)

func (s *DockerSuite) TestInspectAPIContainerResponse(c *check.C) ***REMOVED***
	out, _ := dockerCmd(c, "run", "-d", "busybox", "true")

	cleanedContainerID := strings.TrimSpace(out)
	keysBase := []string***REMOVED***"Id", "State", "Created", "Path", "Args", "Config", "Image", "NetworkSettings",
		"ResolvConfPath", "HostnamePath", "HostsPath", "LogPath", "Name", "Driver", "MountLabel", "ProcessLabel", "GraphDriver"***REMOVED***

	type acase struct ***REMOVED***
		version string
		keys    []string
	***REMOVED***

	var cases []acase

	if testEnv.OSType == "windows" ***REMOVED***
		cases = []acase***REMOVED***
			***REMOVED***"v1.25", append(keysBase, "Mounts")***REMOVED***,
		***REMOVED***

	***REMOVED*** else ***REMOVED***
		cases = []acase***REMOVED***
			***REMOVED***"v1.20", append(keysBase, "Mounts")***REMOVED***,
			***REMOVED***"v1.19", append(keysBase, "Volumes", "VolumesRW")***REMOVED***,
		***REMOVED***
	***REMOVED***

	for _, cs := range cases ***REMOVED***
		body := getInspectBody(c, cs.version, cleanedContainerID)

		var inspectJSON map[string]interface***REMOVED******REMOVED***
		err := json.Unmarshal(body, &inspectJSON)
		c.Assert(err, checker.IsNil, check.Commentf("Unable to unmarshal body for version %s", cs.version))

		for _, key := range cs.keys ***REMOVED***
			_, ok := inspectJSON[key]
			c.Check(ok, checker.True, check.Commentf("%s does not exist in response for version %s", key, cs.version))
		***REMOVED***

		//Issue #6830: type not properly converted to JSON/back
		_, ok := inspectJSON["Path"].(bool)
		c.Assert(ok, checker.False, check.Commentf("Path of `true` should not be converted to boolean `true` via JSON marshalling"))
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestInspectAPIContainerVolumeDriverLegacy(c *check.C) ***REMOVED***
	// No legacy implications for Windows
	testRequires(c, DaemonIsLinux)
	out, _ := dockerCmd(c, "run", "-d", "busybox", "true")

	cleanedContainerID := strings.TrimSpace(out)

	cases := []string***REMOVED***"v1.19", "v1.20"***REMOVED***
	for _, version := range cases ***REMOVED***
		body := getInspectBody(c, version, cleanedContainerID)

		var inspectJSON map[string]interface***REMOVED******REMOVED***
		err := json.Unmarshal(body, &inspectJSON)
		c.Assert(err, checker.IsNil, check.Commentf("Unable to unmarshal body for version %s", version))

		config, ok := inspectJSON["Config"]
		c.Assert(ok, checker.True, check.Commentf("Unable to find 'Config'"))
		cfg := config.(map[string]interface***REMOVED******REMOVED***)
		_, ok = cfg["VolumeDriver"]
		c.Assert(ok, checker.True, check.Commentf("API version %s expected to include VolumeDriver in 'Config'", version))
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestInspectAPIContainerVolumeDriver(c *check.C) ***REMOVED***
	out, _ := dockerCmd(c, "run", "-d", "--volume-driver", "local", "busybox", "true")

	cleanedContainerID := strings.TrimSpace(out)

	body := getInspectBody(c, "v1.25", cleanedContainerID)

	var inspectJSON map[string]interface***REMOVED******REMOVED***
	err := json.Unmarshal(body, &inspectJSON)
	c.Assert(err, checker.IsNil, check.Commentf("Unable to unmarshal body for version 1.25"))

	config, ok := inspectJSON["Config"]
	c.Assert(ok, checker.True, check.Commentf("Unable to find 'Config'"))
	cfg := config.(map[string]interface***REMOVED******REMOVED***)
	_, ok = cfg["VolumeDriver"]
	c.Assert(ok, checker.False, check.Commentf("API version 1.25 expected to not include VolumeDriver in 'Config'"))

	config, ok = inspectJSON["HostConfig"]
	c.Assert(ok, checker.True, check.Commentf("Unable to find 'HostConfig'"))
	cfg = config.(map[string]interface***REMOVED******REMOVED***)
	_, ok = cfg["VolumeDriver"]
	c.Assert(ok, checker.True, check.Commentf("API version 1.25 expected to include VolumeDriver in 'HostConfig'"))
***REMOVED***

func (s *DockerSuite) TestInspectAPIImageResponse(c *check.C) ***REMOVED***
	dockerCmd(c, "tag", "busybox:latest", "busybox:mytag")
	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	imageJSON, _, err := cli.ImageInspectWithRaw(context.Background(), "busybox")
	c.Assert(err, checker.IsNil)

	c.Assert(imageJSON.RepoTags, checker.HasLen, 2)
	assert.Contains(c, imageJSON.RepoTags, "busybox:latest")
	assert.Contains(c, imageJSON.RepoTags, "busybox:mytag")
***REMOVED***

// #17131, #17139, #17173
func (s *DockerSuite) TestInspectAPIEmptyFieldsInConfigPre121(c *check.C) ***REMOVED***
	// Not relevant on Windows
	testRequires(c, DaemonIsLinux)
	out, _ := dockerCmd(c, "run", "-d", "busybox", "true")

	cleanedContainerID := strings.TrimSpace(out)

	cases := []string***REMOVED***"v1.19", "v1.20"***REMOVED***
	for _, version := range cases ***REMOVED***
		body := getInspectBody(c, version, cleanedContainerID)

		var inspectJSON map[string]interface***REMOVED******REMOVED***
		err := json.Unmarshal(body, &inspectJSON)
		c.Assert(err, checker.IsNil, check.Commentf("Unable to unmarshal body for version %s", version))
		config, ok := inspectJSON["Config"]
		c.Assert(ok, checker.True, check.Commentf("Unable to find 'Config'"))
		cfg := config.(map[string]interface***REMOVED******REMOVED***)
		for _, f := range []string***REMOVED***"MacAddress", "NetworkDisabled", "ExposedPorts"***REMOVED*** ***REMOVED***
			_, ok := cfg[f]
			c.Check(ok, checker.True, check.Commentf("API version %s expected to include %s in 'Config'", version, f))
		***REMOVED***
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestInspectAPIBridgeNetworkSettings120(c *check.C) ***REMOVED***
	// Not relevant on Windows, and besides it doesn't have any bridge network settings
	testRequires(c, DaemonIsLinux)
	out, _ := dockerCmd(c, "run", "-d", "busybox", "top")
	containerID := strings.TrimSpace(out)
	waitRun(containerID)

	body := getInspectBody(c, "v1.20", containerID)

	var inspectJSON v1p20.ContainerJSON
	err := json.Unmarshal(body, &inspectJSON)
	c.Assert(err, checker.IsNil)

	settings := inspectJSON.NetworkSettings
	c.Assert(settings.IPAddress, checker.Not(checker.HasLen), 0)
***REMOVED***

func (s *DockerSuite) TestInspectAPIBridgeNetworkSettings121(c *check.C) ***REMOVED***
	// Windows doesn't have any bridge network settings
	testRequires(c, DaemonIsLinux)
	out, _ := dockerCmd(c, "run", "-d", "busybox", "top")
	containerID := strings.TrimSpace(out)
	waitRun(containerID)

	body := getInspectBody(c, "v1.21", containerID)

	var inspectJSON types.ContainerJSON
	err := json.Unmarshal(body, &inspectJSON)
	c.Assert(err, checker.IsNil)

	settings := inspectJSON.NetworkSettings
	c.Assert(settings.IPAddress, checker.Not(checker.HasLen), 0)
	c.Assert(settings.Networks["bridge"], checker.Not(checker.IsNil))
	c.Assert(settings.IPAddress, checker.Equals, settings.Networks["bridge"].IPAddress)
***REMOVED***
