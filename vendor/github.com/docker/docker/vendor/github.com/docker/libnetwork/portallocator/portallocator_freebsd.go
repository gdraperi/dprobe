package portallocator

import (
	"bytes"
	"fmt"
	"os/exec"
)

func getDynamicPortRange() (start int, end int, err error) ***REMOVED***
	portRangeKernelSysctl := []string***REMOVED***"net.inet.ip.portrange.hifirst", "net.ip.portrange.hilast"***REMOVED***
	portRangeFallback := fmt.Sprintf("using fallback port range %d-%d", DefaultPortRangeStart, DefaultPortRangeEnd)
	portRangeLowCmd := exec.Command("/sbin/sysctl", portRangeKernelSysctl[0])
	var portRangeLowOut bytes.Buffer
	portRangeLowCmd.Stdout = &portRangeLowOut
	cmdErr := portRangeLowCmd.Run()
	if cmdErr != nil ***REMOVED***
		return 0, 0, fmt.Errorf("port allocator - sysctl net.inet.ip.portrange.hifirst failed - %s: %v", portRangeFallback, err)
	***REMOVED***
	n, err := fmt.Sscanf(portRangeLowOut.String(), "%d", &start)
	if n != 1 || err != nil ***REMOVED***
		if err == nil ***REMOVED***
			err = fmt.Errorf("unexpected count of parsed numbers (%d)", n)
		***REMOVED***
		return 0, 0, fmt.Errorf("port allocator - failed to parse system ephemeral port range start from %s - %s: %v", portRangeLowOut.String(), portRangeFallback, err)
	***REMOVED***

	portRangeHighCmd := exec.Command("/sbin/sysctl", portRangeKernelSysctl[1])
	var portRangeHighOut bytes.Buffer
	portRangeHighCmd.Stdout = &portRangeHighOut
	cmdErr = portRangeHighCmd.Run()
	if cmdErr != nil ***REMOVED***
		return 0, 0, fmt.Errorf("port allocator - sysctl net.inet.ip.portrange.hilast failed - %s: %v", portRangeFallback, err)
	***REMOVED***
	n, err = fmt.Sscanf(portRangeHighOut.String(), "%d", &end)
	if n != 1 || err != nil ***REMOVED***
		if err == nil ***REMOVED***
			err = fmt.Errorf("unexpected count of parsed numbers (%d)", n)
		***REMOVED***
		return 0, 0, fmt.Errorf("port allocator - failed to parse system ephemeral port range end from %s - %s: %v", portRangeHighOut.String(), portRangeFallback, err)
	***REMOVED***
	return start, end, nil
***REMOVED***
