package instructions

import "strings"

// handleJSONArgs parses command passed to CMD, ENTRYPOINT, RUN and SHELL instruction in Dockerfile
// for exec form it returns untouched args slice
// for shell form it returns concatenated args as the first element of a slice
func handleJSONArgs(args []string, attributes map[string]bool) []string ***REMOVED***
	if len(args) == 0 ***REMOVED***
		return []string***REMOVED******REMOVED***
	***REMOVED***

	if attributes != nil && attributes["json"] ***REMOVED***
		return args
	***REMOVED***

	// literal string command, not an exec array
	return []string***REMOVED***strings.Join(args, " ")***REMOVED***
***REMOVED***
