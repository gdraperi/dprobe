package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseNameValOldFormat(t *testing.T) ***REMOVED***
	directive := Directive***REMOVED******REMOVED***
	node, err := parseNameVal("foo bar", "LABEL", &directive)
	assert.NoError(t, err)

	expected := &Node***REMOVED***
		Value: "foo",
		Next:  &Node***REMOVED***Value: "bar"***REMOVED***,
	***REMOVED***
	assert.Equal(t, expected, node)
***REMOVED***

func TestParseNameValNewFormat(t *testing.T) ***REMOVED***
	directive := Directive***REMOVED******REMOVED***
	node, err := parseNameVal("foo=bar thing=star", "LABEL", &directive)
	assert.NoError(t, err)

	expected := &Node***REMOVED***
		Value: "foo",
		Next: &Node***REMOVED***
			Value: "bar",
			Next: &Node***REMOVED***
				Value: "thing",
				Next: &Node***REMOVED***
					Value: "star",
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	assert.Equal(t, expected, node)
***REMOVED***

func TestNodeFromLabels(t *testing.T) ***REMOVED***
	labels := map[string]string***REMOVED***
		"foo":   "bar",
		"weird": "first' second",
	***REMOVED***
	expected := &Node***REMOVED***
		Value:    "label",
		Original: `LABEL "foo"='bar' "weird"='first' second'`,
		Next: &Node***REMOVED***
			Value: "foo",
			Next: &Node***REMOVED***
				Value: "'bar'",
				Next: &Node***REMOVED***
					Value: "weird",
					Next: &Node***REMOVED***
						Value: "'first' second'",
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	node := NodeFromLabels(labels)
	assert.Equal(t, expected, node)

***REMOVED***

func TestParseNameValWithoutVal(t *testing.T) ***REMOVED***
	directive := Directive***REMOVED******REMOVED***
	// In Config.Env, a variable without `=` is removed from the environment. (#31634)
	// However, in Dockerfile, we don't allow "unsetting" an environment variable. (#11922)
	_, err := parseNameVal("foo", "ENV", &directive)
	assert.Error(t, err, "ENV must have two arguments")
***REMOVED***
