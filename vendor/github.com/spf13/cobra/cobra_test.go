package cobra

import (
	"testing"
	"text/template"
)

func TestAddTemplateFunctions(t *testing.T) ***REMOVED***
	AddTemplateFunc("t", func() bool ***REMOVED*** return true ***REMOVED***)
	AddTemplateFuncs(template.FuncMap***REMOVED***
		"f": func() bool ***REMOVED*** return false ***REMOVED***,
		"h": func() string ***REMOVED*** return "Hello," ***REMOVED***,
		"w": func() string ***REMOVED*** return "world." ***REMOVED******REMOVED***)

	c := &Command***REMOVED******REMOVED***
	c.SetUsageTemplate(`***REMOVED******REMOVED***if t***REMOVED******REMOVED******REMOVED******REMOVED***h***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***if f***REMOVED******REMOVED******REMOVED******REMOVED***h***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED*** ***REMOVED******REMOVED***w***REMOVED******REMOVED***`)

	const expected = "Hello, world."
	if got := c.UsageString(); got != expected ***REMOVED***
		t.Errorf("Expected UsageString: %v\nGot: %v", expected, got)
	***REMOVED***
***REMOVED***
