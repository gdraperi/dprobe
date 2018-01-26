package network

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/docker/docker/integration-cli/daemon"
	"github.com/docker/docker/integration-cli/request"
	"github.com/gotestyourself/gotestyourself/poll"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

const defaultSwarmPort = 2477
const dockerdBinary = "dockerd"

func TestInspectNetwork(t *testing.T) ***REMOVED***
	defer setupTest(t)()
	d := newSwarm(t)
	defer d.Stop(t)
	client, err := request.NewClientForHost(d.Sock())
	require.NoError(t, err)

	overlayName := "overlay1"
	networkCreate := types.NetworkCreate***REMOVED***
		CheckDuplicate: true,
		Driver:         "overlay",
	***REMOVED***

	netResp, err := client.NetworkCreate(context.Background(), overlayName, networkCreate)
	require.NoError(t, err)
	overlayID := netResp.ID

	var instances uint64 = 4
	serviceName := "TestService"
	serviceSpec := swarmServiceSpec(serviceName, instances)
	serviceSpec.TaskTemplate.Networks = append(serviceSpec.TaskTemplate.Networks, swarm.NetworkAttachmentConfig***REMOVED***Target: overlayName***REMOVED***)

	serviceResp, err := client.ServiceCreate(context.Background(), serviceSpec, types.ServiceCreateOptions***REMOVED***
		QueryRegistry: false,
	***REMOVED***)
	require.NoError(t, err)

	pollSettings := func(config *poll.Settings) ***REMOVED***
		if runtime.GOARCH == "arm" ***REMOVED***
			config.Timeout = 30 * time.Second
			config.Delay = 100 * time.Millisecond
		***REMOVED***
	***REMOVED***

	serviceID := serviceResp.ID
	poll.WaitOn(t, serviceRunningTasksCount(client, serviceID, instances), pollSettings)

	_, _, err = client.ServiceInspectWithRaw(context.Background(), serviceID, types.ServiceInspectOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	// Test inspect verbose with full NetworkID
	networkVerbose, err := client.NetworkInspect(context.Background(), overlayID, types.NetworkInspectOptions***REMOVED***
		Verbose: true,
	***REMOVED***)
	require.NoError(t, err)
	require.True(t, validNetworkVerbose(networkVerbose, serviceName, instances))

	// Test inspect verbose with partial NetworkID
	networkVerbose, err = client.NetworkInspect(context.Background(), overlayID[0:11], types.NetworkInspectOptions***REMOVED***
		Verbose: true,
	***REMOVED***)
	require.NoError(t, err)
	require.True(t, validNetworkVerbose(networkVerbose, serviceName, instances))

	// Test inspect verbose with Network name and swarm scope
	networkVerbose, err = client.NetworkInspect(context.Background(), overlayName, types.NetworkInspectOptions***REMOVED***
		Verbose: true,
		Scope:   "swarm",
	***REMOVED***)
	require.NoError(t, err)
	require.True(t, validNetworkVerbose(networkVerbose, serviceName, instances))

	err = client.ServiceRemove(context.Background(), serviceID)
	require.NoError(t, err)

	poll.WaitOn(t, serviceIsRemoved(client, serviceID), pollSettings)
	poll.WaitOn(t, noTasks(client), pollSettings)

	serviceResp, err = client.ServiceCreate(context.Background(), serviceSpec, types.ServiceCreateOptions***REMOVED***
		QueryRegistry: false,
	***REMOVED***)
	require.NoError(t, err)

	serviceID2 := serviceResp.ID
	poll.WaitOn(t, serviceRunningTasksCount(client, serviceID2, instances), pollSettings)

	err = client.ServiceRemove(context.Background(), serviceID2)
	require.NoError(t, err)

	poll.WaitOn(t, serviceIsRemoved(client, serviceID2), pollSettings)
	poll.WaitOn(t, noTasks(client), pollSettings)

	err = client.NetworkRemove(context.Background(), overlayID)
	require.NoError(t, err)

	poll.WaitOn(t, networkIsRemoved(client, overlayID), poll.WithTimeout(1*time.Minute), poll.WithDelay(10*time.Second))
***REMOVED***

func newSwarm(t *testing.T) *daemon.Swarm ***REMOVED***
	d := &daemon.Swarm***REMOVED***
		Daemon: daemon.New(t, "", dockerdBinary, daemon.Config***REMOVED***
			Experimental: testEnv.DaemonInfo.ExperimentalBuild,
		***REMOVED***),
		// TODO: better method of finding an unused port
		Port: defaultSwarmPort,
	***REMOVED***
	// TODO: move to a NewSwarm constructor
	d.ListenAddr = fmt.Sprintf("0.0.0.0:%d", d.Port)

	// avoid networking conflicts
	args := []string***REMOVED***"--iptables=false", "--swarm-default-advertise-addr=lo"***REMOVED***
	d.StartWithBusybox(t, args...)

	require.NoError(t, d.Init(swarm.InitRequest***REMOVED******REMOVED***))
	return d
***REMOVED***

func swarmServiceSpec(name string, replicas uint64) swarm.ServiceSpec ***REMOVED***
	return swarm.ServiceSpec***REMOVED***
		Annotations: swarm.Annotations***REMOVED***
			Name: name,
		***REMOVED***,
		TaskTemplate: swarm.TaskSpec***REMOVED***
			ContainerSpec: &swarm.ContainerSpec***REMOVED***
				Image:   "busybox:latest",
				Command: []string***REMOVED***"/bin/top"***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		Mode: swarm.ServiceMode***REMOVED***
			Replicated: &swarm.ReplicatedService***REMOVED***
				Replicas: &replicas,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func serviceRunningTasksCount(client client.ServiceAPIClient, serviceID string, instances uint64) func(log poll.LogT) poll.Result ***REMOVED***
	return func(log poll.LogT) poll.Result ***REMOVED***
		filter := filters.NewArgs()
		filter.Add("service", serviceID)
		tasks, err := client.TaskList(context.Background(), types.TaskListOptions***REMOVED***
			Filters: filter,
		***REMOVED***)
		switch ***REMOVED***
		case err != nil:
			return poll.Error(err)
		case len(tasks) == int(instances):
			for _, task := range tasks ***REMOVED***
				if task.Status.State != swarm.TaskStateRunning ***REMOVED***
					return poll.Continue("waiting for tasks to enter run state")
				***REMOVED***
			***REMOVED***
			return poll.Success()
		default:
			return poll.Continue("task count at %d waiting for %d", len(tasks), instances)
		***REMOVED***
	***REMOVED***
***REMOVED***

func networkIsRemoved(client client.NetworkAPIClient, networkID string) func(log poll.LogT) poll.Result ***REMOVED***
	return func(log poll.LogT) poll.Result ***REMOVED***
		_, err := client.NetworkInspect(context.Background(), networkID, types.NetworkInspectOptions***REMOVED******REMOVED***)
		if err == nil ***REMOVED***
			return poll.Continue("waiting for network %s to be removed", networkID)
		***REMOVED***
		return poll.Success()
	***REMOVED***
***REMOVED***

func serviceIsRemoved(client client.ServiceAPIClient, serviceID string) func(log poll.LogT) poll.Result ***REMOVED***
	return func(log poll.LogT) poll.Result ***REMOVED***
		filter := filters.NewArgs()
		filter.Add("service", serviceID)
		_, err := client.TaskList(context.Background(), types.TaskListOptions***REMOVED***
			Filters: filter,
		***REMOVED***)
		if err == nil ***REMOVED***
			return poll.Continue("waiting for service %s to be deleted", serviceID)
		***REMOVED***
		return poll.Success()
	***REMOVED***
***REMOVED***

func noTasks(client client.ServiceAPIClient) func(log poll.LogT) poll.Result ***REMOVED***
	return func(log poll.LogT) poll.Result ***REMOVED***
		filter := filters.NewArgs()
		tasks, err := client.TaskList(context.Background(), types.TaskListOptions***REMOVED***
			Filters: filter,
		***REMOVED***)
		switch ***REMOVED***
		case err != nil:
			return poll.Error(err)
		case len(tasks) == 0:
			return poll.Success()
		default:
			return poll.Continue("task count at %d waiting for 0", len(tasks))
		***REMOVED***
	***REMOVED***
***REMOVED***

// Check to see if Service and Tasks info are part of the inspect verbose response
func validNetworkVerbose(network types.NetworkResource, service string, instances uint64) bool ***REMOVED***
	if service, ok := network.Services[service]; ok ***REMOVED***
		if len(service.Tasks) == int(instances) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
