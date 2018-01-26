package logger

import (
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types/plugins/logdriver"
	"github.com/docker/docker/pkg/plugingetter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// pluginAdapter takes a plugin and implements the Logger interface for logger
// instances
type pluginAdapter struct ***REMOVED***
	driverName   string
	id           string
	plugin       logPlugin
	basePath     string
	fifoPath     string
	capabilities Capability
	logInfo      Info

	// synchronize access to the log stream and shared buffer
	mu     sync.Mutex
	enc    logdriver.LogEntryEncoder
	stream io.WriteCloser
	// buf is shared for each `Log()` call to reduce allocations.
	// buf must be protected by mutex
	buf logdriver.LogEntry
***REMOVED***

func (a *pluginAdapter) Log(msg *Message) error ***REMOVED***
	a.mu.Lock()

	a.buf.Line = msg.Line
	a.buf.TimeNano = msg.Timestamp.UnixNano()
	a.buf.Partial = msg.Partial
	a.buf.Source = msg.Source

	err := a.enc.Encode(&a.buf)
	a.buf.Reset()

	a.mu.Unlock()

	PutMessage(msg)
	return err
***REMOVED***

func (a *pluginAdapter) Name() string ***REMOVED***
	return a.driverName
***REMOVED***

func (a *pluginAdapter) Close() error ***REMOVED***
	a.mu.Lock()
	defer a.mu.Unlock()

	if err := a.plugin.StopLogging(strings.TrimPrefix(a.fifoPath, a.basePath)); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := a.stream.Close(); err != nil ***REMOVED***
		logrus.WithError(err).Error("error closing plugin fifo")
	***REMOVED***
	if err := os.Remove(a.fifoPath); err != nil && !os.IsNotExist(err) ***REMOVED***
		logrus.WithError(err).Error("error cleaning up plugin fifo")
	***REMOVED***

	// may be nil, especially for unit tests
	if pluginGetter != nil ***REMOVED***
		pluginGetter.Get(a.Name(), extName, plugingetter.Release)
	***REMOVED***
	return nil
***REMOVED***

type pluginAdapterWithRead struct ***REMOVED***
	*pluginAdapter
***REMOVED***

func (a *pluginAdapterWithRead) ReadLogs(config ReadConfig) *LogWatcher ***REMOVED***
	watcher := NewLogWatcher()

	go func() ***REMOVED***
		defer close(watcher.Msg)
		stream, err := a.plugin.ReadLogs(a.logInfo, config)
		if err != nil ***REMOVED***
			watcher.Err <- errors.Wrap(err, "error getting log reader")
			return
		***REMOVED***
		defer stream.Close()

		dec := logdriver.NewLogEntryDecoder(stream)
		for ***REMOVED***
			select ***REMOVED***
			case <-watcher.WatchClose():
				return
			default:
			***REMOVED***

			var buf logdriver.LogEntry
			if err := dec.Decode(&buf); err != nil ***REMOVED***
				if err == io.EOF ***REMOVED***
					return
				***REMOVED***
				select ***REMOVED***
				case watcher.Err <- errors.Wrap(err, "error decoding log message"):
				case <-watcher.WatchClose():
				***REMOVED***
				return
			***REMOVED***

			msg := &Message***REMOVED***
				Timestamp: time.Unix(0, buf.TimeNano),
				Line:      buf.Line,
				Source:    buf.Source,
			***REMOVED***

			// plugin should handle this, but check just in case
			if !config.Since.IsZero() && msg.Timestamp.Before(config.Since) ***REMOVED***
				continue
			***REMOVED***
			if !config.Until.IsZero() && msg.Timestamp.After(config.Until) ***REMOVED***
				return
			***REMOVED***

			select ***REMOVED***
			case watcher.Msg <- msg:
			case <-watcher.WatchClose():
				// make sure the message we consumed is sent
				watcher.Msg <- msg
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return watcher
***REMOVED***
