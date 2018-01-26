package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// PluginInstall installs a plugin
func (cli *Client) PluginInstall(ctx context.Context, name string, options types.PluginInstallOptions) (rc io.ReadCloser, err error) ***REMOVED***
	query := url.Values***REMOVED******REMOVED***
	if _, err := reference.ParseNormalizedNamed(options.RemoteRef); err != nil ***REMOVED***
		return nil, errors.Wrap(err, "invalid remote reference")
	***REMOVED***
	query.Set("remote", options.RemoteRef)

	privileges, err := cli.checkPluginPermissions(ctx, query, options)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// set name for plugin pull, if empty should default to remote reference
	query.Set("name", name)

	resp, err := cli.tryPluginPull(ctx, query, privileges, options.RegistryAuth)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	name = resp.header.Get("Docker-Plugin-Name")

	pr, pw := io.Pipe()
	go func() ***REMOVED*** // todo: the client should probably be designed more around the actual api
		_, err := io.Copy(pw, resp.body)
		if err != nil ***REMOVED***
			pw.CloseWithError(err)
			return
		***REMOVED***
		defer func() ***REMOVED***
			if err != nil ***REMOVED***
				delResp, _ := cli.delete(ctx, "/plugins/"+name, nil, nil)
				ensureReaderClosed(delResp)
			***REMOVED***
		***REMOVED***()
		if len(options.Args) > 0 ***REMOVED***
			if err := cli.PluginSet(ctx, name, options.Args); err != nil ***REMOVED***
				pw.CloseWithError(err)
				return
			***REMOVED***
		***REMOVED***

		if options.Disabled ***REMOVED***
			pw.Close()
			return
		***REMOVED***

		enableErr := cli.PluginEnable(ctx, name, types.PluginEnableOptions***REMOVED***Timeout: 0***REMOVED***)
		pw.CloseWithError(enableErr)
	***REMOVED***()
	return pr, nil
***REMOVED***

func (cli *Client) tryPluginPrivileges(ctx context.Context, query url.Values, registryAuth string) (serverResponse, error) ***REMOVED***
	headers := map[string][]string***REMOVED***"X-Registry-Auth": ***REMOVED***registryAuth***REMOVED******REMOVED***
	return cli.get(ctx, "/plugins/privileges", query, headers)
***REMOVED***

func (cli *Client) tryPluginPull(ctx context.Context, query url.Values, privileges types.PluginPrivileges, registryAuth string) (serverResponse, error) ***REMOVED***
	headers := map[string][]string***REMOVED***"X-Registry-Auth": ***REMOVED***registryAuth***REMOVED******REMOVED***
	return cli.post(ctx, "/plugins/pull", query, privileges, headers)
***REMOVED***

func (cli *Client) checkPluginPermissions(ctx context.Context, query url.Values, options types.PluginInstallOptions) (types.PluginPrivileges, error) ***REMOVED***
	resp, err := cli.tryPluginPrivileges(ctx, query, options.RegistryAuth)
	if resp.statusCode == http.StatusUnauthorized && options.PrivilegeFunc != nil ***REMOVED***
		// todo: do inspect before to check existing name before checking privileges
		newAuthHeader, privilegeErr := options.PrivilegeFunc()
		if privilegeErr != nil ***REMOVED***
			ensureReaderClosed(resp)
			return nil, privilegeErr
		***REMOVED***
		options.RegistryAuth = newAuthHeader
		resp, err = cli.tryPluginPrivileges(ctx, query, options.RegistryAuth)
	***REMOVED***
	if err != nil ***REMOVED***
		ensureReaderClosed(resp)
		return nil, err
	***REMOVED***

	var privileges types.PluginPrivileges
	if err := json.NewDecoder(resp.body).Decode(&privileges); err != nil ***REMOVED***
		ensureReaderClosed(resp)
		return nil, err
	***REMOVED***
	ensureReaderClosed(resp)

	if !options.AcceptAllPermissions && options.AcceptPermissionsFunc != nil && len(privileges) > 0 ***REMOVED***
		accept, err := options.AcceptPermissionsFunc(privileges)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if !accept ***REMOVED***
			return nil, pluginPermissionDenied***REMOVED***options.RemoteRef***REMOVED***
		***REMOVED***
	***REMOVED***
	return privileges, nil
***REMOVED***
