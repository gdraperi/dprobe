package cldr_test

import (
	"fmt"

	"golang.org/x/text/unicode/cldr"
)

func ExampleSlice() ***REMOVED***
	var dr *cldr.CLDR // assume this is initialized

	x, _ := dr.LDML("en")
	cs := x.Collations.Collation
	// remove all but the default
	cldr.MakeSlice(&cs).Filter(func(e cldr.Elem) bool ***REMOVED***
		return e.GetCommon().Type != x.Collations.Default()
	***REMOVED***)
	for i, c := range cs ***REMOVED***
		fmt.Println(i, c.Type)
	***REMOVED***
***REMOVED***
