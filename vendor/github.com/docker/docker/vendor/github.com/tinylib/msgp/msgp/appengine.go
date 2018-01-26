// +build appengine

package msgp

// let's just assume appengine
// uses 64-bit hardware...
const smallint = false

func UnsafeString(b []byte) string ***REMOVED***
	return string(b)
***REMOVED***

func UnsafeBytes(s string) []byte ***REMOVED***
	return []byte(s)
***REMOVED***
