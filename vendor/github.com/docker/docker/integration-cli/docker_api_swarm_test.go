// +build !windows

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/helpers"
	"github.com/cloudflare/cfssl/initca"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/daemon"
	"github.com/docker/docker/integration-cli/request"
	"github.com/docker/swarmkit/ca"
	"github.com/go-check/check"
	"golang.org/x/net/context"
)

var defaultReconciliationTimeout = 30 * time.Second

func (s *DockerSwarmSuite) TestAPISwarmInit(c *check.C) ***REMOVED***
	// todo: should find a better way to verify that components are running than /info
	d1 := s.AddDaemon(c, true, true)
	info, err := d1.SwarmInfo()
	c.Assert(err, checker.IsNil)
	c.Assert(info.ControlAvailable, checker.True)
	c.Assert(info.LocalNodeState, checker.Equals, swarm.LocalNodeStateActive)
	c.Assert(info.Cluster.RootRotationInProgress, checker.False)

	d2 := s.AddDaemon(c, true, false)
	info, err = d2.SwarmInfo()
	c.Assert(err, checker.IsNil)
	c.Assert(info.ControlAvailable, checker.False)
	c.Assert(info.LocalNodeState, checker.Equals, swarm.LocalNodeStateActive)

	// Leaving cluster
	c.Assert(d2.Leave(false), checker.IsNil)

	info, err = d2.SwarmInfo()
	c.Assert(err, checker.IsNil)
	c.Assert(info.ControlAvailable, checker.False)
	c.Assert(info.LocalNodeState, checker.Equals, swarm.LocalNodeStateInactive)

	c.Assert(d2.Join(swarm.JoinRequest***REMOVED***JoinToken: d1.JoinTokens(c).Worker, RemoteAddrs: []string***REMOVED***d1.ListenAddr***REMOVED******REMOVED***), checker.IsNil)

	info, err = d2.SwarmInfo()
	c.Assert(err, checker.IsNil)
	c.Assert(info.ControlAvailable, checker.False)
	c.Assert(info.LocalNodeState, checker.Equals, swarm.LocalNodeStateActive)

	// Current state restoring after restarts
	d1.Stop(c)
	d2.Stop(c)

	d1.Start(c)
	d2.Start(c)

	info, err = d1.SwarmInfo()
	c.Assert(err, checker.IsNil)
	c.Assert(info.ControlAvailable, checker.True)
	c.Assert(info.LocalNodeState, checker.Equals, swarm.LocalNodeStateActive)

	info, err = d2.SwarmInfo()
	c.Assert(err, checker.IsNil)
	c.Assert(info.ControlAvailable, checker.False)
	c.Assert(info.LocalNodeState, checker.Equals, swarm.LocalNodeStateActive)
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmJoinToken(c *check.C) ***REMOVED***
	d1 := s.AddDaemon(c, false, false)
	c.Assert(d1.Init(swarm.InitRequest***REMOVED******REMOVED***), checker.IsNil)

	// todo: error message differs depending if some components of token are valid

	d2 := s.AddDaemon(c, false, false)
	err := d2.Join(swarm.JoinRequest***REMOVED***RemoteAddrs: []string***REMOVED***d1.ListenAddr***REMOVED******REMOVED***)
	c.Assert(err, checker.NotNil)
	c.Assert(err.Error(), checker.Contains, "join token is necessary")
	info, err := d2.SwarmInfo()
	c.Assert(err, checker.IsNil)
	c.Assert(info.LocalNodeState, checker.Equals, swarm.LocalNodeStateInactive)

	err = d2.Join(swarm.JoinRequest***REMOVED***JoinToken: "foobaz", RemoteAddrs: []string***REMOVED***d1.ListenAddr***REMOVED******REMOVED***)
	c.Assert(err, checker.NotNil)
	c.Assert(err.Error(), checker.Contains, "invalid join token")
	info, err = d2.SwarmInfo()
	c.Assert(err, checker.IsNil)
	c.Assert(info.LocalNodeState, checker.Equals, swarm.LocalNodeStateInactive)

	workerToken := d1.JoinTokens(c).Worker

	c.Assert(d2.Join(swarm.JoinRequest***REMOVED***JoinToken: workerToken, RemoteAddrs: []string***REMOVED***d1.ListenAddr***REMOVED******REMOVED***), checker.IsNil)
	info, err = d2.SwarmInfo()
	c.Assert(err, checker.IsNil)
	c.Assert(info.LocalNodeState, checker.Equals, swarm.LocalNodeStateActive)
	c.Assert(d2.Leave(false), checker.IsNil)
	info, err = d2.SwarmInfo()
	c.Assert(err, checker.IsNil)
	c.Assert(info.LocalNodeState, checker.Equals, swarm.LocalNodeStateInactive)

	// change tokens
	d1.RotateTokens(c)

	err = d2.Join(swarm.JoinRequest***REMOVED***JoinToken: workerToken, RemoteAddrs: []string***REMOVED***d1.ListenAddr***REMOVED******REMOVED***)
	c.Assert(err, checker.NotNil)
	c.Assert(err.Error(), checker.Contains, "join token is necessary")
	info, err = d2.SwarmInfo()
	c.Assert(err, checker.IsNil)
	c.Assert(info.LocalNodeState, checker.Equals, swarm.LocalNodeStateInactive)

	workerToken = d1.JoinTokens(c).Worker

	c.Assert(d2.Join(swarm.JoinRequest***REMOVED***JoinToken: workerToken, RemoteAddrs: []string***REMOVED***d1.ListenAddr***REMOVED******REMOVED***), checker.IsNil)
	info, err = d2.SwarmInfo()
	c.Assert(err, checker.IsNil)
	c.Assert(info.LocalNodeState, checker.Equals, swarm.LocalNodeStateActive)
	c.Assert(d2.Leave(false), checker.IsNil)
	info, err = d2.SwarmInfo()
	c.Assert(err, checker.IsNil)
	c.Assert(info.LocalNodeState, checker.Equals, swarm.LocalNodeStateInactive)

	// change spec, don't change tokens
	d1.UpdateSwarm(c, func(s *swarm.Spec) ***REMOVED******REMOVED***)

	err = d2.Join(swarm.JoinRequest***REMOVED***RemoteAddrs: []string***REMOVED***d1.ListenAddr***REMOVED******REMOVED***)
	c.Assert(err, checker.NotNil)
	c.Assert(err.Error(), checker.Contains, "join token is necessary")
	info, err = d2.SwarmInfo()
	c.Assert(err, checker.IsNil)
	c.Assert(info.LocalNodeState, checker.Equals, swarm.LocalNodeStateInactive)

	c.Assert(d2.Join(swarm.JoinRequest***REMOVED***JoinToken: workerToken, RemoteAddrs: []string***REMOVED***d1.ListenAddr***REMOVED******REMOVED***), checker.IsNil)
	info, err = d2.SwarmInfo()
	c.Assert(err, checker.IsNil)
	c.Assert(info.LocalNodeState, checker.Equals, swarm.LocalNodeStateActive)
	c.Assert(d2.Leave(false), checker.IsNil)
	info, err = d2.SwarmInfo()
	c.Assert(err, checker.IsNil)
	c.Assert(info.LocalNodeState, checker.Equals, swarm.LocalNodeStateInactive)
