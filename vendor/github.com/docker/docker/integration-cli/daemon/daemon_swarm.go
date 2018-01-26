package daemon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/go-check/check"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// Swarm is a test daemon with helpers for participating in a swarm.
type Swarm struct ***REMOVED***
	*Daemon
	swarm.Info
	Port       int
	ListenAddr string
***REMOVED***

// Init initializes a new swarm cluster.
func (d *Swarm) Init(req swarm.InitRequest) error ***REMOVED***
	if req.ListenAddr == "" ***REMOVED***
		req.ListenAddr = d.ListenAddr
	***REMOVED***
	cli, err := d.NewClient()
	if err != nil ***REMOVED***
		return fmt.Errorf("initializing swarm: failed to create client %v", err)
	***REMOVED***
	defer cli.Close()
	_, err = cli.SwarmInit(context.Background(), req)
	if err != nil ***REMOVED***
		return fmt.Errorf("initializing swarm: %v", err)
	***REMOVED***
	info, err := d.SwarmInfo()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	d.Info = info
	return nil
***REMOVED***

// Join joins a daemon to an existing cluster.
func (d *Swarm) Join(req swarm.JoinRequest) error ***REMOVED***
	if req.ListenAddr == "" ***REMOVED***
		req.ListenAddr = d.ListenAddr
	***REMOVED***
	cli, err := d.NewClient()
	if err != nil ***REMOVED***
		return fmt.Errorf("joining swarm: failed to create client %v", err)
	***REMOVED***
	defer cli.Close()
	err = cli.SwarmJoin(context.Background(), req)
	if err != nil ***REMOVED***
		return fmt.Errorf("joining swarm: %v", err)
	***REMOVED***
	info, err := d.SwarmInfo()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	d.Info = info
	return nil
***REMOVED***

// Leave forces daemon to leave current cluster.
func (d *Swarm) Leave(force bool) error ***REMOVED***
	cli, err := d.NewClient()
	if err != nil ***REMOVED***
		return fmt.Errorf("leaving swarm: failed to create client %v", err)
	***REMOVED***
	defer cli.Close()
	err = cli.SwarmLeave(context.Background(), force)
	if err != nil ***REMOVED***
		err = fmt.Errorf("leaving swarm: %v", err)
	***REMOVED***
	return err
***REMOVED***

// SwarmInfo returns the swarm information of the daemon
func (d *Swarm) SwarmInfo() (swarm.Info, error) ***REMOVED***
	cli, err := d.NewClient()
	if err != nil ***REMOVED***
		return swarm.Info***REMOVED******REMOVED***, fmt.Errorf("get swarm info: %v", err)
	***REMOVED***

	info, err := cli.Info(context.Background())
	if err != nil ***REMOVED***
		return swarm.Info***REMOVED******REMOVED***, fmt.Errorf("get swarm info: %v", err)
	***REMOVED***

	return info.Swarm, nil
***REMOVED***

// Unlock tries to unlock a locked swarm
func (d *Swarm) Unlock(req swarm.UnlockRequest) error ***REMOVED***
	cli, err := d.NewClient()
	if err != nil ***REMOVED***
		return fmt.Errorf("unlocking swarm: failed to create client %v", err)
	***REMOVED***
	defer cli.Close()
	err = cli.SwarmUnlock(context.Background(), req)
	if err != nil ***REMOVED***
		err = errors.Wrap(err, "unlocking swarm")
	***REMOVED***
	return err
***REMOVED***

// ServiceConstructor defines a swarm service constructor function
type ServiceConstructor func(*swarm.Service)

// NodeConstructor defines a swarm node constructor
type NodeConstructor func(*swarm.Node)

// SecretConstructor defines a swarm secret constructor
type SecretConstructor func(*swarm.Secret)

// ConfigConstructor defines a swarm config constructor
type ConfigConstructor func(*swarm.Config)

// SpecConstructor defines a swarm spec constructor
type SpecConstructor func(*swarm.Spec)

