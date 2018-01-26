package image

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/integration/util/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommitInheritsEnv(t *testing.T) ***REMOVED***
	defer setupTest(t)()
	client := request.NewAPIClient(t)
	ctx := context.Background()

	createResp1, err := client.ContainerCreate(ctx, &container.Config***REMOVED***Image: "busybox"***REMOVED***, nil, nil, "")
	require.NoError(t, err)

	commitResp1, err := client.ContainerCommit(ctx, createResp1.ID, types.ContainerCommitOptions***REMOVED***
		Changes:   []string***REMOVED***"ENV PATH=/bin"***REMOVED***,
		Reference: "test-commit-image",
	***REMOVED***)
	require.NoError(t, err)

	image1, _, err := client.ImageInspectWithRaw(ctx, commitResp1.ID)
	require.NoError(t, err)

	expectedEnv1 := []string***REMOVED***"PATH=/bin"***REMOVED***
	assert.Equal(t, expectedEnv1, image1.Config.Env)

	createResp2, err := client.ContainerCreate(ctx, &container.Config***REMOVED***Image: image1.ID***REMOVED***, nil, nil, "")
	require.NoError(t, err)

	commitResp2, err := client.ContainerCommit(ctx, createResp2.ID, types.ContainerCommitOptions***REMOVED***
		Changes:   []string***REMOVED***"ENV PATH=/usr/bin:$PATH"***REMOVED***,
		Reference: "test-commit-image",
	***REMOVED***)
	require.NoError(t, err)

	image2, _, err := client.ImageInspectWithRaw(ctx, commitResp2.ID)
	require.NoError(t, err)
	expectedEnv2 := []string***REMOVED***"PATH=/usr/bin:/bin"***REMOVED***
	assert.Equal(t, expectedEnv2, image2.Config.Env)
***REMOVED***