***REMOVED***

func (s *DockerSwarmSuite) TestUpdateSwarmAddExternalCA(c *check.C) ***REMOVED***
	d1 := s.AddDaemon(c, false, false)
	c.Assert(d1.Init(swarm.InitRequest***REMOVED******REMOVED***), checker.IsNil)
	d1.UpdateSwarm(c, func(s *swarm.Spec) ***REMOVED***
		s.CAConfig.ExternalCAs = []*swarm.ExternalCA***REMOVED***
			***REMOVED***
				Protocol: swarm.ExternalCAProtocolCFSSL,
				URL:      "https://thishasnoca.org",
			***REMOVED***,
			***REMOVED***
				Protocol: swarm.ExternalCAProtocolCFSSL,
				URL:      "https://thishasacacert.org",
				CACert:   "cacert",
			***REMOVED***,
		***REMOVED***
	***REMOVED***)
	info, err := d1.SwarmInfo()
	c.Assert(err, checker.IsNil)
	c.Assert(info.Cluster.Spec.CAConfig.ExternalCAs, checker.HasLen, 2)
	c.Assert(info.Cluster.Spec.CAConfig.ExternalCAs[0].CACert, checker.Equals, "")
	c.Assert(info.Cluster.Spec.CAConfig.ExternalCAs[1].CACert, checker.Equals, "cacert")
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmCAHash(c *check.C) ***REMOVED***
	d1 := s.AddDaemon(c, true, true)
	d2 := s.AddDaemon(c, false, false)
	splitToken := strings.Split(d1.JoinTokens(c).Worker, "-")
	splitToken[2] = "1kxftv4ofnc6mt30lmgipg6ngf9luhwqopfk1tz6bdmnkubg0e"
	replacementToken := strings.Join(splitToken, "-")
	err := d2.Join(swarm.JoinRequest***REMOVED***JoinToken: replacementToken, RemoteAddrs: []string***REMOVED***d1.ListenAddr***REMOVED******REMOVED***)
	c.Assert(err, checker.NotNil)
	c.Assert(err.Error(), checker.Contains, "remote CA does not match fingerprint")
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmPromoteDemote(c *check.C) ***REMOVED***
	d1 := s.AddDaemon(c, false, false)
	c.Assert(d1.Init(swarm.InitRequest***REMOVED******REMOVED***), checker.IsNil)
	d2 := s.AddDaemon(c, true, false)

	info, err := d2.SwarmInfo()
	c.Assert(err, checker.IsNil)
	c.Assert(info.ControlAvailable, checker.False)
	c.Assert(info.LocalNodeState, checker.Equals, swarm.LocalNodeStateActive)

	d1.UpdateNode(c, d2.NodeID, func(n *swarm.Node) ***REMOVED***
		n.Spec.Role = swarm.NodeRoleManager
	***REMOVED***)

	waitAndAssert(c, defaultReconciliationTimeout, d2.CheckControlAvailable, checker.True)

	d1.UpdateNode(c, d2.NodeID, func(n *swarm.Node) ***REMOVED***
		n.Spec.Role = swarm.NodeRoleWorker
	***REMOVED***)

	waitAndAssert(c, defaultReconciliationTimeout, d2.CheckControlAvailable, checker.False)

	// Wait for the role to change to worker in the cert. This is partially
	// done because it's something worth testing in its own right, and
	// partially because changing the role from manager to worker and then
	// back to manager quickly might cause the node to pause for awhile
	// while waiting for the role to change to worker, and the test can
	// time out during this interval.
	waitAndAssert(c, defaultReconciliationTimeout, func(c *check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
		certBytes, err := ioutil.ReadFile(filepath.Join(d2.Folder, "root", "swarm", "certificates", "swarm-node.crt"))
		if err != nil ***REMOVED***
			return "", check.Commentf("error: %v", err)
		***REMOVED***
		certs, err := helpers.ParseCertificatesPEM(certBytes)
		if err == nil && len(certs) > 0 && len(certs[0].Subject.OrganizationalUnit) > 0 ***REMOVED***
			return certs[0].Subject.OrganizationalUnit[0], nil
		***REMOVED***
		return "", check.Commentf("could not get organizational unit from certificate")
	***REMOVED***, checker.Equals, "swarm-worker")

	// Demoting last node should fail
	node := d1.GetNode(c, d1.NodeID)
	node.Spec.Role = swarm.NodeRoleWorker
	url := fmt.Sprintf("/nodes/%s/update?version=%d", node.ID, node.Version.Index)
	res, body, err := request.DoOnHost(d1.Sock(), url, request.Method("POST"), request.JSONBody(node.Spec))
	c.Assert(err, checker.IsNil)
	b, err := request.ReadBody(body)
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusBadRequest, check.Commentf("output: %q", string(b)))

	// The warning specific to demoting the last manager is best-effort and
	// won't appear until the Role field of the demoted manager has been
	// updated.
	// Yes, I know this looks silly, but checker.Matches is broken, since
	// it anchors the regexp contrary to the documentation, and this makes
	// it impossible to match something that includes a line break.
	if !strings.Contains(string(b), "last manager of the swarm") ***REMOVED***
		c.Assert(string(b), checker.Contains, "this would result in a loss of quorum")
	***REMOVED***
	info, err = d1.SwarmInfo()
	c.Assert(err, checker.IsNil)
	c.Assert(info.LocalNodeState, checker.Equals, swarm.LocalNodeStateActive)
	c.Assert(info.ControlAvailable, checker.True)

	// Promote already demoted node
	d1.UpdateNode(c, d2.NodeID, func(n *swarm.Node) ***REMOVED***
		n.Spec.Role = swarm.NodeRoleManager
	***REMOVED***)

	waitAndAssert(c, defaultReconciliationTimeout, d2.CheckControlAvailable, checker.True)
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmLeaderProxy(c *check.C) ***REMOVED***
	// add three managers, one of these is leader
	d1 := s.AddDaemon(c, true, true)
	d2 := s.AddDaemon(c, true, true)
	d3 := s.AddDaemon(c, true, true)

	// start a service by hitting each of the 3 managers
	d1.CreateService(c, simpleTestService, func(s *swarm.Service) ***REMOVED***
		s.Spec.Name = "test1"
	***REMOVED***)
	d2.CreateService(c, simpleTestService, func(s *swarm.Service) ***REMOVED***
		s.Spec.Name = "test2"
	***REMOVED***)
	d3.CreateService(c, simpleTestService, func(s *swarm.Service) ***REMOVED***
		s.Spec.Name = "test3"
	***REMOVED***)

	// 3 services should be started now, because the requests were proxied to leader
	// query each node and make sure it returns 3 services
	for _, d := range []*daemon.Swarm***REMOVED***d1, d2, d3***REMOVED*** ***REMOVED***
		services := d.ListServices(c)
		c.Assert(services, checker.HasLen, 3)
	***REMOVED***
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmLeaderElection(c *check.C) ***REMOVED***
	// Create 3 nodes
	d1 := s.AddDaemon(c, true, true)
	d2 := s.AddDaemon(c, true, true)
	d3 := s.AddDaemon(c, true, true)

	// assert that the first node we made is the leader, and the other two are followers
	c.Assert(d1.GetNode(c, d1.NodeID).ManagerStatus.Leader, checker.True)
	c.Assert(d1.GetNode(c, d2.NodeID).ManagerStatus.Leader, checker.False)
	c.Assert(d1.GetNode(c, d3.NodeID).ManagerStatus.Leader, checker.False)

	d1.Stop(c)

	var (
		leader    *daemon.Swarm   // keep track of leader
		followers []*daemon.Swarm // keep track of followers
	)
	checkLeader := func(nodes ...*daemon.Swarm) checkF ***REMOVED***
		return func(c *check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
			// clear these out before each run
			leader = nil
			followers = nil
			for _, d := range nodes ***REMOVED***
				if d.GetNode(c, d.NodeID).ManagerStatus.Leader ***REMOVED***
					leader = d
				***REMOVED*** else ***REMOVED***
					followers = append(followers, d)
				***REMOVED***
			***REMOVED***

			if leader == nil ***REMOVED***
				return false, check.Commentf("no leader elected")
			***REMOVED***

			return true, check.Commentf("elected %v", leader.ID())
		***REMOVED***
	***REMOVED***

	// wait for an election to occur
	waitAndAssert(c, defaultReconciliationTimeout, checkLeader(d2, d3), checker.True)

	// assert that we have a new leader
	c.Assert(leader, checker.NotNil)

	// Keep track of the current leader, since we want that to be chosen.
	stableleader := leader

	// add the d1, the initial leader, back
	d1.Start(c)

	// TODO(stevvooe): may need to wait for rejoin here

	// wait for possible election
	waitAndAssert(c, defaultReconciliationTimeout, checkLeader(d1, d2, d3), checker.True)
	// pick out the leader and the followers again

	// verify that we still only have 1 leader and 2 followers
	c.Assert(leader, checker.NotNil)
	c.Assert(followers, checker.HasLen, 2)
	// and that after we added d1 back, the leader hasn't changed
	c.Assert(leader.NodeID, checker.Equals, stableleader.NodeID)
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmRaftQuorum(c *check.C) ***REMOVED***
	d1 := s.AddDaemon(c, true, true)
	d2 := s.AddDaemon(c, true, true)
	d3 := s.AddDaemon(c, true, true)

	d1.CreateService(c, simpleTestService)

	d2.Stop(c)

	// make sure there is a leader
	waitAndAssert(c, defaultReconciliationTimeout, d1.CheckLeader, checker.IsNil)

	d1.CreateService(c, simpleTestService, func(s *swarm.Service) ***REMOVED***
		s.Spec.Name = "top1"
	***REMOVED***)

	d3.Stop(c)

	var service swarm.Service
	simpleTestService(&service)
	service.Spec.Name = "top2"
	cli, err := d1.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	// d1 will eventually step down from leader because there is no longer an active quorum, wait for that to happen
	waitAndAssert(c, defaultReconciliationTimeout, func(c *check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
		_, err = cli.ServiceCreate(context.Background(), service.Spec, types.ServiceCreateOptions***REMOVED******REMOVED***)
		return err.Error(), nil
	***REMOVED***, checker.Contains, "Make sure more than half of the managers are online.")

	d2.Start(c)

	// make sure there is a leader
	waitAndAssert(c, defaultReconciliationTimeout, d1.CheckLeader, checker.IsNil)

	d1.CreateService(c, simpleTestService, func(s *swarm.Service) ***REMOVED***
		s.Spec.Name = "top3"
	***REMOVED***)
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmLeaveRemovesContainer(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, true, true)

	instances := 2
	d.CreateService(c, simpleTestService, setInstances(instances))

	id, err := d.Cmd("run", "-d", "busybox", "top")
	c.Assert(err, checker.IsNil)
	id = strings.TrimSpace(id)

	waitAndAssert(c, defaultReconciliationTimeout, d.CheckActiveContainerCount, checker.Equals, instances+1)

	c.Assert(d.Leave(false), checker.NotNil)
	c.Assert(d.Leave(true), checker.IsNil)

	waitAndAssert(c, defaultReconciliationTimeout, d.CheckActiveContainerCount, checker.Equals, 1)

	id2, err := d.Cmd("ps", "-q")
	c.Assert(err, checker.IsNil)
	c.Assert(id, checker.HasPrefix, strings.TrimSpace(id2))
***REMOVED***

// #23629
func (s *DockerSwarmSuite) TestAPISwarmLeaveOnPendingJoin(c *check.C) ***REMOVED***
	testRequires(c, Network)
	s.AddDaemon(c, true, true)
	d2 := s.AddDaemon(c, false, false)

	id, err := d2.Cmd("run", "-d", "busybox", "top")
	c.Assert(err, checker.IsNil)
	id = strings.TrimSpace(id)

	err = d2.Join(swarm.JoinRequest***REMOVED***
		RemoteAddrs: []string***REMOVED***"123.123.123.123:1234"***REMOVED***,
	***REMOVED***)
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), checker.Contains, "Timeout was reached")

	info, err := d2.SwarmInfo()
	c.Assert(err, checker.IsNil)
	c.Assert(info.LocalNodeState, checker.Equals, swarm.LocalNodeStatePending)

	c.Assert(d2.Leave(true), checker.IsNil)

	waitAndAssert(c, defaultReconciliationTimeout, d2.CheckActiveContainerCount, checker.Equals, 1)

	id2, err := d2.Cmd("ps", "-q")
	c.Assert(err, checker.IsNil)
	c.Assert(id, checker.HasPrefix, strings.TrimSpace(id2))
