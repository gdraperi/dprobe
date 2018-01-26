// +build linux,cgo,libdm_no_deferred_remove

package devicemapper

// LibraryDeferredRemovalSupport tells if the feature is enabled in the build
const LibraryDeferredRemovalSupport = false

func dmTaskDeferredRemoveFct(task *cdmTask) int ***REMOVED***
	// Error. Nobody should be calling it.
	return -1
***REMOVED***

func dmTaskGetInfoWithDeferredFct(task *cdmTask, info *Info) int ***REMOVED***
	return -1
***REMOVED***
