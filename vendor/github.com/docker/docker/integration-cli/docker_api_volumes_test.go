package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types/filters"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/go-check/check"
	"golang.org/x/net/context"
)

func (s *DockerSuite) TestVolumesAPIList(c *check.C) ***REMOVED***
	prefix, _ := getPrefixAndSlashFromDaemonPlatform()
	cid, _ := dockerCmd(c, "run", "-d", "-v", prefix+"/foo", "busybox")

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	container, err := cli.ContainerInspect(context.Background(), strings.TrimSpace(cid))
	c.Assert(err, checker.IsNil)
	vname := container.Mounts[0].Name

	volumes, err := cli.VolumeList(context.Background(), filters.Args***REMOVED******REMOVED***)
	c.Assert(err, checker.IsNil)

	found := false
	for _, vol := range volumes.Volumes ***REMOVED***
		if vol.Name == vname ***REMOVED***
			found = true
			break
		***REMOVED***
	***REMOVED***
	c.Assert(found, checker.Equals, true)
***REMOVED***

func (s *DockerSuite) TestVolumesAPICreate(c *check.C) ***REMOVED***
	config := volumetypes.VolumesCreateBody***REMOVED***
		Name: "test",
	***REMOVED***

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	vol, err := cli.VolumeCreate(context.Background(), config)
	c.Assert(err, check.IsNil)

	c.Assert(filepath.Base(filepath.Dir(vol.Mountpoint)), checker.Equals, config.Name)
***REMOVED***

func (s *DockerSuite) TestVolumesAPIRemove(c *check.C) ***REMOVED***
	prefix, _ := getPrefixAndSlashFromDaemonPlatform()
	cid, _ := dockerCmd(c, "run", "-d", "-v", prefix+"/foo", "--name=test", "busybox")

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	container, err := cli.ContainerInspect(context.Background(), strings.TrimSpace(cid))
	c.Assert(err, checker.IsNil)
	vname := container.Mounts[0].Name

	err = cli.VolumeRemove(context.Background(), vname, false)
	c.Assert(err.Error(), checker.Contains, "volume is in use")

	dockerCmd(c, "rm", "-f", "test")
	err = cli.VolumeRemove(context.Background(), vname, false)
	c.Assert(err, checker.IsNil)
***REMOVED***

func (s *DockerSuite) TestVolumesAPIInspect(c *check.C) ***REMOVED***
	config := volumetypes.VolumesCreateBody***REMOVED***
		Name: "test",
	***REMOVED***

	// sampling current time minus a minute so to now have false positive in case of delays
	now := time.Now().Truncate(time.Minute)

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	_, err = cli.VolumeCreate(context.Background(), config)
	c.Assert(err, check.IsNil)

	vol, err := cli.VolumeInspect(context.Background(), config.Name)
	c.Assert(err, checker.IsNil)
	c.Assert(vol.Name, checker.Equals, config.Name)

	// comparing CreatedAt field time for the new volume to now. Removing a minute from both to avoid false positive
	testCreatedAt, err := time.Parse(time.RFC3339, strings.TrimSpace(vol.CreatedAt))
	c.Assert(err, check.IsNil)
	testCreatedAt = testCreatedAt.Truncate(time.Minute)
	if !testCreatedAt.Equal(now) ***REMOVED***
		c.Assert(fmt.Errorf("Time Volume is CreatedAt not equal to current time"), check.NotNil)
	***REMOVED***
***REMOVED***
