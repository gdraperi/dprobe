package shakers

import (
	"reflect"
	"time"

	"github.com/go-check/check"
)

// As a commodity, we bring all check.Checker variables into the current namespace to avoid having
// to think about check.X versus checker.X.
var (
	DeepEquals   = check.DeepEquals
	ErrorMatches = check.ErrorMatches
	FitsTypeOf   = check.FitsTypeOf
	HasLen       = check.HasLen
	Implements   = check.Implements
	IsNil        = check.IsNil
	Matches      = check.Matches
	Not          = check.Not
	NotNil       = check.NotNil
	PanicMatches = check.PanicMatches
	Panics       = check.Panics
)

// Equaler is an interface implemented if the type has a Equal method.
// This is used to compare struct using shakers.Equals.
type Equaler interface ***REMOVED***
	Equal(Equaler) bool
***REMOVED***

// Equals checker verifies the obtained value is equal to the specified one.
// It's is smart in a wait that it supports several *types* (built-in, Equaler,
// time.Time)
//
//    c.Assert(myStruct, Equals, aStruct, check.Commentf("bouuuhh"))
//    c.Assert(myTime, Equals, aTime, check.Commentf("bouuuhh"))
//
var Equals check.Checker = &equalChecker***REMOVED***
	&check.CheckerInfo***REMOVED***
		Name:   "Equals",
		Params: []string***REMOVED***"obtained", "expected"***REMOVED***,
	***REMOVED***,
***REMOVED***

type equalChecker struct ***REMOVED***
	*check.CheckerInfo
***REMOVED***

func (checker *equalChecker) Check(params []interface***REMOVED******REMOVED***, names []string) (bool, string) ***REMOVED***
	return isEqual(params[0], params[1])
***REMOVED***

func isEqual(obtained, expected interface***REMOVED******REMOVED***) (bool, string) ***REMOVED***
	switch obtained.(type) ***REMOVED***
	case time.Time:
		return timeEquals(obtained, expected)
	case Equaler:
		return equalerEquals(obtained, expected)
	default:
		if reflect.TypeOf(obtained) != reflect.TypeOf(expected) ***REMOVED***
			return false, "obtained value and expected value have not the same type."
		***REMOVED***
		return obtained == expected, ""
	***REMOVED***
***REMOVED***

func equalerEquals(obtained, expected interface***REMOVED******REMOVED***) (bool, string) ***REMOVED***
	expectedEqualer, ok := expected.(Equaler)
	if !ok ***REMOVED***
		return false, "expected value must be an Equaler - implementing Equal(Equaler)."
	***REMOVED***
	obtainedEqualer, ok := obtained.(Equaler)
	if !ok ***REMOVED***
		return false, "obtained value must be an Equaler - implementing Equal(Equaler)."
	***REMOVED***
	return obtainedEqualer.Equal(expectedEqualer), ""
***REMOVED***

// GreaterThan checker verifies the obtained value is greater than the specified one.
// It's is smart in a wait that it supports several *types* (built-in, time.Time)
//
//    c.Assert(myTime, GreaterThan, aTime, check.Commentf("bouuuhh"))
//    c.Assert(myInt, GreaterThan, 2, check.Commentf("bouuuhh"))
//
var GreaterThan check.Checker = &greaterThanChecker***REMOVED***
	&check.CheckerInfo***REMOVED***
		Name:   "GreaterThan",
		Params: []string***REMOVED***"obtained", "expected"***REMOVED***,
	***REMOVED***,
***REMOVED***

type greaterThanChecker struct ***REMOVED***
	*check.CheckerInfo
***REMOVED***

func (checker *greaterThanChecker) Check(params []interface***REMOVED******REMOVED***, names []string) (bool, string) ***REMOVED***
	return greaterThan(params[0], params[1])
***REMOVED***

func greaterThan(obtained, expected interface***REMOVED******REMOVED***) (bool, string) ***REMOVED***
	if _, ok := obtained.(time.Time); ok ***REMOVED***
		return isAfter(obtained, expected)
	***REMOVED***
	if reflect.TypeOf(obtained) != reflect.TypeOf(expected) ***REMOVED***
		return false, "obtained value and expected value have not the same type."
	***REMOVED***
	switch v := obtained.(type) ***REMOVED***
	case float32:
		return v > expected.(float32), ""
	case float64:
		return v > expected.(float64), ""
	case int:
		return v > expected.(int), ""
	case int8:
		return v > expected.(int8), ""
	case int16:
		return v > expected.(int16), ""
	case int32:
		return v > expected.(int32), ""
	case int64:
		return v > expected.(int64), ""
	case uint:
		return v > expected.(uint), ""
	case uint8:
		return v > expected.(uint8), ""
	case uint16:
		return v > expected.(uint16), ""
	case uint32:
		return v > expected.(uint32), ""
	case uint64:
		return v > expected.(uint64), ""
	default:
		return false, "obtained value type not supported."
	***REMOVED***
***REMOVED***

// GreaterOrEqualThan checker verifies the obtained value is greater or equal than the specified one.
// It's is smart in a wait that it supports several *types* (built-in, time.Time)
//
//    c.Assert(myTime, GreaterOrEqualThan, aTime, check.Commentf("bouuuhh"))
//    c.Assert(myInt, GreaterOrEqualThan, 2, check.Commentf("bouuuhh"))
//
var GreaterOrEqualThan check.Checker = &greaterOrEqualThanChecker***REMOVED***
	&check.CheckerInfo***REMOVED***
		Name:   "GreaterOrEqualThan",
		Params: []string***REMOVED***"obtained", "expected"***REMOVED***,
	***REMOVED***,
***REMOVED***

type greaterOrEqualThanChecker struct ***REMOVED***
	*check.CheckerInfo
