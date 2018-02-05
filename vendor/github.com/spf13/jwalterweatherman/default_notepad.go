// Copyright © 2016 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package jwalterweatherman

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	TRACE    *log.Logger
	DEBUG    *log.Logger
	INFO     *log.Logger
	WARN     *log.Logger
	ERROR    *log.Logger
	CRITICAL *log.Logger
	FATAL    *log.Logger

	LOG      *log.Logger
	FEEDBACK *Feedback

	defaultNotepad *Notepad
)

func reloadDefaultNotepad() ***REMOVED***
	TRACE = defaultNotepad.TRACE
	DEBUG = defaultNotepad.DEBUG
	INFO = defaultNotepad.INFO
	WARN = defaultNotepad.WARN
	ERROR = defaultNotepad.ERROR
	CRITICAL = defaultNotepad.CRITICAL
	FATAL = defaultNotepad.FATAL

	LOG = defaultNotepad.LOG
	FEEDBACK = defaultNotepad.FEEDBACK
***REMOVED***

func init() ***REMOVED***
	defaultNotepad = NewNotepad(LevelError, LevelWarn, os.Stdout, ioutil.Discard, "", log.Ldate|log.Ltime)
	reloadDefaultNotepad()
***REMOVED***

// SetLogThreshold set the log threshold for the default notepad. Trace by default.
func SetLogThreshold(threshold Threshold) ***REMOVED***
	defaultNotepad.SetLogThreshold(threshold)
	reloadDefaultNotepad()
***REMOVED***

// SetLogOutput set the log output for the default notepad. Discarded by default.
func SetLogOutput(handle io.Writer) ***REMOVED***
	defaultNotepad.SetLogOutput(handle)
	reloadDefaultNotepad()
***REMOVED***

// SetStdoutThreshold set the standard output threshold for the default notepad.
// Info by default.
func SetStdoutThreshold(threshold Threshold) ***REMOVED***
	defaultNotepad.SetStdoutThreshold(threshold)
	reloadDefaultNotepad()
***REMOVED***

// SetPrefix set the prefix for the default logger. Empty by default.
func SetPrefix(prefix string) ***REMOVED***
	defaultNotepad.SetPrefix(prefix)
	reloadDefaultNotepad()
***REMOVED***

// SetFlags set the flags for the default logger. "log.Ldate | log.Ltime" by default.
func SetFlags(flags int) ***REMOVED***
	defaultNotepad.SetFlags(flags)
	reloadDefaultNotepad()
***REMOVED***

// Level returns the current global log threshold.
func LogThreshold() Threshold ***REMOVED***
	return defaultNotepad.logThreshold
***REMOVED***

// Level returns the current global output threshold.
func StdoutThreshold() Threshold ***REMOVED***
	return defaultNotepad.stdoutThreshold
***REMOVED***

// GetStdoutThreshold returns the defined Treshold for the log logger.
func GetLogThreshold() Threshold ***REMOVED***
	return defaultNotepad.GetLogThreshold()
***REMOVED***

// GetStdoutThreshold returns the Treshold for the stdout logger.
func GetStdoutThreshold() Threshold ***REMOVED***
	return defaultNotepad.GetStdoutThreshold()
***REMOVED***

// LogCountForLevel returns the number of log invocations for a given threshold.
func LogCountForLevel(l Threshold) uint64 ***REMOVED***
	return defaultNotepad.LogCountForLevel(l)
***REMOVED***

// LogCountForLevelsGreaterThanorEqualTo returns the number of log invocations
// greater than or equal to a given threshold.
func LogCountForLevelsGreaterThanorEqualTo(threshold Threshold) uint64 ***REMOVED***
	return defaultNotepad.LogCountForLevelsGreaterThanorEqualTo(threshold)
***REMOVED***

// ResetLogCounters resets the invocation counters for all levels.
func ResetLogCounters() ***REMOVED***
	defaultNotepad.ResetLogCounters()
***REMOVED***
