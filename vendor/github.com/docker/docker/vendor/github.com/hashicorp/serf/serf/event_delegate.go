package serf

import (
	"github.com/hashicorp/memberlist"
)

type eventDelegate struct ***REMOVED***
	serf *Serf
***REMOVED***

func (e *eventDelegate) NotifyJoin(n *memberlist.Node) ***REMOVED***
	e.serf.handleNodeJoin(n)
***REMOVED***

func (e *eventDelegate) NotifyLeave(n *memberlist.Node) ***REMOVED***
	e.serf.handleNodeLeave(n)
***REMOVED***

func (e *eventDelegate) NotifyUpdate(n *memberlist.Node) ***REMOVED***
	e.serf.handleNodeUpdate(n)
***REMOVED***
