package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

// SecretInspectWithRaw returns the secret information with raw data
func (cli *Client) SecretInspectWithRaw(ctx context.Context, id string) (swarm.Secret, []byte, error) ***REMOVED***
	if err := cli.NewVersionError("1.25", "secret inspect"); err != nil ***REMOVED***
		return swarm.Secret***REMOVED******REMOVED***, nil, err
	***REMOVED***
	resp, err := cli.get(ctx, "/secrets/"+id, nil, nil)
	if err != nil ***REMOVED***
		return swarm.Secret***REMOVED******REMOVED***, nil, wrapResponseError(err, resp, "secret", id)
	***REMOVED***
	defer ensureReaderClosed(resp)

	body, err := ioutil.ReadAll(resp.body)
	if err != nil ***REMOVED***
		return swarm.Secret***REMOVED******REMOVED***, nil, err
	***REMOVED***

	var secret swarm.Secret
	rdr := bytes.NewReader(body)
	err = json.NewDecoder(rdr).Decode(&secret)

	return secret, body, err
***REMOVED***
