// +build !windows

package main

import (
	"strings"

	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/go-check/check"
)

func (s *DockerSwarmSuite) TestConfigList(c *check.C) ***REMOVED***
	testRequires(c, SameHostDaemon)
	d := s.AddDaemon(c, true, true)

	testName0 := "test0"
	testName1 := "test1"

	// create config test0
	id0 := d.CreateConfig(c, swarm.ConfigSpec***REMOVED***
		Annotations: swarm.Annotations***REMOVED***
			Name:   testName0,
			Labels: map[string]string***REMOVED***"type": "test"***REMOVED***,
		***REMOVED***,
		Data: []byte("TESTINGDATA0"),
	***REMOVED***)
	c.Assert(id0, checker.Not(checker.Equals), "", check.Commentf("configs: %s", id0))

	config := d.GetConfig(c, id0)
	c.Assert(config.Spec.Name, checker.Equals, testName0)

	// create config test1
	id1 := d.CreateConfig(c, swarm.ConfigSpec***REMOVED***
		Annotations: swarm.Annotations***REMOVED***
			Name:   testName1,
			Labels: map[string]string***REMOVED***"type": "production"***REMOVED***,
		***REMOVED***,
		Data: []byte("TESTINGDATA1"),
	***REMOVED***)
	c.Assert(id1, checker.Not(checker.Equals), "", check.Commentf("configs: %s", id1))

	config = d.GetConfig(c, id1)
	c.Assert(config.Spec.Name, checker.Equals, testName1)

	// test by command `docker config ls`
	out, err := d.Cmd("config", "ls")
	c.Assert(err, checker.IsNil, check.Commentf(out))
	c.Assert(strings.TrimSpace(out), checker.Contains, testName0)
	c.Assert(strings.TrimSpace(out), checker.Contains, testName1)

	// test filter by name `docker config ls --filter name=xxx`
	args := []string***REMOVED***
		"config",
		"ls",
		"--filter",
		"name=test0",
	***REMOVED***
	out, err = d.Cmd(args...)
	c.Assert(err, checker.IsNil, check.Commentf(out))

	c.Assert(strings.TrimSpace(out), checker.Contains, testName0)
	c.Assert(strings.TrimSpace(out), checker.Not(checker.Contains), testName1)

	// test filter by id `docker config ls --filter id=xxx`
	args = []string***REMOVED***
		"config",
		"ls",
		"--filter",
		"id=" + id1,
	***REMOVED***
	out, err = d.Cmd(args...)
	c.Assert(err, checker.IsNil, check.Commentf(out))

	c.Assert(strings.TrimSpace(out), checker.Not(checker.Contains), testName0)
	c.Assert(strings.TrimSpace(out), checker.Contains, testName1)

	// test filter by label `docker config ls --filter label=xxx`
	args = []string***REMOVED***
		"config",
		"ls",
		"--filter",
		"label=type",
	***REMOVED***
	out, err = d.Cmd(args...)
	c.Assert(err, checker.IsNil, check.Commentf(out))

	c.Assert(strings.TrimSpace(out), checker.Contains, testName0)
	c.Assert(strings.TrimSpace(out), checker.Contains, testName1)

	args = []string***REMOVED***
		"config",
		"ls",
		"--filter",
		"label=type=test",
	***REMOVED***
	out, err = d.Cmd(args...)
	c.Assert(err, checker.IsNil, check.Commentf(out))

	c.Assert(strings.TrimSpace(out), checker.Contains, testName0)
	c.Assert(strings.TrimSpace(out), checker.Not(checker.Contains), testName1)

	args = []string***REMOVED***
		"config",
		"ls",
		"--filter",
		"label=type=production",
	***REMOVED***
	out, err = d.Cmd(args...)
	c.Assert(err, checker.IsNil, check.Commentf(out))

	c.Assert(strings.TrimSpace(out), checker.Not(checker.Contains), testName0)
	c.Assert(strings.TrimSpace(out), checker.Contains, testName1)

	// test invalid filter `docker config ls --filter noexisttype=xxx`
	args = []string***REMOVED***
		"config",
		"ls",
		"--filter",
		"noexisttype=test0",
	***REMOVED***
	out, err = d.Cmd(args...)
	c.Assert(err, checker.NotNil, check.Commentf(out))

	c.Assert(strings.TrimSpace(out), checker.Contains, "Error response from daemon: Invalid filter 'noexisttype'")
***REMOVED***
