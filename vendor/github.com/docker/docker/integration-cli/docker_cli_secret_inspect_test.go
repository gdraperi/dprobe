// +build !windows

package main

import (
	"encoding/json"

	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/go-check/check"
)

func (s *DockerSwarmSuite) TestSecretInspect(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, true, true)

	testName := "test_secret"
	id := d.CreateSecret(c, swarm.SecretSpec***REMOVED***
		Annotations: swarm.Annotations***REMOVED***
			Name: testName,
		***REMOVED***,
		Data: []byte("TESTINGDATA"),
	***REMOVED***)
	c.Assert(id, checker.Not(checker.Equals), "", check.Commentf("secrets: %s", id))

	secret := d.GetSecret(c, id)
	c.Assert(secret.Spec.Name, checker.Equals, testName)

	out, err := d.Cmd("secret", "inspect", testName)
	c.Assert(err, checker.IsNil, check.Commentf(out))

	var secrets []swarm.Secret
	c.Assert(json.Unmarshal([]byte(out), &secrets), checker.IsNil)
	c.Assert(secrets, checker.HasLen, 1)
***REMOVED***

func (s *DockerSwarmSuite) TestSecretInspectMultiple(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, true, true)

	testNames := []string***REMOVED***
		"test0",
		"test1",
	***REMOVED***
	for _, n := range testNames ***REMOVED***
		id := d.CreateSecret(c, swarm.SecretSpec***REMOVED***
			Annotations: swarm.Annotations***REMOVED***
				Name: n,
			***REMOVED***,
			Data: []byte("TESTINGDATA"),
		***REMOVED***)
		c.Assert(id, checker.Not(checker.Equals), "", check.Commentf("secrets: %s", id))

		secret := d.GetSecret(c, id)
		c.Assert(secret.Spec.Name, checker.Equals, n)

	***REMOVED***

	args := []string***REMOVED***
		"secret",
		"inspect",
	***REMOVED***
	args = append(args, testNames...)
	out, err := d.Cmd(args...)
	c.Assert(err, checker.IsNil, check.Commentf(out))

	var secrets []swarm.Secret
	c.Assert(json.Unmarshal([]byte(out), &secrets), checker.IsNil)
	c.Assert(secrets, checker.HasLen, 2)
***REMOVED***
