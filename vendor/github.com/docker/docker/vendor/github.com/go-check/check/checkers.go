package check

import (
	"fmt"
	"reflect"
	"regexp"
)

// -----------------------------------------------------------------------
// CommentInterface and Commentf helper, to attach extra information to checks.

type comment struct ***REMOVED***
	format string
	args   []interface***REMOVED******REMOVED***
***REMOVED***

// Commentf returns an infomational value to use with Assert or Check calls.
// If the checker test fails, the provided arguments will be passed to
// fmt.Sprintf, and will be presented next to the logged failure.
//
// For example:
//
//     c.Assert(v, Equals, 42, Commentf("Iteration #%d failed.", i))
//
// Note that if the comment is constant, a better option is to
// simply use a normal comment right above or next to the line, as
// it will also get printed with any errors:
//
//     c.Assert(l, Equals, 8192) // Ensure buffer size is correct (bug #123)
//
func Commentf(format string, args ...interface***REMOVED******REMOVED***) CommentInterface ***REMOVED***
	return &comment***REMOVED***format, args***REMOVED***
***REMOVED***

// CommentInterface must be implemented by types that attach extra
// information to failed checks. See the Commentf function for details.
type CommentInterface interface ***REMOVED***
	CheckCommentString() string
***REMOVED***

func (c *comment) CheckCommentString() string ***REMOVED***
	return fmt.Sprintf(c.format, c.args...)
***REMOVED***

// -----------------------------------------------------------------------
// The Checker interface.

// The Checker interface must be provided by checkers used with
// the Assert and Check verification methods.
type Checker interface ***REMOVED***
	Info() *CheckerInfo
	Check(params []interface***REMOVED******REMOVED***, names []string) (result bool, error string)
***REMOVED***

// See the Checker interface.
type CheckerInfo struct ***REMOVED***
	Name   string
	Params []string
***REMOVED***

func (info *CheckerInfo) Info() *CheckerInfo ***REMOVED***
	return info
***REMOVED***

// -----------------------------------------------------------------------
// Not checker logic inverter.

// The Not checker inverts the logic of the provided checker.  The
// resulting checker will succeed where the original one failed, and
// vice-versa.
//
// For example:
//
//     c.Assert(a, Not(Equals), b)
//
func Not(checker Checker) Checker ***REMOVED***
	return &notChecker***REMOVED***checker***REMOVED***
***REMOVED***

type notChecker struct ***REMOVED***
	sub Checker
***REMOVED***

func (checker *notChecker) Info() *CheckerInfo ***REMOVED***
	info := *checker.sub.Info()
	info.Name = "Not(" + info.Name + ")"
	return &info
***REMOVED***

func (checker *notChecker) Check(params []interface***REMOVED******REMOVED***, names []string) (result bool, error string) ***REMOVED***
	result, error = checker.sub.Check(params, names)
	result = !result
	return
***REMOVED***

// -----------------------------------------------------------------------
// IsNil checker.

type isNilChecker struct ***REMOVED***
	*CheckerInfo
***REMOVED***

// The IsNil checker tests whether the obtained value is nil.
//
// For example:
//
//    c.Assert(err, IsNil)
//
var IsNil Checker = &isNilChecker***REMOVED***
	&CheckerInfo***REMOVED***Name: "IsNil", Params: []string***REMOVED***"value"***REMOVED******REMOVED***,
***REMOVED***

func (checker *isNilChecker) Check(params []interface***REMOVED******REMOVED***, names []string) (result bool, error string) ***REMOVED***
	return isNil(params[0]), ""
***REMOVED***

func isNil(obtained interface***REMOVED******REMOVED***) (result bool) ***REMOVED***
	if obtained == nil ***REMOVED***
		result = true
	***REMOVED*** else ***REMOVED***
		switch v := reflect.ValueOf(obtained); v.Kind() ***REMOVED***
		case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
			return v.IsNil()
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// -----------------------------------------------------------------------
// NotNil checker. Alias for Not(IsNil), since it's so common.

type notNilChecker struct ***REMOVED***
	*CheckerInfo