***REMOVED***

// #23705
func (s *DockerSwarmSuite) TestAPISwarmRestoreOnPendingJoin(c *check.C) ***REMOVED***
	testRequires(c, Network)
	d := s.AddDaemon(c, false, false)
	err := d.Join(swarm.JoinRequest***REMOVED***
		RemoteAddrs: []string***REMOVED***"123.123.123.123:1234"***REMOVED***,
	***REMOVED***)
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), checker.Contains, "Timeout was reached")

	waitAndAssert(c, defaultReconciliationTimeout, d.CheckLocalNodeState, checker.Equals, swarm.LocalNodeStatePending)

	d.Stop(c)
	d.Start(c)

	info, err := d.SwarmInfo()
	c.Assert(err, checker.IsNil)
	c.Assert(info.LocalNodeState, checker.Equals, swarm.LocalNodeStateInactive)
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmManagerRestore(c *check.C) ***REMOVED***
	d1 := s.AddDaemon(c, true, true)

	instances := 2
	id := d1.CreateService(c, simpleTestService, setInstances(instances))

	d1.GetService(c, id)
	d1.Stop(c)
	d1.Start(c)
	d1.GetService(c, id)

	d2 := s.AddDaemon(c, true, true)
	d2.GetService(c, id)
	d2.Stop(c)
	d2.Start(c)
	d2.GetService(c, id)

	d3 := s.AddDaemon(c, true, true)
	d3.GetService(c, id)
	d3.Stop(c)
	d3.Start(c)
	d3.GetService(c, id)

	d3.Kill()
	time.Sleep(1 * time.Second) // time to handle signal
	d3.Start(c)
	d3.GetService(c, id)
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmScaleNoRollingUpdate(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, true, true)

	instances := 2
	id := d.CreateService(c, simpleTestService, setInstances(instances))

	waitAndAssert(c, defaultReconciliationTimeout, d.CheckActiveContainerCount, checker.Equals, instances)
	containers := d.ActiveContainers()
	instances = 4
	d.UpdateService(c, d.GetService(c, id), setInstances(instances))
	waitAndAssert(c, defaultReconciliationTimeout, d.CheckActiveContainerCount, checker.Equals, instances)
	containers2 := d.ActiveContainers()

