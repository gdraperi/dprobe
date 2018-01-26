package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/docker/docker/integration-cli/daemon"
	"github.com/docker/docker/integration-cli/request"
	"github.com/gotestyourself/gotestyourself/poll"
	"github.com/gotestyourself/gotestyourself/skip"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestInspect(t *testing.T) ***REMOVED***
	skip.IfCondition(t, !testEnv.IsLocalDaemon())
	defer setupTest(t)()
	d := newSwarm(t)
	defer d.Stop(t)
	client, err := request.NewClientForHost(d.Sock())
	require.NoError(t, err)

	var before = time.Now()
	var instances uint64 = 2
	serviceSpec := fullSwarmServiceSpec("test-service-inspect", instances)

	ctx := context.Background()
	resp, err := client.ServiceCreate(ctx, serviceSpec, types.ServiceCreateOptions***REMOVED***
		QueryRegistry: false,
	***REMOVED***)
	require.NoError(t, err)

	id := resp.ID
	poll.WaitOn(t, serviceContainerCount(client, id, instances))

	service, _, err := client.ServiceInspectWithRaw(ctx, id, types.ServiceInspectOptions***REMOVED******REMOVED***)
	require.NoError(t, err)
	assert.Equal(t, serviceSpec, service.Spec)
	assert.Equal(t, uint64(11), service.Meta.Version.Index)
	assert.Equal(t, id, service.ID)
	assert.WithinDuration(t, before, service.CreatedAt, 30*time.Second)
	assert.WithinDuration(t, before, service.UpdatedAt, 30*time.Second)
***REMOVED***

func fullSwarmServiceSpec(name string, replicas uint64) swarm.ServiceSpec ***REMOVED***
	restartDelay := 100 * time.Millisecond
	maxAttempts := uint64(4)

	return swarm.ServiceSpec***REMOVED***
		Annotations: swarm.Annotations***REMOVED***
			Name: name,
			Labels: map[string]string***REMOVED***
				"service-label": "service-label-value",
			***REMOVED***,
		***REMOVED***,
		TaskTemplate: swarm.TaskSpec***REMOVED***
			ContainerSpec: &swarm.ContainerSpec***REMOVED***
				Image:           "busybox:latest",
				Labels:          map[string]string***REMOVED***"container-label": "container-value"***REMOVED***,
				Command:         []string***REMOVED***"/bin/top"***REMOVED***,
				Args:            []string***REMOVED***"-u", "root"***REMOVED***,
				Hostname:        "hostname",
				Env:             []string***REMOVED***"envvar=envvalue"***REMOVED***,
				Dir:             "/work",
				User:            "root",
				StopSignal:      "SIGINT",
				StopGracePeriod: &restartDelay,
				Hosts:           []string***REMOVED***"8.8.8.8  google"***REMOVED***,
				DNSConfig: &swarm.DNSConfig***REMOVED***
					Nameservers: []string***REMOVED***"8.8.8.8"***REMOVED***,
					Search:      []string***REMOVED***"somedomain"***REMOVED***,
				***REMOVED***,
				Isolation: container.IsolationDefault,
			***REMOVED***,
			RestartPolicy: &swarm.RestartPolicy***REMOVED***
				Delay:       &restartDelay,
				Condition:   swarm.RestartPolicyConditionOnFailure,
				MaxAttempts: &maxAttempts,
			***REMOVED***,
			Runtime: swarm.RuntimeContainer,
		***REMOVED***,
		Mode: swarm.ServiceMode***REMOVED***
			Replicated: &swarm.ReplicatedService***REMOVED***
				Replicas: &replicas,
			***REMOVED***,
		***REMOVED***,
		UpdateConfig: &swarm.UpdateConfig***REMOVED***
			Parallelism:     2,
			Delay:           200 * time.Second,
			FailureAction:   swarm.UpdateFailureActionContinue,
			Monitor:         2 * time.Second,
			MaxFailureRatio: 0.2,
			Order:           swarm.UpdateOrderStopFirst,
		***REMOVED***,
		RollbackConfig: &swarm.UpdateConfig***REMOVED***
			Parallelism:     3,
			Delay:           300 * time.Second,
			FailureAction:   swarm.UpdateFailureActionPause,
			Monitor:         3 * time.Second,
			MaxFailureRatio: 0.3,
			Order:           swarm.UpdateOrderStartFirst,
		***REMOVED***,
	***REMOVED***
***REMOVED***

const defaultSwarmPort = 2477

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

func serviceContainerCount(client client.ServiceAPIClient, id string, count uint64) func(log poll.LogT) poll.Result ***REMOVED***
	return func(log poll.LogT) poll.Result ***REMOVED***
		filter := filters.NewArgs()
		filter.Add("service", id)
		tasks, err := client.TaskList(context.Background(), types.TaskListOptions***REMOVED***
			Filters: filter,
		***REMOVED***)
		switch ***REMOVED***
		case err != nil:
			return poll.Error(err)
		case len(tasks) == int(count):
			return poll.Success()
		default:
			return poll.Continue("task count at %d waiting for %d", len(tasks), count)
		***REMOVED***
	***REMOVED***
***REMOVED***
