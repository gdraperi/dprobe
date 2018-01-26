package client

import (
	"io"
	"net/url"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// PluginUpgrade upgrades a plugin
func (cli *Client) PluginUpgrade(ctx context.Context, name string, options types.PluginInstallOptions) (rc io.ReadCloser, err error) ***REMOVED***
	if err := cli.NewVersionError("1.26", "plugin upgrade"); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	query := url.Values***REMOVED******REMOVED***
	if _, err := reference.ParseNormalizedNamed(options.RemoteRef); err != nil ***REMOVED***
		return nil, errors.Wrap(err, "invalid remote reference")
	***REMOVED***
	query.Set("remote", options.RemoteRef)

	privileges, err := cli.checkPluginPermissions(ctx, query, options)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	resp, err := cli.tryPluginUpgrade(ctx, query, privileges, name, options.RegistryAuth)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return resp.body, nil
***REMOVED***

func (cli *Client) tryPluginUpgrade(ctx context.Context, query url.Values, privileges types.PluginPrivileges, name, registryAuth string) (serverResponse, error) ***REMOVED***
	headers := map[string][]string***REMOVED***"X-Registry-Auth": ***REMOVED***registryAuth***REMOVED******REMOVED***
	return cli.post(ctx, "/plugins/"+name+"/upgrade", query, privileges, headers)
***REMOVED***
