package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// PluginInspectWithRaw inspects an existing plugin
func (cli *Client) PluginInspectWithRaw(ctx context.Context, name string) (*types.Plugin, []byte, error) ***REMOVED***
	resp, err := cli.get(ctx, "/plugins/"+name+"/json", nil, nil)
	if err != nil ***REMOVED***
		return nil, nil, wrapResponseError(err, resp, "plugin", name)
	***REMOVED***

	defer ensureReaderClosed(resp)
	body, err := ioutil.ReadAll(resp.body)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	var p types.Plugin
	rdr := bytes.NewReader(body)
	err = json.NewDecoder(rdr).Decode(&p)
	return &p, body, err
***REMOVED***
