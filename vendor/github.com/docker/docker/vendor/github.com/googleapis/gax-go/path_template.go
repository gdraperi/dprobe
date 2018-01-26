// Copyright 2016, Google Inc.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package gax

import (
	"errors"
	"fmt"
	"strings"
)

type matcher interface ***REMOVED***
	match([]string) (int, error)
	String() string
***REMOVED***

type segment struct ***REMOVED***
	matcher
	name string
***REMOVED***

type labelMatcher string

func (ls labelMatcher) match(segments []string) (int, error) ***REMOVED***
	if len(segments) == 0 ***REMOVED***
		return 0, fmt.Errorf("expected %s but no more segments found", ls)
	***REMOVED***
	if segments[0] != string(ls) ***REMOVED***
		return 0, fmt.Errorf("expected %s but got %s", ls, segments[0])
	***REMOVED***
	return 1, nil
***REMOVED***

func (ls labelMatcher) String() string ***REMOVED***
	return string(ls)
***REMOVED***

type wildcardMatcher int

func (wm wildcardMatcher) match(segments []string) (int, error) ***REMOVED***
	if len(segments) == 0 ***REMOVED***
		return 0, errors.New("no more segments found")
	***REMOVED***
	return 1, nil
***REMOVED***

func (wm wildcardMatcher) String() string ***REMOVED***
	return "*"
***REMOVED***

type pathWildcardMatcher int

func (pwm pathWildcardMatcher) match(segments []string) (int, error) ***REMOVED***
	length := len(segments) - int(pwm)
	if length <= 0 ***REMOVED***
		return 0, errors.New("not sufficient segments are supplied for path wildcard")
	***REMOVED***
	return length, nil
***REMOVED***

func (pwm pathWildcardMatcher) String() string ***REMOVED***
	return "**"
***REMOVED***

type ParseError struct ***REMOVED***
	Pos      int
	Template string
	Message  string
***REMOVED***

func (pe ParseError) Error() string ***REMOVED***
	return fmt.Sprintf("at %d of template '%s', %s", pe.Pos, pe.Template, pe.Message)
***REMOVED***

// PathTemplate manages the template to build and match with paths used
// by API services. It holds a template and variable names in it, and
// it can extract matched patterns from a path string or build a path
// string from a binding.
//
// See http.proto in github.com/googleapis/googleapis/ for the details of
// the template syntax.
type PathTemplate struct ***REMOVED***
	segments []segment
***REMOVED***

// NewPathTemplate parses a path template, and returns a PathTemplate
// instance if successful.
func NewPathTemplate(template string) (*PathTemplate, error) ***REMOVED***
	return parsePathTemplate(template)
***REMOVED***

// MustCompilePathTemplate is like NewPathTemplate but panics if the
// expression cannot be parsed. It simplifies safe initialization of
// global variables holding compiled regular expressions.
func MustCompilePathTemplate(template string) *PathTemplate ***REMOVED***
	pt, err := NewPathTemplate(template)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return pt
***REMOVED***

// Match attempts to match the given path with the template, and returns
// the mapping of the variable name to the matched pattern string.
func (pt *PathTemplate) Match(path string) (map[string]string, error) ***REMOVED***
	paths := strings.Split(path, "/")
	values := map[string]string***REMOVED******REMOVED***
	for _, segment := range pt.segments ***REMOVED***
		length, err := segment.match(paths)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if segment.name != "" ***REMOVED***
			value := strings.Join(paths[:length], "/")
			if oldValue, ok := values[segment.name]; ok ***REMOVED***
				values[segment.name] = oldValue + "/" + value
			***REMOVED*** else ***REMOVED***
				values[segment.name] = value
			***REMOVED***
		***REMOVED***
		paths = paths[length:]
	***REMOVED***
	if len(paths) != 0 ***REMOVED***
		return nil, fmt.Errorf("Trailing path %s remains after the matching", strings.Join(paths, "/"))
	***REMOVED***
	return values, nil
***REMOVED***

// Render creates a path string from its template and the binding from
// the variable name to the value.
func (pt *PathTemplate) Render(binding map[string]string) (string, error) ***REMOVED***
	result := make([]string, 0, len(pt.segments))
	var lastVariableName string
	for _, segment := range pt.segments ***REMOVED***
		name := segment.name
		if lastVariableName != "" && name == lastVariableName ***REMOVED***
			continue
		***REMOVED***
		lastVariableName = name
		if name == "" ***REMOVED***
			result = append(result, segment.String())
		***REMOVED*** else if value, ok := binding[name]; ok ***REMOVED***
			result = append(result, value)
		***REMOVED*** else ***REMOVED***
			return "", fmt.Errorf("%s is not found", name)
		***REMOVED***
	***REMOVED***
	built := strings.Join(result, "/")
	return built, nil
***REMOVED***
