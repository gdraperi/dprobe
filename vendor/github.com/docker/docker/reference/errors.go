package reference

type notFoundError string

func (e notFoundError) Error() string ***REMOVED***
	return string(e)
***REMOVED***

func (notFoundError) NotFound() ***REMOVED******REMOVED***

type invalidTagError string

func (e invalidTagError) Error() string ***REMOVED***
	return string(e)
***REMOVED***

func (invalidTagError) InvalidParameter() ***REMOVED******REMOVED***

type conflictingTagError string

func (e conflictingTagError) Error() string ***REMOVED***
	return string(e)
***REMOVED***

func (conflictingTagError) Conflict() ***REMOVED******REMOVED***
