package system

import (
	"context"
	"testing"

	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/integration/util/request"
	"github.com/stretchr/testify/require"
)

func TestEvents(t *testing.T) ***REMOVED***
	defer setupTest(t)()
	ctx := context.Background()
	client := request.NewAPIClient(t)

	container, err := client.ContainerCreate(ctx,
		&container.Config***REMOVED***
			Image:      "busybox",
			Tty:        true,
			WorkingDir: "/root",
			Cmd:        strslice.StrSlice([]string***REMOVED***"top"***REMOVED***),
		***REMOVED***,
		&container.HostConfig***REMOVED******REMOVED***,
		&network.NetworkingConfig***REMOVED******REMOVED***,
		"foo",
	)
	require.NoError(t, err)
	err = client.ContainerStart(ctx, container.ID, types.ContainerStartOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	id, err := client.ContainerExecCreate(ctx, container.ID,
		types.ExecConfig***REMOVED***
			Cmd: strslice.StrSlice([]string***REMOVED***"echo", "hello"***REMOVED***),
		***REMOVED***,
	)
	require.NoError(t, err)

	filters := filters.NewArgs(
		filters.Arg("container", container.ID),
		filters.Arg("event", "exec_die"),
	)
	msg, errors := client.Events(ctx, types.EventsOptions***REMOVED***
		Filters: filters,
	***REMOVED***)

	err = client.ContainerExecStart(ctx, id.ID,
		types.ExecStartCheck***REMOVED***
			Detach: true,
			Tty:    false,
		***REMOVED***,
	)
	require.NoError(t, err)

	select ***REMOVED***
	case m := <-msg:
		require.Equal(t, m.Type, "container")
		require.Equal(t, m.Actor.ID, container.ID)
		require.Equal(t, m.Action, "exec_die")
		require.Equal(t, m.Actor.Attributes["execID"], id.ID)
		require.Equal(t, m.Actor.Attributes["exitCode"], "0")
	case err = <-errors:
		t.Fatal(err)
	case <-time.After(time.Second * 3):
		t.Fatal("timeout hit")
	***REMOVED***

***REMOVED***
