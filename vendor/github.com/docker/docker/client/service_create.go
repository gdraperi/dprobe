package client

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	digest "github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// ServiceCreate creates a new Service.
func (cli *Client) ServiceCreate(ctx context.Context, service swarm.ServiceSpec, options types.ServiceCreateOptions) (types.ServiceCreateResponse, error) ***REMOVED***
	var distErr error

	headers := map[string][]string***REMOVED***
		"version": ***REMOVED***cli.version***REMOVED***,
	***REMOVED***

	if options.EncodedRegistryAuth != "" ***REMOVED***
		headers["X-Registry-Auth"] = []string***REMOVED***options.EncodedRegistryAuth***REMOVED***
	***REMOVED***

	// Make sure containerSpec is not nil when no runtime is set or the runtime is set to container
	if service.TaskTemplate.ContainerSpec == nil && (service.TaskTemplate.Runtime == "" || service.TaskTemplate.Runtime == swarm.RuntimeContainer) ***REMOVED***
		service.TaskTemplate.ContainerSpec = &swarm.ContainerSpec***REMOVED******REMOVED***
	***REMOVED***

	if err := validateServiceSpec(service); err != nil ***REMOVED***
		return types.ServiceCreateResponse***REMOVED******REMOVED***, err
	***REMOVED***

	// ensure that the image is tagged
	var imgPlatforms []swarm.Platform
	if service.TaskTemplate.ContainerSpec != nil ***REMOVED***
		if taggedImg := imageWithTagString(service.TaskTemplate.ContainerSpec.Image); taggedImg != "" ***REMOVED***
			service.TaskTemplate.ContainerSpec.Image = taggedImg
		***REMOVED***
		if options.QueryRegistry ***REMOVED***
			var img string
			img, imgPlatforms, distErr = imageDigestAndPlatforms(ctx, cli, service.TaskTemplate.ContainerSpec.Image, options.EncodedRegistryAuth)
			if img != "" ***REMOVED***
				service.TaskTemplate.ContainerSpec.Image = img
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// ensure that the image is tagged
	if service.TaskTemplate.PluginSpec != nil ***REMOVED***
		if taggedImg := imageWithTagString(service.TaskTemplate.PluginSpec.Remote); taggedImg != "" ***REMOVED***
			service.TaskTemplate.PluginSpec.Remote = taggedImg
		***REMOVED***
		if options.QueryRegistry ***REMOVED***
			var img string
			img, imgPlatforms, distErr = imageDigestAndPlatforms(ctx, cli, service.TaskTemplate.PluginSpec.Remote, options.EncodedRegistryAuth)
			if img != "" ***REMOVED***
				service.TaskTemplate.PluginSpec.Remote = img
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if service.TaskTemplate.Placement == nil && len(imgPlatforms) > 0 ***REMOVED***
		service.TaskTemplate.Placement = &swarm.Placement***REMOVED******REMOVED***
	***REMOVED***
	if len(imgPlatforms) > 0 ***REMOVED***
		service.TaskTemplate.Placement.Platforms = imgPlatforms
	***REMOVED***

	var response types.ServiceCreateResponse
	resp, err := cli.post(ctx, "/services/create", nil, service, headers)
	if err != nil ***REMOVED***
		return response, err
	***REMOVED***

	err = json.NewDecoder(resp.body).Decode(&response)

	if distErr != nil ***REMOVED***
		response.Warnings = append(response.Warnings, digestWarning(service.TaskTemplate.ContainerSpec.Image))
	***REMOVED***

	ensureReaderClosed(resp)
	return response, err
***REMOVED***

func imageDigestAndPlatforms(ctx context.Context, cli DistributionAPIClient, image, encodedAuth string) (string, []swarm.Platform, error) ***REMOVED***
	distributionInspect, err := cli.DistributionInspect(ctx, image, encodedAuth)
	var platforms []swarm.Platform
	if err != nil ***REMOVED***
		return "", nil, err
	***REMOVED***

	imageWithDigest := imageWithDigestString(image, distributionInspect.Descriptor.Digest)

	if len(distributionInspect.Platforms) > 0 ***REMOVED***
		platforms = make([]swarm.Platform, 0, len(distributionInspect.Platforms))
		for _, p := range distributionInspect.Platforms ***REMOVED***
			// clear architecture field for arm. This is a temporary patch to address
			// https://github.com/docker/swarmkit/issues/2294. The issue is that while
			// image manifests report "arm" as the architecture, the node reports
			// something like "armv7l" (includes the variant), which causes arm images
			// to stop working with swarm mode. This patch removes the architecture
			// constraint for arm images to ensure tasks get scheduled.
			arch := p.Architecture
			if strings.ToLower(arch) == "arm" ***REMOVED***
				arch = ""
			***REMOVED***
			platforms = append(platforms, swarm.Platform***REMOVED***
				Architecture: arch,
				OS:           p.OS,
			***REMOVED***)
		***REMOVED***
	***REMOVED***
	return imageWithDigest, platforms, err
***REMOVED***

// imageWithDigestString takes an image string and a digest, and updates
// the image string if it didn't originally contain a digest. It returns
// an empty string if there are no updates.
func imageWithDigestString(image string, dgst digest.Digest) string ***REMOVED***
	namedRef, err := reference.ParseNormalizedNamed(image)
	if err == nil ***REMOVED***
		if _, isCanonical := namedRef.(reference.Canonical); !isCanonical ***REMOVED***
			// ensure that image gets a default tag if none is provided
			img, err := reference.WithDigest(namedRef, dgst)
			if err == nil ***REMOVED***
				return reference.FamiliarString(img)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

// imageWithTagString takes an image string, and returns a tagged image
// string, adding a 'latest' tag if one was not provided. It returns an
// emptry string if a canonical reference was provided
func imageWithTagString(image string) string ***REMOVED***
	namedRef, err := reference.ParseNormalizedNamed(image)
	if err == nil ***REMOVED***
		return reference.FamiliarString(reference.TagNameOnly(namedRef))
	***REMOVED***
	return ""
***REMOVED***

// digestWarning constructs a formatted warning string using the
// image name that could not be pinned by digest. The formatting
// is hardcoded, but could me made smarter in the future
func digestWarning(image string) string ***REMOVED***
	return fmt.Sprintf("image %s could not be accessed on a registry to record\nits digest. Each node will access %s independently,\npossibly leading to different nodes running different\nversions of the image.\n", image, image)
***REMOVED***

func validateServiceSpec(s swarm.ServiceSpec) error ***REMOVED***
	if s.TaskTemplate.ContainerSpec != nil && s.TaskTemplate.PluginSpec != nil ***REMOVED***
		return errors.New("must not specify both a container spec and a plugin spec in the task template")
	***REMOVED***
	if s.TaskTemplate.PluginSpec != nil && s.TaskTemplate.Runtime != swarm.RuntimePlugin ***REMOVED***
		return errors.New("mismatched runtime with plugin spec")
	***REMOVED***
	if s.TaskTemplate.ContainerSpec != nil && (s.TaskTemplate.Runtime != "" && s.TaskTemplate.Runtime != swarm.RuntimeContainer) ***REMOVED***
		return errors.New("mismatched runtime with container spec")
	***REMOVED***
	return nil
***REMOVED***
