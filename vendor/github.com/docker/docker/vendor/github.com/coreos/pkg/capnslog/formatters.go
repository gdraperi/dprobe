// Copyright 2015 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package capnslog

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"runtime"
	"strings"
	"time"
)

type Formatter interface ***REMOVED***
	Format(pkg string, level LogLevel, depth int, entries ...interface***REMOVED******REMOVED***)
	Flush()
***REMOVED***

func NewStringFormatter(w io.Writer) Formatter ***REMOVED***
	return &StringFormatter***REMOVED***
		w: bufio.NewWriter(w),
	***REMOVED***
***REMOVED***

type StringFormatter struct ***REMOVED***
	w *bufio.Writer
***REMOVED***

func (s *StringFormatter) Format(pkg string, l LogLevel, i int, entries ...interface***REMOVED******REMOVED***) ***REMOVED***
	now := time.Now().UTC()
	s.w.WriteString(now.Format(time.RFC3339))
	s.w.WriteByte(' ')
	writeEntries(s.w, pkg, l, i, entries...)
	s.Flush()
***REMOVED***

func writeEntries(w *bufio.Writer, pkg string, _ LogLevel, _ int, entries ...interface***REMOVED******REMOVED***) ***REMOVED***
	if pkg != "" ***REMOVED***
		w.WriteString(pkg + ": ")
	***REMOVED***
	str := fmt.Sprint(entries...)
	endsInNL := strings.HasSuffix(str, "\n")
	w.WriteString(str)
	if !endsInNL ***REMOVED***
		w.WriteString("\n")
	***REMOVED***
***REMOVED***

func (s *StringFormatter) Flush() ***REMOVED***
	s.w.Flush()
***REMOVED***

func NewPrettyFormatter(w io.Writer, debug bool) Formatter ***REMOVED***
	return &PrettyFormatter***REMOVED***
		w:     bufio.NewWriter(w),
		debug: debug,
	***REMOVED***
***REMOVED***

type PrettyFormatter struct ***REMOVED***
	w     *bufio.Writer
	debug bool
***REMOVED***

func (c *PrettyFormatter) Format(pkg string, l LogLevel, depth int, entries ...interface***REMOVED******REMOVED***) ***REMOVED***
	now := time.Now()
	ts := now.Format("2006-01-02 15:04:05")
	c.w.WriteString(ts)
	ms := now.Nanosecond() / 1000
	c.w.WriteString(fmt.Sprintf(".%06d", ms))
	if c.debug ***REMOVED***
		_, file, line, ok := runtime.Caller(depth) // It's always the same number of frames to the user's call.
		if !ok ***REMOVED***
			file = "???"
			line = 1
		***REMOVED*** else ***REMOVED***
			slash := strings.LastIndex(file, "/")
			if slash >= 0 ***REMOVED***
				file = file[slash+1:]
			***REMOVED***
		***REMOVED***
		if line < 0 ***REMOVED***
			line = 0 // not a real line number
		***REMOVED***
		c.w.WriteString(fmt.Sprintf(" [%s:%d]", file, line))
	***REMOVED***
	c.w.WriteString(fmt.Sprint(" ", l.Char(), " | "))
	writeEntries(c.w, pkg, l, depth, entries...)
	c.Flush()
***REMOVED***

func (c *PrettyFormatter) Flush() ***REMOVED***
	c.w.Flush()
***REMOVED***

// LogFormatter emulates the form of the traditional built-in logger.
type LogFormatter struct ***REMOVED***
	logger *log.Logger
	prefix string
***REMOVED***

// NewLogFormatter is a helper to produce a new LogFormatter struct. It uses the
// golang log package to actually do the logging work so that logs look similar.
func NewLogFormatter(w io.Writer, prefix string, flag int) Formatter ***REMOVED***
	return &LogFormatter***REMOVED***
		logger: log.New(w, "", flag), // don't use prefix here
		prefix: prefix,               // save it instead
	***REMOVED***
***REMOVED***

// Format builds a log message for the LogFormatter. The LogLevel is ignored.
func (lf *LogFormatter) Format(pkg string, _ LogLevel, _ int, entries ...interface***REMOVED******REMOVED***) ***REMOVED***
	str := fmt.Sprint(entries...)
	prefix := lf.prefix
	if pkg != "" ***REMOVED***
		prefix = fmt.Sprintf("%s%s: ", prefix, pkg)
	***REMOVED***
	lf.logger.Output(5, fmt.Sprintf("%s%v", prefix, str)) // call depth is 5
***REMOVED***

// Flush is included so that the interface is complete, but is a no-op.
func (lf *LogFormatter) Flush() ***REMOVED***
	// noop
***REMOVED***

// NilFormatter is a no-op log formatter that does nothing.
type NilFormatter struct ***REMOVED***
***REMOVED***

// NewNilFormatter is a helper to produce a new LogFormatter struct. It logs no
// messages so that you can cause part of your logging to be silent.
func NewNilFormatter() Formatter ***REMOVED***
	return &NilFormatter***REMOVED******REMOVED***
***REMOVED***

// Format does nothing.
func (_ *NilFormatter) Format(_ string, _ LogLevel, _ int, _ ...interface***REMOVED******REMOVED***) ***REMOVED***
	// noop
***REMOVED***

// Flush is included so that the interface is complete, but is a no-op.
func (_ *NilFormatter) Flush() ***REMOVED***
	// noop
***REMOVED***
