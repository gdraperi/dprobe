package networkdb

import (
	"net"

	"github.com/docker/go-events"
)

type opType uint8

const (
	opCreate opType = 1 + iota
	opUpdate
	opDelete
)

type event struct ***REMOVED***
	Table     string
	NetworkID string
	Key       string
	Value     []byte
***REMOVED***

// NodeTable represents table event for node join and leave
const NodeTable = "NodeTable"

// NodeAddr represents the value carried for node event in NodeTable
type NodeAddr struct ***REMOVED***
	Addr net.IP
***REMOVED***

// CreateEvent generates a table entry create event to the watchers
type CreateEvent event

// UpdateEvent generates a table entry update event to the watchers
type UpdateEvent event

// DeleteEvent generates a table entry delete event to the watchers
type DeleteEvent event

// Watch creates a watcher with filters for a particular table or
// network or key or any combination of the tuple. If any of the
// filter is an empty string it acts as a wildcard for that
// field. Watch returns a channel of events, where the events will be
// sent.
func (nDB *NetworkDB) Watch(tname, nid, key string) (*events.Channel, func()) ***REMOVED***
	var matcher events.Matcher

	if tname != "" || nid != "" || key != "" ***REMOVED***
		matcher = events.MatcherFunc(func(ev events.Event) bool ***REMOVED***
			var evt event
			switch ev := ev.(type) ***REMOVED***
			case CreateEvent:
				evt = event(ev)
			case UpdateEvent:
				evt = event(ev)
			case DeleteEvent:
				evt = event(ev)
			***REMOVED***

			if tname != "" && evt.Table != tname ***REMOVED***
				return false
			***REMOVED***

			if nid != "" && evt.NetworkID != nid ***REMOVED***
				return false
			***REMOVED***

			if key != "" && evt.Key != key ***REMOVED***
				return false
			***REMOVED***

			return true
		***REMOVED***)
	***REMOVED***

	ch := events.NewChannel(0)
	sink := events.Sink(events.NewQueue(ch))

	if matcher != nil ***REMOVED***
		sink = events.NewFilter(sink, matcher)
	***REMOVED***

	nDB.broadcaster.Add(sink)
	return ch, func() ***REMOVED***
		nDB.broadcaster.Remove(sink)
		ch.Close()
		sink.Close()
	***REMOVED***
***REMOVED***

func makeEvent(op opType, tname, nid, key string, value []byte) events.Event ***REMOVED***
	ev := event***REMOVED***
		Table:     tname,
		NetworkID: nid,
		Key:       key,
		Value:     value,
	***REMOVED***

	switch op ***REMOVED***
	case opCreate:
		return CreateEvent(ev)
	case opUpdate:
		return UpdateEvent(ev)
	case opDelete:
		return DeleteEvent(ev)
	***REMOVED***

	return nil
***REMOVED***
