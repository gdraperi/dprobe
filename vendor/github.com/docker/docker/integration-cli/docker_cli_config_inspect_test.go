// +build !windows

package main

import (
	"encoding/json"

	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/go-check/check"
)

func (s *DockerSwarmSuite) TestConfigInspect(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, true, true)

	testName := "test_config"
	id := d.CreateConfig(c, swarm.ConfigSpec***REMOVED***
		Annotations: swarm.Annotations***REMOVED***
			Name: testName,
		***REMOVED***,
		Data: []byte("TESTINGDATA"),
	***REMOVED***)
	c.Assert(id, checker.Not(checker.Equals), "", check.Commentf("configs: %s", id))

	config := d.GetConfig(c, id)
	c.Assert(config.Spec.Name, checker.Equals, testName)

	out, err := d.Cmd("config", "inspect", testName)
	c.Assert(err, checker.IsNil, check.Commentf(out))

	var configs []swarm.Config
	c.Assert(json.Unmarshal([]byte(out), &configs), checker.IsNil)
	c.Assert(configs, checker.HasLen, 1)
***REMOVED***

func (s *DockerSwarmSuite) TestConfigInspectMultiple(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, true, true)

	testNames := []string***REMOVED***
		"test0",
		"test1",
	***REMOVED***
	for _, n := range testNames ***REMOVED***
		id := d.CreateConfig(c, swarm.ConfigSpec***REMOVED***
			Annotations: swarm.Annotations***REMOVED***
				Name: n,
			***REMOVED***,
			Data: []byte("TESTINGDATA"),
		***REMOVED***)
		c.Assert(id, checker.Not(checker.Equals), "", check.Commentf("configs: %s", id))

		config := d.GetConfig(c, id)
		c.Assert(config.Spec.Name, checker.Equals, n)

	***REMOVED***

	args := []string***REMOVED***
		"config",
		"inspect",
	***REMOVED***
	args = append(args, testNames...)
	out, err := d.Cmd(args...)
	c.Assert(err, checker.IsNil, check.Commentf(out))

	var configs []swarm.Config
	c.Assert(json.Unmarshal([]byte(out), &configs), checker.IsNil)
	c.Assert(configs, checker.HasLen, 2)
***REMOVED***
