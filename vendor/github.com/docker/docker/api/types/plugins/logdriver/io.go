package logdriver

import (
	"encoding/binary"
	"io"
)

const binaryEncodeLen = 4

// LogEntryEncoder encodes a LogEntry to a protobuf stream
// The stream should look like:
//
// [uint32 binary encoded message size][protobuf message]
//
// To decode an entry, read the first 4 bytes to get the size of the entry,
// then read `size` bytes from the stream.
type LogEntryEncoder interface ***REMOVED***
	Encode(*LogEntry) error
***REMOVED***

// NewLogEntryEncoder creates a protobuf stream encoder for log entries.
// This is used to write out  log entries to a stream.
func NewLogEntryEncoder(w io.Writer) LogEntryEncoder ***REMOVED***
	return &logEntryEncoder***REMOVED***
		w:   w,
		buf: make([]byte, 1024),
	***REMOVED***
***REMOVED***

type logEntryEncoder struct ***REMOVED***
	buf []byte
	w   io.Writer
***REMOVED***

func (e *logEntryEncoder) Encode(l *LogEntry) error ***REMOVED***
	n := l.Size()

	total := n + binaryEncodeLen
	if total > len(e.buf) ***REMOVED***
		e.buf = make([]byte, total)
	***REMOVED***
	binary.BigEndian.PutUint32(e.buf, uint32(n))

	if _, err := l.MarshalTo(e.buf[binaryEncodeLen:]); err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err := e.w.Write(e.buf[:total])
	return err
***REMOVED***

// LogEntryDecoder decodes log entries from a stream
// It is expected that the wire format is as defined by LogEntryEncoder.
type LogEntryDecoder interface ***REMOVED***
	Decode(*LogEntry) error
***REMOVED***

// NewLogEntryDecoder creates a new stream decoder for log entries
func NewLogEntryDecoder(r io.Reader) LogEntryDecoder ***REMOVED***
	return &logEntryDecoder***REMOVED***
		lenBuf: make([]byte, binaryEncodeLen),
		buf:    make([]byte, 1024),
		r:      r,
	***REMOVED***
***REMOVED***

type logEntryDecoder struct ***REMOVED***
	r      io.Reader
	lenBuf []byte
	buf    []byte
***REMOVED***

func (d *logEntryDecoder) Decode(l *LogEntry) error ***REMOVED***
	_, err := io.ReadFull(d.r, d.lenBuf)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	size := int(binary.BigEndian.Uint32(d.lenBuf))
	if len(d.buf) < size ***REMOVED***
		d.buf = make([]byte, size)
	***REMOVED***

	if _, err := io.ReadFull(d.r, d.buf[:size]); err != nil ***REMOVED***
		return err
	***REMOVED***
	return l.Unmarshal(d.buf[:size])
***REMOVED***
