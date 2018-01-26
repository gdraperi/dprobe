package fs

import "context"

// Usage of disk information
type Usage struct ***REMOVED***
	Inodes int64
	Size   int64
***REMOVED***

// DiskUsage counts the number of inodes and disk usage for the resources under
// path.
func DiskUsage(roots ...string) (Usage, error) ***REMOVED***
	return diskUsage(roots...)
***REMOVED***

// DiffUsage counts the numbers of inodes and disk usage in the
// diff between the 2 directories. The first path is intended
// as the base directory and the second as the changed directory.
func DiffUsage(ctx context.Context, a, b string) (Usage, error) ***REMOVED***
	return diffUsage(ctx, a, b)
***REMOVED***
