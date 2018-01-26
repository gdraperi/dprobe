package state

import (
	"github.com/docker/go-events"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/watch"
)

// EventCommit delineates a transaction boundary.
type EventCommit struct ***REMOVED***
	Version *api.Version
***REMOVED***

// Matches returns true if this event is a commit event.
func (e EventCommit) Matches(watchEvent events.Event) bool ***REMOVED***
	_, ok := watchEvent.(EventCommit)
	return ok
***REMOVED***

// TaskCheckStateGreaterThan is a TaskCheckFunc for checking task state.
func TaskCheckStateGreaterThan(t1, t2 *api.Task) bool ***REMOVED***
	return t2.Status.State > t1.Status.State
***REMOVED***

// NodeCheckState is a NodeCheckFunc for matching node state.
func NodeCheckState(n1, n2 *api.Node) bool ***REMOVED***
	return n1.Status.State == n2.Status.State
***REMOVED***

// Watch takes a variable number of events to match against. The subscriber
// will receive events that match any of the arguments passed to Watch.
//
// Examples:
//
// // subscribe to all events
// Watch(q)
//
// // subscribe to all UpdateTask events
// Watch(q, EventUpdateTask***REMOVED******REMOVED***)
//
// // subscribe to all task-related events
// Watch(q, EventUpdateTask***REMOVED******REMOVED***, EventCreateTask***REMOVED******REMOVED***, EventDeleteTask***REMOVED******REMOVED***)
//
// // subscribe to UpdateTask for node 123
// Watch(q, EventUpdateTask***REMOVED***Task: &api.Task***REMOVED***NodeID: 123***REMOVED***,
//                         Checks: []TaskCheckFunc***REMOVED***TaskCheckNodeID***REMOVED******REMOVED***)
//
// // subscribe to UpdateTask for node 123, as well as CreateTask
// // for node 123 that also has ServiceID set to "abc"
// Watch(q, EventUpdateTask***REMOVED***Task: &api.Task***REMOVED***NodeID: 123***REMOVED***,
//                         Checks: []TaskCheckFunc***REMOVED***TaskCheckNodeID***REMOVED******REMOVED***,
//         EventCreateTask***REMOVED***Task: &api.Task***REMOVED***NodeID: 123, ServiceID: "abc"***REMOVED***,
//                         Checks: []TaskCheckFunc***REMOVED***TaskCheckNodeID,
//                                                 func(t1, t2 *api.Task) bool ***REMOVED***
//                                                         return t1.ServiceID == t2.ServiceID
//                                             ***REMOVED******REMOVED******REMOVED***)
func Watch(queue *watch.Queue, specifiers ...api.Event) (eventq chan events.Event, cancel func()) ***REMOVED***
	if len(specifiers) == 0 ***REMOVED***
		return queue.Watch()
	***REMOVED***
	return queue.CallbackWatch(Matcher(specifiers...))
***REMOVED***

// Matcher returns an events.Matcher that Matches the specifiers with OR logic.
func Matcher(specifiers ...api.Event) events.MatcherFunc ***REMOVED***
	return events.MatcherFunc(func(event events.Event) bool ***REMOVED***
		for _, s := range specifiers ***REMOVED***
			if s.Matches(event) ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
		return false
	***REMOVED***)
***REMOVED***
