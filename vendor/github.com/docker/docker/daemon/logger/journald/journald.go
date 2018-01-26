// +build linux

// Package journald provides the log driver for forwarding server logs
// to endpoints that receive the systemd format.
package journald

import (
	"fmt"
	"sync"
	"unicode"

	"github.com/coreos/go-systemd/journal"
	"github.com/docker/docker/daemon/logger"
	"github.com/docker/docker/daemon/logger/loggerutils"
	"github.com/sirupsen/logrus"
)

const name = "journald"

type journald struct ***REMOVED***
	mu      sync.Mutex
	vars    map[string]string // additional variables and values to send to the journal along with the log message
	readers readerList
	closed  bool
***REMOVED***

type readerList struct ***REMOVED***
	readers map[*logger.LogWatcher]*logger.LogWatcher
***REMOVED***

func init() ***REMOVED***
	if err := logger.RegisterLogDriver(name, New); err != nil ***REMOVED***
		logrus.Fatal(err)
	***REMOVED***
	if err := logger.RegisterLogOptValidator(name, validateLogOpt); err != nil ***REMOVED***
		logrus.Fatal(err)
	***REMOVED***
***REMOVED***

// sanitizeKeyMode returns the sanitized string so that it could be used in journald.
// In journald log, there are special requirements for fields.
// Fields must be composed of uppercase letters, numbers, and underscores, but must
// not start with an underscore.
func sanitizeKeyMod(s string) string ***REMOVED***
	n := ""
	for _, v := range s ***REMOVED***
		if 'a' <= v && v <= 'z' ***REMOVED***
			v = unicode.ToUpper(v)
		***REMOVED*** else if ('Z' < v || v < 'A') && ('9' < v || v < '0') ***REMOVED***
			v = '_'
		***REMOVED***
		// If (n == "" && v == '_'), then we will skip as this is the beginning with '_'
		if !(n == "" && v == '_') ***REMOVED***
			n += string(v)
		***REMOVED***
	***REMOVED***
	return n
***REMOVED***

// New creates a journald logger using the configuration passed in on
// the context.
func New(info logger.Info) (logger.Logger, error) ***REMOVED***
	if !journal.Enabled() ***REMOVED***
		return nil, fmt.Errorf("journald is not enabled on this host")
	***REMOVED***

	// parse log tag
	tag, err := loggerutils.ParseLogTag(info, loggerutils.DefaultTemplate)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	vars := map[string]string***REMOVED***
		"CONTAINER_ID":      info.ContainerID[:12],
		"CONTAINER_ID_FULL": info.ContainerID,
		"CONTAINER_NAME":    info.Name(),
		"CONTAINER_TAG":     tag,
		"SYSLOG_IDENTIFIER": tag,
	***REMOVED***
	extraAttrs, err := info.ExtraAttributes(sanitizeKeyMod)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	for k, v := range extraAttrs ***REMOVED***
		vars[k] = v
	***REMOVED***
	return &journald***REMOVED***vars: vars, readers: readerList***REMOVED***readers: make(map[*logger.LogWatcher]*logger.LogWatcher)***REMOVED******REMOVED***, nil
***REMOVED***

// We don't actually accept any options, but we have to supply a callback for
// the factory to pass the (probably empty) configuration map to.
func validateLogOpt(cfg map[string]string) error ***REMOVED***
	for key := range cfg ***REMOVED***
		switch key ***REMOVED***
		case "labels":
		case "env":
		case "env-regex":
		case "tag":
		default:
			return fmt.Errorf("unknown log opt '%s' for journald log driver", key)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (s *journald) Log(msg *logger.Message) error ***REMOVED***
	vars := map[string]string***REMOVED******REMOVED***
	for k, v := range s.vars ***REMOVED***
		vars[k] = v
	***REMOVED***
	if msg.Partial ***REMOVED***
		vars["CONTAINER_PARTIAL_MESSAGE"] = "true"
	***REMOVED***

	line := string(msg.Line)
	source := msg.Source
	logger.PutMessage(msg)

	if source == "stderr" ***REMOVED***
		return journal.Send(line, journal.PriErr, vars)
	***REMOVED***
	return journal.Send(line, journal.PriInfo, vars)
***REMOVED***

func (s *journald) Name() string ***REMOVED***
	return name
***REMOVED***
