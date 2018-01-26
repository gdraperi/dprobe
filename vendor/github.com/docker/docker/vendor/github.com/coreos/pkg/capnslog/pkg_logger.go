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
	"fmt"
	"os"
)

type PackageLogger struct ***REMOVED***
	pkg   string
	level LogLevel
***REMOVED***

const calldepth = 2

func (p *PackageLogger) internalLog(depth int, inLevel LogLevel, entries ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Lock()
	defer logger.Unlock()
	if inLevel != CRITICAL && p.level < inLevel ***REMOVED***
		return
	***REMOVED***
	if logger.formatter != nil ***REMOVED***
		logger.formatter.Format(p.pkg, inLevel, depth+1, entries...)
	***REMOVED***
***REMOVED***

func (p *PackageLogger) LevelAt(l LogLevel) bool ***REMOVED***
	logger.Lock()
	defer logger.Unlock()
	return p.level >= l
***REMOVED***

// Log a formatted string at any level between ERROR and TRACE
func (p *PackageLogger) Logf(l LogLevel, format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	p.internalLog(calldepth, l, fmt.Sprintf(format, args...))
***REMOVED***

// Log a message at any level between ERROR and TRACE
func (p *PackageLogger) Log(l LogLevel, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	p.internalLog(calldepth, l, fmt.Sprint(args...))
***REMOVED***

// log stdlib compatibility

func (p *PackageLogger) Println(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	p.internalLog(calldepth, INFO, fmt.Sprintln(args...))
***REMOVED***

func (p *PackageLogger) Printf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	p.Logf(INFO, format, args...)
***REMOVED***

func (p *PackageLogger) Print(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	p.internalLog(calldepth, INFO, fmt.Sprint(args...))
***REMOVED***

// Panic and fatal

func (p *PackageLogger) Panicf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	s := fmt.Sprintf(format, args...)
	p.internalLog(calldepth, CRITICAL, s)
	panic(s)
***REMOVED***

func (p *PackageLogger) Panic(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	s := fmt.Sprint(args...)
	p.internalLog(calldepth, CRITICAL, s)
	panic(s)
***REMOVED***

func (p *PackageLogger) Fatalf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	p.Logf(CRITICAL, format, args...)
	os.Exit(1)
***REMOVED***

func (p *PackageLogger) Fatal(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	s := fmt.Sprint(args...)
	p.internalLog(calldepth, CRITICAL, s)
	os.Exit(1)
***REMOVED***

func (p *PackageLogger) Fatalln(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	s := fmt.Sprintln(args...)
	p.internalLog(calldepth, CRITICAL, s)
	os.Exit(1)
***REMOVED***

// Error Functions

func (p *PackageLogger) Errorf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	p.Logf(ERROR, format, args...)
***REMOVED***

func (p *PackageLogger) Error(entries ...interface***REMOVED******REMOVED***) ***REMOVED***
	p.internalLog(calldepth, ERROR, entries...)
***REMOVED***

// Warning Functions

func (p *PackageLogger) Warningf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	p.Logf(WARNING, format, args...)
***REMOVED***

func (p *PackageLogger) Warning(entries ...interface***REMOVED******REMOVED***) ***REMOVED***
	p.internalLog(calldepth, WARNING, entries...)
***REMOVED***

// Notice Functions

func (p *PackageLogger) Noticef(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	p.Logf(NOTICE, format, args...)
***REMOVED***

func (p *PackageLogger) Notice(entries ...interface***REMOVED******REMOVED***) ***REMOVED***
	p.internalLog(calldepth, NOTICE, entries...)
***REMOVED***

// Info Functions

func (p *PackageLogger) Infof(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	p.Logf(INFO, format, args...)
***REMOVED***

func (p *PackageLogger) Info(entries ...interface***REMOVED******REMOVED***) ***REMOVED***
	p.internalLog(calldepth, INFO, entries...)
***REMOVED***

// Debug Functions

func (p *PackageLogger) Debugf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if p.level < DEBUG ***REMOVED***
		return
	***REMOVED***
	p.Logf(DEBUG, format, args...)
***REMOVED***

func (p *PackageLogger) Debug(entries ...interface***REMOVED******REMOVED***) ***REMOVED***
	if p.level < DEBUG ***REMOVED***
		return
	***REMOVED***
	p.internalLog(calldepth, DEBUG, entries...)
***REMOVED***

// Trace Functions

func (p *PackageLogger) Tracef(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if p.level < TRACE ***REMOVED***
		return
	***REMOVED***
	p.Logf(TRACE, format, args...)
***REMOVED***

func (p *PackageLogger) Trace(entries ...interface***REMOVED******REMOVED***) ***REMOVED***
	if p.level < TRACE ***REMOVED***
		return
	***REMOVED***
	p.internalLog(calldepth, TRACE, entries...)
***REMOVED***

func (p *PackageLogger) Flush() ***REMOVED***
	logger.Lock()
	defer logger.Unlock()
	logger.formatter.Flush()
***REMOVED***
