package jsonlog

import (
	"bytes"
	"encoding/json"
	"time"
	"unicode/utf8"
)

// JSONLogs marshals encoded JSONLog objects
type JSONLogs struct ***REMOVED***
	Log     []byte    `json:"log,omitempty"`
	Stream  string    `json:"stream,omitempty"`
	Created time.Time `json:"time"`

	// json-encoded bytes
	RawAttrs json.RawMessage `json:"attrs,omitempty"`
***REMOVED***

// MarshalJSONBuf is an optimized JSON marshaller that avoids reflection
// and unnecessary allocation.
func (mj *JSONLogs) MarshalJSONBuf(buf *bytes.Buffer) error ***REMOVED***
	var first = true

	buf.WriteString(`***REMOVED***`)
	if len(mj.Log) != 0 ***REMOVED***
		first = false
		buf.WriteString(`"log":`)
		ffjsonWriteJSONBytesAsString(buf, mj.Log)
	***REMOVED***
	if len(mj.Stream) != 0 ***REMOVED***
		if first ***REMOVED***
			first = false
		***REMOVED*** else ***REMOVED***
			buf.WriteString(`,`)
		***REMOVED***
		buf.WriteString(`"stream":`)
		ffjsonWriteJSONBytesAsString(buf, []byte(mj.Stream))
	***REMOVED***
	if len(mj.RawAttrs) > 0 ***REMOVED***
		if first ***REMOVED***
			first = false
		***REMOVED*** else ***REMOVED***
			buf.WriteString(`,`)
		***REMOVED***
		buf.WriteString(`"attrs":`)
		buf.Write(mj.RawAttrs)
	***REMOVED***
	if !first ***REMOVED***
		buf.WriteString(`,`)
	***REMOVED***

	created, err := fastTimeMarshalJSON(mj.Created)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	buf.WriteString(`"time":`)
	buf.WriteString(created)
	buf.WriteString(`***REMOVED***`)
	return nil
***REMOVED***

func ffjsonWriteJSONBytesAsString(buf *bytes.Buffer, s []byte) ***REMOVED***
	const hex = "0123456789abcdef"

	buf.WriteByte('"')
	start := 0
	for i := 0; i < len(s); ***REMOVED***
		if b := s[i]; b < utf8.RuneSelf ***REMOVED***
			if 0x20 <= b && b != '\\' && b != '"' && b != '<' && b != '>' && b != '&' ***REMOVED***
				i++
				continue
			***REMOVED***
			if start < i ***REMOVED***
				buf.Write(s[start:i])
			***REMOVED***
			switch b ***REMOVED***
			case '\\', '"':
				buf.WriteByte('\\')
				buf.WriteByte(b)
			case '\n':
				buf.WriteByte('\\')
				buf.WriteByte('n')
			case '\r':
				buf.WriteByte('\\')
				buf.WriteByte('r')
			default:

				buf.WriteString(`\u00`)
				buf.WriteByte(hex[b>>4])
				buf.WriteByte(hex[b&0xF])
			***REMOVED***
			i++
			start = i
			continue
		***REMOVED***
		c, size := utf8.DecodeRune(s[i:])
		if c == utf8.RuneError && size == 1 ***REMOVED***
			if start < i ***REMOVED***
				buf.Write(s[start:i])
			***REMOVED***
			buf.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		***REMOVED***

		if c == '\u2028' || c == '\u2029' ***REMOVED***
			if start < i ***REMOVED***
				buf.Write(s[start:i])
			***REMOVED***
			buf.WriteString(`\u202`)
			buf.WriteByte(hex[c&0xF])
			i += size
			start = i
			continue
		***REMOVED***
		i += size
	***REMOVED***
	if start < len(s) ***REMOVED***
		buf.Write(s[start:])
	***REMOVED***
	buf.WriteByte('"')
***REMOVED***
