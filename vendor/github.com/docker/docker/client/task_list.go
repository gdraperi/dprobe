package client

import (
	"encoding/json"
	"net/url"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

// TaskList returns the list of tasks.
func (cli *Client) TaskList(ctx context.Context, options types.TaskListOptions) ([]swarm.Task, error) ***REMOVED***
	query := url.Values***REMOVED******REMOVED***

	if options.Filters.Len() > 0 ***REMOVED***
		filterJSON, err := filters.ToJSON(options.Filters)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		query.Set("filters", filterJSON)
	***REMOVED***

	resp, err := cli.get(ctx, "/tasks", query, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var tasks []swarm.Task
	err = json.NewDecoder(resp.body).Decode(&tasks)
	ensureReaderClosed(resp)
	return tasks, err
***REMOVED***
