// +build !windows

package main

import (
	"fmt"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/api/types/swarm/runtime"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/daemon"
	"github.com/docker/docker/integration-cli/fixtures/plugin"
	"github.com/go-check/check"
	"golang.org/x/net/context"
	"golang.org/x/sys/unix"
)

func setPortConfig(portConfig []swarm.PortConfig) daemon.ServiceConstructor ***REMOVED***
	return func(s *swarm.Service) ***REMOVED***
		if s.Spec.EndpointSpec == nil ***REMOVED***
			s.Spec.EndpointSpec = &swarm.EndpointSpec***REMOVED******REMOVED***
		***REMOVED***
		s.Spec.EndpointSpec.Ports = portConfig
	***REMOVED***
***REMOVED***

func (s *DockerSwarmSuite) TestAPIServiceUpdatePort(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, true, true)

	// Create a service with a port mapping of 8080:8081.
	portConfig := []swarm.PortConfig***REMOVED******REMOVED***TargetPort: 8081, PublishedPort: 8080***REMOVED******REMOVED***
	serviceID := d.CreateService(c, simpleTestService, setInstances(1), setPortConfig(portConfig))
	waitAndAssert(c, defaultReconciliationTimeout, d.CheckActiveContainerCount, checker.Equals, 1)

	// Update the service: changed the port mapping from 8080:8081 to 8082:8083.
	updatedPortConfig := []swarm.PortConfig***REMOVED******REMOVED***TargetPort: 8083, PublishedPort: 8082***REMOVED******REMOVED***
	remoteService := d.GetService(c, serviceID)
	d.UpdateService(c, remoteService, setPortConfig(updatedPortConfig))

	// Inspect the service and verify port mapping.
	updatedService := d.GetService(c, serviceID)
	c.Assert(updatedService.Spec.EndpointSpec, check.NotNil)
	c.Assert(len(updatedService.Spec.EndpointSpec.Ports), check.Equals, 1)
	c.Assert(updatedService.Spec.EndpointSpec.Ports[0].TargetPort, check.Equals, uint32(8083))
	c.Assert(updatedService.Spec.EndpointSpec.Ports[0].PublishedPort, check.Equals, uint32(8082))
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmServicesEmptyList(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, true, true)

	services := d.ListServices(c)
	c.Assert(services, checker.NotNil)
	c.Assert(len(services), checker.Equals, 0, check.Commentf("services: %#v", services))
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmServicesCreate(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, true, true)

	instances := 2
	id := d.CreateService(c, simpleTestService, setInstances(instances))
	waitAndAssert(c, defaultReconciliationTimeout, d.CheckActiveContainerCount, checker.Equals, instances)

	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	options := types.ServiceInspectOptions***REMOVED***InsertDefaults: true***REMOVED***

	// insertDefaults inserts UpdateConfig when service is fetched by ID
	resp, _, err := cli.ServiceInspectWithRaw(context.Background(), id, options)
	out := fmt.Sprintf("%+v", resp)
	c.Assert(err, checker.IsNil)
	c.Assert(out, checker.Contains, "UpdateConfig")

	// insertDefaults inserts UpdateConfig when service is fetched by ID
	resp, _, err = cli.ServiceInspectWithRaw(context.Background(), "top", options)
	out = fmt.Sprintf("%+v", resp)
	c.Assert(err, checker.IsNil)
	c.Assert(string(out), checker.Contains, "UpdateConfig")

	service := d.GetService(c, id)
	instances = 5
	d.UpdateService(c, service, setInstances(instances))
	waitAndAssert(c, defaultReconciliationTimeout, d.CheckActiveContainerCount, checker.Equals, instances)

	d.RemoveService(c, service.ID)
	waitAndAssert(c, defaultReconciliationTimeout, d.CheckActiveContainerCount, checker.Equals, 0)
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmServicesMultipleAgents(c *check.C) ***REMOVED***
	d1 := s.AddDaemon(c, true, true)
	d2 := s.AddDaemon(c, true, false)
	d3 := s.AddDaemon(c, true, false)

	time.Sleep(1 * time.Second) // make sure all daemons are ready to accept tasks

	instances := 9
	id := d1.CreateService(c, simpleTestService, setInstances(instances))

	waitAndAssert(c, defaultReconciliationTimeout, d1.CheckActiveContainerCount, checker.GreaterThan, 0)
	waitAndAssert(c, defaultReconciliationTimeout, d2.CheckActiveContainerCount, checker.GreaterThan, 0)
	waitAndAssert(c, defaultReconciliationTimeout, d3.CheckActiveContainerCount, checker.GreaterThan, 0)

	waitAndAssert(c, defaultReconciliationTimeout, reducedCheck(sumAsIntegers, d1.CheckActiveContainerCount, d2.CheckActiveContainerCount, d3.CheckActiveContainerCount), checker.Equals, instances)

	// reconciliation on d2 node down
	d2.Stop(c)

	waitAndAssert(c, defaultReconciliationTimeout, reducedCheck(sumAsIntegers, d1.CheckActiveContainerCount, d3.CheckActiveContainerCount), checker.Equals, instances)

	// test downscaling
	instances = 5
	d1.UpdateService(c, d1.GetService(c, id), setInstances(instances))
	waitAndAssert(c, defaultReconciliationTimeout, reducedCheck(sumAsIntegers, d1.CheckActiveContainerCount, d3.CheckActiveContainerCount), checker.Equals, instances)

