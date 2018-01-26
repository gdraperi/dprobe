// Package shakers provide some checker implementation the go-check.Checker interface.
package shakers

import (
	"fmt"
	"strings"

	"github.com/go-check/check"
)

// Contains checker verifies that obtained value contains a substring.
var Contains check.Checker = &substringChecker***REMOVED***
	&check.CheckerInfo***REMOVED***
		Name:   "Contains",
		Params: []string***REMOVED***"obtained", "substring"***REMOVED***,
	***REMOVED***,
	strings.Contains,
***REMOVED***

// ContainsAny checker verifies that any Unicode code points in chars
// are in the obtained string.
var ContainsAny check.Checker = &substringChecker***REMOVED***
	&check.CheckerInfo***REMOVED***
		Name:   "ContainsAny",
		Params: []string***REMOVED***"obtained", "chars"***REMOVED***,
	***REMOVED***,
	strings.ContainsAny,
***REMOVED***

// HasPrefix checker verifies that obtained value has the specified substring as prefix
var HasPrefix check.Checker = &substringChecker***REMOVED***
	&check.CheckerInfo***REMOVED***
		Name:   "HasPrefix",
		Params: []string***REMOVED***"obtained", "prefix"***REMOVED***,
	***REMOVED***,
	strings.HasPrefix,
***REMOVED***

// HasSuffix checker verifies that obtained value has the specified substring as prefix
var HasSuffix check.Checker = &substringChecker***REMOVED***
	&check.CheckerInfo***REMOVED***
		Name:   "HasSuffix",
		Params: []string***REMOVED***"obtained", "suffix"***REMOVED***,
	***REMOVED***,
	strings.HasSuffix,
***REMOVED***

// EqualFold checker verifies that obtained value is, interpreted as UTF-8 strings, are equal under Unicode case-folding.
var EqualFold check.Checker = &substringChecker***REMOVED***
	&check.CheckerInfo***REMOVED***
		Name:   "EqualFold",
		Params: []string***REMOVED***"obtained", "expected"***REMOVED***,
	***REMOVED***,
	strings.EqualFold,
***REMOVED***

type substringChecker struct ***REMOVED***
	*check.CheckerInfo
	substringFunction func(string, string) bool
***REMOVED***

func (checker *substringChecker) Check(params []interface***REMOVED******REMOVED***, names []string) (bool, string) ***REMOVED***
	obtained := params[0]
	substring := params[1]
	substringStr, ok := substring.(string)
	if !ok ***REMOVED***
		return false, fmt.Sprintf("%s value must be a string.", names[1])
	***REMOVED***
	obtainedString, obtainedIsStr := obtained.(string)
	if !obtainedIsStr ***REMOVED***
		if obtainedWithStringer, obtainedHasStringer := obtained.(fmt.Stringer); obtainedHasStringer ***REMOVED***
			obtainedString, obtainedIsStr = obtainedWithStringer.String(), true
		***REMOVED***
	***REMOVED***
	if obtainedIsStr ***REMOVED***
		return checker.substringFunction(obtainedString, substringStr), ""
	***REMOVED***
	return false, "obtained value is not a string and has no .String()."
***REMOVED***

// IndexAny checker verifies that the index of the first instance of any Unicode code point from chars in the obtained value is equal to expected
var IndexAny check.Checker = &substringCountChecker***REMOVED***
	&check.CheckerInfo***REMOVED***
		Name:   "IndexAny",
		Params: []string***REMOVED***"obtained", "chars", "expected"***REMOVED***,
	***REMOVED***,
	strings.IndexAny,
***REMOVED***

// Index checker verifies that the index of the first instance of sep in the obtained value is equal to expected
var Index check.Checker = &substringCountChecker***REMOVED***
	&check.CheckerInfo***REMOVED***
		Name:   "Index",
		Params: []string***REMOVED***"obtained", "sep", "expected"***REMOVED***,
	***REMOVED***,
	strings.Index,
***REMOVED***

// Count checker verifies that obtained value has the specified number of non-overlapping instances of sep
var Count check.Checker = &substringCountChecker***REMOVED***
	&check.CheckerInfo***REMOVED***
		Name:   "Count",
		Params: []string***REMOVED***"obtained", "sep", "expected"***REMOVED***,
	***REMOVED***,
	strings.Count,
***REMOVED***

type substringCountChecker struct ***REMOVED***
	*check.CheckerInfo
	substringFunction func(string, string) int
***REMOVED***

func (checker *substringCountChecker) Check(params []interface***REMOVED******REMOVED***, names []string) (bool, string) ***REMOVED***
	obtained := params[0]
	substring := params[1]
	expected := params[2]
	substringStr, ok := substring.(string)
	if !ok ***REMOVED***
		return false, fmt.Sprintf("%s value must be a string.", names[1])
	***REMOVED***
	obtainedString, obtainedIsStr := obtained.(string)
	if !obtainedIsStr ***REMOVED***
		if obtainedWithStringer, obtainedHasStringer := obtained.(fmt.Stringer); obtainedHasStringer ***REMOVED***
			obtainedString, obtainedIsStr = obtainedWithStringer.String(), true
		***REMOVED***
	***REMOVED***
	if obtainedIsStr ***REMOVED***
		return checker.substringFunction(obtainedString, substringStr) == expected, ""
	***REMOVED***
	return false, "obtained value is not a string and has no .String()."
***REMOVED***

// IsLower checker verifies that the obtained value is in lower case
var IsLower check.Checker = &stringTransformChecker***REMOVED***
	&check.CheckerInfo***REMOVED***
		Name:   "IsLower",
		Params: []string***REMOVED***"obtained"***REMOVED***,
	***REMOVED***,
	strings.ToLower,
***REMOVED***

// IsUpper checker verifies that the obtained value is in lower case
var IsUpper check.Checker = &stringTransformChecker***REMOVED***
	&check.CheckerInfo***REMOVED***
		Name:   "IsUpper",
		Params: []string***REMOVED***"obtained"***REMOVED***,
	***REMOVED***,
	strings.ToUpper,
***REMOVED***

type stringTransformChecker struct ***REMOVED***
	*check.CheckerInfo
	stringFunction func(string) string
***REMOVED***

func (checker *stringTransformChecker) Check(params []interface***REMOVED******REMOVED***, names []string) (bool, string) ***REMOVED***
	obtained := params[0]
	obtainedString, obtainedIsStr := obtained.(string)
	if !obtainedIsStr ***REMOVED***
		if obtainedWithStringer, obtainedHasStringer := obtained.(fmt.Stringer); obtainedHasStringer ***REMOVED***
			obtainedString, obtainedIsStr = obtainedWithStringer.String(), true
		***REMOVED***
	***REMOVED***
	if obtainedIsStr ***REMOVED***
		return checker.stringFunction(obtainedString) == obtainedString, ""
	***REMOVED***
	return false, "obtained value is not a string and has no .String()."
***REMOVED***
