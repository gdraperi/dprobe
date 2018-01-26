package cluster

import (
	apitypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	types "github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/daemon/cluster/convert"
	swarmapi "github.com/docker/swarmkit/api"
	"golang.org/x/net/context"
)

// GetTasks returns a list of tasks matching the filter options.
func (c *Cluster) GetTasks(options apitypes.TaskListOptions) ([]types.Task, error) ***REMOVED***
	var r *swarmapi.ListTasksResponse

	if err := c.lockedManagerAction(func(ctx context.Context, state nodeState) error ***REMOVED***
		filterTransform := func(filter filters.Args) error ***REMOVED***
			if filter.Contains("service") ***REMOVED***
				serviceFilters := filter.Get("service")
				for _, serviceFilter := range serviceFilters ***REMOVED***
					service, err := getService(ctx, state.controlClient, serviceFilter, false)
					if err != nil ***REMOVED***
						return err
					***REMOVED***
					filter.Del("service", serviceFilter)
					filter.Add("service", service.ID)
				***REMOVED***
			***REMOVED***
			if filter.Contains("node") ***REMOVED***
				nodeFilters := filter.Get("node")
				for _, nodeFilter := range nodeFilters ***REMOVED***
					node, err := getNode(ctx, state.controlClient, nodeFilter)
					if err != nil ***REMOVED***
						return err
					***REMOVED***
					filter.Del("node", nodeFilter)
					filter.Add("node", node.ID)
				***REMOVED***
			***REMOVED***
			if !filter.Contains("runtime") ***REMOVED***
				// default to only showing container tasks
				filter.Add("runtime", "container")
				filter.Add("runtime", "")
			***REMOVED***
			return nil
		***REMOVED***

		filters, err := newListTasksFilters(options.Filters, filterTransform)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		r, err = state.controlClient.ListTasks(
			ctx,
			&swarmapi.ListTasksRequest***REMOVED***Filters: filters***REMOVED***)
		return err
	***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	tasks := make([]types.Task, 0, len(r.Tasks))
	for _, task := range r.Tasks ***REMOVED***
		t, err := convert.TaskFromGRPC(*task)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		tasks = append(tasks, t)
	***REMOVED***
	return tasks, nil
***REMOVED***

// GetTask returns a task by an ID.
func (c *Cluster) GetTask(input string) (types.Task, error) ***REMOVED***
	var task *swarmapi.Task
	if err := c.lockedManagerAction(func(ctx context.Context, state nodeState) error ***REMOVED***
		t, err := getTask(ctx, state.controlClient, input)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		task = t
		return nil
	***REMOVED***); err != nil ***REMOVED***
		return types.Task***REMOVED******REMOVED***, err
	***REMOVED***
	return convert.TaskFromGRPC(*task)
***REMOVED***
