package syslog

import (
	"log/syslog"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLocalhostAddAndPrint(t *testing.T) ***REMOVED***
	log := logrus.New()
	hook, err := NewSyslogHook("udp", "localhost:514", syslog.LOG_INFO, "")

	if err != nil ***REMOVED***
		t.Errorf("Unable to connect to local syslog.")
	***REMOVED***

	log.Hooks.Add(hook)

	for _, level := range hook.Levels() ***REMOVED***
		if len(log.Hooks[level]) != 1 ***REMOVED***
			t.Errorf("SyslogHook was not added. The length of log.Hooks[%v]: %v", level, len(log.Hooks[level]))
		***REMOVED***
	***REMOVED***

	log.Info("Congratulations!")
***REMOVED***
