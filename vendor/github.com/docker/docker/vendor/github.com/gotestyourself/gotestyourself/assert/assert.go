/*Package assert provides assertions and checks for comparing expected values to
actual values. When an assertion or check fails a helpful error message is
printed.

Assert and Check

Assert() and Check() both accept a Comparison, and fail the test when the
comparison fails. The one difference is that Assert() will end the test execution
immediately (using t.FailNow()) whereas Check() will fail the test (using t.Fail()),
return the value of the comparison, then proceed with the rest of the test case.

Example Usage

The example below shows assert used with some common types.


	import (
	    "testing"

	    "github.com/gotestyourself/gotestyourself/assert"
	    is "github.com/gotestyourself/gotestyourself/assert/cmp"
	)

	func TestEverything(t *testing.T) ***REMOVED***
	    // booleans
	    assert.Assert(t, isOk)
	    assert.Assert(t, !missing)

	    // primitives
	    assert.Equal(t, count, 1)
	    assert.Equal(t, msg, "the message")
	    assert.Assert(t, total != 10) // NotEqual

	    // errors
	    assert.NilError(t, closer.Close())
	    assert.Assert(t, is.Error(err, "the exact error message"))
	    assert.Assert(t, is.ErrorContains(err, "includes this"))

	    // complex types
	    assert.Assert(t, is.Len(items, 3))
	    assert.Assert(t, len(sequence) != 0) // NotEmpty
	    assert.Assert(t, is.Contains(mapping, "key"))
	    assert.Assert(t, is.Compare(result, myStruct***REMOVED***Name: "title"***REMOVED***))

	    // pointers and interface
	    assert.Assert(t, is.Nil(ref))
	    assert.Assert(t, ref != nil) // NotNil
	***REMOVED***

Comparisons

https://godoc.org/github.com/gotestyourself/gotestyourself/assert/cmp provides
many common comparisons. Additional comparisons can be written to compare
values in other ways.

Below is an example of a custom comparison using a regex pattern:

	func RegexP(value string, pattern string) func() (bool, string) ***REMOVED***
	    return func() (bool, string) ***REMOVED***
	        re := regexp.MustCompile(pattern)
	        msg := fmt.Sprintf("%q did not match pattern %q", value, pattern)
	        return re.MatchString(value), msg
	***REMOVED***
	***REMOVED***

*/
package assert

import (
	"fmt"

	"github.com/gotestyourself/gotestyourself/assert/cmp"
	"github.com/gotestyourself/gotestyourself/internal/format"
	"github.com/gotestyourself/gotestyourself/internal/source"
)

// BoolOrComparison can be a bool, or Comparison. Other types will panic.
type BoolOrComparison interface***REMOVED******REMOVED***

// Comparison is a function which compares values and returns true if the actual
// value matches the expected value. If the values do not match it returns a message
// with details about why it failed.
//
// https://godoc.org/github.com/gotestyourself/gotestyourself/assert/cmp
// provides many general purpose Comparisons.
type Comparison func() (success bool, message string)

// TestingT is the subset of testing.T used by the assert package.
type TestingT interface ***REMOVED***
	FailNow()
	Fail()
	Log(args ...interface***REMOVED******REMOVED***)
***REMOVED***

type helperT interface ***REMOVED***
	Helper()
***REMOVED***

// stackIndex = Assert()/Check(), assert()
const stackIndex = 2
const comparisonArgPos = 1

const failureMessage = "assertion failed: "

func assert(
	t TestingT,
	failer func(),
	comparison BoolOrComparison,
	msgAndArgs ...interface***REMOVED******REMOVED***,
) bool ***REMOVED***
	if ht, ok := t.(helperT); ok ***REMOVED***
		ht.Helper()
	***REMOVED***
	switch check := comparison.(type) ***REMOVED***
	case bool:
		if check ***REMOVED***
			return true
		***REMOVED***
		source, err := source.GetCondition(stackIndex, comparisonArgPos)
		if err != nil ***REMOVED***
			t.Log(err.Error())
		***REMOVED***

		msg := " is false"
		t.Log(format.WithCustomMessage(failureMessage+source+msg, msgAndArgs...))
		failer()
		return false

	case Comparison:
		return runCompareFunc(failer, t, check, msgAndArgs...)

	case func() (success bool, message string):
		return runCompareFunc(failer, t, check, msgAndArgs...)

	default:
		panic(fmt.Sprintf("comparison arg must be bool or Comparison, not %T", comparison))
	***REMOVED***
***REMOVED***

func runCompareFunc(failer func(), t TestingT, f Comparison, msgAndArgs ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	if ht, ok := t.(helperT); ok ***REMOVED***
		ht.Helper()
	***REMOVED***
	if success, message := f(); !success ***REMOVED***
		t.Log(format.WithCustomMessage(failureMessage+message, msgAndArgs...))
		failer()
		return false
	***REMOVED***
	return true
***REMOVED***

// Assert performs a comparison, marks the test as having failed if the comparison
// returns false, and stops execution immediately.
func Assert(t TestingT, comparison BoolOrComparison, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if ht, ok := t.(helperT); ok ***REMOVED***
		ht.Helper()
	***REMOVED***
	assert(t, t.FailNow, comparison, msgAndArgs...)
***REMOVED***

// Check performs a comparison and marks the test as having failed if the comparison
// returns false. Returns the result of the comparison.
func Check(t TestingT, comparison BoolOrComparison, msgAndArgs ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	if ht, ok := t.(helperT); ok ***REMOVED***
		ht.Helper()
	***REMOVED***
	return assert(t, t.Fail, comparison, msgAndArgs...)
***REMOVED***

// NilError fails the test immediately if the last arg is a non-nil error.
// This is equivalent to Assert(t, cmp.NilError(err)).
func NilError(t TestingT, err error, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if ht, ok := t.(helperT); ok ***REMOVED***
		ht.Helper()
	***REMOVED***
	assert(t, t.FailNow, cmp.NilError(err), msgAndArgs...)
***REMOVED***

// Equal uses the == operator to assert two values are equal and fails the test
// if they are not equal. This is equivalent to Assert(t, cmp.Equal(x, y)).
func Equal(t TestingT, x, y interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if ht, ok := t.(helperT); ok ***REMOVED***
		ht.Helper()
	***REMOVED***
	assert(t, t.FailNow, cmp.Equal(x, y), msgAndArgs...)
***REMOVED***
