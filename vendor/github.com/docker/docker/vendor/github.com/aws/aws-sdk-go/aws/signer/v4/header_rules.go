package v4

import (
	"net/http"
	"strings"
)

// validator houses a set of rule needed for validation of a
// string value
type rules []rule

// rule interface allows for more flexible rules and just simply
// checks whether or not a value adheres to that rule
type rule interface ***REMOVED***
	IsValid(value string) bool
***REMOVED***

// IsValid will iterate through all rules and see if any rules
// apply to the value and supports nested rules
func (r rules) IsValid(value string) bool ***REMOVED***
	for _, rule := range r ***REMOVED***
		if rule.IsValid(value) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// mapRule generic rule for maps
type mapRule map[string]struct***REMOVED******REMOVED***

// IsValid for the map rule satisfies whether it exists in the map
func (m mapRule) IsValid(value string) bool ***REMOVED***
	_, ok := m[value]
	return ok
***REMOVED***

// whitelist is a generic rule for whitelisting
type whitelist struct ***REMOVED***
	rule
***REMOVED***

// IsValid for whitelist checks if the value is within the whitelist
func (w whitelist) IsValid(value string) bool ***REMOVED***
	return w.rule.IsValid(value)
***REMOVED***

// blacklist is a generic rule for blacklisting
type blacklist struct ***REMOVED***
	rule
***REMOVED***

// IsValid for whitelist checks if the value is within the whitelist
func (b blacklist) IsValid(value string) bool ***REMOVED***
	return !b.rule.IsValid(value)
***REMOVED***

type patterns []string

// IsValid for patterns checks each pattern and returns if a match has
// been found
func (p patterns) IsValid(value string) bool ***REMOVED***
	for _, pattern := range p ***REMOVED***
		if strings.HasPrefix(http.CanonicalHeaderKey(value), pattern) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// inclusiveRules rules allow for rules to depend on one another
type inclusiveRules []rule

// IsValid will return true if all rules are true
func (r inclusiveRules) IsValid(value string) bool ***REMOVED***
	for _, rule := range r ***REMOVED***
		if !rule.IsValid(value) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***
