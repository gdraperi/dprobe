package service

import (
	"runtime"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/docker/docker/integration-cli/request"
	"github.com/gotestyourself/gotestyourself/poll"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestCreateServiceMultipleTimes(t *testing.T) ***REMOVED***
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
	serviceSpec := swarmServiceSpec("TestService", instances)
	serviceSpec.TaskTemplate.Networks = append(serviceSpec.TaskTemplate.Networks, swarm.NetworkAttachmentConfig***REMOVED***Target: overlayName***REMOVED***)

	serviceResp, err := client.ServiceCreate(context.Background(), serviceSpec, types.ServiceCreateOptions***REMOVED***
		QueryRegistry: false,
	***REMOVED***)
	require.NoError(t, err)

	pollSettings := func(config *poll.Settings) ***REMOVED***
		// It takes about ~25s to finish the multi services creation in this case per the pratical observation on arm64/arm platform
		if runtime.GOARCH == "arm64" || runtime.GOARCH == "arm" ***REMOVED***
			config.Timeout = 30 * time.Second
			config.Delay = 100 * time.Millisecond
		***REMOVED***
	***REMOVED***

	serviceID := serviceResp.ID
	poll.WaitOn(t, serviceRunningTasksCount(client, serviceID, instances), pollSettings)

	_, _, err = client.ServiceInspectWithRaw(context.Background(), serviceID, types.ServiceInspectOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

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

func TestCreateWithDuplicateNetworkNames(t *testing.T) ***REMOVED***
	defer setupTest(t)()
	d := newSwarm(t)
	defer d.Stop(t)
	client, err := request.NewClientForHost(d.Sock())
	require.NoError(t, err)

	name := "foo"
	networkCreate := types.NetworkCreate***REMOVED***
		CheckDuplicate: false,
		Driver:         "bridge",
	***REMOVED***

	n1, err := client.NetworkCreate(context.Background(), name, networkCreate)
	require.NoError(t, err)

	n2, err := client.NetworkCreate(context.Background(), name, networkCreate)
	require.NoError(t, err)

	// Dupliates with name but with different driver
	networkCreate.Driver = "overlay"
	n3, err := client.NetworkCreate(context.Background(), name, networkCreate)
	require.NoError(t, err)

	// Create Service with the same name
	var instances uint64 = 1
	serviceSpec := swarmServiceSpec("top", instances)

	serviceSpec.TaskTemplate.Networks = append(serviceSpec.TaskTemplate.Networks, swarm.NetworkAttachmentConfig***REMOVED***Target: name***REMOVED***)

	service, err := client.ServiceCreate(context.Background(), serviceSpec, types.ServiceCreateOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	poll.WaitOn(t, serviceRunningTasksCount(client, service.ID, instances))

	resp, _, err := client.ServiceInspectWithRaw(context.Background(), service.ID, types.ServiceInspectOptions***REMOVED******REMOVED***)
	require.NoError(t, err)
	assert.Equal(t, n3.ID, resp.Spec.TaskTemplate.Networks[0].Target)

	// Remove Service
	err = client.ServiceRemove(context.Background(), service.ID)
	require.NoError(t, err)

	// Make sure task has been destroyed.
	poll.WaitOn(t, serviceIsRemoved(client, service.ID))

	// Remove networks
	err = client.NetworkRemove(context.Background(), n3.ID)
	require.NoError(t, err)

	err = client.NetworkRemove(context.Background(), n2.ID)
	require.NoError(t, err)

	err = client.NetworkRemove(context.Background(), n1.ID)
	require.NoError(t, err)

	// Make sure networks have been destroyed.
	poll.WaitOn(t, networkIsRemoved(client, n3.ID), poll.WithTimeout(1*time.Minute), poll.WithDelay(10*time.Second))
	poll.WaitOn(t, networkIsRemoved(client, n2.ID), poll.WithTimeout(1*time.Minute), poll.WithDelay(10*time.Second))
	poll.WaitOn(t, networkIsRemoved(client, n1.ID), poll.WithTimeout(1*time.Minute), poll.WithDelay(10*time.Second))
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

func networkIsRemoved(client client.NetworkAPIClient, networkID string) func(log poll.LogT) poll.Result ***REMOVED***
	return func(log poll.LogT) poll.Result ***REMOVED***
		_, err := client.NetworkInspect(context.Background(), networkID, types.NetworkInspectOptions***REMOVED******REMOVED***)
		if err == nil ***REMOVED***
			return poll.Continue("waiting for network %s to be removed", networkID)
		***REMOVED***
		return poll.Success()
	***REMOVED***
***REMOVED***
