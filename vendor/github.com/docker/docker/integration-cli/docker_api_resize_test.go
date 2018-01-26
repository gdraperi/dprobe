package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/request"
	"github.com/go-check/check"
)

func (s *DockerSuite) TestResizeAPIResponse(c *check.C) ***REMOVED***
	out := runSleepingContainer(c, "-d")
	cleanedContainerID := strings.TrimSpace(out)
	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	options := types.ResizeOptions***REMOVED***
		Height: 40,
		Width:  40,
	***REMOVED***
	err = cli.ContainerResize(context.Background(), cleanedContainerID, options)
	c.Assert(err, check.IsNil)
***REMOVED***

func (s *DockerSuite) TestResizeAPIHeightWidthNoInt(c *check.C) ***REMOVED***
	out := runSleepingContainer(c, "-d")
	cleanedContainerID := strings.TrimSpace(out)

	endpoint := "/containers/" + cleanedContainerID + "/resize?h=foo&w=bar"
	res, _, err := request.Post(endpoint)
	c.Assert(res.StatusCode, check.Equals, http.StatusBadRequest)
	c.Assert(err, check.IsNil)
***REMOVED***

func (s *DockerSuite) TestResizeAPIResponseWhenContainerNotStarted(c *check.C) ***REMOVED***
	out, _ := dockerCmd(c, "run", "-d", "busybox", "true")
	cleanedContainerID := strings.TrimSpace(out)

	// make sure the exited container is not running
	dockerCmd(c, "wait", cleanedContainerID)

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	options := types.ResizeOptions***REMOVED***
		Height: 40,
		Width:  40,
	***REMOVED***

	err = cli.ContainerResize(context.Background(), cleanedContainerID, options)
	c.Assert(err.Error(), checker.Contains, "is not running")
***REMOVED***
