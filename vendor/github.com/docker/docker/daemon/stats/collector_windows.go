package stats

// platformNewStatsCollector performs platform specific initialisation of the
// Collector structure. This is a no-op on Windows.
func platformNewStatsCollector(s *Collector) ***REMOVED***
***REMOVED***

// getSystemCPUUsage returns the host system's cpu usage in
// nanoseconds. An error is returned if the format of the underlying
// file does not match. This is a no-op on Windows.
func (s *Collector) getSystemCPUUsage() (uint64, error) ***REMOVED***
	return 0, nil
***REMOVED***

func (s *Collector) getNumberOnlineCPUs() (uint32, error) ***REMOVED***
	return 0, nil
***REMOVED***