***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmServicesCreateGlobal(c *check.C) ***REMOVED***
	d1 := s.AddDaemon(c, true, true)
	d2 := s.AddDaemon(c, true, false)
	d3 := s.AddDaemon(c, true, false)

	d1.CreateService(c, simpleTestService, setGlobalMode)

	waitAndAssert(c, defaultReconciliationTimeout, d1.CheckActiveContainerCount, checker.Equals, 1)
	waitAndAssert(c, defaultReconciliationTimeout, d2.CheckActiveContainerCount, checker.Equals, 1)
	waitAndAssert(c, defaultReconciliationTimeout, d3.CheckActiveContainerCount, checker.Equals, 1)

	d4 := s.AddDaemon(c, true, false)
	d5 := s.AddDaemon(c, true, false)

	waitAndAssert(c, defaultReconciliationTimeout, d4.CheckActiveContainerCount, checker.Equals, 1)
	waitAndAssert(c, defaultReconciliationTimeout, d5.CheckActiveContainerCount, checker.Equals, 1)
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmServicesUpdate(c *check.C) ***REMOVED***
	const nodeCount = 3
	var daemons [nodeCount]*daemon.Swarm
	for i := 0; i < nodeCount; i++ ***REMOVED***
		daemons[i] = s.AddDaemon(c, true, i == 0)
	***REMOVED***
	// wait for nodes ready
	waitAndAssert(c, 5*time.Second, daemons[0].CheckNodeReadyCount, checker.Equals, nodeCount)

	// service image at start
	image1 := "busybox:latest"
	// target image in update
	image2 := "busybox:test"

	// create a different tag
	for _, d := range daemons ***REMOVED***
		out, err := d.Cmd("tag", image1, image2)
		c.Assert(err, checker.IsNil, check.Commentf(out))
	***REMOVED***

	// create service
	instances := 5
	parallelism := 2
	rollbackParallelism := 3
	id := daemons[0].CreateService(c, serviceForUpdate, setInstances(instances))

	// wait for tasks ready
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckRunningTaskImages, checker.DeepEquals,
		map[string]int***REMOVED***image1: instances***REMOVED***)

	// issue service update
	service := daemons[0].GetService(c, id)
	daemons[0].UpdateService(c, service, setImage(image2))

	// first batch
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckRunningTaskImages, checker.DeepEquals,
		map[string]int***REMOVED***image1: instances - parallelism, image2: parallelism***REMOVED***)

	// 2nd batch
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckRunningTaskImages, checker.DeepEquals,
		map[string]int***REMOVED***image1: instances - 2*parallelism, image2: 2 * parallelism***REMOVED***)

	// 3nd batch
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckRunningTaskImages, checker.DeepEquals,
		map[string]int***REMOVED***image2: instances***REMOVED***)

	// Roll back to the previous version. This uses the CLI because
	// rollback used to be a client-side operation.
	out, err := daemons[0].Cmd("service", "update", "--detach", "--rollback", id)
	c.Assert(err, checker.IsNil, check.Commentf(out))

	// first batch
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckRunningTaskImages, checker.DeepEquals,
		map[string]int***REMOVED***image2: instances - rollbackParallelism, image1: rollbackParallelism***REMOVED***)

	// 2nd batch
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckRunningTaskImages, checker.DeepEquals,
		map[string]int***REMOVED***image1: instances***REMOVED***)
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmServicesUpdateStartFirst(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, true, true)

	// service image at start
	image1 := "busybox:latest"
	// target image in update
	image2 := "testhealth:latest"

	// service started from this image won't pass health check
	_, _, err := d.BuildImageWithOut(image2,
		`FROM busybox
		HEALTHCHECK --interval=1s --timeout=30s --retries=1024 \
		  CMD cat /status`,
		true)
	c.Check(err, check.IsNil)

	// create service
	instances := 5
	parallelism := 2
	rollbackParallelism := 3
	id := d.CreateService(c, serviceForUpdate, setInstances(instances), setUpdateOrder(swarm.UpdateOrderStartFirst), setRollbackOrder(swarm.UpdateOrderStartFirst))

	checkStartingTasks := func(expected int) []swarm.Task ***REMOVED***
		var startingTasks []swarm.Task
		waitAndAssert(c, defaultReconciliationTimeout, func(c *check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
			tasks := d.GetServiceTasks(c, id)
			startingTasks = nil
			for _, t := range tasks ***REMOVED***
				if t.Status.State == swarm.TaskStateStarting ***REMOVED***
					startingTasks = append(startingTasks, t)
				***REMOVED***
			***REMOVED***
			return startingTasks, nil
		***REMOVED***, checker.HasLen, expected)

		return startingTasks
	***REMOVED***

	makeTasksHealthy := func(tasks []swarm.Task) ***REMOVED***
		for _, t := range tasks ***REMOVED***
			containerID := t.Status.ContainerStatus.ContainerID
			d.Cmd("exec", containerID, "touch", "/status")
		***REMOVED***
	***REMOVED***

	// wait for tasks ready
	waitAndAssert(c, defaultReconciliationTimeout, d.CheckRunningTaskImages, checker.DeepEquals,
		map[string]int***REMOVED***image1: instances***REMOVED***)

	// issue service update
	service := d.GetService(c, id)
	d.UpdateService(c, service, setImage(image2))

	// first batch

	// The old tasks should be running, and the new ones should be starting.
	startingTasks := checkStartingTasks(parallelism)

	waitAndAssert(c, defaultReconciliationTimeout, d.CheckRunningTaskImages, checker.DeepEquals,
		map[string]int***REMOVED***image1: instances***REMOVED***)

	// make it healthy
	makeTasksHealthy(startingTasks)

	waitAndAssert(c, defaultReconciliationTimeout, d.CheckRunningTaskImages, checker.DeepEquals,
		map[string]int***REMOVED***image1: instances - parallelism, image2: parallelism***REMOVED***)

	// 2nd batch

	// The old tasks should be running, and the new ones should be starting.
	startingTasks = checkStartingTasks(parallelism)

	waitAndAssert(c, defaultReconciliationTimeout, d.CheckRunningTaskImages, checker.DeepEquals,
		map[string]int***REMOVED***image1: instances - parallelism, image2: parallelism***REMOVED***)

	// make it healthy
	makeTasksHealthy(startingTasks)

	waitAndAssert(c, defaultReconciliationTimeout, d.CheckRunningTaskImages, checker.DeepEquals,
		map[string]int***REMOVED***image1: instances - 2*parallelism, image2: 2 * parallelism***REMOVED***)

	// 3nd batch

	// The old tasks should be running, and the new ones should be starting.
	startingTasks = checkStartingTasks(1)

	waitAndAssert(c, defaultReconciliationTimeout, d.CheckRunningTaskImages, checker.DeepEquals,
		map[string]int***REMOVED***image1: instances - 2*parallelism, image2: 2 * parallelism***REMOVED***)

	// make it healthy
	makeTasksHealthy(startingTasks)

	waitAndAssert(c, defaultReconciliationTimeout, d.CheckRunningTaskImages, checker.DeepEquals,
		map[string]int***REMOVED***image2: instances***REMOVED***)

	// Roll back to the previous version. This uses the CLI because
	// rollback is a client-side operation.
	out, err := d.Cmd("service", "update", "--detach", "--rollback", id)
	c.Assert(err, checker.IsNil, check.Commentf(out))

	// first batch
	waitAndAssert(c, defaultReconciliationTimeout, d.CheckRunningTaskImages, checker.DeepEquals,
		map[string]int***REMOVED***image2: instances - rollbackParallelism, image1: rollbackParallelism***REMOVED***)

	// 2nd batch
	waitAndAssert(c, defaultReconciliationTimeout, d.CheckRunningTaskImages, checker.DeepEquals,
		map[string]int***REMOVED***image1: instances***REMOVED***)
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmServicesFailedUpdate(c *check.C) ***REMOVED***
	const nodeCount = 3
	var daemons [nodeCount]*daemon.Swarm
	for i := 0; i < nodeCount; i++ ***REMOVED***
		daemons[i] = s.AddDaemon(c, true, i == 0)
	***REMOVED***
	// wait for nodes ready
	waitAndAssert(c, 5*time.Second, daemons[0].CheckNodeReadyCount, checker.Equals, nodeCount)

	// service image at start
	image1 := "busybox:latest"
	// target image in update
	image2 := "busybox:badtag"

	// create service
	instances := 5
	id := daemons[0].CreateService(c, serviceForUpdate, setInstances(instances))

	// wait for tasks ready
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckRunningTaskImages, checker.DeepEquals,
		map[string]int***REMOVED***image1: instances***REMOVED***)

	// issue service update
	service := daemons[0].GetService(c, id)
	daemons[0].UpdateService(c, service, setImage(image2), setFailureAction(swarm.UpdateFailureActionPause), setMaxFailureRatio(0.25), setParallelism(1))

	// should update 2 tasks and then pause
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckServiceUpdateState(id), checker.Equals, swarm.UpdateStatePaused)
	v, _ := daemons[0].CheckServiceRunningTasks(id)(c)
	c.Assert(v, checker.Equals, instances-2)

	// Roll back to the previous version. This uses the CLI because
	// rollback used to be a client-side operation.
	out, err := daemons[0].Cmd("service", "update", "--detach", "--rollback", id)
	c.Assert(err, checker.IsNil, check.Commentf(out))

	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckRunningTaskImages, checker.DeepEquals,
		map[string]int***REMOVED***image1: instances***REMOVED***)
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmServiceConstraintRole(c *check.C) ***REMOVED***
	const nodeCount = 3
	var daemons [nodeCount]*daemon.Swarm
	for i := 0; i < nodeCount; i++ ***REMOVED***
		daemons[i] = s.AddDaemon(c, true, i == 0)
	***REMOVED***
	// wait for nodes ready
	waitAndAssert(c, 5*time.Second, daemons[0].CheckNodeReadyCount, checker.Equals, nodeCount)

	// create service
	constraints := []string***REMOVED***"node.role==worker"***REMOVED***
	instances := 3
	id := daemons[0].CreateService(c, simpleTestService, setConstraints(constraints), setInstances(instances))
	// wait for tasks ready
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckServiceRunningTasks(id), checker.Equals, instances)
	// validate tasks are running on worker nodes
	tasks := daemons[0].GetServiceTasks(c, id)
	for _, task := range tasks ***REMOVED***
		node := daemons[0].GetNode(c, task.NodeID)
		c.Assert(node.Spec.Role, checker.Equals, swarm.NodeRoleWorker)
	***REMOVED***
	//remove service
	daemons[0].RemoveService(c, id)

	// create service
	constraints = []string***REMOVED***"node.role!=worker"***REMOVED***
	id = daemons[0].CreateService(c, simpleTestService, setConstraints(constraints), setInstances(instances))
	// wait for tasks ready
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckServiceRunningTasks(id), checker.Equals, instances)
	tasks = daemons[0].GetServiceTasks(c, id)
	// validate tasks are running on manager nodes
	for _, task := range tasks ***REMOVED***
		node := daemons[0].GetNode(c, task.NodeID)
		c.Assert(node.Spec.Role, checker.Equals, swarm.NodeRoleManager)
	***REMOVED***
	//remove service
	daemons[0].RemoveService(c, id)

	// create service
	constraints = []string***REMOVED***"node.role==nosuchrole"***REMOVED***
	id = daemons[0].CreateService(c, simpleTestService, setConstraints(constraints), setInstances(instances))
	// wait for tasks created
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckServiceTasks(id), checker.Equals, instances)
	// let scheduler try
	time.Sleep(250 * time.Millisecond)
	// validate tasks are not assigned to any node
	tasks = daemons[0].GetServiceTasks(c, id)
	for _, task := range tasks ***REMOVED***
		c.Assert(task.NodeID, checker.Equals, "")
	***REMOVED***
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmServiceConstraintLabel(c *check.C) ***REMOVED***
	const nodeCount = 3
	var daemons [nodeCount]*daemon.Swarm
	for i := 0; i < nodeCount; i++ ***REMOVED***
		daemons[i] = s.AddDaemon(c, true, i == 0)
	***REMOVED***
	// wait for nodes ready
	waitAndAssert(c, 5*time.Second, daemons[0].CheckNodeReadyCount, checker.Equals, nodeCount)
	nodes := daemons[0].ListNodes(c)
	c.Assert(len(nodes), checker.Equals, nodeCount)

	// add labels to nodes
	daemons[0].UpdateNode(c, nodes[0].ID, func(n *swarm.Node) ***REMOVED***
		n.Spec.Annotations.Labels = map[string]string***REMOVED***
			"security": "high",
		***REMOVED***
	***REMOVED***)
	for i := 1; i < nodeCount; i++ ***REMOVED***
		daemons[0].UpdateNode(c, nodes[i].ID, func(n *swarm.Node) ***REMOVED***
			n.Spec.Annotations.Labels = map[string]string***REMOVED***
				"security": "low",
			***REMOVED***
		***REMOVED***)
	***REMOVED***

	// create service
	instances := 3
	constraints := []string***REMOVED***"node.labels.security==high"***REMOVED***
	id := daemons[0].CreateService(c, simpleTestService, setConstraints(constraints), setInstances(instances))
	// wait for tasks ready
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckServiceRunningTasks(id), checker.Equals, instances)
	tasks := daemons[0].GetServiceTasks(c, id)
	// validate all tasks are running on nodes[0]
	for _, task := range tasks ***REMOVED***
		c.Assert(task.NodeID, checker.Equals, nodes[0].ID)
	***REMOVED***
	//remove service
	daemons[0].RemoveService(c, id)

	// create service
	constraints = []string***REMOVED***"node.labels.security!=high"***REMOVED***
	id = daemons[0].CreateService(c, simpleTestService, setConstraints(constraints), setInstances(instances))
	// wait for tasks ready
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckServiceRunningTasks(id), checker.Equals, instances)
	tasks = daemons[0].GetServiceTasks(c, id)
	// validate all tasks are NOT running on nodes[0]
	for _, task := range tasks ***REMOVED***
		c.Assert(task.NodeID, checker.Not(checker.Equals), nodes[0].ID)
	***REMOVED***
	//remove service
	daemons[0].RemoveService(c, id)

	constraints = []string***REMOVED***"node.labels.security==medium"***REMOVED***
	id = daemons[0].CreateService(c, simpleTestService, setConstraints(constraints), setInstances(instances))
	// wait for tasks created
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckServiceTasks(id), checker.Equals, instances)
	// let scheduler try
	time.Sleep(250 * time.Millisecond)
	tasks = daemons[0].GetServiceTasks(c, id)
	// validate tasks are not assigned
	for _, task := range tasks ***REMOVED***
		c.Assert(task.NodeID, checker.Equals, "")
	***REMOVED***
	//remove service
	daemons[0].RemoveService(c, id)

	// multiple constraints
	constraints = []string***REMOVED***
		"node.labels.security==high",
		fmt.Sprintf("node.id==%s", nodes[1].ID),
	***REMOVED***
	id = daemons[0].CreateService(c, simpleTestService, setConstraints(constraints), setInstances(instances))
	// wait for tasks created
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckServiceTasks(id), checker.Equals, instances)
	// let scheduler try
	time.Sleep(250 * time.Millisecond)
	tasks = daemons[0].GetServiceTasks(c, id)
	// validate tasks are not assigned
	for _, task := range tasks ***REMOVED***
		c.Assert(task.NodeID, checker.Equals, "")
	***REMOVED***
	// make nodes[1] fulfills the constraints
	daemons[0].UpdateNode(c, nodes[1].ID, func(n *swarm.Node) ***REMOVED***
		n.Spec.Annotations.Labels = map[string]string***REMOVED***
			"security": "high",
		***REMOVED***
	***REMOVED***)
	// wait for tasks ready
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckServiceRunningTasks(id), checker.Equals, instances)
	tasks = daemons[0].GetServiceTasks(c, id)
	for _, task := range tasks ***REMOVED***
		c.Assert(task.NodeID, checker.Equals, nodes[1].ID)
	***REMOVED***
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmServicePlacementPrefs(c *check.C) ***REMOVED***
	const nodeCount = 3
	var daemons [nodeCount]*daemon.Swarm
	for i := 0; i < nodeCount; i++ ***REMOVED***
		daemons[i] = s.AddDaemon(c, true, i == 0)
	***REMOVED***
	// wait for nodes ready
	waitAndAssert(c, 5*time.Second, daemons[0].CheckNodeReadyCount, checker.Equals, nodeCount)
	nodes := daemons[0].ListNodes(c)
	c.Assert(len(nodes), checker.Equals, nodeCount)

	// add labels to nodes
	daemons[0].UpdateNode(c, nodes[0].ID, func(n *swarm.Node) ***REMOVED***
		n.Spec.Annotations.Labels = map[string]string***REMOVED***
			"rack": "a",
		***REMOVED***
	***REMOVED***)
	for i := 1; i < nodeCount; i++ ***REMOVED***
		daemons[0].UpdateNode(c, nodes[i].ID, func(n *swarm.Node) ***REMOVED***
			n.Spec.Annotations.Labels = map[string]string***REMOVED***
				"rack": "b",
			***REMOVED***
		***REMOVED***)
	***REMOVED***

	// create service
	instances := 4
	prefs := []swarm.PlacementPreference***REMOVED******REMOVED***Spread: &swarm.SpreadOver***REMOVED***SpreadDescriptor: "node.labels.rack"***REMOVED******REMOVED******REMOVED***
	id := daemons[0].CreateService(c, simpleTestService, setPlacementPrefs(prefs), setInstances(instances))
	// wait for tasks ready
	waitAndAssert(c, defaultReconciliationTimeout, daemons[0].CheckServiceRunningTasks(id), checker.Equals, instances)
	tasks := daemons[0].GetServiceTasks(c, id)
	// validate all tasks are running on nodes[0]
	tasksOnNode := make(map[string]int)
	for _, task := range tasks ***REMOVED***
		tasksOnNode[task.NodeID]++
	***REMOVED***
	c.Assert(tasksOnNode[nodes[0].ID], checker.Equals, 2)
	c.Assert(tasksOnNode[nodes[1].ID], checker.Equals, 1)
	c.Assert(tasksOnNode[nodes[2].ID], checker.Equals, 1)
