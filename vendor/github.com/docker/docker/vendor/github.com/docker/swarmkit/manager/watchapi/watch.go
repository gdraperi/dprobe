package watchapi

import (
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/manager/state"
	"github.com/docker/swarmkit/manager/state/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Watch starts a stream that returns any changes to objects that match
// the specified selectors. When the stream begins, it immediately sends
// an empty message back to the client. It is important to wait for
// this message before taking any actions that depend on an established
// stream of changes for consistency.
func (s *Server) Watch(request *api.WatchRequest, stream api.Watch_WatchServer) error ***REMOVED***
	ctx := stream.Context()

	s.mu.Lock()
	pctx := s.pctx
	s.mu.Unlock()
	if pctx == nil ***REMOVED***
		return errNotRunning
	***REMOVED***

	watchArgs, err := api.ConvertWatchArgs(request.Entries)
	if err != nil ***REMOVED***
		return status.Errorf(codes.InvalidArgument, "%s", err.Error())
	***REMOVED***

	watchArgs = append(watchArgs, state.EventCommit***REMOVED******REMOVED***)
	watch, cancel, err := store.WatchFrom(s.store, request.ResumeFrom, watchArgs...)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer cancel()

	// TODO(aaronl): Send current version in this WatchMessage?
	if err := stream.Send(&api.WatchMessage***REMOVED******REMOVED***); err != nil ***REMOVED***
		return err
	***REMOVED***

	var events []*api.WatchMessage_Event
	for ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return ctx.Err()
		case <-pctx.Done():
			return pctx.Err()
		case event := <-watch:
			if commitEvent, ok := event.(state.EventCommit); ok && len(events) > 0 ***REMOVED***
				if err := stream.Send(&api.WatchMessage***REMOVED***Events: events, Version: commitEvent.Version***REMOVED***); err != nil ***REMOVED***
					return err
				***REMOVED***
				events = nil
			***REMOVED*** else if eventMessage := api.WatchMessageEvent(event.(api.Event)); eventMessage != nil ***REMOVED***
				if !request.IncludeOldObject ***REMOVED***
					eventMessage.OldObject = nil
				***REMOVED***
				events = append(events, eventMessage)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