***REMOVED***

// The NotNil checker verifies that the obtained value is not nil.
//
// For example:
//
//     c.Assert(iface, NotNil)
//
// This is an alias for Not(IsNil), made available since it's a
// fairly common check.
//
var NotNil Checker = &notNilChecker***REMOVED***
	&CheckerInfo***REMOVED***Name: "NotNil", Params: []string***REMOVED***"value"***REMOVED******REMOVED***,
***REMOVED***

func (checker *notNilChecker) Check(params []interface***REMOVED******REMOVED***, names []string) (result bool, error string) ***REMOVED***
	return !isNil(params[0]), ""
***REMOVED***

// -----------------------------------------------------------------------
// Equals checker.

type equalsChecker struct ***REMOVED***
	*CheckerInfo
***REMOVED***

// The Equals checker verifies that the obtained value is equal to
// the expected value, according to usual Go semantics for ==.
//
// For example:
//
//     c.Assert(value, Equals, 42)
//
var Equals Checker = &equalsChecker***REMOVED***
	&CheckerInfo***REMOVED***Name: "Equals", Params: []string***REMOVED***"obtained", "expected"***REMOVED******REMOVED***,
***REMOVED***

func (checker *equalsChecker) Check(params []interface***REMOVED******REMOVED***, names []string) (result bool, error string) ***REMOVED***
	defer func() ***REMOVED***
		if v := recover(); v != nil ***REMOVED***
			result = false
			error = fmt.Sprint(v)
		***REMOVED***
	***REMOVED***()
	return params[0] == params[1], ""
***REMOVED***

// -----------------------------------------------------------------------
// DeepEquals checker.

type deepEqualsChecker struct ***REMOVED***
	*CheckerInfo
***REMOVED***

// The DeepEquals checker verifies that the obtained value is deep-equal to
// the expected value.  The check will work correctly even when facing
// slices, interfaces, and values of different types (which always fail
// the test).
//
// For example:
//
//     c.Assert(value, DeepEquals, 42)
//     c.Assert(array, DeepEquals, []string***REMOVED***"hi", "there"***REMOVED***)
//
var DeepEquals Checker = &deepEqualsChecker***REMOVED***
	&CheckerInfo***REMOVED***Name: "DeepEquals", Params: []string***REMOVED***"obtained", "expected"***REMOVED******REMOVED***,
***REMOVED***

func (checker *deepEqualsChecker) Check(params []interface***REMOVED******REMOVED***, names []string) (result bool, error string) ***REMOVED***
	return reflect.DeepEqual(params[0], params[1]), ""
***REMOVED***

// -----------------------------------------------------------------------
// HasLen checker.

type hasLenChecker struct ***REMOVED***
	*CheckerInfo
***REMOVED***

// The HasLen checker verifies that the obtained value has the
// provided length. In many cases this is superior to using Equals
// in conjuction with the len function because in case the check
// fails the value itself will be printed, instead of its length,
// providing more details for figuring the problem.
//
// For example:
//
//     c.Assert(list, HasLen, 5)
//
var HasLen Checker = &hasLenChecker***REMOVED***
	&CheckerInfo***REMOVED***Name: "HasLen", Params: []string***REMOVED***"obtained", "n"***REMOVED******REMOVED***,
***REMOVED***

func (checker *hasLenChecker) Check(params []interface***REMOVED******REMOVED***, names []string) (result bool, error string) ***REMOVED***
	n, ok := params[1].(int)
	if !ok ***REMOVED***
		return false, "n must be an int"
	***REMOVED***
	value := reflect.ValueOf(params[0])
	switch value.Kind() ***REMOVED***
	case reflect.Map, reflect.Array, reflect.Slice, reflect.Chan, reflect.String:
	default:
		return false, "obtained value type has no length"
	***REMOVED***
	return value.Len() == n, ""
***REMOVED***

// -----------------------------------------------------------------------
// ErrorMatches checker.

type errorMatchesChecker struct ***REMOVED***
	*CheckerInfo
***REMOVED***

