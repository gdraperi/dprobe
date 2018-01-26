// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package grpc_prometheus

import (
	"strings"

	"google.golang.org/grpc/codes"
)

var (
	allCodes = []codes.Code***REMOVED***
		codes.OK, codes.Canceled, codes.Unknown, codes.InvalidArgument, codes.DeadlineExceeded, codes.NotFound,
		codes.AlreadyExists, codes.PermissionDenied, codes.Unauthenticated, codes.ResourceExhausted,
		codes.FailedPrecondition, codes.Aborted, codes.OutOfRange, codes.Unimplemented, codes.Internal,
		codes.Unavailable, codes.DataLoss,
	***REMOVED***
)

func splitMethodName(fullMethodName string) (string, string) ***REMOVED***
	fullMethodName = strings.TrimPrefix(fullMethodName, "/") // remove leading slash
	if i := strings.Index(fullMethodName, "/"); i >= 0 ***REMOVED***
		return fullMethodName[:i], fullMethodName[i+1:]
	***REMOVED***
	return "unknown", "unknown"
***REMOVED***
