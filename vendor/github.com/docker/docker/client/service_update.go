package client

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

// ServiceUpdate updates a Service.
func (cli *Client) ServiceUpdate(ctx context.Context, serviceID string, version swarm.Version, service swarm.ServiceSpec, options types.ServiceUpdateOptions) (types.ServiceUpdateResponse, error) ***REMOVED***
	var (
		query   = url.Values***REMOVED******REMOVED***
		distErr error
	)

	headers := map[string][]string***REMOVED***
		"version": ***REMOVED***cli.version***REMOVED***,
	***REMOVED***

	if options.EncodedRegistryAuth != "" ***REMOVED***
		headers["X-Registry-Auth"] = []string***REMOVED***options.EncodedRegistryAuth***REMOVED***
	***REMOVED***

	if options.RegistryAuthFrom != "" ***REMOVED***
		query.Set("registryAuthFrom", options.RegistryAuthFrom)
	***REMOVED***

	if options.Rollback != "" ***REMOVED***
		query.Set("rollback", options.Rollback)
	***REMOVED***

	query.Set("version", strconv.FormatUint(version.Index, 10))

	if err := validateServiceSpec(service); err != nil ***REMOVED***
		return types.ServiceUpdateResponse***REMOVED******REMOVED***, err
	***REMOVED***

	var imgPlatforms []swarm.Platform
	// ensure that the image is tagged
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

	var response types.ServiceUpdateResponse
	resp, err := cli.post(ctx, "/services/"+serviceID+"/update", query, service, headers)
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
