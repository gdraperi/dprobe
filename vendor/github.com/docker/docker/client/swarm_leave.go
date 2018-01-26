package client

import (
	"net/url"

	"golang.org/x/net/context"
)

// SwarmLeave leaves the swarm.
func (cli *Client) SwarmLeave(ctx context.Context, force bool) error ***REMOVED***
	query := url.Values***REMOVED******REMOVED***
	if force ***REMOVED***
		query.Set("force", "1")
	***REMOVED***
	resp, err := cli.post(ctx, "/swarm/leave", query, nil, nil)
	ensureReaderClosed(resp)
	return err
***REMOVED***
