package serf

import (
	"github.com/hashicorp/memberlist"
)

type conflictDelegate struct ***REMOVED***
	serf *Serf
***REMOVED***

func (c *conflictDelegate) NotifyConflict(existing, other *memberlist.Node) ***REMOVED***
	c.serf.handleNodeConflict(existing, other)
***REMOVED***