// The ErrorMatches checker verifies that the error value
// is non nil and matches the regular expression provided.
//
// For example:
//
//     c.Assert(err, ErrorMatches, "perm.*denied")
//
var ErrorMatches Checker = errorMatchesChecker***REMOVED***
	&CheckerInfo***REMOVED***Name: "ErrorMatches", Params: []string***REMOVED***"value", "regex"***REMOVED******REMOVED***,
***REMOVED***

func (checker errorMatchesChecker) Check(params []interface***REMOVED******REMOVED***, names []string) (result bool, errStr string) ***REMOVED***
	if params[0] == nil ***REMOVED***
		return false, "Error value is nil"
	***REMOVED***
	err, ok := params[0].(error)
	if !ok ***REMOVED***
		return false, "Value is not an error"
	***REMOVED***
	params[0] = err.Error()
	names[0] = "error"
	return matches(params[0], params[1])
***REMOVED***

// -----------------------------------------------------------------------
// Matches checker.

type matchesChecker struct ***REMOVED***
	*CheckerInfo
***REMOVED***

// The Matches checker verifies that the string provided as the obtained
// value (or the string resulting from obtained.String()) matches the
// regular expression provided.
//
// For example:
//
//     c.Assert(err, Matches, "perm.*denied")
//
var Matches Checker = &matchesChecker***REMOVED***
	&CheckerInfo***REMOVED***Name: "Matches", Params: []string***REMOVED***"value", "regex"***REMOVED******REMOVED***,
***REMOVED***

func (checker *matchesChecker) Check(params []interface***REMOVED******REMOVED***, names []string) (result bool, error string) ***REMOVED***
	return matches(params[0], params[1])
***REMOVED***

func matches(value, regex interface***REMOVED******REMOVED***) (result bool, error string) ***REMOVED***
	reStr, ok := regex.(string)
	if !ok ***REMOVED***
		return false, "Regex must be a string"
	***REMOVED***
	valueStr, valueIsStr := value.(string)
	if !valueIsStr ***REMOVED***
		if valueWithStr, valueHasStr := value.(fmt.Stringer); valueHasStr ***REMOVED***
			valueStr, valueIsStr = valueWithStr.String(), true
		***REMOVED***
	***REMOVED***
	if valueIsStr ***REMOVED***
		matches, err := regexp.MatchString("^"+reStr+"$", valueStr)
		if err != nil ***REMOVED***
			return false, "Can't compile regex: " + err.Error()
		***REMOVED***
		return matches, ""
	***REMOVED***
	return false, "Obtained value is not a string and has no .String()"
***REMOVED***

// -----------------------------------------------------------------------
// Panics checker.

type panicsChecker struct ***REMOVED***
	*CheckerInfo
***REMOVED***

// The Panics checker verifies that calling the provided zero-argument
// function will cause a panic which is deep-equal to the provided value.
//
// For example:
//
//     c.Assert(func() ***REMOVED*** f(1, 2) ***REMOVED***, Panics, &SomeErrorType***REMOVED***"BOOM"***REMOVED***).
//
//
var Panics Checker = &panicsChecker***REMOVED***
	&CheckerInfo***REMOVED***Name: "Panics", Params: []string***REMOVED***"function", "expected"***REMOVED******REMOVED***,
***REMOVED***

func (checker *panicsChecker) Check(params []interface***REMOVED******REMOVED***, names []string) (result bool, error string) ***REMOVED***
	f := reflect.ValueOf(params[0])
	if f.Kind() != reflect.Func || f.Type().NumIn() != 0 ***REMOVED***
		return false, "Function must take zero arguments"
	***REMOVED***
	defer func() ***REMOVED***
		// If the function has not panicked, then don't do the check.
		if error != "" ***REMOVED***
			return
		***REMOVED***
		params[0] = recover()
		names[0] = "panic"
		result = reflect.DeepEqual(params[0], params[1])
	***REMOVED***()
	f.Call(nil)
	return false, "Function has not panicked"
***REMOVED***

type panicMatchesChecker struct ***REMOVED***
	*CheckerInfo
***REMOVED***

