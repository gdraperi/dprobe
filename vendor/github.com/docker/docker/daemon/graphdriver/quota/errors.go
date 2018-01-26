package quota

import "github.com/docker/docker/errdefs"

var (
	_ errdefs.ErrNotImplemented = (*errQuotaNotSupported)(nil)
)

// ErrQuotaNotSupported indicates if were found the FS didn't have projects quotas available
var ErrQuotaNotSupported = errQuotaNotSupported***REMOVED******REMOVED***

type errQuotaNotSupported struct ***REMOVED***
***REMOVED***

func (e errQuotaNotSupported) NotImplemented() ***REMOVED******REMOVED***

func (e errQuotaNotSupported) Error() string ***REMOVED***
	return "Filesystem does not support, or has not enabled quotas"
***REMOVED***
