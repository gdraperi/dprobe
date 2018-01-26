// Package logentries provides the log driver for forwarding server logs
// to logentries endpoints.
package logentries

import (
	"fmt"
	"strconv"

	"github.com/bsphere/le_go"
	"github.com/docker/docker/daemon/logger"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type logentries struct ***REMOVED***
	tag           string
	containerID   string
	containerName string
	writer        *le_go.Logger
	extra         map[string]string
	lineOnly      bool
***REMOVED***

const (
	name     = "logentries"
	token    = "logentries-token"
	lineonly = "line-only"
)

func init() ***REMOVED***
	if err := logger.RegisterLogDriver(name, New); err != nil ***REMOVED***
		logrus.Fatal(err)
	***REMOVED***
	if err := logger.RegisterLogOptValidator(name, ValidateLogOpt); err != nil ***REMOVED***
		logrus.Fatal(err)
	***REMOVED***
***REMOVED***

// New creates a logentries logger using the configuration passed in on
// the context. The supported context configuration variable is
// logentries-token.
func New(info logger.Info) (logger.Logger, error) ***REMOVED***
	logrus.WithField("container", info.ContainerID).
		WithField("token", info.Config[token]).
		WithField("line-only", info.Config[lineonly]).
		Debug("logging driver logentries configured")

	log, err := le_go.Connect(info.Config[token])
	if err != nil ***REMOVED***
		return nil, errors.Wrap(err, "error connecting to logentries")
	***REMOVED***
	var lineOnly bool
	if info.Config[lineonly] != "" ***REMOVED***
		if lineOnly, err = strconv.ParseBool(info.Config[lineonly]); err != nil ***REMOVED***
			return nil, errors.Wrap(err, "error parsing lineonly option")
		***REMOVED***
	***REMOVED***
	return &logentries***REMOVED***
		containerID:   info.ContainerID,
		containerName: info.ContainerName,
		writer:        log,
		lineOnly:      lineOnly,
	***REMOVED***, nil
***REMOVED***

func (f *logentries) Log(msg *logger.Message) error ***REMOVED***
	if !f.lineOnly ***REMOVED***
		data := map[string]string***REMOVED***
			"container_id":   f.containerID,
			"container_name": f.containerName,
			"source":         msg.Source,
			"log":            string(msg.Line),
		***REMOVED***
		for k, v := range f.extra ***REMOVED***
			data[k] = v
		***REMOVED***
		ts := msg.Timestamp
		logger.PutMessage(msg)
		f.writer.Println(f.tag, ts, data)
	***REMOVED*** else ***REMOVED***
		line := string(msg.Line)
		logger.PutMessage(msg)
		f.writer.Println(line)
	***REMOVED***
	return nil
***REMOVED***

func (f *logentries) Close() error ***REMOVED***
	return f.writer.Close()
***REMOVED***

func (f *logentries) Name() string ***REMOVED***
	return name
***REMOVED***

// ValidateLogOpt looks for logentries specific log option logentries-address.
func ValidateLogOpt(cfg map[string]string) error ***REMOVED***
	for key := range cfg ***REMOVED***
		switch key ***REMOVED***
		case "env":
		case "env-regex":
		case "labels":
		case "tag":
		case key:
		default:
			return fmt.Errorf("unknown log opt '%s' for logentries log driver", key)
		***REMOVED***
	***REMOVED***

	if cfg[token] == "" ***REMOVED***
		return fmt.Errorf("Missing logentries token")
	***REMOVED***

	return nil
***REMOVED***
