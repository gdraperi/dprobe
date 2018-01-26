package client

import (
	"net/url"
	"strconv"

	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

// SecretUpdate attempts to update a Secret
func (cli *Client) SecretUpdate(ctx context.Context, id string, version swarm.Version, secret swarm.SecretSpec) error ***REMOVED***
	if err := cli.NewVersionError("1.25", "secret update"); err != nil ***REMOVED***
		return err
	***REMOVED***
	query := url.Values***REMOVED******REMOVED***
	query.Set("version", strconv.FormatUint(version.Index, 10))
	resp, err := cli.post(ctx, "/secrets/"+id+"/update", query, secret, nil)
	ensureReaderClosed(resp)
	return err
***REMOVED***
