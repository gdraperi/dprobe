package main

// this is a file 'a.go'

import (
	superos "os"
)

func B() superos.Error ***REMOVED***
	return nil
***REMOVED***

// notice how changing type of a return function in one file,
// the inferred type of a variable in another file changes also

func (t *Tester) SetC() ***REMOVED***
	t.c = 31337
***REMOVED***

func (t *Tester) SetD() ***REMOVED***
	t.d = 31337
***REMOVED***

// support for multifile packages, including correct namespace handling
