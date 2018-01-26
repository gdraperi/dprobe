package libcontainerd

import (
	"time"

	"github.com/Microsoft/hcsshim"
	opengcs "github.com/Microsoft/opengcs/client"
)

// Summary contains a ProcessList item from HCS to support `top`
type Summary hcsshim.ProcessListItem

// Stats contains statistics from HCS
type Stats struct ***REMOVED***
	Read     time.Time
	HCSStats *hcsshim.Statistics
***REMOVED***

func interfaceToStats(read time.Time, v interface***REMOVED******REMOVED***) *Stats ***REMOVED***
	return &Stats***REMOVED***
		HCSStats: v.(*hcsshim.Statistics),
		Read:     read,
	***REMOVED***
***REMOVED***

// Resources defines updatable container resource values.
type Resources struct***REMOVED******REMOVED***

// LCOWOption is a CreateOption required for LCOW configuration
type LCOWOption struct ***REMOVED***
	Config *opengcs.Config
***REMOVED***

// Checkpoint holds the details of a checkpoint (not supported in windows)
type Checkpoint struct ***REMOVED***
	Name string
***REMOVED***

// Checkpoints contains the details of a checkpoint
type Checkpoints struct ***REMOVED***
	Checkpoints []*Checkpoint
***REMOVED***
