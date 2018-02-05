// Copyright © 2016 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package jwalterweatherman

import (
	"fmt"
	"io"
	"log"
)

type Threshold int

func (t Threshold) String() string ***REMOVED***
	return prefixes[t]
***REMOVED***

const (
	LevelTrace Threshold = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelCritical
	LevelFatal
)

var prefixes map[Threshold]string = map[Threshold]string***REMOVED***
	LevelTrace:    "TRACE",
	LevelDebug:    "DEBUG",
	LevelInfo:     "INFO",
	LevelWarn:     "WARN",
	LevelError:    "ERROR",
	LevelCritical: "CRITICAL",
	LevelFatal:    "FATAL",
***REMOVED***

// Notepad is where you leave a note!
type Notepad struct ***REMOVED***
	TRACE    *log.Logger
	DEBUG    *log.Logger
	INFO     *log.Logger
	WARN     *log.Logger
	ERROR    *log.Logger
	CRITICAL *log.Logger
	FATAL    *log.Logger

	LOG      *log.Logger
	FEEDBACK *Feedback

	loggers         [7]**log.Logger
	logHandle       io.Writer
	outHandle       io.Writer
	logThreshold    Threshold
	stdoutThreshold Threshold
	prefix          string
	flags           int

	// One per Threshold
	logCounters [7]*logCounter
***REMOVED***

// NewNotepad create a new notepad.
func NewNotepad(outThreshold Threshold, logThreshold Threshold, outHandle, logHandle io.Writer, prefix string, flags int) *Notepad ***REMOVED***
	n := &Notepad***REMOVED******REMOVED***

	n.loggers = [7]**log.Logger***REMOVED***&n.TRACE, &n.DEBUG, &n.INFO, &n.WARN, &n.ERROR, &n.CRITICAL, &n.FATAL***REMOVED***
	n.outHandle = outHandle
	n.logHandle = logHandle
	n.stdoutThreshold = outThreshold
	n.logThreshold = logThreshold

	if len(prefix) != 0 ***REMOVED***
		n.prefix = "[" + prefix + "] "
	***REMOVED*** else ***REMOVED***
		n.prefix = ""
	***REMOVED***

	n.flags = flags

	n.LOG = log.New(n.logHandle,
		"LOG:   ",
		n.flags)
	n.FEEDBACK = &Feedback***REMOVED***out: log.New(outHandle, "", 0), log: n.LOG***REMOVED***

	n.init()
	return n
***REMOVED***

// init creates the loggers for each level depending on the notepad thresholds.
func (n *Notepad) init() ***REMOVED***
	logAndOut := io.MultiWriter(n.outHandle, n.logHandle)

	for t, logger := range n.loggers ***REMOVED***
		threshold := Threshold(t)
		counter := &logCounter***REMOVED******REMOVED***
		n.logCounters[t] = counter
		prefix := n.prefix + threshold.String() + " "

		switch ***REMOVED***
		case threshold >= n.logThreshold && threshold >= n.stdoutThreshold:
			*logger = log.New(io.MultiWriter(counter, logAndOut), prefix, n.flags)

		case threshold >= n.logThreshold:
			*logger = log.New(io.MultiWriter(counter, n.logHandle), prefix, n.flags)

		case threshold >= n.stdoutThreshold:
			*logger = log.New(io.MultiWriter(counter, n.outHandle), prefix, n.flags)

		default:
			// counter doesn't care about prefix and flags, so don't use them
			// for performance.
			*logger = log.New(counter, "", 0)
		***REMOVED***
	***REMOVED***
***REMOVED***

// SetLogThreshold changes the threshold above which messages are written to the
// log file.
func (n *Notepad) SetLogThreshold(threshold Threshold) ***REMOVED***
	n.logThreshold = threshold
	n.init()
***REMOVED***

// SetLogOutput changes the file where log messages are written.
func (n *Notepad) SetLogOutput(handle io.Writer) ***REMOVED***
	n.logHandle = handle
	n.init()
***REMOVED***

// GetStdoutThreshold returns the defined Treshold for the log logger.
func (n *Notepad) GetLogThreshold() Threshold ***REMOVED***
	return n.logThreshold
***REMOVED***

// SetStdoutThreshold changes the threshold above which messages are written to the
// standard output.
func (n *Notepad) SetStdoutThreshold(threshold Threshold) ***REMOVED***
	n.stdoutThreshold = threshold
	n.init()
***REMOVED***

// GetStdoutThreshold returns the Treshold for the stdout logger.
func (n *Notepad) GetStdoutThreshold() Threshold ***REMOVED***
	return n.stdoutThreshold
***REMOVED***

// SetPrefix changes the prefix used by the notepad. Prefixes are displayed between
// brackets at the beginning of the line. An empty prefix won't be displayed at all.
func (n *Notepad) SetPrefix(prefix string) ***REMOVED***
	if len(prefix) != 0 ***REMOVED***
		n.prefix = "[" + prefix + "] "
	***REMOVED*** else ***REMOVED***
		n.prefix = ""
	***REMOVED***
	n.init()
***REMOVED***

// SetFlags choose which flags the logger will display (after prefix and message
// level). See the package log for more informations on this.
func (n *Notepad) SetFlags(flags int) ***REMOVED***
	n.flags = flags
	n.init()
***REMOVED***

// Feedback writes plainly to the outHandle while
// logging with the standard extra information (date, file, etc).
type Feedback struct ***REMOVED***
	out *log.Logger
	log *log.Logger
***REMOVED***

func (fb *Feedback) Println(v ...interface***REMOVED******REMOVED***) ***REMOVED***
	fb.output(fmt.Sprintln(v...))
***REMOVED***

func (fb *Feedback) Printf(format string, v ...interface***REMOVED******REMOVED***) ***REMOVED***
	fb.output(fmt.Sprintf(format, v...))
***REMOVED***

func (fb *Feedback) Print(v ...interface***REMOVED******REMOVED***) ***REMOVED***
	fb.output(fmt.Sprint(v...))
***REMOVED***

func (fb *Feedback) output(s string) ***REMOVED***
	if fb.out != nil ***REMOVED***
		fb.out.Output(2, s)
	***REMOVED***
	if fb.log != nil ***REMOVED***
		fb.log.Output(2, s)
	***REMOVED***
***REMOVED***
