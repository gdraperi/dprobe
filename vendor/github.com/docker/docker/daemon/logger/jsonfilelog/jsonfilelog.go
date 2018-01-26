// Package jsonfilelog provides the default Logger implementation for
// Docker logging. This logger logs to files on the host server in the
// JSON format.
package jsonfilelog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"github.com/docker/docker/daemon/logger"
	"github.com/docker/docker/daemon/logger/jsonfilelog/jsonlog"
	"github.com/docker/docker/daemon/logger/loggerutils"
	units "github.com/docker/go-units"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Name is the name of the file that the jsonlogger logs to.
const Name = "json-file"

// JSONFileLogger is Logger implementation for default Docker logging.
type JSONFileLogger struct ***REMOVED***
	mu      sync.Mutex
	closed  bool
	writer  *loggerutils.LogFile
	readers map[*logger.LogWatcher]struct***REMOVED******REMOVED*** // stores the active log followers
	tag     string                          // tag values requested by the user to log
***REMOVED***

func init() ***REMOVED***
	if err := logger.RegisterLogDriver(Name, New); err != nil ***REMOVED***
		logrus.Fatal(err)
	***REMOVED***
	if err := logger.RegisterLogOptValidator(Name, ValidateLogOpt); err != nil ***REMOVED***
		logrus.Fatal(err)
	***REMOVED***
***REMOVED***

// New creates new JSONFileLogger which writes to filename passed in
// on given context.
func New(info logger.Info) (logger.Logger, error) ***REMOVED***
	var capval int64 = -1
	if capacity, ok := info.Config["max-size"]; ok ***REMOVED***
		var err error
		capval, err = units.FromHumanSize(capacity)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	var maxFiles = 1
	if maxFileString, ok := info.Config["max-file"]; ok ***REMOVED***
		var err error
		maxFiles, err = strconv.Atoi(maxFileString)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if maxFiles < 1 ***REMOVED***
			return nil, fmt.Errorf("max-file cannot be less than 1")
		***REMOVED***
	***REMOVED***

	attrs, err := info.ExtraAttributes(nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// no default template. only use a tag if the user asked for it
	tag, err := loggerutils.ParseLogTag(info, "")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if tag != "" ***REMOVED***
		attrs["tag"] = tag
	***REMOVED***

	var extra []byte
	if len(attrs) > 0 ***REMOVED***
		var err error
		extra, err = json.Marshal(attrs)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	buf := bytes.NewBuffer(nil)
	marshalFunc := func(msg *logger.Message) ([]byte, error) ***REMOVED***
		if err := marshalMessage(msg, extra, buf); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		b := buf.Bytes()
		buf.Reset()
		return b, nil
	***REMOVED***

	writer, err := loggerutils.NewLogFile(info.LogPath, capval, maxFiles, marshalFunc, decodeFunc)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &JSONFileLogger***REMOVED***
		writer:  writer,
		readers: make(map[*logger.LogWatcher]struct***REMOVED******REMOVED***),
		tag:     tag,
	***REMOVED***, nil
***REMOVED***

// Log converts logger.Message to jsonlog.JSONLog and serializes it to file.
func (l *JSONFileLogger) Log(msg *logger.Message) error ***REMOVED***
	l.mu.Lock()
	err := l.writer.WriteLogEntry(msg)
	l.mu.Unlock()
	return err
***REMOVED***

func marshalMessage(msg *logger.Message, extra json.RawMessage, buf *bytes.Buffer) error ***REMOVED***
	logLine := msg.Line
	if !msg.Partial ***REMOVED***
		logLine = append(msg.Line, '\n')
	***REMOVED***
	err := (&jsonlog.JSONLogs***REMOVED***
		Log:      logLine,
		Stream:   msg.Source,
		Created:  msg.Timestamp,
		RawAttrs: extra,
	***REMOVED***).MarshalJSONBuf(buf)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "error writing log message to buffer")
	***REMOVED***
	err = buf.WriteByte('\n')
	return errors.Wrap(err, "error finalizing log buffer")
***REMOVED***

// ValidateLogOpt looks for json specific log options max-file & max-size.
func ValidateLogOpt(cfg map[string]string) error ***REMOVED***
	for key := range cfg ***REMOVED***
		switch key ***REMOVED***
		case "max-file":
		case "max-size":
		case "labels":
		case "env":
		case "env-regex":
		case "tag":
		default:
			return fmt.Errorf("unknown log opt '%s' for json-file log driver", key)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// LogPath returns the location the given json logger logs to.
func (l *JSONFileLogger) LogPath() string ***REMOVED***
	return l.writer.LogPath()
***REMOVED***

// Close closes underlying file and signals all readers to stop.
func (l *JSONFileLogger) Close() error ***REMOVED***
	l.mu.Lock()
	l.closed = true
	err := l.writer.Close()
	for r := range l.readers ***REMOVED***
		r.Close()
		delete(l.readers, r)
	***REMOVED***
	l.mu.Unlock()
	return err
***REMOVED***

// Name returns name of this logger.
func (l *JSONFileLogger) Name() string ***REMOVED***
	return Name
***REMOVED***
