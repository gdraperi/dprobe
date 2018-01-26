package store

type orCombinator struct ***REMOVED***
	bys []By
***REMOVED***

func (b orCombinator) isBy() ***REMOVED***
***REMOVED***

// Or returns a combinator that applies OR logic on all the supplied By
// arguments.
func Or(bys ...By) By ***REMOVED***
	return orCombinator***REMOVED***bys: bys***REMOVED***
***REMOVED***