// CreateServiceWithOptions creates a swarm service given the specified service constructors
// and auth config
func (d *Swarm) CreateServiceWithOptions(c *check.C, opts types.ServiceCreateOptions, f ...ServiceConstructor) string ***REMOVED***
	var service swarm.Service
	for _, fn := range f ***REMOVED***
		fn(&service)
	***REMOVED***

	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := cli.ServiceCreate(ctx, service.Spec, opts)
	c.Assert(err, checker.IsNil)
	return res.ID
***REMOVED***

// CreateService creates a swarm service given the specified service constructor
func (d *Swarm) CreateService(c *check.C, f ...ServiceConstructor) string ***REMOVED***
	return d.CreateServiceWithOptions(c, types.ServiceCreateOptions***REMOVED******REMOVED***, f...)
***REMOVED***

// GetService returns the swarm service corresponding to the specified id
func (d *Swarm) GetService(c *check.C, id string) *swarm.Service ***REMOVED***
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	service, _, err := cli.ServiceInspectWithRaw(context.Background(), id, types.ServiceInspectOptions***REMOVED******REMOVED***)
	c.Assert(err, checker.IsNil)
	return &service
***REMOVED***

// GetServiceTasks returns the swarm tasks for the specified service
func (d *Swarm) GetServiceTasks(c *check.C, service string) []swarm.Task ***REMOVED***
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	filterArgs := filters.NewArgs()
	filterArgs.Add("desired-state", "running")
	filterArgs.Add("service", service)

	options := types.TaskListOptions***REMOVED***
		Filters: filterArgs,
	***REMOVED***

	tasks, err := cli.TaskList(context.Background(), options)
	c.Assert(err, checker.IsNil)
	return tasks
***REMOVED***

