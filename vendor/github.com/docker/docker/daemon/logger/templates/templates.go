package templates

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"
)

// basicFunctions are the set of initial
// functions provided to every template.
var basicFunctions = template.FuncMap***REMOVED***
	"json": func(v interface***REMOVED******REMOVED***) string ***REMOVED***
		buf := &bytes.Buffer***REMOVED******REMOVED***
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		enc.Encode(v)
		// Remove the trailing new line added by the encoder
		return strings.TrimSpace(buf.String())
	***REMOVED***,
	"split":    strings.Split,
	"join":     strings.Join,
	"title":    strings.Title,
	"lower":    strings.ToLower,
	"upper":    strings.ToUpper,
	"pad":      padWithSpace,
	"truncate": truncateWithLength,
***REMOVED***

// NewParse creates a new tagged template with the basic functions
// and parses the given format.
func NewParse(tag, format string) (*template.Template, error) ***REMOVED***
	return template.New(tag).Funcs(basicFunctions).Parse(format)
***REMOVED***

// padWithSpace adds whitespace to the input if the input is non-empty
func padWithSpace(source string, prefix, suffix int) string ***REMOVED***
	if source == "" ***REMOVED***
		return source
	***REMOVED***
	return strings.Repeat(" ", prefix) + source + strings.Repeat(" ", suffix)
***REMOVED***

// truncateWithLength truncates the source string up to the length provided by the input
func truncateWithLength(source string, length int) string ***REMOVED***
	if len(source) < length ***REMOVED***
		return source
	***REMOVED***
	return source[:length]
***REMOVED***
