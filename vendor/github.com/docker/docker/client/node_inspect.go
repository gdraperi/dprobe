package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

// NodeInspectWithRaw returns the node information.
func (cli *Client) NodeInspectWithRaw(ctx context.Context, nodeID string) (swarm.Node, []byte, error) ***REMOVED***
	serverResp, err := cli.get(ctx, "/nodes/"+nodeID, nil, nil)
	if err != nil ***REMOVED***
		return swarm.Node***REMOVED******REMOVED***, nil, wrapResponseError(err, serverResp, "node", nodeID)
	***REMOVED***
	defer ensureReaderClosed(serverResp)

	body, err := ioutil.ReadAll(serverResp.body)
	if err != nil ***REMOVED***
		return swarm.Node***REMOVED******REMOVED***, nil, err
	***REMOVED***

	var response swarm.Node
	rdr := bytes.NewReader(body)
	err = json.NewDecoder(rdr).Decode(&response)
	return response, body, err
***REMOVED***