loop0:
	for _, c1 := range containers ***REMOVED***
		for _, c2 := range containers2 ***REMOVED***
			if c1 == c2 ***REMOVED***
				continue loop0
			***REMOVED***
		***REMOVED***
		c.Errorf("container %v not found in new set %#v", c1, containers2)
	***REMOVED***
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmInvalidAddress(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, false, false)
	req := swarm.InitRequest***REMOVED***
		ListenAddr: "",
	***REMOVED***
	res, _, err := request.DoOnHost(d.Sock(), "/swarm/init", request.Method("POST"), request.JSONBody(req))
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusBadRequest)

	req2 := swarm.JoinRequest***REMOVED***
		ListenAddr:  "0.0.0.0:2377",
		RemoteAddrs: []string***REMOVED***""***REMOVED***,
	***REMOVED***
	res, _, err = request.DoOnHost(d.Sock(), "/swarm/join", request.Method("POST"), request.JSONBody(req2))
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusBadRequest)
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmForceNewCluster(c *check.C) ***REMOVED***
	d1 := s.AddDaemon(c, true, true)
	d2 := s.AddDaemon(c, true, true)

	instances := 2
	id := d1.CreateService(c, simpleTestService, setInstances(instances))
	waitAndAssert(c, defaultReconciliationTimeout, reducedCheck(sumAsIntegers, d1.CheckActiveContainerCount, d2.CheckActiveContainerCount), checker.Equals, instances)

	// drain d2, all containers should move to d1
	d1.UpdateNode(c, d2.NodeID, func(n *swarm.Node) ***REMOVED***
		n.Spec.Availability = swarm.NodeAvailabilityDrain
	***REMOVED***)
	waitAndAssert(c, defaultReconciliationTimeout, d1.CheckActiveContainerCount, checker.Equals, instances)
	waitAndAssert(c, defaultReconciliationTimeout, d2.CheckActiveContainerCount, checker.Equals, 0)

	d2.Stop(c)

	c.Assert(d1.Init(swarm.InitRequest***REMOVED***
		ForceNewCluster: true,
		Spec:            swarm.Spec***REMOVED******REMOVED***,
	***REMOVED***), checker.IsNil)

	waitAndAssert(c, defaultReconciliationTimeout, d1.CheckActiveContainerCount, checker.Equals, instances)

	d3 := s.AddDaemon(c, true, true)
	info, err := d3.SwarmInfo()
	c.Assert(err, checker.IsNil)
	c.Assert(info.ControlAvailable, checker.True)
	c.Assert(info.LocalNodeState, checker.Equals, swarm.LocalNodeStateActive)

	instances = 4
	d3.UpdateService(c, d3.GetService(c, id), setInstances(instances))

	waitAndAssert(c, defaultReconciliationTimeout, reducedCheck(sumAsIntegers, d1.CheckActiveContainerCount, d3.CheckActiveContainerCount), checker.Equals, instances)
