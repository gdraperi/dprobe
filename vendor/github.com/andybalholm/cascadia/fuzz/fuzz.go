package fuzz

import "github.com/andybalholm/cascadia"

// Fuzz is the entrypoint used by the go-fuzz framework
func Fuzz(data []byte) int ***REMOVED***
	sel, err := cascadia.Compile(string(data))
	if err != nil ***REMOVED***
		if sel != nil ***REMOVED***
			panic("sel != nil on error")
		***REMOVED***
		return 0
	***REMOVED***
	return 1
***REMOVED***
