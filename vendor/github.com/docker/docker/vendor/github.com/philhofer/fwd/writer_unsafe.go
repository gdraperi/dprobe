// +build !appengine

package fwd

import (
	"reflect"
	"unsafe"
)

// unsafe cast string as []byte
func unsafestr(b string) []byte ***REMOVED***
	l := len(b)
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader***REMOVED***
		Len:  l,
		Cap:  l,
		Data: (*reflect.StringHeader)(unsafe.Pointer(&b)).Data,
	***REMOVED***))
***REMOVED***
