// +build !windows

package stats

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/opencontainers/runc/libcontainer/system"
)

/*
#include <unistd.h>
*/
import "C"

// platformNewStatsCollector performs platform specific initialisation of the
// Collector structure.
func platformNewStatsCollector(s *Collector) ***REMOVED***
	s.clockTicksPerSecond = uint64(system.GetClockTicks())
***REMOVED***

const nanoSecondsPerSecond = 1e9

// getSystemCPUUsage returns the host system's cpu usage in
// nanoseconds. An error is returned if the format of the underlying
// file does not match.
//
// Uses /proc/stat defined by POSIX. Looks for the cpu
// statistics line and then sums up the first seven fields
// provided. See `man 5 proc` for details on specific field
// information.
func (s *Collector) getSystemCPUUsage() (uint64, error) ***REMOVED***
	var line string
	f, err := os.Open("/proc/stat")
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	defer func() ***REMOVED***
		s.bufReader.Reset(nil)
		f.Close()
	***REMOVED***()
	s.bufReader.Reset(f)
	err = nil
	for err == nil ***REMOVED***
		line, err = s.bufReader.ReadString('\n')
		if err != nil ***REMOVED***
			break
		***REMOVED***
		parts := strings.Fields(line)
		switch parts[0] ***REMOVED***
		case "cpu":
			if len(parts) < 8 ***REMOVED***
				return 0, fmt.Errorf("invalid number of cpu fields")
			***REMOVED***
			var totalClockTicks uint64
			for _, i := range parts[1:8] ***REMOVED***
				v, err := strconv.ParseUint(i, 10, 64)
				if err != nil ***REMOVED***
					return 0, fmt.Errorf("Unable to convert value %s to int: %s", i, err)
				***REMOVED***
				totalClockTicks += v
			***REMOVED***
			return (totalClockTicks * nanoSecondsPerSecond) /
				s.clockTicksPerSecond, nil
		***REMOVED***
	***REMOVED***
	return 0, fmt.Errorf("invalid stat format. Error trying to parse the '/proc/stat' file")
***REMOVED***

func (s *Collector) getNumberOnlineCPUs() (uint32, error) ***REMOVED***
	i, err := C.sysconf(C._SC_NPROCESSORS_ONLN)
	// According to POSIX - errno is undefined after successful
	// sysconf, and can be non-zero in several cases, so look for
	// error in returned value not in errno.
	// (https://sourceware.org/bugzilla/show_bug.cgi?id=21536)
	if i == -1 ***REMOVED***
		return 0, err
	***REMOVED***
	return uint32(i), nil
***REMOVED***
