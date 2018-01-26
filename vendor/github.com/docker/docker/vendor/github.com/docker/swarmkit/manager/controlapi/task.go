package controlapi

import (
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/api/naming"
	"github.com/docker/swarmkit/manager/orchestrator"
	"github.com/docker/swarmkit/manager/state/store"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetTask returns a Task given a TaskID.
// - Returns `InvalidArgument` if TaskID is not provided.
// - Returns `NotFound` if the Task is not found.
func (s *Server) GetTask(ctx context.Context, request *api.GetTaskRequest) (*api.GetTaskResponse, error) ***REMOVED***
	if request.TaskID == "" ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, errInvalidArgument.Error())
	***REMOVED***

	var task *api.Task
	s.store.View(func(tx store.ReadTx) ***REMOVED***
		task = store.GetTask(tx, request.TaskID)
	***REMOVED***)
	if task == nil ***REMOVED***
		return nil, status.Errorf(codes.NotFound, "task %s not found", request.TaskID)
	***REMOVED***
	return &api.GetTaskResponse***REMOVED***
		Task: task,
	***REMOVED***, nil
***REMOVED***

// RemoveTask removes a Task referenced by TaskID.
// - Returns `InvalidArgument` if TaskID is not provided.
// - Returns `NotFound` if the Task is not found.
// - Returns an error if the deletion fails.
func (s *Server) RemoveTask(ctx context.Context, request *api.RemoveTaskRequest) (*api.RemoveTaskResponse, error) ***REMOVED***
	if request.TaskID == "" ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, errInvalidArgument.Error())
	***REMOVED***

	err := s.store.Update(func(tx store.Tx) error ***REMOVED***
		return store.DeleteTask(tx, request.TaskID)
	***REMOVED***)
	if err != nil ***REMOVED***
		if err == store.ErrNotExist ***REMOVED***
			return nil, status.Errorf(codes.NotFound, "task %s not found", request.TaskID)
		***REMOVED***
		return nil, err
	***REMOVED***
	return &api.RemoveTaskResponse***REMOVED******REMOVED***, nil
***REMOVED***

func filterTasks(candidates []*api.Task, filters ...func(*api.Task) bool) []*api.Task ***REMOVED***
	result := []*api.Task***REMOVED******REMOVED***

	for _, c := range candidates ***REMOVED***
		match := true
		for _, f := range filters ***REMOVED***
			if !f(c) ***REMOVED***
				match = false
				break
			***REMOVED***
		***REMOVED***
		if match ***REMOVED***
			result = append(result, c)
		***REMOVED***
	***REMOVED***

	return result
***REMOVED***

// ListTasks returns a list of all tasks.
func (s *Server) ListTasks(ctx context.Context, request *api.ListTasksRequest) (*api.ListTasksResponse, error) ***REMOVED***
	var (
		tasks []*api.Task
		err   error
	)

	s.store.View(func(tx store.ReadTx) ***REMOVED***
		switch ***REMOVED***
		case request.Filters != nil && len(request.Filters.Names) > 0:
			tasks, err = store.FindTasks(tx, buildFilters(store.ByName, request.Filters.Names))
		case request.Filters != nil && len(request.Filters.NamePrefixes) > 0:
			tasks, err = store.FindTasks(tx, buildFilters(store.ByNamePrefix, request.Filters.NamePrefixes))
		case request.Filters != nil && len(request.Filters.IDPrefixes) > 0:
			tasks, err = store.FindTasks(tx, buildFilters(store.ByIDPrefix, request.Filters.IDPrefixes))
		case request.Filters != nil && len(request.Filters.ServiceIDs) > 0:
			tasks, err = store.FindTasks(tx, buildFilters(store.ByServiceID, request.Filters.ServiceIDs))
		case request.Filters != nil && len(request.Filters.Runtimes) > 0:
			tasks, err = store.FindTasks(tx, buildFilters(store.ByRuntime, request.Filters.Runtimes))
		case request.Filters != nil && len(request.Filters.NodeIDs) > 0:
			tasks, err = store.FindTasks(tx, buildFilters(store.ByNodeID, request.Filters.NodeIDs))
		case request.Filters != nil && len(request.Filters.DesiredStates) > 0:
			filters := make([]store.By, 0, len(request.Filters.DesiredStates))
			for _, v := range request.Filters.DesiredStates ***REMOVED***
				filters = append(filters, store.ByDesiredState(v))
			***REMOVED***
			tasks, err = store.FindTasks(tx, store.Or(filters...))
		default:
			tasks, err = store.FindTasks(tx, store.All)
		***REMOVED***

		if err != nil || request.Filters == nil ***REMOVED***
			return
		***REMOVED***

		tasks = filterTasks(tasks,
			func(e *api.Task) bool ***REMOVED***
				return filterContains(naming.Task(e), request.Filters.Names)
			***REMOVED***,
			func(e *api.Task) bool ***REMOVED***
				return filterContainsPrefix(naming.Task(e), request.Filters.NamePrefixes)
			***REMOVED***,
			func(e *api.Task) bool ***REMOVED***
				return filterContainsPrefix(e.ID, request.Filters.IDPrefixes)
			***REMOVED***,
			func(e *api.Task) bool ***REMOVED***
				return filterMatchLabels(e.ServiceAnnotations.Labels, request.Filters.Labels)
			***REMOVED***,
			func(e *api.Task) bool ***REMOVED***
				return filterContains(e.ServiceID, request.Filters.ServiceIDs)
			***REMOVED***,
			func(e *api.Task) bool ***REMOVED***
				return filterContains(e.NodeID, request.Filters.NodeIDs)
			***REMOVED***,
			func(e *api.Task) bool ***REMOVED***
				if len(request.Filters.Runtimes) == 0 ***REMOVED***
					return true
				***REMOVED***
				r, err := naming.Runtime(e.Spec)
				if err != nil ***REMOVED***
					return false
				***REMOVED***
				return filterContains(r, request.Filters.Runtimes)
			***REMOVED***,
			func(e *api.Task) bool ***REMOVED***
				if len(request.Filters.DesiredStates) == 0 ***REMOVED***
					return true
				***REMOVED***
				for _, c := range request.Filters.DesiredStates ***REMOVED***
					if c == e.DesiredState ***REMOVED***
						return true
					***REMOVED***
				***REMOVED***
				return false
			***REMOVED***,
			func(e *api.Task) bool ***REMOVED***
				if !request.Filters.UpToDate ***REMOVED***
					return true
				***REMOVED***

				service := store.GetService(tx, e.ServiceID)
				if service == nil ***REMOVED***
					return false
				***REMOVED***

				return !orchestrator.IsTaskDirty(service, e)
			***REMOVED***,
		)
	***REMOVED***)

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &api.ListTasksResponse***REMOVED***
		Tasks: tasks,
	***REMOVED***, nil
***REMOVED***
