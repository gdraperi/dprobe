package environment

import (
	"regexp"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

type testingT interface ***REMOVED***
	require.TestingT
	logT
	Fatalf(string, ...interface***REMOVED******REMOVED***)
***REMOVED***

type logT interface ***REMOVED***
	Logf(string, ...interface***REMOVED******REMOVED***)
***REMOVED***

// Clean the environment, preserving protected objects (images, containers, ...)
// and removing everything else. It's meant to run after any tests so that they don't
// depend on each others.
func (e *Execution) Clean(t testingT) ***REMOVED***
	client := e.APIClient()

	platform := e.OSType
	if (platform != "windows") || (platform == "windows" && e.DaemonInfo.Isolation == "hyperv") ***REMOVED***
		unpauseAllContainers(t, client)
	***REMOVED***
	deleteAllContainers(t, client, e.protectedElements.containers)
	deleteAllImages(t, client, e.protectedElements.images)
	deleteAllVolumes(t, client, e.protectedElements.volumes)
	deleteAllNetworks(t, client, platform, e.protectedElements.networks)
	if platform == "linux" ***REMOVED***
		deleteAllPlugins(t, client, e.protectedElements.plugins)
	***REMOVED***
***REMOVED***

func unpauseAllContainers(t assert.TestingT, client client.ContainerAPIClient) ***REMOVED***
	ctx := context.Background()
	containers := getPausedContainers(ctx, t, client)
	if len(containers) > 0 ***REMOVED***
		for _, container := range containers ***REMOVED***
			err := client.ContainerUnpause(ctx, container.ID)
			assert.NoError(t, err, "failed to unpause container %s", container.ID)
		***REMOVED***
	***REMOVED***
***REMOVED***

func getPausedContainers(ctx context.Context, t assert.TestingT, client client.ContainerAPIClient) []types.Container ***REMOVED***
	filter := filters.NewArgs()
	filter.Add("status", "paused")
	containers, err := client.ContainerList(ctx, types.ContainerListOptions***REMOVED***
		Filters: filter,
		Quiet:   true,
		All:     true,
	***REMOVED***)
	assert.NoError(t, err, "failed to list containers")
	return containers
***REMOVED***

var alreadyExists = regexp.MustCompile(`Error response from daemon: removal of container (\w+) is already in progress`)

func deleteAllContainers(t assert.TestingT, apiclient client.ContainerAPIClient, protectedContainers map[string]struct***REMOVED******REMOVED***) ***REMOVED***
	ctx := context.Background()
	containers := getAllContainers(ctx, t, apiclient)
	if len(containers) == 0 ***REMOVED***
		return
	***REMOVED***

	for _, container := range containers ***REMOVED***
		if _, ok := protectedContainers[container.ID]; ok ***REMOVED***
			continue
		***REMOVED***
		err := apiclient.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions***REMOVED***
			Force:         true,
			RemoveVolumes: true,
		***REMOVED***)
		if err == nil || client.IsErrNotFound(err) || alreadyExists.MatchString(err.Error()) || isErrNotFoundSwarmClassic(err) ***REMOVED***
			continue
		***REMOVED***
		assert.NoError(t, err, "failed to remove %s", container.ID)
	***REMOVED***
***REMOVED***

func getAllContainers(ctx context.Context, t assert.TestingT, client client.ContainerAPIClient) []types.Container ***REMOVED***
	containers, err := client.ContainerList(ctx, types.ContainerListOptions***REMOVED***
		Quiet: true,
		All:   true,
	***REMOVED***)
	assert.NoError(t, err, "failed to list containers")
	return containers
***REMOVED***

