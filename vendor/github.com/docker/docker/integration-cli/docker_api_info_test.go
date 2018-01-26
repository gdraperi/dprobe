package main

import (
	"net/http"

	"fmt"

	"github.com/docker/docker/client"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/request"
	"github.com/go-check/check"
	"golang.org/x/net/context"
)

func (s *DockerSuite) TestInfoAPI(c *check.C) ***REMOVED***
	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	info, err := cli.Info(context.Background())
	c.Assert(err, checker.IsNil)

	// always shown fields
	stringsToCheck := []string***REMOVED***
		"ID",
		"Containers",
		"ContainersRunning",
		"ContainersPaused",
		"ContainersStopped",
		"Images",
		"LoggingDriver",
		"OperatingSystem",
		"NCPU",
		"OSType",
		"Architecture",
		"MemTotal",
		"KernelVersion",
		"Driver",
		"ServerVersion",
		"SecurityOptions"***REMOVED***

	out := fmt.Sprintf("%+v", info)
	for _, linePrefix := range stringsToCheck ***REMOVED***
		c.Assert(out, checker.Contains, linePrefix)
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestInfoAPIVersioned(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux) // Windows only supports 1.25 or later

	res, body, err := request.Get("/v1.20/info")
	c.Assert(res.StatusCode, checker.Equals, http.StatusOK)
	c.Assert(err, checker.IsNil)

	b, err := request.ReadBody(body)
	c.Assert(err, checker.IsNil)

	out := string(b)
	c.Assert(out, checker.Contains, "ExecutionDriver")
	c.Assert(out, checker.Contains, "not supported")
***REMOVED***
