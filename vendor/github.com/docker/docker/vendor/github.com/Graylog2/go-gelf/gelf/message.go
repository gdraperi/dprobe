package gelf

import (
	"bytes"
	"encoding/json"
	"time"
)

// Message represents the contents of the GELF message.  It is gzipped
// before sending.
type Message struct ***REMOVED***
	Version  string                 `json:"version"`
	Host     string                 `json:"host"`
	Short    string                 `json:"short_message"`
	Full     string                 `json:"full_message,omitempty"`
	TimeUnix float64                `json:"timestamp"`
	Level    int32                  `json:"level,omitempty"`
	Facility string                 `json:"facility,omitempty"`
	Extra    map[string]interface***REMOVED******REMOVED*** `json:"-"`
	RawExtra json.RawMessage        `json:"-"`
***REMOVED***

// Syslog severity levels
const (
	LOG_EMERG   = 0
	LOG_ALERT   = 1
	LOG_CRIT    = 2
	LOG_ERR     = 3
	LOG_WARNING = 4
	LOG_NOTICE  = 5
	LOG_INFO    = 6
	LOG_DEBUG   = 7
)

func (m *Message) MarshalJSONBuf(buf *bytes.Buffer) error ***REMOVED***
	b, err := json.Marshal(m)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	// write up until the final ***REMOVED***
	if _, err = buf.Write(b[:len(b)-1]); err != nil ***REMOVED***
		return err
	***REMOVED***
	if len(m.Extra) > 0 ***REMOVED***
		eb, err := json.Marshal(m.Extra)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		// merge serialized message + serialized extra map
		if err = buf.WriteByte(','); err != nil ***REMOVED***
			return err
		***REMOVED***
		// write serialized extra bytes, without enclosing quotes
		if _, err = buf.Write(eb[1 : len(eb)-1]); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if len(m.RawExtra) > 0 ***REMOVED***
		if err := buf.WriteByte(','); err != nil ***REMOVED***
			return err
		***REMOVED***

		// write serialized extra bytes, without enclosing quotes
		if _, err = buf.Write(m.RawExtra[1 : len(m.RawExtra)-1]); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// write final closing quotes
	return buf.WriteByte('***REMOVED***')
***REMOVED***

func (m *Message) UnmarshalJSON(data []byte) error ***REMOVED***
	i := make(map[string]interface***REMOVED******REMOVED***, 16)
	if err := json.Unmarshal(data, &i); err != nil ***REMOVED***
		return err
	***REMOVED***
	for k, v := range i ***REMOVED***
		if k[0] == '_' ***REMOVED***
			if m.Extra == nil ***REMOVED***
				m.Extra = make(map[string]interface***REMOVED******REMOVED***, 1)
			***REMOVED***
			m.Extra[k] = v
			continue
		***REMOVED***
		switch k ***REMOVED***
		case "version":
			m.Version = v.(string)
		case "host":
			m.Host = v.(string)
		case "short_message":
			m.Short = v.(string)
		case "full_message":
			m.Full = v.(string)
		case "timestamp":
			m.TimeUnix = v.(float64)
		case "level":
			m.Level = int32(v.(float64))
		case "facility":
			m.Facility = v.(string)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (m *Message) toBytes() (messageBytes []byte, err error) ***REMOVED***
	buf := newBuffer()
	defer bufPool.Put(buf)
	if err = m.MarshalJSONBuf(buf); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	messageBytes = buf.Bytes()
	return messageBytes, nil
***REMOVED***

func constructMessage(p []byte, hostname string, facility string, file string, line int) (m *Message) ***REMOVED***
	// remove trailing and leading whitespace
	p = bytes.TrimSpace(p)

	// If there are newlines in the message, use the first line
	// for the short message and set the full message to the
	// original input.  If the input has no newlines, stick the
	// whole thing in Short.
	short := p
	full := []byte("")
	if i := bytes.IndexRune(p, '\n'); i > 0 ***REMOVED***
		short = p[:i]
		full = p
	***REMOVED***

	m = &Message***REMOVED***
		Version:  "1.1",
		Host:     hostname,
		Short:    string(short),
		Full:     string(full),
		TimeUnix: float64(time.Now().Unix()),
		Level:    6, // info
		Facility: facility,
		Extra: map[string]interface***REMOVED******REMOVED******REMOVED***
			"_file": file,
			"_line": line,
		***REMOVED***,
	***REMOVED***

	return m
***REMOVED***