***REMOVED***

func (checker *greaterOrEqualThanChecker) Check(params []interface***REMOVED******REMOVED***, names []string) (bool, string) ***REMOVED***
	return greaterOrEqualThan(params[0], params[1])
***REMOVED***

func greaterOrEqualThan(obtained, expected interface***REMOVED******REMOVED***) (bool, string) ***REMOVED***
	if _, ok := obtained.(time.Time); ok ***REMOVED***
		return isAfter(obtained, expected)
	***REMOVED***
	if reflect.TypeOf(obtained) != reflect.TypeOf(expected) ***REMOVED***
		return false, "obtained value and expected value have not the same type."
	***REMOVED***
	switch v := obtained.(type) ***REMOVED***
	case float32:
		return v >= expected.(float32), ""
	case float64:
		return v >= expected.(float64), ""
	case int:
		return v >= expected.(int), ""
	case int8:
		return v >= expected.(int8), ""
	case int16:
		return v >= expected.(int16), ""
	case int32:
		return v >= expected.(int32), ""
	case int64:
		return v >= expected.(int64), ""
	case uint:
		return v >= expected.(uint), ""
	case uint8:
		return v >= expected.(uint8), ""
	case uint16:
		return v >= expected.(uint16), ""
	case uint32:
		return v >= expected.(uint32), ""
	case uint64:
		return v >= expected.(uint64), ""
	default:
		return false, "obtained value type not supported."
	***REMOVED***
***REMOVED***

// LessThan checker verifies the obtained value is less than the specified one.
// It's is smart in a wait that it supports several *types* (built-in, time.Time)
//
//    c.Assert(myTime, LessThan, aTime, check.Commentf("bouuuhh"))
//    c.Assert(myInt, LessThan, 2, check.Commentf("bouuuhh"))
//
var LessThan check.Checker = &lessThanChecker***REMOVED***
	&check.CheckerInfo***REMOVED***
		Name:   "LessThan",
		Params: []string***REMOVED***"obtained", "expected"***REMOVED***,
	***REMOVED***,
***REMOVED***

type lessThanChecker struct ***REMOVED***
	*check.CheckerInfo
***REMOVED***

func (checker *lessThanChecker) Check(params []interface***REMOVED******REMOVED***, names []string) (bool, string) ***REMOVED***
	return lessThan(params[0], params[1])
***REMOVED***

func lessThan(obtained, expected interface***REMOVED******REMOVED***) (bool, string) ***REMOVED***
	if _, ok := obtained.(time.Time); ok ***REMOVED***
		return isBefore(obtained, expected)
	***REMOVED***
	if reflect.TypeOf(obtained) != reflect.TypeOf(expected) ***REMOVED***
		return false, "obtained value and expected value have not the same type."
	***REMOVED***
	switch v := obtained.(type) ***REMOVED***
	case float32:
		return v < expected.(float32), ""
	case float64:
		return v < expected.(float64), ""
	case int:
		return v < expected.(int), ""
	case int8:
		return v < expected.(int8), ""
	case int16:
		return v < expected.(int16), ""
	case int32:
		return v < expected.(int32), ""
	case int64:
		return v < expected.(int64), ""
	case uint:
		return v < expected.(uint), ""
	case uint8:
		return v < expected.(uint8), ""
	case uint16:
		return v < expected.(uint16), ""
	case uint32:
		return v < expected.(uint32), ""
	case uint64:
		return v < expected.(uint64), ""
	default:
		return false, "obtained value type not supported."
	***REMOVED***
***REMOVED***

// LessOrEqualThan checker verifies the obtained value is less or equal than the specified one.
// It's is smart in a wait that it supports several *types* (built-in, time.Time)
//
//    c.Assert(myTime, LessThan, aTime, check.Commentf("bouuuhh"))
//    c.Assert(myInt, LessThan, 2, check.Commentf("bouuuhh"))
//
var LessOrEqualThan check.Checker = &lessOrEqualThanChecker***REMOVED***
	&check.CheckerInfo***REMOVED***
		Name:   "LessOrEqualThan",
		Params: []string***REMOVED***"obtained", "expected"***REMOVED***,
	***REMOVED***,
***REMOVED***

type lessOrEqualThanChecker struct ***REMOVED***
	*check.CheckerInfo
***REMOVED***

func (checker *lessOrEqualThanChecker) Check(params []interface***REMOVED******REMOVED***, names []string) (bool, string) ***REMOVED***
	return lessOrEqualThan(params[0], params[1])
***REMOVED***

func lessOrEqualThan(obtained, expected interface***REMOVED******REMOVED***) (bool, string) ***REMOVED***
	if _, ok := obtained.(time.Time); ok ***REMOVED***
		return isBefore(obtained, expected)
	***REMOVED***
	if reflect.TypeOf(obtained) != reflect.TypeOf(expected) ***REMOVED***
		return false, "obtained value and expected value have not the same type."
	***REMOVED***
	switch v := obtained.(type) ***REMOVED***
	case float32:
		return v <= expected.(float32), ""
	case float64:
		return v <= expected.(float64), ""
	case int:
		return v <= expected.(int), ""
	case int8:
		return v <= expected.(int8), ""
	case int16:
		return v <= expected.(int16), ""
	case int32:
		return v <= expected.(int32), ""
	case int64:
		return v <= expected.(int64), ""
	case uint:
		return v <= expected.(uint), ""
	case uint8:
		return v <= expected.(uint8), ""
	case uint16:
		return v <= expected.(uint16), ""
	case uint32:
		return v <= expected.(uint32), ""
	case uint64:
		return v <= expected.(uint64), ""
	default:
		return false, "obtained value type not supported."
	***REMOVED***
***REMOVED***
