package reference

import "path"

// IsNameOnly returns true if reference only contains a repo name.
func IsNameOnly(ref Named) bool ***REMOVED***
	if _, ok := ref.(NamedTagged); ok ***REMOVED***
		return false
	***REMOVED***
	if _, ok := ref.(Canonical); ok ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

// FamiliarName returns the familiar name string
// for the given named, familiarizing if needed.
func FamiliarName(ref Named) string ***REMOVED***
	if nn, ok := ref.(normalizedNamed); ok ***REMOVED***
		return nn.Familiar().Name()
	***REMOVED***
	return ref.Name()
***REMOVED***

// FamiliarString returns the familiar string representation
// for the given reference, familiarizing if needed.
func FamiliarString(ref Reference) string ***REMOVED***
	if nn, ok := ref.(normalizedNamed); ok ***REMOVED***
		return nn.Familiar().String()
	***REMOVED***
	return ref.String()
***REMOVED***

// FamiliarMatch reports whether ref matches the specified pattern.
// See https://godoc.org/path#Match for supported patterns.
func FamiliarMatch(pattern string, ref Reference) (bool, error) ***REMOVED***
	matched, err := path.Match(pattern, FamiliarString(ref))
	if namedRef, isNamed := ref.(Named); isNamed && !matched ***REMOVED***
		matched, _ = path.Match(pattern, FamiliarName(namedRef))
	***REMOVED***
	return matched, err
***REMOVED***
