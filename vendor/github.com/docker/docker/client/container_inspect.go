package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/url"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// ContainerInspect returns the container information.
func (cli *Client) ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error) ***REMOVED***
	serverResp, err := cli.get(ctx, "/containers/"+containerID+"/json", nil, nil)
	if err != nil ***REMOVED***
		return types.ContainerJSON***REMOVED******REMOVED***, wrapResponseError(err, serverResp, "container", containerID)
	***REMOVED***

	var response types.ContainerJSON
	err = json.NewDecoder(serverResp.body).Decode(&response)
	ensureReaderClosed(serverResp)
	return response, err
***REMOVED***

// ContainerInspectWithRaw returns the container information and its raw representation.
func (cli *Client) ContainerInspectWithRaw(ctx context.Context, containerID string, getSize bool) (types.ContainerJSON, []byte, error) ***REMOVED***
	query := url.Values***REMOVED******REMOVED***
	if getSize ***REMOVED***
		query.Set("size", "1")
	***REMOVED***
	serverResp, err := cli.get(ctx, "/containers/"+containerID+"/json", query, nil)
	if err != nil ***REMOVED***
		return types.ContainerJSON***REMOVED******REMOVED***, nil, wrapResponseError(err, serverResp, "container", containerID)
	***REMOVED***
	defer ensureReaderClosed(serverResp)

	body, err := ioutil.ReadAll(serverResp.body)
	if err != nil ***REMOVED***
		return types.ContainerJSON***REMOVED******REMOVED***, nil, err
	***REMOVED***

	var response types.ContainerJSON
	rdr := bytes.NewReader(body)
	err = json.NewDecoder(rdr).Decode(&response)
	return response, body, err
***REMOVED***
