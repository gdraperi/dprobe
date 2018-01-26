/*Package skip provides functions for skipping based on a condition.
 */
package skip

import (
	"fmt"
	"path"
	"reflect"
	"runtime"
	"strings"

	"github.com/gotestyourself/gotestyourself/internal/format"
	"github.com/gotestyourself/gotestyourself/internal/source"
)

type skipT interface ***REMOVED***
	Skip(args ...interface***REMOVED******REMOVED***)
	Log(args ...interface***REMOVED******REMOVED***)
***REMOVED***

type helperT interface ***REMOVED***
	Helper()
***REMOVED***

// BoolOrCheckFunc can be a bool or func() bool, other types will panic
type BoolOrCheckFunc interface***REMOVED******REMOVED***

// If skips the test if the check function returns true. The skip message will
// contain the name of the check function. Extra message text can be passed as a
// format string with args
func If(t skipT, condition BoolOrCheckFunc, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if ht, ok := t.(helperT); ok ***REMOVED***
		ht.Helper()
	***REMOVED***
	switch check := condition.(type) ***REMOVED***
	case bool:
		ifCondition(t, check, msgAndArgs...)
	case func() bool:
		if check() ***REMOVED***
			t.Skip(format.WithCustomMessage(getFunctionName(check), msgAndArgs...))
		***REMOVED***
	default:
		panic(fmt.Sprintf("invalid type for condition arg: %T", check))
	***REMOVED***
***REMOVED***

func getFunctionName(function func() bool) string ***REMOVED***
	funcPath := runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name()
	return strings.SplitN(path.Base(funcPath), ".", 2)[1]
***REMOVED***

// IfCondition skips the test if the condition is true. The skip message will
// contain the source of the expression passed as the condition. Extra message
// text can be passed as a format string with args.
//
// Deprecated: Use If() which now accepts bool arguments
func IfCondition(t skipT, condition bool, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if ht, ok := t.(helperT); ok ***REMOVED***
		ht.Helper()
	***REMOVED***
	ifCondition(t, condition, msgAndArgs...)
***REMOVED***

func ifCondition(t skipT, condition bool, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if ht, ok := t.(helperT); ok ***REMOVED***
		ht.Helper()
	***REMOVED***
	if !condition ***REMOVED***
		return
	***REMOVED***
	const (
		stackIndex = 2
		argPos     = 1
	)
	source, err := source.GetCondition(stackIndex, argPos)
	if err != nil ***REMOVED***
		t.Log(err.Error())
		t.Skip(format.Message(msgAndArgs...))
	***REMOVED***
	t.Skip(format.WithCustomMessage(source, msgAndArgs...))
***REMOVED***
