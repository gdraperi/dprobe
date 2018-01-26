package logger

import (
	"github.com/docker/docker/api/types/backend"
)

func (m *Message) copy() *Message ***REMOVED***
	msg := &Message***REMOVED***
		Source:    m.Source,
		Partial:   m.Partial,
		Timestamp: m.Timestamp,
	***REMOVED***

	if m.Attrs != nil ***REMOVED***
		msg.Attrs = make([]backend.LogAttr, len(m.Attrs))
		copy(msg.Attrs, m.Attrs)
	***REMOVED***

	msg.Line = append(make([]byte, 0, len(m.Line)), m.Line...)
	return msg
***REMOVED***
