package build

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/integration-cli/cli/build/fakecontext"
	"github.com/docker/docker/integration/util/request"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildWithRemoveAndForceRemove(t *testing.T) ***REMOVED***
	defer setupTest(t)()
	t.Parallel()
	cases := []struct ***REMOVED***
		name                           string
		dockerfile                     string
		numberOfIntermediateContainers int
		rm                             bool
		forceRm                        bool
	***REMOVED******REMOVED***
		***REMOVED***
			name: "successful build with no removal",
			dockerfile: `FROM busybox
			RUN exit 0
			RUN exit 0`,
			numberOfIntermediateContainers: 2,
			rm:      false,
			forceRm: false,
		***REMOVED***,
		***REMOVED***
			name: "successful build with remove",
			dockerfile: `FROM busybox
			RUN exit 0
			RUN exit 0`,
			numberOfIntermediateContainers: 0,
			rm:      true,
			forceRm: false,
		***REMOVED***,
		***REMOVED***
			name: "successful build with remove and force remove",
			dockerfile: `FROM busybox
			RUN exit 0
			RUN exit 0`,
			numberOfIntermediateContainers: 0,
			rm:      true,
			forceRm: true,
		***REMOVED***,
		***REMOVED***
			name: "failed build with no removal",
			dockerfile: `FROM busybox
			RUN exit 0
			RUN exit 1`,
			numberOfIntermediateContainers: 2,
			rm:      false,
			forceRm: false,
		***REMOVED***,
		***REMOVED***
			name: "failed build with remove",
			dockerfile: `FROM busybox
			RUN exit 0
			RUN exit 1`,
			numberOfIntermediateContainers: 1,
			rm:      true,
			forceRm: false,
		***REMOVED***,
		***REMOVED***
			name: "failed build with remove and force remove",
			dockerfile: `FROM busybox
			RUN exit 0
			RUN exit 1`,
			numberOfIntermediateContainers: 0,
			rm:      true,
			forceRm: true,
		***REMOVED***,
	***REMOVED***

	client := request.NewAPIClient(t)
	ctx := context.Background()
	for _, c := range cases ***REMOVED***
		t.Run(c.name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			dockerfile := []byte(c.dockerfile)

			buff := bytes.NewBuffer(nil)
			tw := tar.NewWriter(buff)
			require.NoError(t, tw.WriteHeader(&tar.Header***REMOVED***
				Name: "Dockerfile",
				Size: int64(len(dockerfile)),
			***REMOVED***))
			_, err := tw.Write(dockerfile)
			require.NoError(t, err)
			require.NoError(t, tw.Close())
			resp, err := client.ImageBuild(ctx, buff, types.ImageBuildOptions***REMOVED***Remove: c.rm, ForceRemove: c.forceRm, NoCache: true***REMOVED***)
			require.NoError(t, err)
			defer resp.Body.Close()
			filter, err := buildContainerIdsFilter(resp.Body)
			require.NoError(t, err)
			remainingContainers, err := client.ContainerList(ctx, types.ContainerListOptions***REMOVED***Filters: filter, All: true***REMOVED***)
			require.NoError(t, err)
			require.Equal(t, c.numberOfIntermediateContainers, len(remainingContainers), "Expected %v remaining intermediate containers, got %v", c.numberOfIntermediateContainers, len(remainingContainers))
		***REMOVED***)
	***REMOVED***
***REMOVED***

