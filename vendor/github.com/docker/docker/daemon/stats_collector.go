package daemon

import (
	"runtime"
	"time"

	"github.com/docker/docker/daemon/stats"
	"github.com/docker/docker/pkg/system"
)

// newStatsCollector returns a new statsCollector that collections
// stats for a registered container at the specified interval.
// The collector allows non-running containers to be added
// and will start processing stats when they are started.
func (daemon *Daemon) newStatsCollector(interval time.Duration) *stats.Collector ***REMOVED***
	// FIXME(vdemeester) move this elsewhere
	if runtime.GOOS == "linux" ***REMOVED***
		meminfo, err := system.ReadMemInfo()
		if err == nil && meminfo.MemTotal > 0 ***REMOVED***
			daemon.machineMemory = uint64(meminfo.MemTotal)
		***REMOVED***
	***REMOVED***
	s := stats.NewCollector(daemon, interval)
	go s.Run()
	return s
***REMOVED***
