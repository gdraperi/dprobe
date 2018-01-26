package environment

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	dclient "github.com/docker/docker/client"
	"github.com/stretchr/testify/require"
)

var frozenImages = []string***REMOVED***"busybox:latest", "hello-world:frozen", "debian:jessie"***REMOVED***

type protectedElements struct ***REMOVED***
	containers map[string]struct***REMOVED******REMOVED***
	images     map[string]struct***REMOVED******REMOVED***
	networks   map[string]struct***REMOVED******REMOVED***
	plugins    map[string]struct***REMOVED******REMOVED***
	volumes    map[string]struct***REMOVED******REMOVED***
***REMOVED***

func newProtectedElements() protectedElements ***REMOVED***
	return protectedElements***REMOVED***
		containers: map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***,
		images:     map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***,
		networks:   map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***,
		plugins:    map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***,
		volumes:    map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

// ProtectAll protects the existing environment (containers, images, networks,
// volumes, and, on Linux, plugins) from being cleaned up at the end of test
// runs
func ProtectAll(t testingT, testEnv *Execution) ***REMOVED***
	ProtectContainers(t, testEnv)
	ProtectImages(t, testEnv)
	ProtectNetworks(t, testEnv)
	ProtectVolumes(t, testEnv)
	if testEnv.OSType == "linux" ***REMOVED***
		ProtectPlugins(t, testEnv)
	***REMOVED***
***REMOVED***

// ProtectContainer adds the specified container(s) to be protected in case of
// clean
func (e *Execution) ProtectContainer(t testingT, containers ...string) ***REMOVED***
	for _, container := range containers ***REMOVED***
		e.protectedElements.containers[container] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
***REMOVED***

// ProtectContainers protects existing containers from being cleaned up at the
// end of test runs
func ProtectContainers(t testingT, testEnv *Execution) ***REMOVED***
	containers := getExistingContainers(t, testEnv)
	testEnv.ProtectContainer(t, containers...)
***REMOVED***

func getExistingContainers(t require.TestingT, testEnv *Execution) []string ***REMOVED***
	client := testEnv.APIClient()
	containerList, err := client.ContainerList(context.Background(), types.ContainerListOptions***REMOVED***
		All: true,
	***REMOVED***)
	require.NoError(t, err, "failed to list containers")

	containers := []string***REMOVED******REMOVED***
	for _, container := range containerList ***REMOVED***
		containers = append(containers, container.ID)
	***REMOVED***
	return containers
***REMOVED***

// ProtectImage adds the specified image(s) to be protected in case of clean
func (e *Execution) ProtectImage(t testingT, images ...string) ***REMOVED***
	for _, image := range images ***REMOVED***
		e.protectedElements.images[image] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
***REMOVED***

// ProtectImages protects existing images and on linux frozen images from being
// cleaned up at the end of test runs
func ProtectImages(t testingT, testEnv *Execution) ***REMOVED***
	images := getExistingImages(t, testEnv)

	if testEnv.OSType == "linux" ***REMOVED***
		images = append(images, frozenImages...)
	***REMOVED***
	testEnv.ProtectImage(t, images...)
***REMOVED***

func getExistingImages(t require.TestingT, testEnv *Execution) []string ***REMOVED***
	client := testEnv.APIClient()
	filter := filters.NewArgs()
	filter.Add("dangling", "false")
	imageList, err := client.ImageList(context.Background(), types.ImageListOptions***REMOVED***
		All:     true,
		Filters: filter,
	***REMOVED***)
	require.NoError(t, err, "failed to list images")

	images := []string***REMOVED******REMOVED***
	for _, image := range imageList ***REMOVED***
		images = append(images, tagsFromImageSummary(image)...)
	***REMOVED***
	return images
***REMOVED***

func tagsFromImageSummary(image types.ImageSummary) []string ***REMOVED***
	result := []string***REMOVED******REMOVED***
	for _, tag := range image.RepoTags ***REMOVED***
		if tag != "<none>:<none>" ***REMOVED***
			result = append(result, tag)
		***REMOVED***
	***REMOVED***
	for _, digest := range image.RepoDigests ***REMOVED***
		if digest != "<none>@<none>" ***REMOVED***
			result = append(result, digest)
		***REMOVED***
	***REMOVED***
	return result
***REMOVED***

// ProtectNetwork adds the specified network(s) to be protected in case of
// clean
func (e *Execution) ProtectNetwork(t testingT, networks ...string) ***REMOVED***
	for _, network := range networks ***REMOVED***
		e.protectedElements.networks[network] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
***REMOVED***

// ProtectNetworks protects existing networks from being cleaned up at the end
// of test runs
func ProtectNetworks(t testingT, testEnv *Execution) ***REMOVED***
	networks := getExistingNetworks(t, testEnv)
	testEnv.ProtectNetwork(t, networks...)
***REMOVED***

func getExistingNetworks(t require.TestingT, testEnv *Execution) []string ***REMOVED***
	client := testEnv.APIClient()
	networkList, err := client.NetworkList(context.Background(), types.NetworkListOptions***REMOVED******REMOVED***)
	require.NoError(t, err, "failed to list networks")

	networks := []string***REMOVED******REMOVED***
	for _, network := range networkList ***REMOVED***
		networks = append(networks, network.ID)
	***REMOVED***
	return networks
***REMOVED***

// ProtectPlugin adds the specified plugin(s) to be protected in case of clean
func (e *Execution) ProtectPlugin(t testingT, plugins ...string) ***REMOVED***
	for _, plugin := range plugins ***REMOVED***
		e.protectedElements.plugins[plugin] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
***REMOVED***

// ProtectPlugins protects existing plugins from being cleaned up at the end of
// test runs
func ProtectPlugins(t testingT, testEnv *Execution) ***REMOVED***
	plugins := getExistingPlugins(t, testEnv)
	testEnv.ProtectPlugin(t, plugins...)
***REMOVED***

func getExistingPlugins(t require.TestingT, testEnv *Execution) []string ***REMOVED***
	client := testEnv.APIClient()
	pluginList, err := client.PluginList(context.Background(), filters.Args***REMOVED******REMOVED***)
	// Docker EE does not allow cluster-wide plugin management.
	if dclient.IsErrNotImplemented(err) ***REMOVED***
		return []string***REMOVED******REMOVED***
	***REMOVED***
	require.NoError(t, err, "failed to list plugins")

	plugins := []string***REMOVED******REMOVED***
	for _, plugin := range pluginList ***REMOVED***
		plugins = append(plugins, plugin.Name)
	***REMOVED***
	return plugins
***REMOVED***

// ProtectVolume adds the specified volume(s) to be protected in case of clean
func (e *Execution) ProtectVolume(t testingT, volumes ...string) ***REMOVED***
	for _, volume := range volumes ***REMOVED***
		e.protectedElements.volumes[volume] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
***REMOVED***

// ProtectVolumes protects existing volumes from being cleaned up at the end of
// test runs
func ProtectVolumes(t testingT, testEnv *Execution) ***REMOVED***
	volumes := getExistingVolumes(t, testEnv)
	testEnv.ProtectVolume(t, volumes...)
***REMOVED***

func getExistingVolumes(t require.TestingT, testEnv *Execution) []string ***REMOVED***
	client := testEnv.APIClient()
	volumeList, err := client.VolumeList(context.Background(), filters.Args***REMOVED******REMOVED***)
	require.NoError(t, err, "failed to list volumes")

	volumes := []string***REMOVED******REMOVED***
	for _, volume := range volumeList.Volumes ***REMOVED***
		volumes = append(volumes, volume.Name)
	***REMOVED***
	return volumes
***REMOVED***
