package service

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/integration-cli/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDockerNetworkConnectAlias(t *testing.T) ***REMOVED***
	defer setupTest(t)()
	d := newSwarm(t)
	defer d.Stop(t)
	client, err := request.NewClientForHost(d.Sock())
	require.NoError(t, err)
	ctx := context.Background()

	name := "test-alias"
	_, err = client.NetworkCreate(ctx, name, types.NetworkCreate***REMOVED***
		Driver:     "overlay",
		Attachable: true,
	***REMOVED***)
	require.NoError(t, err)
	_, err = client.ContainerCreate(ctx,
		&container.Config***REMOVED***
			Cmd:   []string***REMOVED***"top"***REMOVED***,
			Image: "busybox",
		***REMOVED***,
		&container.HostConfig***REMOVED******REMOVED***,
		&network.NetworkingConfig***REMOVED***
			EndpointsConfig: map[string]*network.EndpointSettings***REMOVED***
				name: ***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"ng1",
	)
	require.NoError(t, err)
	err = client.NetworkConnect(ctx, name, "ng1", &network.EndpointSettings***REMOVED***
		Aliases: []string***REMOVED***
			"aaa",
		***REMOVED***,
	***REMOVED***)
	require.NoError(t, err)

	err = client.ContainerStart(ctx, "ng1", types.ContainerStartOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	ng1, err := client.ContainerInspect(ctx, "ng1")
	require.NoError(t, err)
	assert.Equal(t, len(ng1.NetworkSettings.Networks[name].Aliases), 2)
	assert.Equal(t, ng1.NetworkSettings.Networks[name].Aliases[0], "aaa")

	_, err = client.ContainerCreate(ctx,
		&container.Config***REMOVED***
			Cmd:   []string***REMOVED***"top"***REMOVED***,
			Image: "busybox",
		***REMOVED***,
		&container.HostConfig***REMOVED******REMOVED***,
		&network.NetworkingConfig***REMOVED***
			EndpointsConfig: map[string]*network.EndpointSettings***REMOVED***
				name: ***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"ng2",
	)
	require.NoError(t, err)
	err = client.NetworkConnect(ctx, name, "ng2", &network.EndpointSettings***REMOVED***
		Aliases: []string***REMOVED***
			"bbb",
		***REMOVED***,
	***REMOVED***)
	require.NoError(t, err)

	err = client.ContainerStart(ctx, "ng2", types.ContainerStartOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	ng2, err := client.ContainerInspect(ctx, "ng2")
	require.NoError(t, err)
	assert.Equal(t, len(ng2.NetworkSettings.Networks[name].Aliases), 2)
	assert.Equal(t, ng2.NetworkSettings.Networks[name].Aliases[0], "bbb")
***REMOVED***