***REMOVED***

func (s *DockerSwarmSuite) TestAPISwarmServicesStateReporting(c *check.C) ***REMOVED***
	testRequires(c, SameHostDaemon)
	testRequires(c, DaemonIsLinux)

	d1 := s.AddDaemon(c, true, true)
	d2 := s.AddDaemon(c, true, true)
	d3 := s.AddDaemon(c, true, false)

	time.Sleep(1 * time.Second) // make sure all daemons are ready to accept

	instances := 9
	d1.CreateService(c, simpleTestService, setInstances(instances))

	waitAndAssert(c, defaultReconciliationTimeout, reducedCheck(sumAsIntegers, d1.CheckActiveContainerCount, d2.CheckActiveContainerCount, d3.CheckActiveContainerCount), checker.Equals, instances)

	getContainers := func() map[string]*daemon.Swarm ***REMOVED***
		m := make(map[string]*daemon.Swarm)
		for _, d := range []*daemon.Swarm***REMOVED***d1, d2, d3***REMOVED*** ***REMOVED***
			for _, id := range d.ActiveContainers() ***REMOVED***
				m[id] = d
			***REMOVED***
		***REMOVED***
		return m
	***REMOVED***

	containers := getContainers()
	c.Assert(containers, checker.HasLen, instances)
	var toRemove string
	for i := range containers ***REMOVED***
		toRemove = i
	***REMOVED***

	_, err := containers[toRemove].Cmd("stop", toRemove)
	c.Assert(err, checker.IsNil)

	waitAndAssert(c, defaultReconciliationTimeout, reducedCheck(sumAsIntegers, d1.CheckActiveContainerCount, d2.CheckActiveContainerCount, d3.CheckActiveContainerCount), checker.Equals, instances)

	containers2 := getContainers()
	c.Assert(containers2, checker.HasLen, instances)
	for i := range containers ***REMOVED***
		if i == toRemove ***REMOVED***
			c.Assert(containers2[i], checker.IsNil)
		***REMOVED*** else ***REMOVED***
			c.Assert(containers2[i], checker.NotNil)
		***REMOVED***
	***REMOVED***

	containers = containers2
	for i := range containers ***REMOVED***
		toRemove = i
	***REMOVED***

	// try with killing process outside of docker
	pidStr, err := containers[toRemove].Cmd("inspect", "-f", "***REMOVED******REMOVED***.State.Pid***REMOVED******REMOVED***", toRemove)
	c.Assert(err, checker.IsNil)
	pid, err := strconv.Atoi(strings.TrimSpace(pidStr))
	c.Assert(err, checker.IsNil)
	c.Assert(unix.Kill(pid, unix.SIGKILL), checker.IsNil)

	time.Sleep(time.Second) // give some time to handle the signal

	waitAndAssert(c, defaultReconciliationTimeout, reducedCheck(sumAsIntegers, d1.CheckActiveContainerCount, d2.CheckActiveContainerCount, d3.CheckActiveContainerCount), checker.Equals, instances)

	containers2 = getContainers()
	c.Assert(containers2, checker.HasLen, instances)
	for i := range containers ***REMOVED***
		if i == toRemove ***REMOVED***
			c.Assert(containers2[i], checker.IsNil)
		***REMOVED*** else ***REMOVED***
			c.Assert(containers2[i], checker.NotNil)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Test plugins deployed via swarm services