***REMOVED***

func simpleTestService(s *swarm.Service) ***REMOVED***
	ureplicas := uint64(1)
	restartDelay := time.Duration(100 * time.Millisecond)

	s.Spec = swarm.ServiceSpec***REMOVED***
		TaskTemplate: swarm.TaskSpec***REMOVED***
			ContainerSpec: &swarm.ContainerSpec***REMOVED***
				Image:   "busybox:latest",
				Command: []string***REMOVED***"/bin/top"***REMOVED***,
			***REMOVED***,
			RestartPolicy: &swarm.RestartPolicy***REMOVED***
				Delay: &restartDelay,
			***REMOVED***,
		***REMOVED***,
		Mode: swarm.ServiceMode***REMOVED***
			Replicated: &swarm.ReplicatedService***REMOVED***
				Replicas: &ureplicas,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	s.Spec.Name = "top"
***REMOVED***

func serviceForUpdate(s *swarm.Service) ***REMOVED***
	ureplicas := uint64(1)
	restartDelay := time.Duration(100 * time.Millisecond)

	s.Spec = swarm.ServiceSpec***REMOVED***
		TaskTemplate: swarm.TaskSpec***REMOVED***
			ContainerSpec: &swarm.ContainerSpec***REMOVED***
				Image:   "busybox:latest",
				Command: []string***REMOVED***"/bin/top"***REMOVED***,
			***REMOVED***,
			RestartPolicy: &swarm.RestartPolicy***REMOVED***
				Delay: &restartDelay,
			***REMOVED***,
		***REMOVED***,
		Mode: swarm.ServiceMode***REMOVED***
			Replicated: &swarm.ReplicatedService***REMOVED***
				Replicas: &ureplicas,
			***REMOVED***,
		***REMOVED***,
		UpdateConfig: &swarm.UpdateConfig***REMOVED***
			Parallelism:   2,
			Delay:         4 * time.Second,
			FailureAction: swarm.UpdateFailureActionContinue,
		***REMOVED***,
		RollbackConfig: &swarm.UpdateConfig***REMOVED***
			Parallelism:   3,
			Delay:         4 * time.Second,
			FailureAction: swarm.UpdateFailureActionContinue,
		***REMOVED***,
	***REMOVED***
	s.Spec.Name = "updatetest"
***REMOVED***

func setInstances(replicas int) daemon.ServiceConstructor ***REMOVED***
	ureplicas := uint64(replicas)
	return func(s *swarm.Service) ***REMOVED***
		s.Spec.Mode = swarm.ServiceMode***REMOVED***
			Replicated: &swarm.ReplicatedService***REMOVED***
				Replicas: &ureplicas,
			***REMOVED***,
		***REMOVED***
	***REMOVED***
***REMOVED***

func setUpdateOrder(order string) daemon.ServiceConstructor ***REMOVED***
	return func(s *swarm.Service) ***REMOVED***
		if s.Spec.UpdateConfig == nil ***REMOVED***
			s.Spec.UpdateConfig = &swarm.UpdateConfig***REMOVED******REMOVED***
		***REMOVED***
		s.Spec.UpdateConfig.Order = order
	***REMOVED***
***REMOVED***

func setRollbackOrder(order string) daemon.ServiceConstructor ***REMOVED***
	return func(s *swarm.Service) ***REMOVED***
		if s.Spec.RollbackConfig == nil ***REMOVED***
			s.Spec.RollbackConfig = &swarm.UpdateConfig***REMOVED******REMOVED***
		***REMOVED***
		s.Spec.RollbackConfig.Order = order
	***REMOVED***
