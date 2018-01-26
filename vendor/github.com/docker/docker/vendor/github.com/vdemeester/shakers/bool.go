package shakers

import (
	"github.com/go-check/check"
)

// True checker verifies the obtained value is true
//
//    c.Assert(myBool, True)
//
var True check.Checker = &boolChecker***REMOVED***
	&check.CheckerInfo***REMOVED***
		Name:   "True",
		Params: []string***REMOVED***"obtained"***REMOVED***,
	***REMOVED***,
	true,
***REMOVED***

// False checker verifies the obtained value is false
//
//    c.Assert(myBool, False)
//
var False check.Checker = &boolChecker***REMOVED***
	&check.CheckerInfo***REMOVED***
		Name:   "False",
		Params: []string***REMOVED***"obtained"***REMOVED***,
	***REMOVED***,
	false,
***REMOVED***

type boolChecker struct ***REMOVED***
	*check.CheckerInfo
	expected bool
***REMOVED***

func (checker *boolChecker) Check(params []interface***REMOVED******REMOVED***, names []string) (bool, string) ***REMOVED***
	return is(checker.expected, params[0])
***REMOVED***

func is(expected bool, obtained interface***REMOVED******REMOVED***) (bool, string) ***REMOVED***
	obtainedBool, ok := obtained.(bool)
	if !ok ***REMOVED***
		return false, "obtained value must be a bool."
	***REMOVED***
	return obtainedBool == expected, ""
***REMOVED***
