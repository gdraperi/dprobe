package client

import (
	"net/url"
	"strconv"

	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

// ConfigUpdate attempts to update a Config
func (cli *Client) ConfigUpdate(ctx context.Context, id string, version swarm.Version, config swarm.ConfigSpec) error ***REMOVED***
	if err := cli.NewVersionError("1.30", "config update"); err != nil ***REMOVED***
		return err
	***REMOVED***
	query := url.Values***REMOVED******REMOVED***
	query.Set("version", strconv.FormatUint(version.Index, 10))
	resp, err := cli.post(ctx, "/configs/"+id+"/update", query, config, nil)
	ensureReaderClosed(resp)
	return err
***REMOVED***
