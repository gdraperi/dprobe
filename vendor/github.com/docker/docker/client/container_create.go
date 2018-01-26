package client

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/versions"
	"golang.org/x/net/context"
)

type configWrapper struct ***REMOVED***
	*container.Config
	HostConfig       *container.HostConfig
	NetworkingConfig *network.NetworkingConfig
***REMOVED***

// ContainerCreate creates a new container based in the given configuration.
// It can be associated with a name, but it's not mandatory.
func (cli *Client) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (container.ContainerCreateCreatedBody, error) ***REMOVED***
	var response container.ContainerCreateCreatedBody

	if err := cli.NewVersionError("1.25", "stop timeout"); config != nil && config.StopTimeout != nil && err != nil ***REMOVED***
		return response, err
	***REMOVED***

	// When using API 1.24 and under, the client is responsible for removing the container
	if hostConfig != nil && versions.LessThan(cli.ClientVersion(), "1.25") ***REMOVED***
		hostConfig.AutoRemove = false
	***REMOVED***

	query := url.Values***REMOVED******REMOVED***
	if containerName != "" ***REMOVED***
		query.Set("name", containerName)
	***REMOVED***

	body := configWrapper***REMOVED***
		Config:           config,
		HostConfig:       hostConfig,
		NetworkingConfig: networkingConfig,
	***REMOVED***

	serverResp, err := cli.post(ctx, "/containers/create", query, body, nil)
	if err != nil ***REMOVED***
		if serverResp.statusCode == 404 && strings.Contains(err.Error(), "No such image") ***REMOVED***
			return response, objectNotFoundError***REMOVED***object: "image", id: config.Image***REMOVED***
		***REMOVED***
		return response, err
	***REMOVED***

	err = json.NewDecoder(serverResp.body).Decode(&response)
	ensureReaderClosed(serverResp)
	return response, err
***REMOVED***
