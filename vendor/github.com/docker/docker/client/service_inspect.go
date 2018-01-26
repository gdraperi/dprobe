package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

// ServiceInspectWithRaw returns the service information and the raw data.
func (cli *Client) ServiceInspectWithRaw(ctx context.Context, serviceID string, opts types.ServiceInspectOptions) (swarm.Service, []byte, error) ***REMOVED***
	query := url.Values***REMOVED******REMOVED***
	query.Set("insertDefaults", fmt.Sprintf("%v", opts.InsertDefaults))
	serverResp, err := cli.get(ctx, "/services/"+serviceID, query, nil)
	if err != nil ***REMOVED***
		return swarm.Service***REMOVED******REMOVED***, nil, wrapResponseError(err, serverResp, "service", serviceID)
	***REMOVED***
	defer ensureReaderClosed(serverResp)

	body, err := ioutil.ReadAll(serverResp.body)
	if err != nil ***REMOVED***
		return swarm.Service***REMOVED******REMOVED***, nil, err
	***REMOVED***

	var response swarm.Service
	rdr := bytes.NewReader(body)
	err = json.NewDecoder(rdr).Decode(&response)
	return response, body, err
***REMOVED***
