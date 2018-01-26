package jsonfilelog

import (
	"encoding/json"
	"io"

	"github.com/docker/docker/api/types/backend"
	"github.com/docker/docker/daemon/logger"
	"github.com/docker/docker/daemon/logger/jsonfilelog/jsonlog"
)

const maxJSONDecodeRetry = 20000

// ReadLogs implements the logger's LogReader interface for the logs
// created by this driver.
func (l *JSONFileLogger) ReadLogs(config logger.ReadConfig) *logger.LogWatcher ***REMOVED***
	logWatcher := logger.NewLogWatcher()

	go l.readLogs(logWatcher, config)
	return logWatcher
***REMOVED***

func (l *JSONFileLogger) readLogs(watcher *logger.LogWatcher, config logger.ReadConfig) ***REMOVED***
	defer close(watcher.Msg)

	l.mu.Lock()
	l.readers[watcher] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	l.mu.Unlock()

	l.writer.ReadLogs(config, watcher)

	l.mu.Lock()
	delete(l.readers, watcher)
	l.mu.Unlock()
***REMOVED***

func decodeLogLine(dec *json.Decoder, l *jsonlog.JSONLog) (*logger.Message, error) ***REMOVED***
	l.Reset()
	if err := dec.Decode(l); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var attrs []backend.LogAttr
	if len(l.Attrs) != 0 ***REMOVED***
		attrs = make([]backend.LogAttr, 0, len(l.Attrs))
		for k, v := range l.Attrs ***REMOVED***
			attrs = append(attrs, backend.LogAttr***REMOVED***Key: k, Value: v***REMOVED***)
		***REMOVED***
	***REMOVED***
	msg := &logger.Message***REMOVED***
		Source:    l.Stream,
		Timestamp: l.Created,
		Line:      []byte(l.Log),
		Attrs:     attrs,
	***REMOVED***
	return msg, nil
***REMOVED***

// decodeFunc is used to create a decoder for the log file reader
func decodeFunc(rdr io.Reader) func() (*logger.Message, error) ***REMOVED***
	l := &jsonlog.JSONLog***REMOVED******REMOVED***
	dec := json.NewDecoder(rdr)
	return func() (msg *logger.Message, err error) ***REMOVED***
		for retries := 0; retries < maxJSONDecodeRetry; retries++ ***REMOVED***
			msg, err = decodeLogLine(dec, l)
			if err == nil ***REMOVED***
				break
			***REMOVED***

			// try again, could be due to a an incomplete json object as we read
			if _, ok := err.(*json.SyntaxError); ok ***REMOVED***
				dec = json.NewDecoder(rdr)
				retries++
				continue
			***REMOVED***

			// io.ErrUnexpectedEOF is returned from json.Decoder when there is
			// remaining data in the parser's buffer while an io.EOF occurs.
			// If the json logger writes a partial json log entry to the disk
			// while at the same time the decoder tries to decode it, the race condition happens.
			if err == io.ErrUnexpectedEOF ***REMOVED***
				reader := io.MultiReader(dec.Buffered(), rdr)
				dec = json.NewDecoder(reader)
				retries++
			***REMOVED***
		***REMOVED***
		return msg, err
	***REMOVED***
***REMOVED***
