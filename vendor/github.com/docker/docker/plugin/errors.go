package plugin

import "fmt"

type errNotFound string

func (name errNotFound) Error() string ***REMOVED***
	return fmt.Sprintf("plugin %q not found", string(name))
***REMOVED***

func (errNotFound) NotFound() ***REMOVED******REMOVED***

type errAmbiguous string

func (name errAmbiguous) Error() string ***REMOVED***
	return fmt.Sprintf("multiple plugins found for %q", string(name))
***REMOVED***

func (name errAmbiguous) InvalidParameter() ***REMOVED******REMOVED***

type errDisabled string

func (name errDisabled) Error() string ***REMOVED***
	return fmt.Sprintf("plugin %s found but disabled", string(name))
***REMOVED***

func (name errDisabled) Conflict() ***REMOVED******REMOVED***

type invalidFilter struct ***REMOVED***
	filter string
	value  []string
***REMOVED***

func (e invalidFilter) Error() string ***REMOVED***
	msg := "Invalid filter '" + e.filter
	if len(e.value) > 0 ***REMOVED***
		msg += fmt.Sprintf("=%s", e.value)
	***REMOVED***
	return msg + "'"
***REMOVED***

func (invalidFilter) InvalidParameter() ***REMOVED******REMOVED***

type inUseError string

func (e inUseError) Error() string ***REMOVED***
	return "plugin " + string(e) + " is in use"
***REMOVED***

func (inUseError) Conflict() ***REMOVED******REMOVED***

type enabledError string

func (e enabledError) Error() string ***REMOVED***
	return "plugin " + string(e) + " is enabled"
***REMOVED***

func (enabledError) Conflict() ***REMOVED******REMOVED***

type alreadyExistsError string

func (e alreadyExistsError) Error() string ***REMOVED***
	return "plugin " + string(e) + " already exists"
***REMOVED***

func (alreadyExistsError) Conflict() ***REMOVED******REMOVED***
