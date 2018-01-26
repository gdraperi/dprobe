package system

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/docker/go-units"
)

// ReadMemInfo retrieves memory statistics of the host system and returns a
// MemInfo type.
func ReadMemInfo() (*MemInfo, error) ***REMOVED***
	file, err := os.Open("/proc/meminfo")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer file.Close()
	return parseMemInfo(file)
***REMOVED***

// parseMemInfo parses the /proc/meminfo file into
// a MemInfo object given an io.Reader to the file.
// Throws error if there are problems reading from the file
func parseMemInfo(reader io.Reader) (*MemInfo, error) ***REMOVED***
	meminfo := &MemInfo***REMOVED******REMOVED***
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() ***REMOVED***
		// Expected format: ["MemTotal:", "1234", "kB"]
		parts := strings.Fields(scanner.Text())

		// Sanity checks: Skip malformed entries.
		if len(parts) < 3 || parts[2] != "kB" ***REMOVED***
			continue
		***REMOVED***

		// Convert to bytes.
		size, err := strconv.Atoi(parts[1])
		if err != nil ***REMOVED***
			continue
		***REMOVED***
		bytes := int64(size) * units.KiB

		switch parts[0] ***REMOVED***
		case "MemTotal:":
			meminfo.MemTotal = bytes
		case "MemFree:":
			meminfo.MemFree = bytes
		case "SwapTotal:":
			meminfo.SwapTotal = bytes
		case "SwapFree:":
			meminfo.SwapFree = bytes
		***REMOVED***

	***REMOVED***

	// Handle errors that may have occurred during the reading of the file.
	if err := scanner.Err(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return meminfo, nil
***REMOVED***
