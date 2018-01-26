package template

import (
	"strings"
	"text/template"
)

// funcMap defines functions for our template system.
var funcMap = template.FuncMap***REMOVED***
	"join": func(s ...string) string ***REMOVED***
		// first arg is sep, remaining args are strings to join
		return strings.Join(s[1:], s[0])
	***REMOVED***,
***REMOVED***

func newTemplate(s string, extraFuncs template.FuncMap) (*template.Template, error) ***REMOVED***
	tmpl := template.New("expansion").Option("missingkey=error").Funcs(funcMap)
	if len(extraFuncs) != 0 ***REMOVED***
		tmpl = tmpl.Funcs(extraFuncs)
	***REMOVED***
	return tmpl.Parse(s)
***REMOVED***