func deleteAllImages(t testingT, apiclient client.ImageAPIClient, protectedImages map[string]struct***REMOVED******REMOVED***) ***REMOVED***
	images, err := apiclient.ImageList(context.Background(), types.ImageListOptions***REMOVED******REMOVED***)
	assert.NoError(t, err, "failed to list images")

	ctx := context.Background()
	for _, image := range images ***REMOVED***
		tags := tagsFromImageSummary(image)
		if len(tags) == 0 ***REMOVED***
			t.Logf("Removing image %s", image.ID)
			removeImage(ctx, t, apiclient, image.ID)
			continue
		***REMOVED***
		for _, tag := range tags ***REMOVED***
			if _, ok := protectedImages[tag]; !ok ***REMOVED***
				t.Logf("Removing image %s", tag)
				removeImage(ctx, t, apiclient, tag)
				continue
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func removeImage(ctx context.Context, t assert.TestingT, apiclient client.ImageAPIClient, ref string) ***REMOVED***
	_, err := apiclient.ImageRemove(ctx, ref, types.ImageRemoveOptions***REMOVED***
		Force: true,
	***REMOVED***)
	if client.IsErrNotFound(err) ***REMOVED***
		return
	***REMOVED***
	assert.NoError(t, err, "failed to remove image %s", ref)
***REMOVED***

func deleteAllVolumes(t assert.TestingT, c client.VolumeAPIClient, protectedVolumes map[string]struct***REMOVED******REMOVED***) ***REMOVED***
	volumes, err := c.VolumeList(context.Background(), filters.Args***REMOVED******REMOVED***)
	assert.NoError(t, err, "failed to list volumes")

	for _, v := range volumes.Volumes ***REMOVED***
		if _, ok := protectedVolumes[v.Name]; ok ***REMOVED***
			continue
		***REMOVED***
		err := c.VolumeRemove(context.Background(), v.Name, true)
		// Docker EE may list volumes that no longer exist.
		if isErrNotFoundSwarmClassic(err) ***REMOVED***
			continue
		***REMOVED***
		assert.NoError(t, err, "failed to remove volume %s", v.Name)
	***REMOVED***
***REMOVED***

func deleteAllNetworks(t assert.TestingT, c client.NetworkAPIClient, daemonPlatform string, protectedNetworks map[string]struct***REMOVED******REMOVED***) ***REMOVED***
	networks, err := c.NetworkList(context.Background(), types.NetworkListOptions***REMOVED******REMOVED***)
	assert.NoError(t, err, "failed to list networks")

	for _, n := range networks ***REMOVED***
		if n.Name == "bridge" || n.Name == "none" || n.Name == "host" ***REMOVED***
			continue
		***REMOVED***
		if _, ok := protectedNetworks[n.ID]; ok ***REMOVED***
			continue
		***REMOVED***
		if daemonPlatform == "windows" && strings.ToLower(n.Name) == "nat" ***REMOVED***
			// nat is a pre-defined network on Windows and cannot be removed
			continue
		***REMOVED***
		err := c.NetworkRemove(context.Background(), n.ID)
		assert.NoError(t, err, "failed to remove network %s", n.ID)
	***REMOVED***
***REMOVED***

func deleteAllPlugins(t assert.TestingT, c client.PluginAPIClient, protectedPlugins map[string]struct***REMOVED******REMOVED***) ***REMOVED***
	plugins, err := c.PluginList(context.Background(), filters.Args***REMOVED******REMOVED***)
	// Docker EE does not allow cluster-wide plugin management.
	if client.IsErrNotImplemented(err) ***REMOVED***
		return
	***REMOVED***
	assert.NoError(t, err, "failed to list plugins")

	for _, p := range plugins ***REMOVED***
		if _, ok := protectedPlugins[p.Name]; ok ***REMOVED***
			continue
		***REMOVED***
		err := c.PluginRemove(context.Background(), p.Name, types.PluginRemoveOptions***REMOVED***Force: true***REMOVED***)
		assert.NoError(t, err, "failed to remove plugin %s", p.ID)
	***REMOVED***
***REMOVED***

// Swarm classic aggregates node errors and returns a 500 so we need to check
// the error string instead of just IsErrNotFound().
func isErrNotFoundSwarmClassic(err error) bool ***REMOVED***
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "no such")
***REMOVED***