***REMOVED***

func setImage(image string) daemon.ServiceConstructor ***REMOVED***
	return func(s *swarm.Service) ***REMOVED***
		if s.Spec.TaskTemplate.ContainerSpec == nil ***REMOVED***
			s.Spec.TaskTemplate.ContainerSpec = &swarm.ContainerSpec***REMOVED******REMOVED***
		***REMOVED***
		s.Spec.TaskTemplate.ContainerSpec.Image = image
	***REMOVED***
***REMOVED***

func setFailureAction(failureAction string) daemon.ServiceConstructor ***REMOVED***
	return func(s *swarm.Service) ***REMOVED***
		s.Spec.UpdateConfig.FailureAction = failureAction
	***REMOVED***
***REMOVED***

func setMaxFailureRatio(maxFailureRatio float32) daemon.ServiceConstructor ***REMOVED***
	return func(s *swarm.Service) ***REMOVED***
		s.Spec.UpdateConfig.MaxFailureRatio = maxFailureRatio
	***REMOVED***
***REMOVED***

func setParallelism(parallelism uint64) daemon.ServiceConstructor ***REMOVED***
	return func(s *swarm.Service) ***REMOVED***
		s.Spec.UpdateConfig.Parallelism = parallelism
	***REMOVED***
***REMOVED***

func setConstraints(constraints []string) daemon.ServiceConstructor ***REMOVED***
	return func(s *swarm.Service) ***REMOVED***
		if s.Spec.TaskTemplate.Placement == nil ***REMOVED***
			s.Spec.TaskTemplate.Placement = &swarm.Placement***REMOVED******REMOVED***
		***REMOVED***
		s.Spec.TaskTemplate.Placement.Constraints = constraints
	***REMOVED***
***REMOVED***

func setPlacementPrefs(prefs []swarm.PlacementPreference) daemon.ServiceConstructor ***REMOVED***
	return func(s *swarm.Service) ***REMOVED***
		if s.Spec.TaskTemplate.Placement == nil ***REMOVED***
			s.Spec.TaskTemplate.Placement = &swarm.Placement***REMOVED******REMOVED***
		***REMOVED***
		s.Spec.TaskTemplate.Placement.Preferences = prefs
	***REMOVED***
***REMOVED***

func setGlobalMode(s *swarm.Service) ***REMOVED***
	s.Spec.Mode = swarm.ServiceMode***REMOVED***
		Global: &swarm.GlobalService***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

