package store

import (
	"strings"
)

// CreateEndpoints creates a list of endpoints given the right scheme
func CreateEndpoints(addrs []string, scheme string) (entries []string) ***REMOVED***
	for _, addr := range addrs ***REMOVED***
		entries = append(entries, scheme+"://"+addr)
	***REMOVED***
	return entries
***REMOVED***

// Normalize the key for each store to the form:
//
//     /path/to/key
//
func Normalize(key string) string ***REMOVED***
	return "/" + join(SplitKey(key))
***REMOVED***

// GetDirectory gets the full directory part of
// the key to the form:
//
//     /path/to/
//
func GetDirectory(key string) string ***REMOVED***
	parts := SplitKey(key)
	parts = parts[:len(parts)-1]
	return "/" + join(parts)
***REMOVED***

// SplitKey splits the key to extract path informations
func SplitKey(key string) (path []string) ***REMOVED***
	if strings.Contains(key, "/") ***REMOVED***
		path = strings.Split(key, "/")
	***REMOVED*** else ***REMOVED***
		path = []string***REMOVED***key***REMOVED***
	***REMOVED***
	return path
***REMOVED***

// join the path parts with '/'
func join(parts []string) string ***REMOVED***
	return strings.Join(parts, "/")
***REMOVED***
