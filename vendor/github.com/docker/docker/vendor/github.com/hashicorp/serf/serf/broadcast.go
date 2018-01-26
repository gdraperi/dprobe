package serf

import (
	"github.com/hashicorp/memberlist"
)

// broadcast is an implementation of memberlist.Broadcast and is used
// to manage broadcasts across the memberlist channel that are related
// only to Serf.
type broadcast struct ***REMOVED***
	msg    []byte
	notify chan<- struct***REMOVED******REMOVED***
***REMOVED***

func (b *broadcast) Invalidates(other memberlist.Broadcast) bool ***REMOVED***
	return false
***REMOVED***

func (b *broadcast) Message() []byte ***REMOVED***
	return b.msg
***REMOVED***

func (b *broadcast) Finished() ***REMOVED***
	if b.notify != nil ***REMOVED***
		close(b.notify)
	***REMOVED***
***REMOVED***
