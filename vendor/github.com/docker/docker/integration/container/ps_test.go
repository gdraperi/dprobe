package container

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/integration/util/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPsFilter(t *testing.T) ***REMOVED***
	defer setupTest(t)()
	client := request.NewAPIClient(t)
	ctx := context.Background()

	createContainerForFilter := func(ctx context.Context, name string) string ***REMOVED***
		body, err := client.ContainerCreate(ctx,
			&container.Config***REMOVED***
				Cmd:   []string***REMOVED***"top"***REMOVED***,
				Image: "busybox",
			***REMOVED***,
			&container.HostConfig***REMOVED******REMOVED***,
			&network.NetworkingConfig***REMOVED******REMOVED***,
			name,
		)
		require.NoError(t, err)
		return body.ID
	***REMOVED***

	prev := createContainerForFilter(ctx, "prev")
	createContainerForFilter(ctx, "top")
	next := createContainerForFilter(ctx, "next")

	containerIDs := func(containers []types.Container) []string ***REMOVED***
		entries := []string***REMOVED******REMOVED***
		for _, container := range containers ***REMOVED***
			entries = append(entries, container.ID)
		***REMOVED***
		return entries
	***REMOVED***

	f1 := filters.NewArgs()
	f1.Add("since", "top")
	q1, err := client.ContainerList(ctx, types.ContainerListOptions***REMOVED***
		All:     true,
		Filters: f1,
	***REMOVED***)
	require.NoError(t, err)
	assert.Contains(t, containerIDs(q1), next)

	f2 := filters.NewArgs()
	f2.Add("before", "top")
	q2, err := client.ContainerList(ctx, types.ContainerListOptions***REMOVED***
		All:     true,
		Filters: f2,
	***REMOVED***)
	require.NoError(t, err)
	assert.Contains(t, containerIDs(q2), prev)
***REMOVED***
