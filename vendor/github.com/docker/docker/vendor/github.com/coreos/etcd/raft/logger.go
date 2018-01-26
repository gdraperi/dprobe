// Copyright 2015 The etcd Authors
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

package raft

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Logger interface ***REMOVED***
	Debug(v ...interface***REMOVED******REMOVED***)
	Debugf(format string, v ...interface***REMOVED******REMOVED***)

	Error(v ...interface***REMOVED******REMOVED***)
	Errorf(format string, v ...interface***REMOVED******REMOVED***)

	Info(v ...interface***REMOVED******REMOVED***)
	Infof(format string, v ...interface***REMOVED******REMOVED***)

	Warning(v ...interface***REMOVED******REMOVED***)
	Warningf(format string, v ...interface***REMOVED******REMOVED***)

	Fatal(v ...interface***REMOVED******REMOVED***)
	Fatalf(format string, v ...interface***REMOVED******REMOVED***)

	Panic(v ...interface***REMOVED******REMOVED***)
	Panicf(format string, v ...interface***REMOVED******REMOVED***)
***REMOVED***

func SetLogger(l Logger) ***REMOVED*** raftLogger = l ***REMOVED***

var (
	defaultLogger = &DefaultLogger***REMOVED***Logger: log.New(os.Stderr, "raft", log.LstdFlags)***REMOVED***
	discardLogger = &DefaultLogger***REMOVED***Logger: log.New(ioutil.Discard, "", 0)***REMOVED***
	raftLogger    = Logger(defaultLogger)
)

const (
	calldepth = 2
)

// DefaultLogger is a default implementation of the Logger interface.
type DefaultLogger struct ***REMOVED***
	*log.Logger
	debug bool
***REMOVED***

func (l *DefaultLogger) EnableTimestamps() ***REMOVED***
	l.SetFlags(l.Flags() | log.Ldate | log.Ltime)
***REMOVED***

func (l *DefaultLogger) EnableDebug() ***REMOVED***
	l.debug = true
***REMOVED***

func (l *DefaultLogger) Debug(v ...interface***REMOVED******REMOVED***) ***REMOVED***
	if l.debug ***REMOVED***
		l.Output(calldepth, header("DEBUG", fmt.Sprint(v...)))
	***REMOVED***
***REMOVED***

func (l *DefaultLogger) Debugf(format string, v ...interface***REMOVED******REMOVED***) ***REMOVED***
	if l.debug ***REMOVED***
		l.Output(calldepth, header("DEBUG", fmt.Sprintf(format, v...)))
	***REMOVED***
***REMOVED***

func (l *DefaultLogger) Info(v ...interface***REMOVED******REMOVED***) ***REMOVED***
	l.Output(calldepth, header("INFO", fmt.Sprint(v...)))
***REMOVED***

func (l *DefaultLogger) Infof(format string, v ...interface***REMOVED******REMOVED***) ***REMOVED***
	l.Output(calldepth, header("INFO", fmt.Sprintf(format, v...)))
***REMOVED***

func (l *DefaultLogger) Error(v ...interface***REMOVED******REMOVED***) ***REMOVED***
	l.Output(calldepth, header("ERROR", fmt.Sprint(v...)))
***REMOVED***

func (l *DefaultLogger) Errorf(format string, v ...interface***REMOVED******REMOVED***) ***REMOVED***
	l.Output(calldepth, header("ERROR", fmt.Sprintf(format, v...)))
***REMOVED***

func (l *DefaultLogger) Warning(v ...interface***REMOVED******REMOVED***) ***REMOVED***
	l.Output(calldepth, header("WARN", fmt.Sprint(v...)))
***REMOVED***

func (l *DefaultLogger) Warningf(format string, v ...interface***REMOVED******REMOVED***) ***REMOVED***
	l.Output(calldepth, header("WARN", fmt.Sprintf(format, v...)))
***REMOVED***

func (l *DefaultLogger) Fatal(v ...interface***REMOVED******REMOVED***) ***REMOVED***
	l.Output(calldepth, header("FATAL", fmt.Sprint(v...)))
	os.Exit(1)
***REMOVED***

func (l *DefaultLogger) Fatalf(format string, v ...interface***REMOVED******REMOVED***) ***REMOVED***
	l.Output(calldepth, header("FATAL", fmt.Sprintf(format, v...)))
	os.Exit(1)
***REMOVED***

func (l *DefaultLogger) Panic(v ...interface***REMOVED******REMOVED***) ***REMOVED***
	l.Logger.Panic(v)
***REMOVED***

func (l *DefaultLogger) Panicf(format string, v ...interface***REMOVED******REMOVED***) ***REMOVED***
	l.Logger.Panicf(format, v...)
***REMOVED***

func header(lvl, msg string) string ***REMOVED***
	return fmt.Sprintf("%s: %s", lvl, msg)
***REMOVED***
