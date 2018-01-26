package requirement

import (
	"fmt"
	"path"
	"reflect"
	"runtime"
	"strings"
)

// SkipT is the interface required to skip tests
type SkipT interface ***REMOVED***
	Skip(reason string)
***REMOVED***

// Test represent a function that can be used as a requirement validation.
type Test func() bool

// Is checks if the environment satisfies the requirements
// for the test to run or skips the tests.
func Is(s SkipT, requirements ...Test) ***REMOVED***
	for _, r := range requirements ***REMOVED***
		isValid := r()
		if !isValid ***REMOVED***
			requirementFunc := runtime.FuncForPC(reflect.ValueOf(r).Pointer()).Name()
			s.Skip(fmt.Sprintf("unmatched requirement %s", extractRequirement(requirementFunc)))
		***REMOVED***
	***REMOVED***
***REMOVED***

func extractRequirement(requirementFunc string) string ***REMOVED***
	requirement := path.Base(requirementFunc)
	return strings.SplitN(requirement, ".", 2)[1]
***REMOVED***
