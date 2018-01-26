// +build !windows

package main

import (
	"time"

	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/daemon"
	"github.com/go-check/check"
)

func (s *DockerSwarmSuite) TestAPISwarmListNodes(c *check.C) ***REMOVED***
	d1 := s.AddDaemon(c, true, true)
	d2 := s.AddDaemon(c, true, false)
	d3 := s.AddDaemon(c, true, false)

	nodes := d1.ListNodes(c)
	c.Assert(len(nodes), checker.Equals, 3, check.Commentf("nodes: %#v", nodes))

loop0:
	for _, n := range nodes ***REMOVED***
		for _, d := range []*daemon.Swarm***REMOVED***d1, d2, d3***REMOVED*** ***REMOVED***
			if n.ID == d.NodeID ***REMOVED***
				continue loop0
			***REMOVED***
		***REMOVED***
		c.Errorf("unknown nodeID %v", n.ID)
	***REMOVED***
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmNodeUpdate(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, true, true)

	nodes := d.ListNodes(c)

	d.UpdateNode(c, nodes[0].ID, func(n *swarm.Node) ***REMOVED***
		n.Spec.Availability = swarm.NodeAvailabilityPause
	***REMOVED***)

	n := d.GetNode(c, nodes[0].ID)
	c.Assert(n.Spec.Availability, checker.Equals, swarm.NodeAvailabilityPause)
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmNodeRemove(c *check.C) ***REMOVED***
	testRequires(c, Network)
	d1 := s.AddDaemon(c, true, true)
	d2 := s.AddDaemon(c, true, false)
	_ = s.AddDaemon(c, true, false)

	nodes := d1.ListNodes(c)
	c.Assert(len(nodes), checker.Equals, 3, check.Commentf("nodes: %#v", nodes))

	// Getting the info so we can take the NodeID
	d2Info, err := d2.SwarmInfo()
	c.Assert(err, checker.IsNil)

	// forceful removal of d2 should work
	d1.RemoveNode(c, d2Info.NodeID, true)

	nodes = d1.ListNodes(c)
	c.Assert(len(nodes), checker.Equals, 2, check.Commentf("nodes: %#v", nodes))

	// Restart the node that was removed
	d2.Restart(c)

	// Give some time for the node to rejoin
	time.Sleep(1 * time.Second)

	// Make sure the node didn't rejoin
	nodes = d1.ListNodes(c)
	c.Assert(len(nodes), checker.Equals, 2, check.Commentf("nodes: %#v", nodes))
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmNodeDrainPause(c *check.C) ***REMOVED***
	d1 := s.AddDaemon(c, true, true)
	d2 := s.AddDaemon(c, true, false)

	time.Sleep(1 * time.Second) // make sure all daemons are ready to accept tasks

	// start a service, expect balanced distribution
	instances := 8
	id := d1.CreateService(c, simpleTestService, setInstances(instances))

	waitAndAssert(c, defaultReconciliationTimeout, d1.CheckActiveContainerCount, checker.GreaterThan, 0)
	waitAndAssert(c, defaultReconciliationTimeout, d2.CheckActiveContainerCount, checker.GreaterThan, 0)
	waitAndAssert(c, defaultReconciliationTimeout, reducedCheck(sumAsIntegers, d1.CheckActiveContainerCount, d2.CheckActiveContainerCount), checker.Equals, instances)

	// drain d2, all containers should move to d1
	d1.UpdateNode(c, d2.NodeID, func(n *swarm.Node) ***REMOVED***
		n.Spec.Availability = swarm.NodeAvailabilityDrain
	***REMOVED***)
	waitAndAssert(c, defaultReconciliationTimeout, d1.CheckActiveContainerCount, checker.Equals, instances)
	waitAndAssert(c, defaultReconciliationTimeout, d2.CheckActiveContainerCount, checker.Equals, 0)

	// set d2 back to active
	d1.UpdateNode(c, d2.NodeID, func(n *swarm.Node) ***REMOVED***
		n.Spec.Availability = swarm.NodeAvailabilityActive
	***REMOVED***)

	instances = 1
	d1.UpdateService(c, d1.GetService(c, id), setInstances(instances))

	waitAndAssert(c, defaultReconciliationTimeout*2, reducedCheck(sumAsIntegers, d1.CheckActiveContainerCount, d2.CheckActiveContainerCount), checker.Equals, instances)

	instances = 8
	d1.UpdateService(c, d1.GetService(c, id), setInstances(instances))

	// drained node first so we don't get any old containers
	waitAndAssert(c, defaultReconciliationTimeout, d2.CheckActiveContainerCount, checker.GreaterThan, 0)
	waitAndAssert(c, defaultReconciliationTimeout, d1.CheckActiveContainerCount, checker.GreaterThan, 0)
	waitAndAssert(c, defaultReconciliationTimeout*2, reducedCheck(sumAsIntegers, d1.CheckActiveContainerCount, d2.CheckActiveContainerCount), checker.Equals, instances)

	d2ContainerCount := len(d2.ActiveContainers())

	// set d2 to paused, scale service up, only d1 gets new tasks
	d1.UpdateNode(c, d2.NodeID, func(n *swarm.Node) ***REMOVED***
		n.Spec.Availability = swarm.NodeAvailabilityPause
	***REMOVED***)

	instances = 14
	d1.UpdateService(c, d1.GetService(c, id), setInstances(instances))

	waitAndAssert(c, defaultReconciliationTimeout, d1.CheckActiveContainerCount, checker.Equals, instances-d2ContainerCount)
	waitAndAssert(c, defaultReconciliationTimeout, d2.CheckActiveContainerCount, checker.Equals, d2ContainerCount)

***REMOVED***
