package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// ImageInspectWithRaw returns the image information and its raw representation.
func (cli *Client) ImageInspectWithRaw(ctx context.Context, imageID string) (types.ImageInspect, []byte, error) ***REMOVED***
	serverResp, err := cli.get(ctx, "/images/"+imageID+"/json", nil, nil)
	if err != nil ***REMOVED***
		return types.ImageInspect***REMOVED******REMOVED***, nil, wrapResponseError(err, serverResp, "image", imageID)
	***REMOVED***
	defer ensureReaderClosed(serverResp)

	body, err := ioutil.ReadAll(serverResp.body)
	if err != nil ***REMOVED***
		return types.ImageInspect***REMOVED******REMOVED***, nil, err
	***REMOVED***

	var response types.ImageInspect
	rdr := bytes.NewReader(body)
	err = json.NewDecoder(rdr).Decode(&response)
	return response, body, err
***REMOVED***
