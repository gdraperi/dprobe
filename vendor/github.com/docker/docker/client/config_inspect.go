package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

// ConfigInspectWithRaw returns the config information with raw data
func (cli *Client) ConfigInspectWithRaw(ctx context.Context, id string) (swarm.Config, []byte, error) ***REMOVED***
	if err := cli.NewVersionError("1.30", "config inspect"); err != nil ***REMOVED***
		return swarm.Config***REMOVED******REMOVED***, nil, err
	***REMOVED***
	resp, err := cli.get(ctx, "/configs/"+id, nil, nil)
	if err != nil ***REMOVED***
		return swarm.Config***REMOVED******REMOVED***, nil, wrapResponseError(err, resp, "config", id)
	***REMOVED***
	defer ensureReaderClosed(resp)

	body, err := ioutil.ReadAll(resp.body)
	if err != nil ***REMOVED***
		return swarm.Config***REMOVED******REMOVED***, nil, err
	***REMOVED***

	var config swarm.Config
	rdr := bytes.NewReader(body)
	err = json.NewDecoder(rdr).Decode(&config)

	return config, body, err
***REMOVED***
