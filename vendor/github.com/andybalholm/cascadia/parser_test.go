package cascadia

import (
	"testing"
)

var identifierTests = map[string]string***REMOVED***
	"x":         "x",
	"96":        "",
	"-x":        "-x",
	`r\e9 sumé`: "résumé",
	`a\"b`:      `a"b`,
***REMOVED***

func TestParseIdentifier(t *testing.T) ***REMOVED***
	for source, want := range identifierTests ***REMOVED***
		p := &parser***REMOVED***s: source***REMOVED***
		got, err := p.parseIdentifier()

		if err != nil ***REMOVED***
			if want == "" ***REMOVED***
				// It was supposed to be an error.
				continue
			***REMOVED***
			t.Errorf("parsing %q: got error (%s), want %q", source, err, want)
			continue
		***REMOVED***

		if want == "" ***REMOVED***
			if err == nil ***REMOVED***
				t.Errorf("parsing %q: got %q, want error", source, got)
			***REMOVED***
			continue
		***REMOVED***

		if p.i < len(source) ***REMOVED***
			t.Errorf("parsing %q: %d bytes left over", source, len(source)-p.i)
			continue
		***REMOVED***

		if got != want ***REMOVED***
			t.Errorf("parsing %q: got %q, want %q", source, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

var stringTests = map[string]string***REMOVED***
	`"x"`:         "x",
	`'x'`:         "x",
	`'x`:          "",
	"'x\\\r\nx'":  "xx",
	`"r\e9 sumé"`: "résumé",
	`"a\"b"`:      `a"b`,
***REMOVED***

func TestParseString(t *testing.T) ***REMOVED***
	for source, want := range stringTests ***REMOVED***
		p := &parser***REMOVED***s: source***REMOVED***
		got, err := p.parseString()

		if err != nil ***REMOVED***
			if want == "" ***REMOVED***
				// It was supposed to be an error.
				continue
			***REMOVED***
			t.Errorf("parsing %q: got error (%s), want %q", source, err, want)
			continue
		***REMOVED***

		if want == "" ***REMOVED***
			if err == nil ***REMOVED***
				t.Errorf("parsing %q: got %q, want error", source, got)
			***REMOVED***
			continue
		***REMOVED***

		if p.i < len(source) ***REMOVED***
			t.Errorf("parsing %q: %d bytes left over", source, len(source)-p.i)
			continue
		***REMOVED***

		if got != want ***REMOVED***
			t.Errorf("parsing %q: got %q, want %q", source, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***
