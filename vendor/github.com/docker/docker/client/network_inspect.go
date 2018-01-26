package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/url"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// NetworkInspect returns the information for a specific network configured in the docker host.
func (cli *Client) NetworkInspect(ctx context.Context, networkID string, options types.NetworkInspectOptions) (types.NetworkResource, error) ***REMOVED***
	networkResource, _, err := cli.NetworkInspectWithRaw(ctx, networkID, options)
	return networkResource, err
***REMOVED***

// NetworkInspectWithRaw returns the information for a specific network configured in the docker host and its raw representation.
func (cli *Client) NetworkInspectWithRaw(ctx context.Context, networkID string, options types.NetworkInspectOptions) (types.NetworkResource, []byte, error) ***REMOVED***
	var (
		networkResource types.NetworkResource
		resp            serverResponse
		err             error
	)
	query := url.Values***REMOVED******REMOVED***
	if options.Verbose ***REMOVED***
		query.Set("verbose", "true")
	***REMOVED***
	if options.Scope != "" ***REMOVED***
		query.Set("scope", options.Scope)
	***REMOVED***
	resp, err = cli.get(ctx, "/networks/"+networkID, query, nil)
	if err != nil ***REMOVED***
		return networkResource, nil, wrapResponseError(err, resp, "network", networkID)
	***REMOVED***
	defer ensureReaderClosed(resp)

	body, err := ioutil.ReadAll(resp.body)
	if err != nil ***REMOVED***
		return networkResource, nil, err
	***REMOVED***
	rdr := bytes.NewReader(body)
	err = json.NewDecoder(rdr).Decode(&networkResource)
	return networkResource, body, err
***REMOVED***