// The PanicMatches checker verifies that calling the provided zero-argument
// function will cause a panic with an error value matching
// the regular expression provided.
//
// For example:
//
//     c.Assert(func() ***REMOVED*** f(1, 2) ***REMOVED***, PanicMatches, `open.*: no such file or directory`).
//
//
var PanicMatches Checker = &panicMatchesChecker***REMOVED***
	&CheckerInfo***REMOVED***Name: "PanicMatches", Params: []string***REMOVED***"function", "expected"***REMOVED******REMOVED***,
***REMOVED***

func (checker *panicMatchesChecker) Check(params []interface***REMOVED******REMOVED***, names []string) (result bool, errmsg string) ***REMOVED***
	f := reflect.ValueOf(params[0])
	if f.Kind() != reflect.Func || f.Type().NumIn() != 0 ***REMOVED***
		return false, "Function must take zero arguments"
	***REMOVED***
	defer func() ***REMOVED***
		// If the function has not panicked, then don't do the check.
		if errmsg != "" ***REMOVED***
			return
		***REMOVED***
		obtained := recover()
		names[0] = "panic"
		if e, ok := obtained.(error); ok ***REMOVED***
			params[0] = e.Error()
		***REMOVED*** else if _, ok := obtained.(string); ok ***REMOVED***
			params[0] = obtained
		***REMOVED*** else ***REMOVED***
			errmsg = "Panic value is not a string or an error"
			return
		***REMOVED***
		result, errmsg = matches(params[0], params[1])
	***REMOVED***()
	f.Call(nil)
	return false, "Function has not panicked"
***REMOVED***

// -----------------------------------------------------------------------
// FitsTypeOf checker.

type fitsTypeChecker struct ***REMOVED***
	*CheckerInfo
***REMOVED***

// The FitsTypeOf checker verifies that the obtained value is
// assignable to a variable with the same type as the provided
// sample value.
//
// For example:
//
//     c.Assert(value, FitsTypeOf, int64(0))
//     c.Assert(value, FitsTypeOf, os.Error(nil))
//
var FitsTypeOf Checker = &fitsTypeChecker***REMOVED***
	&CheckerInfo***REMOVED***Name: "FitsTypeOf", Params: []string***REMOVED***"obtained", "sample"***REMOVED******REMOVED***,
***REMOVED***

func (checker *fitsTypeChecker) Check(params []interface***REMOVED******REMOVED***, names []string) (result bool, error string) ***REMOVED***
	obtained := reflect.ValueOf(params[0])
	sample := reflect.ValueOf(params[1])
	if !obtained.IsValid() ***REMOVED***
		return false, ""
	***REMOVED***
	if !sample.IsValid() ***REMOVED***
		return false, "Invalid sample value"
	***REMOVED***
	return obtained.Type().AssignableTo(sample.Type()), ""
***REMOVED***

// -----------------------------------------------------------------------
// Implements checker.

type implementsChecker struct ***REMOVED***
	*CheckerInfo
***REMOVED***

// The Implements checker verifies that the obtained value
// implements the interface specified via a pointer to an interface
// variable.
//
// For example:
//
//     var e os.Error
//     c.Assert(err, Implements, &e)
//
var Implements Checker = &implementsChecker***REMOVED***
	&CheckerInfo***REMOVED***Name: "Implements", Params: []string***REMOVED***"obtained", "ifaceptr"***REMOVED******REMOVED***,
***REMOVED***

func (checker *implementsChecker) Check(params []interface***REMOVED******REMOVED***, names []string) (result bool, error string) ***REMOVED***
	obtained := reflect.ValueOf(params[0])
	ifaceptr := reflect.ValueOf(params[1])
	if !obtained.IsValid() ***REMOVED***
		return false, ""
	***REMOVED***
	if !ifaceptr.IsValid() || ifaceptr.Kind() != reflect.Ptr || ifaceptr.Elem().Kind() != reflect.Interface ***REMOVED***
		return false, "ifaceptr should be a pointer to an interface variable"
	***REMOVED***
	return obtained.Type().Implements(ifaceptr.Elem().Type()), ""
***REMOVED***
