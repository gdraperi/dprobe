// +build !windows,!nacl,!plan9

package syslog

import (
	"fmt"
	"log/syslog"
	"os"

	"github.com/sirupsen/logrus"
)

// SyslogHook to send logs via syslog.
type SyslogHook struct ***REMOVED***
	Writer        *syslog.Writer
	SyslogNetwork string
	SyslogRaddr   string
***REMOVED***

// Creates a hook to be added to an instance of logger. This is called with
// `hook, err := NewSyslogHook("udp", "localhost:514", syslog.LOG_DEBUG, "")`
// `if err == nil ***REMOVED*** log.Hooks.Add(hook) ***REMOVED***`
func NewSyslogHook(network, raddr string, priority syslog.Priority, tag string) (*SyslogHook, error) ***REMOVED***
	w, err := syslog.Dial(network, raddr, priority, tag)
	return &SyslogHook***REMOVED***w, network, raddr***REMOVED***, err
***REMOVED***

func (hook *SyslogHook) Fire(entry *logrus.Entry) error ***REMOVED***
	line, err := entry.String()
	if err != nil ***REMOVED***
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	***REMOVED***

	switch entry.Level ***REMOVED***
	case logrus.PanicLevel:
		return hook.Writer.Crit(line)
	case logrus.FatalLevel:
		return hook.Writer.Crit(line)
	case logrus.ErrorLevel:
		return hook.Writer.Err(line)
	case logrus.WarnLevel:
		return hook.Writer.Warning(line)
	case logrus.InfoLevel:
		return hook.Writer.Info(line)
	case logrus.DebugLevel:
		return hook.Writer.Debug(line)
	default:
		return nil
	***REMOVED***
***REMOVED***

func (hook *SyslogHook) Levels() []logrus.Level ***REMOVED***
	return logrus.AllLevels
***REMOVED***
