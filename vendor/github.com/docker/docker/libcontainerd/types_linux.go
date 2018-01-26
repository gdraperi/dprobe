package libcontainerd

import (
	"time"

	"github.com/containerd/cgroups"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

// Summary is not used on linux
type Summary struct***REMOVED******REMOVED***

// Stats holds metrics properties as returned by containerd
type Stats struct ***REMOVED***
	Read    time.Time
	Metrics *cgroups.Metrics
***REMOVED***

func interfaceToStats(read time.Time, v interface***REMOVED******REMOVED***) *Stats ***REMOVED***
	return &Stats***REMOVED***
		Metrics: v.(*cgroups.Metrics),
		Read:    read,
	***REMOVED***
***REMOVED***

// Resources defines updatable container resource values. TODO: it must match containerd upcoming API
type Resources specs.LinuxResources

// Checkpoints contains the details of a checkpoint
type Checkpoints struct***REMOVED******REMOVED***