// CheckServiceTasksInState returns the number of tasks with a matching state,
// and optional message substring.
func (d *Swarm) CheckServiceTasksInState(service string, state swarm.TaskState, message string) func(*check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
	return func(c *check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
		tasks := d.GetServiceTasks(c, service)
		var count int
		for _, task := range tasks ***REMOVED***
			if task.Status.State == state ***REMOVED***
				if message == "" || strings.Contains(task.Status.Message, message) ***REMOVED***
					count++
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return count, nil
	***REMOVED***
***REMOVED***

// CheckServiceTasksInStateWithError returns the number of tasks with a matching state,
// and optional message substring.
func (d *Swarm) CheckServiceTasksInStateWithError(service string, state swarm.TaskState, errorMessage string) func(*check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
	return func(c *check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
		tasks := d.GetServiceTasks(c, service)
		var count int
		for _, task := range tasks ***REMOVED***
			if task.Status.State == state ***REMOVED***
				if errorMessage == "" || strings.Contains(task.Status.Err, errorMessage) ***REMOVED***
					count++
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return count, nil
	***REMOVED***
***REMOVED***

// CheckServiceRunningTasks returns the number of running tasks for the specified service
func (d *Swarm) CheckServiceRunningTasks(service string) func(*check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
	return d.CheckServiceTasksInState(service, swarm.TaskStateRunning, "")
***REMOVED***

// CheckServiceUpdateState returns the current update state for the specified service
func (d *Swarm) CheckServiceUpdateState(service string) func(*check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
	return func(c *check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
		service := d.GetService(c, service)
		if service.UpdateStatus == nil ***REMOVED***
			return "", nil
		***REMOVED***
		return service.UpdateStatus.State, nil
	***REMOVED***
***REMOVED***

// CheckPluginRunning returns the runtime state of the plugin
func (d *Swarm) CheckPluginRunning(plugin string) func(c *check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
	return func(c *check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
		status, out, err := d.SockRequest("GET", "/plugins/"+plugin+"/json", nil)
		c.Assert(err, checker.IsNil, check.Commentf(string(out)))
		if status != http.StatusOK ***REMOVED***
			return false, nil
		***REMOVED***

		var p types.Plugin
		c.Assert(json.Unmarshal(out, &p), checker.IsNil, check.Commentf(string(out)))

		return p.Enabled, check.Commentf("%+v", p)
	***REMOVED***
***REMOVED***

// CheckPluginImage returns the runtime state of the plugin
func (d *Swarm) CheckPluginImage(plugin string) func(c *check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
	return func(c *check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
		status, out, err := d.SockRequest("GET", "/plugins/"+plugin+"/json", nil)
		c.Assert(err, checker.IsNil, check.Commentf(string(out)))
		if status != http.StatusOK ***REMOVED***
			return false, nil
		***REMOVED***

		var p types.Plugin
		c.Assert(json.Unmarshal(out, &p), checker.IsNil, check.Commentf(string(out)))
		return p.PluginReference, check.Commentf("%+v", p)
	***REMOVED***
***REMOVED***

// CheckServiceTasks returns the number of tasks for the specified service
func (d *Swarm) CheckServiceTasks(service string) func(*check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
	return func(c *check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
		tasks := d.GetServiceTasks(c, service)
		return len(tasks), nil
	***REMOVED***
***REMOVED***

// CheckRunningTaskNetworks returns the number of times each network is referenced from a task.
func (d *Swarm) CheckRunningTaskNetworks(c *check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	filterArgs := filters.NewArgs()
	filterArgs.Add("desired-state", "running")

	options := types.TaskListOptions***REMOVED***
		Filters: filterArgs,
	***REMOVED***

	tasks, err := cli.TaskList(context.Background(), options)
	c.Assert(err, checker.IsNil)

	result := make(map[string]int)
	for _, task := range tasks ***REMOVED***
		for _, network := range task.Spec.Networks ***REMOVED***
			result[network.Target]++
		***REMOVED***
	***REMOVED***
	return result, nil
***REMOVED***

// CheckRunningTaskImages returns the times each image is running as a task.
func (d *Swarm) CheckRunningTaskImages(c *check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	filterArgs := filters.NewArgs()
	filterArgs.Add("desired-state", "running")

	options := types.TaskListOptions***REMOVED***
		Filters: filterArgs,
	***REMOVED***

	tasks, err := cli.TaskList(context.Background(), options)
	c.Assert(err, checker.IsNil)

	result := make(map[string]int)
	for _, task := range tasks ***REMOVED***
		if task.Status.State == swarm.TaskStateRunning && task.Spec.ContainerSpec != nil ***REMOVED***
			result[task.Spec.ContainerSpec.Image]++
		***REMOVED***
	***REMOVED***
	return result, nil
***REMOVED***

// CheckNodeReadyCount returns the number of ready node on the swarm
func (d *Swarm) CheckNodeReadyCount(c *check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
	nodes := d.ListNodes(c)
	var readyCount int
	for _, node := range nodes ***REMOVED***
		if node.Status.State == swarm.NodeStateReady ***REMOVED***
			readyCount++
		***REMOVED***
	***REMOVED***
	return readyCount, nil
***REMOVED***

// GetTask returns the swarm task identified by the specified id
func (d *Swarm) GetTask(c *check.C, id string) swarm.Task ***REMOVED***
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	task, _, err := cli.TaskInspectWithRaw(context.Background(), id)
	c.Assert(err, checker.IsNil)
	return task
***REMOVED***

// UpdateService updates a swarm service with the specified service constructor
func (d *Swarm) UpdateService(c *check.C, service *swarm.Service, f ...ServiceConstructor) ***REMOVED***
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	for _, fn := range f ***REMOVED***
		fn(service)
	***REMOVED***

	_, err = cli.ServiceUpdate(context.Background(), service.ID, service.Version, service.Spec, types.ServiceUpdateOptions***REMOVED******REMOVED***)
	c.Assert(err, checker.IsNil)
***REMOVED***

// RemoveService removes the specified service
func (d *Swarm) RemoveService(c *check.C, id string) ***REMOVED***
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	err = cli.ServiceRemove(context.Background(), id)
	c.Assert(err, checker.IsNil)
***REMOVED***

// GetNode returns a swarm node identified by the specified id
func (d *Swarm) GetNode(c *check.C, id string) *swarm.Node ***REMOVED***
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	node, _, err := cli.NodeInspectWithRaw(context.Background(), id)
	c.Assert(err, checker.IsNil)
	c.Assert(node.ID, checker.Equals, id)
	return &node
***REMOVED***

// RemoveNode removes the specified node
func (d *Swarm) RemoveNode(c *check.C, id string, force bool) ***REMOVED***
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	options := types.NodeRemoveOptions***REMOVED***
		Force: force,
	***REMOVED***
	err = cli.NodeRemove(context.Background(), id, options)
	c.Assert(err, checker.IsNil)
***REMOVED***

// UpdateNode updates a swarm node with the specified node constructor
func (d *Swarm) UpdateNode(c *check.C, id string, f ...NodeConstructor) ***REMOVED***
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	for i := 0; ; i++ ***REMOVED***
		node := d.GetNode(c, id)
		for _, fn := range f ***REMOVED***
			fn(node)
		***REMOVED***

		err = cli.NodeUpdate(context.Background(), node.ID, node.Version, node.Spec)
		if i < 10 && err != nil && strings.Contains(err.Error(), "update out of sequence") ***REMOVED***
			time.Sleep(100 * time.Millisecond)
			continue
		***REMOVED***
		c.Assert(err, checker.IsNil)
		return
	***REMOVED***
***REMOVED***

// ListNodes returns the list of the current swarm nodes
func (d *Swarm) ListNodes(c *check.C) []swarm.Node ***REMOVED***
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	nodes, err := cli.NodeList(context.Background(), types.NodeListOptions***REMOVED******REMOVED***)
	c.Assert(err, checker.IsNil)

	return nodes
***REMOVED***

// ListServices returns the list of the current swarm services
func (d *Swarm) ListServices(c *check.C) []swarm.Service ***REMOVED***
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	services, err := cli.ServiceList(context.Background(), types.ServiceListOptions***REMOVED******REMOVED***)
	c.Assert(err, checker.IsNil)
	return services
***REMOVED***

// CreateSecret creates a secret given the specified spec
func (d *Swarm) CreateSecret(c *check.C, secretSpec swarm.SecretSpec) string ***REMOVED***
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	scr, err := cli.SecretCreate(context.Background(), secretSpec)
	c.Assert(err, checker.IsNil)

	return scr.ID
***REMOVED***

// ListSecrets returns the list of the current swarm secrets
func (d *Swarm) ListSecrets(c *check.C) []swarm.Secret ***REMOVED***
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	secrets, err := cli.SecretList(context.Background(), types.SecretListOptions***REMOVED******REMOVED***)
	c.Assert(err, checker.IsNil)
	return secrets
***REMOVED***

// GetSecret returns a swarm secret identified by the specified id
func (d *Swarm) GetSecret(c *check.C, id string) *swarm.Secret ***REMOVED***
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	secret, _, err := cli.SecretInspectWithRaw(context.Background(), id)
	c.Assert(err, checker.IsNil)
	return &secret
***REMOVED***

// DeleteSecret removes the swarm secret identified by the specified id
func (d *Swarm) DeleteSecret(c *check.C, id string) ***REMOVED***
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	err = cli.SecretRemove(context.Background(), id)
	c.Assert(err, checker.IsNil)
***REMOVED***

// UpdateSecret updates the swarm secret identified by the specified id
// Currently, only label update is supported.
func (d *Swarm) UpdateSecret(c *check.C, id string, f ...SecretConstructor) ***REMOVED***
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	secret := d.GetSecret(c, id)
	for _, fn := range f ***REMOVED***
		fn(secret)
	***REMOVED***

	err = cli.SecretUpdate(context.Background(), secret.ID, secret.Version, secret.Spec)

	c.Assert(err, checker.IsNil)
***REMOVED***

// CreateConfig creates a config given the specified spec
func (d *Swarm) CreateConfig(c *check.C, configSpec swarm.ConfigSpec) string ***REMOVED***
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	scr, err := cli.ConfigCreate(context.Background(), configSpec)
	c.Assert(err, checker.IsNil)
	return scr.ID
***REMOVED***

// ListConfigs returns the list of the current swarm configs
func (d *Swarm) ListConfigs(c *check.C) []swarm.Config ***REMOVED***
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	configs, err := cli.ConfigList(context.Background(), types.ConfigListOptions***REMOVED******REMOVED***)
	c.Assert(err, checker.IsNil)
	return configs
***REMOVED***

// GetConfig returns a swarm config identified by the specified id
func (d *Swarm) GetConfig(c *check.C, id string) *swarm.Config ***REMOVED***
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	config, _, err := cli.ConfigInspectWithRaw(context.Background(), id)
	c.Assert(err, checker.IsNil)
	return &config
***REMOVED***

// DeleteConfig removes the swarm config identified by the specified id
func (d *Swarm) DeleteConfig(c *check.C, id string) ***REMOVED***
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	err = cli.ConfigRemove(context.Background(), id)
	c.Assert(err, checker.IsNil)
***REMOVED***

// UpdateConfig updates the swarm config identified by the specified id
// Currently, only label update is supported.
func (d *Swarm) UpdateConfig(c *check.C, id string, f ...ConfigConstructor) ***REMOVED***
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	config := d.GetConfig(c, id)
	for _, fn := range f ***REMOVED***
		fn(config)
	***REMOVED***

	err = cli.ConfigUpdate(context.Background(), config.ID, config.Version, config.Spec)
	c.Assert(err, checker.IsNil)
***REMOVED***

// GetSwarm returns the current swarm object
func (d *Swarm) GetSwarm(c *check.C) swarm.Swarm ***REMOVED***
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	sw, err := cli.SwarmInspect(context.Background())
	c.Assert(err, checker.IsNil)
	return sw
***REMOVED***

// UpdateSwarm updates the current swarm object with the specified spec constructors
func (d *Swarm) UpdateSwarm(c *check.C, f ...SpecConstructor) ***REMOVED***
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	sw := d.GetSwarm(c)
	for _, fn := range f ***REMOVED***
		fn(&sw.Spec)
	***REMOVED***

	err = cli.SwarmUpdate(context.Background(), sw.Version, sw.Spec, swarm.UpdateFlags***REMOVED******REMOVED***)
	c.Assert(err, checker.IsNil)
***REMOVED***

// RotateTokens update the swarm to rotate tokens
func (d *Swarm) RotateTokens(c *check.C) ***REMOVED***
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	sw, err := cli.SwarmInspect(context.Background())
	c.Assert(err, checker.IsNil)

	flags := swarm.UpdateFlags***REMOVED***
		RotateManagerToken: true,
		RotateWorkerToken:  true,
	***REMOVED***

	err = cli.SwarmUpdate(context.Background(), sw.Version, sw.Spec, flags)
	c.Assert(err, checker.IsNil)
***REMOVED***

// JoinTokens returns the current swarm join tokens
func (d *Swarm) JoinTokens(c *check.C) swarm.JoinTokens ***REMOVED***
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	sw, err := cli.SwarmInspect(context.Background())
	c.Assert(err, checker.IsNil)
	return sw.JoinTokens
***REMOVED***

// CheckLocalNodeState returns the current swarm node state
func (d *Swarm) CheckLocalNodeState(c *check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
	info, err := d.SwarmInfo()
	c.Assert(err, checker.IsNil)
	return info.LocalNodeState, nil
***REMOVED***

// CheckControlAvailable returns the current swarm control available
func (d *Swarm) CheckControlAvailable(c *check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
	info, err := d.SwarmInfo()
	c.Assert(err, checker.IsNil)
	c.Assert(info.LocalNodeState, checker.Equals, swarm.LocalNodeStateActive)
	return info.ControlAvailable, nil
***REMOVED***

// CheckLeader returns whether there is a leader on the swarm or not
func (d *Swarm) CheckLeader(c *check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
	cli, err := d.NewClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	errList := check.Commentf("could not get node list")

	ls, err := cli.NodeList(context.Background(), types.NodeListOptions***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		return err, errList
	***REMOVED***

	for _, node := range ls ***REMOVED***
		if node.ManagerStatus != nil && node.ManagerStatus.Leader ***REMOVED***
			return nil, nil
		***REMOVED***
	***REMOVED***
	return fmt.Errorf("no leader"), check.Commentf("could not find leader")
***REMOVED***

// CmdRetryOutOfSequence tries the specified command against the current daemon for 10 times
func (d *Swarm) CmdRetryOutOfSequence(args ...string) (string, error) ***REMOVED***
	for i := 0; ; i++ ***REMOVED***
		out, err := d.Cmd(args...)
		if err != nil ***REMOVED***
			if strings.Contains(out, "update out of sequence") ***REMOVED***
				if i < 10 ***REMOVED***
					continue
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return out, err
	***REMOVED***
***REMOVED***
