package multierror

// Append is a helper function that will append more errors
// onto an Error in order to create a larger multi-error.
//
// If err is not a multierror.Error, then it will be turned into
// one. If any of the errs are multierr.Error, they will be flattened
// one level into err.
func Append(err error, errs ...error) *Error ***REMOVED***
	switch err := err.(type) ***REMOVED***
	case *Error:
		// Typed nils can reach here, so initialize if we are nil
		if err == nil ***REMOVED***
			err = new(Error)
		***REMOVED***

		err.Errors = append(err.Errors, errs...)
		return err
	default:
		newErrs := make([]error, 0, len(errs)+1)
		if err != nil ***REMOVED***
			newErrs = append(newErrs, err)
		***REMOVED***
		newErrs = append(newErrs, errs...)

		return &Error***REMOVED***
			Errors: newErrs,
		***REMOVED***
	***REMOVED***
***REMOVED***
