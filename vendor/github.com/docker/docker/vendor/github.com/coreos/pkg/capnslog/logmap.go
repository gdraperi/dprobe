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
	"errors"
	"strings"
	"sync"
)

// LogLevel is the set of all log levels.
type LogLevel int8

const (
	// CRITICAL is the lowest log level; only errors which will end the program will be propagated.
	CRITICAL LogLevel = iota - 1
	// ERROR is for errors that are not fatal but lead to troubling behavior.
	ERROR
	// WARNING is for errors which are not fatal and not errors, but are unusual. Often sourced from misconfigurations.
	WARNING
	// NOTICE is for normal but significant conditions.
	NOTICE
	// INFO is a log level for common, everyday log updates.
	INFO
	// DEBUG is the default hidden level for more verbose updates about internal processes.
	DEBUG
	// TRACE is for (potentially) call by call tracing of programs.
	TRACE
)

// Char returns a single-character representation of the log level.
func (l LogLevel) Char() string ***REMOVED***
	switch l ***REMOVED***
	case CRITICAL:
		return "C"
	case ERROR:
		return "E"
	case WARNING:
		return "W"
	case NOTICE:
		return "N"
	case INFO:
		return "I"
	case DEBUG:
		return "D"
	case TRACE:
		return "T"
	default:
		panic("Unhandled loglevel")
	***REMOVED***
***REMOVED***

// String returns a multi-character representation of the log level.
func (l LogLevel) String() string ***REMOVED***
	switch l ***REMOVED***
	case CRITICAL:
		return "CRITICAL"
	case ERROR:
		return "ERROR"
	case WARNING:
		return "WARNING"
	case NOTICE:
		return "NOTICE"
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	case TRACE:
		return "TRACE"
	default:
		panic("Unhandled loglevel")
	***REMOVED***
***REMOVED***

// Update using the given string value. Fulfills the flag.Value interface.
func (l *LogLevel) Set(s string) error ***REMOVED***
	value, err := ParseLevel(s)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	*l = value
	return nil
***REMOVED***

// ParseLevel translates some potential loglevel strings into their corresponding levels.
func ParseLevel(s string) (LogLevel, error) ***REMOVED***
	switch s ***REMOVED***
	case "CRITICAL", "C":
		return CRITICAL, nil
	case "ERROR", "0", "E":
		return ERROR, nil
	case "WARNING", "1", "W":
		return WARNING, nil
	case "NOTICE", "2", "N":
		return NOTICE, nil
	case "INFO", "3", "I":
		return INFO, nil
	case "DEBUG", "4", "D":
		return DEBUG, nil
	case "TRACE", "5", "T":
		return TRACE, nil
	***REMOVED***
	return CRITICAL, errors.New("couldn't parse log level " + s)
***REMOVED***

type RepoLogger map[string]*PackageLogger

type loggerStruct struct ***REMOVED***
	sync.Mutex
	repoMap   map[string]RepoLogger
	formatter Formatter
***REMOVED***

// logger is the global logger
var logger = new(loggerStruct)

// SetGlobalLogLevel sets the log level for all packages in all repositories
// registered with capnslog.
func SetGlobalLogLevel(l LogLevel) ***REMOVED***
	logger.Lock()
	defer logger.Unlock()
	for _, r := range logger.repoMap ***REMOVED***
		r.setRepoLogLevelInternal(l)
	***REMOVED***
***REMOVED***

// GetRepoLogger may return the handle to the repository's set of packages' loggers.
func GetRepoLogger(repo string) (RepoLogger, error) ***REMOVED***
	logger.Lock()
	defer logger.Unlock()
	r, ok := logger.repoMap[repo]
	if !ok ***REMOVED***
		return nil, errors.New("no packages registered for repo " + repo)
	***REMOVED***
	return r, nil
***REMOVED***

// MustRepoLogger returns the handle to the repository's packages' loggers.
func MustRepoLogger(repo string) RepoLogger ***REMOVED***
	r, err := GetRepoLogger(repo)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return r
***REMOVED***

// SetRepoLogLevel sets the log level for all packages in the repository.
func (r RepoLogger) SetRepoLogLevel(l LogLevel) ***REMOVED***
	logger.Lock()
	defer logger.Unlock()
	r.setRepoLogLevelInternal(l)
***REMOVED***

func (r RepoLogger) setRepoLogLevelInternal(l LogLevel) ***REMOVED***
	for _, v := range r ***REMOVED***
		v.level = l
	***REMOVED***
***REMOVED***

// ParseLogLevelConfig parses a comma-separated string of "package=loglevel", in
// order, and returns a map of the results, for use in SetLogLevel.
func (r RepoLogger) ParseLogLevelConfig(conf string) (map[string]LogLevel, error) ***REMOVED***
	setlist := strings.Split(conf, ",")
	out := make(map[string]LogLevel)
	for _, setstring := range setlist ***REMOVED***
		setting := strings.Split(setstring, "=")
		if len(setting) != 2 ***REMOVED***
			return nil, errors.New("oddly structured `pkg=level` option: " + setstring)
		***REMOVED***
		l, err := ParseLevel(setting[1])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		out[setting[0]] = l
	***REMOVED***
	return out, nil
***REMOVED***

// SetLogLevel takes a map of package names within a repository to their desired
// loglevel, and sets the levels appropriately. Unknown packages are ignored.
// "*" is a special package name that corresponds to all packages, and will be
// processed first.
func (r RepoLogger) SetLogLevel(m map[string]LogLevel) ***REMOVED***
	logger.Lock()
	defer logger.Unlock()
	if l, ok := m["*"]; ok ***REMOVED***
		r.setRepoLogLevelInternal(l)
	***REMOVED***
	for k, v := range m ***REMOVED***
		l, ok := r[k]
		if !ok ***REMOVED***
			continue
		***REMOVED***
		l.level = v
	***REMOVED***
***REMOVED***

// SetFormatter sets the formatting function for all logs.
func SetFormatter(f Formatter) ***REMOVED***
	logger.Lock()
	defer logger.Unlock()
	logger.formatter = f
***REMOVED***

// NewPackageLogger creates a package logger object.
// This should be defined as a global var in your package, referencing your repo.
func NewPackageLogger(repo string, pkg string) (p *PackageLogger) ***REMOVED***
	logger.Lock()
	defer logger.Unlock()
	if logger.repoMap == nil ***REMOVED***
		logger.repoMap = make(map[string]RepoLogger)
	***REMOVED***
	r, rok := logger.repoMap[repo]
	if !rok ***REMOVED***
		logger.repoMap[repo] = make(RepoLogger)
		r = logger.repoMap[repo]
	***REMOVED***
	p, pok := r[pkg]
	if !pok ***REMOVED***
		r[pkg] = &PackageLogger***REMOVED***
			pkg:   pkg,
			level: INFO,
		***REMOVED***
		p = r[pkg]
	***REMOVED***
	return
***REMOVED***
