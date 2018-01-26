package httputils

import (
	"net/http"
	"strconv"
	"strings"
)

// BoolValue transforms a form value in different formats into a boolean type.
func BoolValue(r *http.Request, k string) bool ***REMOVED***
	s := strings.ToLower(strings.TrimSpace(r.FormValue(k)))
	return !(s == "" || s == "0" || s == "no" || s == "false" || s == "none")
***REMOVED***

// BoolValueOrDefault returns the default bool passed if the query param is
// missing, otherwise it's just a proxy to boolValue above.
func BoolValueOrDefault(r *http.Request, k string, d bool) bool ***REMOVED***
	if _, ok := r.Form[k]; !ok ***REMOVED***
		return d
	***REMOVED***
	return BoolValue(r, k)
***REMOVED***

// Int64ValueOrZero parses a form value into an int64 type.
// It returns 0 if the parsing fails.
func Int64ValueOrZero(r *http.Request, k string) int64 ***REMOVED***
	val, err := Int64ValueOrDefault(r, k, 0)
	if err != nil ***REMOVED***
		return 0
	***REMOVED***
	return val
***REMOVED***

// Int64ValueOrDefault parses a form value into an int64 type. If there is an
// error, returns the error. If there is no value returns the default value.
func Int64ValueOrDefault(r *http.Request, field string, def int64) (int64, error) ***REMOVED***
	if r.Form.Get(field) != "" ***REMOVED***
		value, err := strconv.ParseInt(r.Form.Get(field), 10, 64)
		return value, err
	***REMOVED***
	return def, nil
***REMOVED***

// ArchiveOptions stores archive information for different operations.
type ArchiveOptions struct ***REMOVED***
	Name string
	Path string
***REMOVED***

type badParameterError struct ***REMOVED***
	param string
***REMOVED***

func (e badParameterError) Error() string ***REMOVED***
	return "bad parameter: " + e.param + "cannot be empty"
***REMOVED***

func (e badParameterError) InvalidParameter() ***REMOVED******REMOVED***

// ArchiveFormValues parses form values and turns them into ArchiveOptions.
// It fails if the archive name and path are not in the request.
func ArchiveFormValues(r *http.Request, vars map[string]string) (ArchiveOptions, error) ***REMOVED***
	if err := ParseForm(r); err != nil ***REMOVED***
		return ArchiveOptions***REMOVED******REMOVED***, err
	***REMOVED***

	name := vars["name"]
	if name == "" ***REMOVED***
		return ArchiveOptions***REMOVED******REMOVED***, badParameterError***REMOVED***"name"***REMOVED***
	***REMOVED***
	path := r.Form.Get("path")
	if path == "" ***REMOVED***
		return ArchiveOptions***REMOVED******REMOVED***, badParameterError***REMOVED***"path"***REMOVED***
	***REMOVED***
	return ArchiveOptions***REMOVED***name, path***REMOVED***, nil
***REMOVED***
