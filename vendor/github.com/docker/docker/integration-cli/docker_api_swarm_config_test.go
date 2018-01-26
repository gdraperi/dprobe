// +build !windows

package main

import (
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/go-check/check"
	"golang.org/x/net/context"
)

func (s *DockerSwarmSuite) TestAPISwarmConfigsEmptyList(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, true, true)

	configs := d.ListConfigs(c)
	c.Assert(configs, checker.NotNil)
	c.Assert(len(configs), checker.Equals, 0, check.Commentf("configs: %#v", configs))
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmConfigsCreate(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, true, true)

	testName := "test_config"
	id := d.CreateConfig(c, swarm.ConfigSpec***REMOVED***
		Annotations: swarm.Annotations***REMOVED***
			Name: testName,
		***REMOVED***,
		Data: []byte("TESTINGDATA"),
	***REMOVED***)
	c.Assert(id, checker.Not(checker.Equals), "", check.Commentf("configs: %s", id))

	configs := d.ListConfigs(c)
	c.Assert(len(configs), checker.Equals, 1, check.Commentf("configs: %#v", configs))
	name := configs[0].Spec.Annotations.Name
	c.Assert(name, checker.Equals, testName, check.Commentf("configs: %s", name))
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmConfigsDelete(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, true, true)

	testName := "test_config"
	id := d.CreateConfig(c, swarm.ConfigSpec***REMOVED***Annotations: swarm.Annotations***REMOVED***
		Name: testName,
	***REMOVED***,
		Data: []byte("TESTINGDATA"),
	***REMOVED***)
	c.Assert(id, checker.Not(checker.Equals), "", check.Commentf("configs: %s", id))

	config := d.GetConfig(c, id)
	c.Assert(config.ID, checker.Equals, id, check.Commentf("config: %v", config))

	d.DeleteConfig(c, config.ID)

	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	_, _, err = cli.ConfigInspectWithRaw(context.Background(), id)
	c.Assert(err.Error(), checker.Contains, "No such config")
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmConfigsUpdate(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, true, true)

	testName := "test_config"
	id := d.CreateConfig(c, swarm.ConfigSpec***REMOVED***
		Annotations: swarm.Annotations***REMOVED***
			Name: testName,
			Labels: map[string]string***REMOVED***
				"test": "test1",
			***REMOVED***,
		***REMOVED***,
		Data: []byte("TESTINGDATA"),
	***REMOVED***)
	c.Assert(id, checker.Not(checker.Equals), "", check.Commentf("configs: %s", id))

	config := d.GetConfig(c, id)
	c.Assert(config.ID, checker.Equals, id, check.Commentf("config: %v", config))

	// test UpdateConfig with full ID
	d.UpdateConfig(c, id, func(s *swarm.Config) ***REMOVED***
		s.Spec.Labels = map[string]string***REMOVED***
			"test": "test1",
		***REMOVED***
	***REMOVED***)

	config = d.GetConfig(c, id)
	c.Assert(config.Spec.Labels["test"], checker.Equals, "test1", check.Commentf("config: %v", config))

	// test UpdateConfig with full name
	d.UpdateConfig(c, config.Spec.Name, func(s *swarm.Config) ***REMOVED***
		s.Spec.Labels = map[string]string***REMOVED***
			"test": "test2",
		***REMOVED***
	***REMOVED***)

	config = d.GetConfig(c, id)
	c.Assert(config.Spec.Labels["test"], checker.Equals, "test2", check.Commentf("config: %v", config))

	// test UpdateConfig with prefix ID
	d.UpdateConfig(c, id[:1], func(s *swarm.Config) ***REMOVED***
		s.Spec.Labels = map[string]string***REMOVED***
			"test": "test3",
		***REMOVED***
	***REMOVED***)

	config = d.GetConfig(c, id)
	c.Assert(config.Spec.Labels["test"], checker.Equals, "test3", check.Commentf("config: %v", config))

	// test UpdateConfig in updating Data which is not supported in daemon
	// this test will produce an error in func UpdateConfig
	config = d.GetConfig(c, id)
	config.Spec.Data = []byte("TESTINGDATA2")

	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	expected := "only updates to Labels are allowed"

	err = cli.ConfigUpdate(context.Background(), config.ID, config.Version, config.Spec)
	c.Assert(err.Error(), checker.Contains, expected)
***REMOVED***
