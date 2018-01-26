package container

import (
	"context"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/docker/integration/util/request"
	"github.com/gotestyourself/gotestyourself/poll"
	"github.com/stretchr/testify/require"
)

// TestHealthCheckWorkdir verifies that health-checks inherit the containers'
// working-dir.
func TestHealthCheckWorkdir(t *testing.T) ***REMOVED***
	defer setupTest(t)()
	ctx := context.Background()
	client := request.NewAPIClient(t)

	c, err := client.ContainerCreate(ctx,
		&container.Config***REMOVED***
			Image:      "busybox",
			Tty:        true,
			WorkingDir: "/foo",
			Cmd:        strslice.StrSlice([]string***REMOVED***"top"***REMOVED***),
			Healthcheck: &container.HealthConfig***REMOVED***
				Test:     []string***REMOVED***"CMD-SHELL", "if [ \"$PWD\" = \"/foo\" ]; then exit 0; else exit 1; fi;"***REMOVED***,
				Interval: 50 * time.Millisecond,
				Retries:  3,
			***REMOVED***,
		***REMOVED***,
		&container.HostConfig***REMOVED******REMOVED***,
		&network.NetworkingConfig***REMOVED******REMOVED***,
		"healthtest",
	)
	require.NoError(t, err)
	err = client.ContainerStart(ctx, c.ID, types.ContainerStartOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	poll.WaitOn(t, pollForHealthStatus(ctx, client, c.ID, types.Healthy), poll.WithDelay(100*time.Millisecond))
***REMOVED***

func pollForHealthStatus(ctx context.Context, client client.APIClient, containerID string, healthStatus string) func(log poll.LogT) poll.Result ***REMOVED***
	return func(log poll.LogT) poll.Result ***REMOVED***
		inspect, err := client.ContainerInspect(ctx, containerID)

		switch ***REMOVED***
		case err != nil:
			return poll.Error(err)
		case inspect.State.Health.Status == healthStatus:
			return poll.Success()
		default:
			return poll.Continue("waiting for container to become %s", healthStatus)
		***REMOVED***
	***REMOVED***
***REMOVED***
