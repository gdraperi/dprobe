package main

import (
	"github.com/docker/docker/integration-cli/daemon"
	"github.com/go-check/check"
)

func (s *DockerSwarmSuite) getDaemon(c *check.C, nodeID string) *daemon.Swarm ***REMOVED***
	s.daemonsLock.Lock()
	defer s.daemonsLock.Unlock()
	for _, d := range s.daemons ***REMOVED***
		if d.NodeID == nodeID ***REMOVED***
			return d
		***REMOVED***
	***REMOVED***
	c.Fatalf("could not find node with id: %s", nodeID)
	return nil
***REMOVED***

// nodeCmd executes a command on a given node via the normal docker socket
func (s *DockerSwarmSuite) nodeCmd(c *check.C, id string, args ...string) (string, error) ***REMOVED***
	return s.getDaemon(c, id).Cmd(args...)
***REMOVED***