func (s *DockerSwarmSuite) TestAPISwarmServicesPlugin(c *check.C) ***REMOVED***
	testRequires(c, ExperimentalDaemon, DaemonIsLinux, IsAmd64)

	reg := setupRegistry(c, false, "", "")
	defer reg.Close()

	repo := path.Join(privateRegistryURL, "swarm", "test:v1")
	repo2 := path.Join(privateRegistryURL, "swarm", "test:v2")
	name := "test"

	err := plugin.CreateInRegistry(context.Background(), repo, nil)
	c.Assert(err, checker.IsNil, check.Commentf("failed to create plugin"))
	err = plugin.CreateInRegistry(context.Background(), repo2, nil)
	c.Assert(err, checker.IsNil, check.Commentf("failed to create plugin"))

	d1 := s.AddDaemon(c, true, true)
	d2 := s.AddDaemon(c, true, true)
	d3 := s.AddDaemon(c, true, false)

	makePlugin := func(repo, name string, constraints []string) func(*swarm.Service) ***REMOVED***
		return func(s *swarm.Service) ***REMOVED***
			s.Spec.TaskTemplate.Runtime = "plugin"
			s.Spec.TaskTemplate.PluginSpec = &runtime.PluginSpec***REMOVED***
				Name:   name,
				Remote: repo,
			***REMOVED***
			if constraints != nil ***REMOVED***
				s.Spec.TaskTemplate.Placement = &swarm.Placement***REMOVED***
					Constraints: constraints,
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	id := d1.CreateService(c, makePlugin(repo, name, nil))
	waitAndAssert(c, defaultReconciliationTimeout, d1.CheckPluginRunning(name), checker.True)
	waitAndAssert(c, defaultReconciliationTimeout, d2.CheckPluginRunning(name), checker.True)
	waitAndAssert(c, defaultReconciliationTimeout, d3.CheckPluginRunning(name), checker.True)

	service := d1.GetService(c, id)
	d1.UpdateService(c, service, makePlugin(repo2, name, nil))
	waitAndAssert(c, defaultReconciliationTimeout, d1.CheckPluginImage(name), checker.Equals, repo2)
	waitAndAssert(c, defaultReconciliationTimeout, d2.CheckPluginImage(name), checker.Equals, repo2)
	waitAndAssert(c, defaultReconciliationTimeout, d3.CheckPluginImage(name), checker.Equals, repo2)
	waitAndAssert(c, defaultReconciliationTimeout, d1.CheckPluginRunning(name), checker.True)
	waitAndAssert(c, defaultReconciliationTimeout, d2.CheckPluginRunning(name), checker.True)
	waitAndAssert(c, defaultReconciliationTimeout, d3.CheckPluginRunning(name), checker.True)

	d1.RemoveService(c, id)
	waitAndAssert(c, defaultReconciliationTimeout, d1.CheckPluginRunning(name), checker.False)
	waitAndAssert(c, defaultReconciliationTimeout, d2.CheckPluginRunning(name), checker.False)
	waitAndAssert(c, defaultReconciliationTimeout, d3.CheckPluginRunning(name), checker.False)

	// constrain to managers only
	id = d1.CreateService(c, makePlugin(repo, name, []string***REMOVED***"node.role==manager"***REMOVED***))
	waitAndAssert(c, defaultReconciliationTimeout, d1.CheckPluginRunning(name), checker.True)
	waitAndAssert(c, defaultReconciliationTimeout, d2.CheckPluginRunning(name), checker.True)
	waitAndAssert(c, defaultReconciliationTimeout, d3.CheckPluginRunning(name), checker.False) // Not a manager, not running it
	d1.RemoveService(c, id)
	waitAndAssert(c, defaultReconciliationTimeout, d1.CheckPluginRunning(name), checker.False)
	waitAndAssert(c, defaultReconciliationTimeout, d2.CheckPluginRunning(name), checker.False)
	waitAndAssert(c, defaultReconciliationTimeout, d3.CheckPluginRunning(name), checker.False)

	// with no name
	id = d1.CreateService(c, makePlugin(repo, "", nil))
	waitAndAssert(c, defaultReconciliationTimeout, d1.CheckPluginRunning(repo), checker.True)
	waitAndAssert(c, defaultReconciliationTimeout, d2.CheckPluginRunning(repo), checker.True)
	waitAndAssert(c, defaultReconciliationTimeout, d3.CheckPluginRunning(repo), checker.True)
	d1.RemoveService(c, id)
	waitAndAssert(c, defaultReconciliationTimeout, d1.CheckPluginRunning(repo), checker.False)
	waitAndAssert(c, defaultReconciliationTimeout, d2.CheckPluginRunning(repo), checker.False)
	waitAndAssert(c, defaultReconciliationTimeout, d3.CheckPluginRunning(repo), checker.False)
***REMOVED***