func buildContainerIdsFilter(buildOutput io.Reader) (filters.Args, error) ***REMOVED***
	const intermediateContainerPrefix = " ---> Running in "
	filter := filters.NewArgs()

	dec := json.NewDecoder(buildOutput)
	for ***REMOVED***
		m := jsonmessage.JSONMessage***REMOVED******REMOVED***
		err := dec.Decode(&m)
		if err == io.EOF ***REMOVED***
			return filter, nil
		***REMOVED***
		if err != nil ***REMOVED***
			return filter, err
		***REMOVED***
		if ix := strings.Index(m.Stream, intermediateContainerPrefix); ix != -1 ***REMOVED***
			filter.Add("id", strings.TrimSpace(m.Stream[ix+len(intermediateContainerPrefix):]))
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestBuildMultiStageParentConfig(t *testing.T) ***REMOVED***
	dockerfile := `
		FROM busybox AS stage0
		ENV WHO=parent
		WORKDIR /foo

		FROM stage0
		ENV WHO=sibling1
		WORKDIR sub1

		FROM stage0
		WORKDIR sub2
	`
	ctx := context.Background()
	source := fakecontext.New(t, "", fakecontext.WithDockerfile(dockerfile))
	defer source.Close()

	apiclient := testEnv.APIClient()
	resp, err := apiclient.ImageBuild(ctx,
		source.AsTarReader(t),
		types.ImageBuildOptions***REMOVED***
			Remove:      true,
			ForceRemove: true,
			Tags:        []string***REMOVED***"build1"***REMOVED***,
		***REMOVED***)
	require.NoError(t, err)
	_, err = io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	require.NoError(t, err)

	image, _, err := apiclient.ImageInspectWithRaw(ctx, "build1")
	require.NoError(t, err)

	assert.Equal(t, "/foo/sub2", image.Config.WorkingDir)
	assert.Contains(t, image.Config.Env, "WHO=parent")
***REMOVED***

func TestBuildWithEmptyLayers(t *testing.T) ***REMOVED***
	dockerfile := `
		FROM    busybox
		COPY    1/ /target/
		COPY    2/ /target/
		COPY    3/ /target/
	`
	ctx := context.Background()
	source := fakecontext.New(t, "",
		fakecontext.WithDockerfile(dockerfile),
		fakecontext.WithFile("1/a", "asdf"),
		fakecontext.WithFile("2/a", "asdf"),
		fakecontext.WithFile("3/a", "asdf"))
	defer source.Close()

	apiclient := testEnv.APIClient()
	resp, err := apiclient.ImageBuild(ctx,
		source.AsTarReader(t),
		types.ImageBuildOptions***REMOVED***
			Remove:      true,
			ForceRemove: true,
		***REMOVED***)
	require.NoError(t, err)
	_, err = io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	require.NoError(t, err)
***REMOVED***

// TestBuildMultiStageOnBuild checks that ONBUILD commands are applied to
// multiple subsequent stages
// #35652
func TestBuildMultiStageOnBuild(t *testing.T) ***REMOVED***
	defer setupTest(t)()
	// test both metadata and layer based commands as they may be implemented differently
	dockerfile := `FROM busybox AS stage1
ONBUILD RUN echo 'foo' >somefile
ONBUILD ENV bar=baz

FROM stage1
RUN cat somefile # fails if ONBUILD RUN fails

FROM stage1
RUN cat somefile`

	ctx := context.Background()
	source := fakecontext.New(t, "",
		fakecontext.WithDockerfile(dockerfile))
	defer source.Close()

	apiclient := testEnv.APIClient()
	resp, err := apiclient.ImageBuild(ctx,
		source.AsTarReader(t),
		types.ImageBuildOptions***REMOVED***
			Remove:      true,
			ForceRemove: true,
		***REMOVED***)

	out := bytes.NewBuffer(nil)
	require.NoError(t, err)
	_, err = io.Copy(out, resp.Body)
	resp.Body.Close()
	require.NoError(t, err)

	assert.Contains(t, out.String(), "Successfully built")

	imageIDs, err := getImageIDsFromBuild(out.Bytes())
	require.NoError(t, err)
	assert.Equal(t, 3, len(imageIDs))

	image, _, err := apiclient.ImageInspectWithRaw(context.Background(), imageIDs[2])
	require.NoError(t, err)
	assert.Contains(t, image.Config.Env, "bar=baz")
***REMOVED***

type buildLine struct ***REMOVED***
	Stream string
	Aux    struct ***REMOVED***
		ID string
	***REMOVED***
***REMOVED***

func getImageIDsFromBuild(output []byte) ([]string, error) ***REMOVED***
	ids := []string***REMOVED******REMOVED***
	for _, line := range bytes.Split(output, []byte("\n")) ***REMOVED***
		if len(line) == 0 ***REMOVED***
			continue
		***REMOVED***
		entry := buildLine***REMOVED******REMOVED***
		if err := json.Unmarshal(line, &entry); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if entry.Aux.ID != "" ***REMOVED***
			ids = append(ids, entry.Aux.ID)
		***REMOVED***
	***REMOVED***
	return ids, nil
***REMOVED***
