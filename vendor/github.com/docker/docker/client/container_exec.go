package client

import (
	"encoding/json"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// ContainerExecCreate creates a new exec configuration to run an exec process.
func (cli *Client) ContainerExecCreate(ctx context.Context, container string, config types.ExecConfig) (types.IDResponse, error) ***REMOVED***
	var response types.IDResponse

	if err := cli.NewVersionError("1.25", "env"); len(config.Env) != 0 && err != nil ***REMOVED***
		return response, err
	***REMOVED***

	resp, err := cli.post(ctx, "/containers/"+container+"/exec", nil, config, nil)
	if err != nil ***REMOVED***
		return response, err
	***REMOVED***
	err = json.NewDecoder(resp.body).Decode(&response)
	ensureReaderClosed(resp)
	return response, err
***REMOVED***

// ContainerExecStart starts an exec process already created in the docker host.
func (cli *Client) ContainerExecStart(ctx context.Context, execID string, config types.ExecStartCheck) error ***REMOVED***
	resp, err := cli.post(ctx, "/exec/"+execID+"/start", nil, config, nil)
	ensureReaderClosed(resp)
	return err
***REMOVED***

// ContainerExecAttach attaches a connection to an exec process in the server.
// It returns a types.HijackedConnection with the hijacked connection
// and the a reader to get output. It's up to the called to close
// the hijacked connection by calling types.HijackedResponse.Close.
func (cli *Client) ContainerExecAttach(ctx context.Context, execID string, config types.ExecStartCheck) (types.HijackedResponse, error) ***REMOVED***
	headers := map[string][]string***REMOVED***"Content-Type": ***REMOVED***"application/json"***REMOVED******REMOVED***
	return cli.postHijacked(ctx, "/exec/"+execID+"/start", nil, config, headers)
***REMOVED***

// ContainerExecInspect returns information about a specific exec process on the docker host.
func (cli *Client) ContainerExecInspect(ctx context.Context, execID string) (types.ContainerExecInspect, error) ***REMOVED***
	var response types.ContainerExecInspect
	resp, err := cli.get(ctx, "/exec/"+execID+"/json", nil, nil)
	if err != nil ***REMOVED***
		return response, err
	***REMOVED***

	err = json.NewDecoder(resp.body).Decode(&response)
	ensureReaderClosed(resp)
	return response, err
***REMOVED***
