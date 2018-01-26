package sts

import "github.com/aws/aws-sdk-go/aws/request"

func init() ***REMOVED***
	initRequest = func(r *request.Request) ***REMOVED***
		switch r.Operation.Name ***REMOVED***
		case opAssumeRoleWithSAML, opAssumeRoleWithWebIdentity:
			r.Handlers.Sign.Clear() // these operations are unsigned
		***REMOVED***
	***REMOVED***
***REMOVED***