func checkClusterHealth(c *check.C, cl []*daemon.Swarm, managerCount, workerCount int) ***REMOVED***
	var totalMCount, totalWCount int

	for _, d := range cl ***REMOVED***
		var (
			info swarm.Info
			err  error
		)

		// check info in a waitAndAssert, because if the cluster doesn't have a leader, `info` will return an error
		checkInfo := func(c *check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
			info, err = d.SwarmInfo()
			return err, check.Commentf("cluster not ready in time")
		***REMOVED***
		waitAndAssert(c, defaultReconciliationTimeout, checkInfo, checker.IsNil)
		if !info.ControlAvailable ***REMOVED***
			totalWCount++
			continue
		***REMOVED***

		var leaderFound bool
		totalMCount++
		var mCount, wCount int

		for _, n := range d.ListNodes(c) ***REMOVED***
			waitReady := func(c *check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
				if n.Status.State == swarm.NodeStateReady ***REMOVED***
					return true, nil
				***REMOVED***
				nn := d.GetNode(c, n.ID)
				n = *nn
				return n.Status.State == swarm.NodeStateReady, check.Commentf("state of node %s, reported by %s", n.ID, d.Info.NodeID)
			***REMOVED***
			waitAndAssert(c, defaultReconciliationTimeout, waitReady, checker.True)

			waitActive := func(c *check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
				if n.Spec.Availability == swarm.NodeAvailabilityActive ***REMOVED***
					return true, nil
				***REMOVED***
				nn := d.GetNode(c, n.ID)
				n = *nn
				return n.Spec.Availability == swarm.NodeAvailabilityActive, check.Commentf("availability of node %s, reported by %s", n.ID, d.Info.NodeID)
			***REMOVED***
			waitAndAssert(c, defaultReconciliationTimeout, waitActive, checker.True)

			if n.Spec.Role == swarm.NodeRoleManager ***REMOVED***
				c.Assert(n.ManagerStatus, checker.NotNil, check.Commentf("manager status of node %s (manager), reported by %s", n.ID, d.Info.NodeID))
				if n.ManagerStatus.Leader ***REMOVED***
					leaderFound = true
				***REMOVED***
				mCount++
			***REMOVED*** else ***REMOVED***
				c.Assert(n.ManagerStatus, checker.IsNil, check.Commentf("manager status of node %s (worker), reported by %s", n.ID, d.Info.NodeID))
				wCount++
			***REMOVED***
		***REMOVED***
		c.Assert(leaderFound, checker.True, check.Commentf("lack of leader reported by node %s", info.NodeID))
		c.Assert(mCount, checker.Equals, managerCount, check.Commentf("managers count reported by node %s", info.NodeID))
		c.Assert(wCount, checker.Equals, workerCount, check.Commentf("workers count reported by node %s", info.NodeID))
	***REMOVED***
	c.Assert(totalMCount, checker.Equals, managerCount)
	c.Assert(totalWCount, checker.Equals, workerCount)
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmRestartCluster(c *check.C) ***REMOVED***
	mCount, wCount := 5, 1

	var nodes []*daemon.Swarm
	for i := 0; i < mCount; i++ ***REMOVED***
		manager := s.AddDaemon(c, true, true)
		info, err := manager.SwarmInfo()
		c.Assert(err, checker.IsNil)
		c.Assert(info.ControlAvailable, checker.True)
		c.Assert(info.LocalNodeState, checker.Equals, swarm.LocalNodeStateActive)
		nodes = append(nodes, manager)
	***REMOVED***

	for i := 0; i < wCount; i++ ***REMOVED***
		worker := s.AddDaemon(c, true, false)
		info, err := worker.SwarmInfo()
		c.Assert(err, checker.IsNil)
		c.Assert(info.ControlAvailable, checker.False)
		c.Assert(info.LocalNodeState, checker.Equals, swarm.LocalNodeStateActive)
		nodes = append(nodes, worker)
	***REMOVED***

	// stop whole cluster
	***REMOVED***
		var wg sync.WaitGroup
		wg.Add(len(nodes))
		errs := make(chan error, len(nodes))

		for _, d := range nodes ***REMOVED***
			go func(daemon *daemon.Swarm) ***REMOVED***
				defer wg.Done()
				if err := daemon.StopWithError(); err != nil ***REMOVED***
					errs <- err
				***REMOVED***
				// FIXME(vdemeester) This is duplicated…
				if root := os.Getenv("DOCKER_REMAP_ROOT"); root != "" ***REMOVED***
					daemon.Root = filepath.Dir(daemon.Root)
				***REMOVED***
			***REMOVED***(d)
		***REMOVED***
		wg.Wait()
		close(errs)
		for err := range errs ***REMOVED***
			c.Assert(err, check.IsNil)
		***REMOVED***
	***REMOVED***

	// start whole cluster
	***REMOVED***
		var wg sync.WaitGroup
		wg.Add(len(nodes))
		errs := make(chan error, len(nodes))

		for _, d := range nodes ***REMOVED***
			go func(daemon *daemon.Swarm) ***REMOVED***
				defer wg.Done()
				if err := daemon.StartWithError("--iptables=false"); err != nil ***REMOVED***
					errs <- err
				***REMOVED***
			***REMOVED***(d)
		***REMOVED***
		wg.Wait()
		close(errs)
		for err := range errs ***REMOVED***
			c.Assert(err, check.IsNil)
		***REMOVED***
	***REMOVED***

	checkClusterHealth(c, nodes, mCount, wCount)
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmServicesUpdateWithName(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, true, true)

	instances := 2
	id := d.CreateService(c, simpleTestService, setInstances(instances))
	waitAndAssert(c, defaultReconciliationTimeout, d.CheckActiveContainerCount, checker.Equals, instances)

	service := d.GetService(c, id)
	instances = 5

	setInstances(instances)(service)
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()
	_, err = cli.ServiceUpdate(context.Background(), service.Spec.Name, service.Version, service.Spec, types.ServiceUpdateOptions***REMOVED******REMOVED***)
	c.Assert(err, checker.IsNil)
	waitAndAssert(c, defaultReconciliationTimeout, d.CheckActiveContainerCount, checker.Equals, instances)
***REMOVED***

// Unlocking an unlocked swarm results in an error
func (s *DockerSwarmSuite) TestAPISwarmUnlockNotLocked(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, true, true)
	err := d.Unlock(swarm.UnlockRequest***REMOVED***UnlockKey: "wrong-key"***REMOVED***)
	c.Assert(err, checker.NotNil)
	c.Assert(err.Error(), checker.Contains, "swarm is not locked")
***REMOVED***

// #29885
func (s *DockerSwarmSuite) TestAPISwarmErrorHandling(c *check.C) ***REMOVED***
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", defaultSwarmPort))
	c.Assert(err, checker.IsNil)
	defer ln.Close()
	d := s.AddDaemon(c, false, false)
	err = d.Init(swarm.InitRequest***REMOVED******REMOVED***)
	c.Assert(err, checker.NotNil)
	c.Assert(err.Error(), checker.Contains, "address already in use")
***REMOVED***

// Test case for 30242, where duplicate networks, with different drivers `bridge` and `overlay`,
// caused both scopes to be `swarm` for `docker network inspect` and `docker network ls`.
// This test makes sure the fixes correctly output scopes instead.
func (s *DockerSwarmSuite) TestAPIDuplicateNetworks(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, true, true)
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	name := "foo"
	networkCreate := types.NetworkCreate***REMOVED***
		CheckDuplicate: false,
	***REMOVED***

	networkCreate.Driver = "bridge"

	n1, err := cli.NetworkCreate(context.Background(), name, networkCreate)
	c.Assert(err, checker.IsNil)

	networkCreate.Driver = "overlay"

	n2, err := cli.NetworkCreate(context.Background(), name, networkCreate)
	c.Assert(err, checker.IsNil)

	r1, err := cli.NetworkInspect(context.Background(), n1.ID, types.NetworkInspectOptions***REMOVED******REMOVED***)
	c.Assert(err, checker.IsNil)
	c.Assert(r1.Scope, checker.Equals, "local")

	r2, err := cli.NetworkInspect(context.Background(), n2.ID, types.NetworkInspectOptions***REMOVED******REMOVED***)
	c.Assert(err, checker.IsNil)
	c.Assert(r2.Scope, checker.Equals, "swarm")
***REMOVED***

// Test case for 30178
func (s *DockerSwarmSuite) TestAPISwarmHealthcheckNone(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, true, true)

	out, err := d.Cmd("network", "create", "-d", "overlay", "lb")
	c.Assert(err, checker.IsNil, check.Commentf(out))

	instances := 1
	d.CreateService(c, simpleTestService, setInstances(instances), func(s *swarm.Service) ***REMOVED***
		if s.Spec.TaskTemplate.ContainerSpec == nil ***REMOVED***
			s.Spec.TaskTemplate.ContainerSpec = &swarm.ContainerSpec***REMOVED******REMOVED***
		***REMOVED***
		s.Spec.TaskTemplate.ContainerSpec.Healthcheck = &container.HealthConfig***REMOVED******REMOVED***
		s.Spec.TaskTemplate.Networks = []swarm.NetworkAttachmentConfig***REMOVED***
			***REMOVED***Target: "lb"***REMOVED***,
		***REMOVED***
	***REMOVED***)

	waitAndAssert(c, defaultReconciliationTimeout, d.CheckActiveContainerCount, checker.Equals, instances)

	containers := d.ActiveContainers()

	out, err = d.Cmd("exec", containers[0], "ping", "-c1", "-W3", "top")
	c.Assert(err, checker.IsNil, check.Commentf(out))
***REMOVED***

func (s *DockerSwarmSuite) TestSwarmRepeatedRootRotation(c *check.C) ***REMOVED***
	m := s.AddDaemon(c, true, true)
	w := s.AddDaemon(c, true, false)

	info, err := m.SwarmInfo()
	c.Assert(err, checker.IsNil)

	currentTrustRoot := info.Cluster.TLSInfo.TrustRoot

	// rotate multiple times
	for i := 0; i < 4; i++ ***REMOVED***
		var cert, key []byte
		if i%2 != 0 ***REMOVED***
			cert, _, key, err = initca.New(&csr.CertificateRequest***REMOVED***
				CN:         "newRoot",
				KeyRequest: csr.NewBasicKeyRequest(),
				CA:         &csr.CAConfig***REMOVED***Expiry: ca.RootCAExpiration***REMOVED***,
			***REMOVED***)
			c.Assert(err, checker.IsNil)
		***REMOVED***
		expectedCert := string(cert)
		m.UpdateSwarm(c, func(s *swarm.Spec) ***REMOVED***
			s.CAConfig.SigningCACert = expectedCert
			s.CAConfig.SigningCAKey = string(key)
			s.CAConfig.ForceRotate++
		***REMOVED***)

		// poll to make sure update succeeds
		var clusterTLSInfo swarm.TLSInfo
		for j := 0; j < 18; j++ ***REMOVED***
			info, err := m.SwarmInfo()
			c.Assert(err, checker.IsNil)

			// the desired CA cert and key is always redacted
			c.Assert(info.Cluster.Spec.CAConfig.SigningCAKey, checker.Equals, "")
			c.Assert(info.Cluster.Spec.CAConfig.SigningCACert, checker.Equals, "")

			clusterTLSInfo = info.Cluster.TLSInfo

			// if root rotation is done and the trust root has changed, we don't have to poll anymore
			if !info.Cluster.RootRotationInProgress && clusterTLSInfo.TrustRoot != currentTrustRoot ***REMOVED***
				break
			***REMOVED***

			// root rotation not done
			time.Sleep(250 * time.Millisecond)
		***REMOVED***
		if cert != nil ***REMOVED***
			c.Assert(clusterTLSInfo.TrustRoot, checker.Equals, expectedCert)
		***REMOVED***
		// could take another second or two for the nodes to trust the new roots after they've all gotten
		// new TLS certificates
		for j := 0; j < 18; j++ ***REMOVED***
			mInfo := m.GetNode(c, m.NodeID).Description.TLSInfo
			wInfo := m.GetNode(c, w.NodeID).Description.TLSInfo

			if mInfo.TrustRoot == clusterTLSInfo.TrustRoot && wInfo.TrustRoot == clusterTLSInfo.TrustRoot ***REMOVED***
				break
			***REMOVED***

			// nodes don't trust root certs yet
			time.Sleep(250 * time.Millisecond)
		***REMOVED***

		c.Assert(m.GetNode(c, m.NodeID).Description.TLSInfo, checker.DeepEquals, clusterTLSInfo)
		c.Assert(m.GetNode(c, w.NodeID).Description.TLSInfo, checker.DeepEquals, clusterTLSInfo)
		currentTrustRoot = clusterTLSInfo.TrustRoot
	***REMOVED***
***REMOVED***

func (s *DockerSwarmSuite) TestAPINetworkInspectWithScope(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, true, true)

	name := "foo"
	networkCreateRequest := types.NetworkCreateRequest***REMOVED***
		Name: name,
	***REMOVED***

	var n types.NetworkCreateResponse
	networkCreateRequest.NetworkCreate.Driver = "overlay"

	status, out, err := d.SockRequest("POST", "/networks/create", networkCreateRequest)
	c.Assert(err, checker.IsNil, check.Commentf(string(out)))
	c.Assert(status, checker.Equals, http.StatusCreated, check.Commentf(string(out)))
	c.Assert(json.Unmarshal(out, &n), checker.IsNil)

	var r types.NetworkResource

	status, body, err := d.SockRequest("GET", "/networks/"+name, nil)
	c.Assert(err, checker.IsNil, check.Commentf(string(out)))
	c.Assert(status, checker.Equals, http.StatusOK, check.Commentf(string(out)))
	c.Assert(json.Unmarshal(body, &r), checker.IsNil)
	c.Assert(r.Scope, checker.Equals, "swarm")
	c.Assert(r.ID, checker.Equals, n.ID)

	v := url.Values***REMOVED******REMOVED***
	v.Set("scope", "local")

	status, _, err = d.SockRequest("GET", "/networks/"+name+"?"+v.Encode(), nil)
	c.Assert(err, checker.IsNil, check.Commentf(string(out)))
	c.Assert(status, checker.Equals, http.StatusNotFound, check.Commentf(string(out)))
***REMOVED***
