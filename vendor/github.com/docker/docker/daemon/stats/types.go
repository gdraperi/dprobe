package stats

import (
	"bufio"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/container"
	"github.com/docker/docker/pkg/pubsub"
)

type supervisor interface ***REMOVED***
	// GetContainerStats collects all the stats related to a container
	GetContainerStats(container *container.Container) (*types.StatsJSON, error)
***REMOVED***

// NewCollector creates a stats collector that will poll the supervisor with the specified interval
func NewCollector(supervisor supervisor, interval time.Duration) *Collector ***REMOVED***
	s := &Collector***REMOVED***
		interval:   interval,
		supervisor: supervisor,
		publishers: make(map[*container.Container]*pubsub.Publisher),
		bufReader:  bufio.NewReaderSize(nil, 128),
	***REMOVED***

	platformNewStatsCollector(s)

	return s
***REMOVED***

// Collector manages and provides container resource stats
type Collector struct ***REMOVED***
	m          sync.Mutex
	supervisor supervisor
	interval   time.Duration
	publishers map[*container.Container]*pubsub.Publisher
	bufReader  *bufio.Reader

	// The following fields are not set on Windows currently.
	clockTicksPerSecond uint64
***REMOVED***
