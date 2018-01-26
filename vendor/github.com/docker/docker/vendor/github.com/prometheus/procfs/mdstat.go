package procfs

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

var (
	statuslineRE = regexp.MustCompile(`(\d+) blocks .*\[(\d+)/(\d+)\] \[[U_]+\]`)
	buildlineRE  = regexp.MustCompile(`\((\d+)/\d+\)`)
)

// MDStat holds info parsed from /proc/mdstat.
type MDStat struct ***REMOVED***
	// Name of the device.
	Name string
	// activity-state of the device.
	ActivityState string
	// Number of active disks.
	DisksActive int64
	// Total number of disks the device consists of.
	DisksTotal int64
	// Number of blocks the device holds.
	BlocksTotal int64
	// Number of blocks on the device that are in sync.
	BlocksSynced int64
***REMOVED***

// ParseMDStat parses an mdstat-file and returns a struct with the relevant infos.
func (fs FS) ParseMDStat() (mdstates []MDStat, err error) ***REMOVED***
	mdStatusFilePath := fs.Path("mdstat")
	content, err := ioutil.ReadFile(mdStatusFilePath)
	if err != nil ***REMOVED***
		return []MDStat***REMOVED******REMOVED***, fmt.Errorf("error parsing %s: %s", mdStatusFilePath, err)
	***REMOVED***

	mdStates := []MDStat***REMOVED******REMOVED***
	lines := strings.Split(string(content), "\n")
	for i, l := range lines ***REMOVED***
		if l == "" ***REMOVED***
			continue
		***REMOVED***
		if l[0] == ' ' ***REMOVED***
			continue
		***REMOVED***
		if strings.HasPrefix(l, "Personalities") || strings.HasPrefix(l, "unused") ***REMOVED***
			continue
		***REMOVED***

		mainLine := strings.Split(l, " ")
		if len(mainLine) < 3 ***REMOVED***
			return mdStates, fmt.Errorf("error parsing mdline: %s", l)
		***REMOVED***
		mdName := mainLine[0]
		activityState := mainLine[2]

		if len(lines) <= i+3 ***REMOVED***
			return mdStates, fmt.Errorf(
				"error parsing %s: too few lines for md device %s",
				mdStatusFilePath,
				mdName,
			)
		***REMOVED***

		active, total, size, err := evalStatusline(lines[i+1])
		if err != nil ***REMOVED***
			return mdStates, fmt.Errorf("error parsing %s: %s", mdStatusFilePath, err)
		***REMOVED***

		// j is the line number of the syncing-line.
		j := i + 2
		if strings.Contains(lines[i+2], "bitmap") ***REMOVED*** // skip bitmap line
			j = i + 3
		***REMOVED***

		// If device is syncing at the moment, get the number of currently
		// synced bytes, otherwise that number equals the size of the device.
		syncedBlocks := size
		if strings.Contains(lines[j], "recovery") || strings.Contains(lines[j], "resync") ***REMOVED***
			syncedBlocks, err = evalBuildline(lines[j])
			if err != nil ***REMOVED***
				return mdStates, fmt.Errorf("error parsing %s: %s", mdStatusFilePath, err)
			***REMOVED***
		***REMOVED***

		mdStates = append(mdStates, MDStat***REMOVED***
			Name:          mdName,
			ActivityState: activityState,
			DisksActive:   active,
			DisksTotal:    total,
			BlocksTotal:   size,
			BlocksSynced:  syncedBlocks,
		***REMOVED***)
	***REMOVED***

	return mdStates, nil
***REMOVED***

func evalStatusline(statusline string) (active, total, size int64, err error) ***REMOVED***
	matches := statuslineRE.FindStringSubmatch(statusline)
	if len(matches) != 4 ***REMOVED***
		return 0, 0, 0, fmt.Errorf("unexpected statusline: %s", statusline)
	***REMOVED***

	size, err = strconv.ParseInt(matches[1], 10, 64)
	if err != nil ***REMOVED***
		return 0, 0, 0, fmt.Errorf("unexpected statusline %s: %s", statusline, err)
	***REMOVED***

	total, err = strconv.ParseInt(matches[2], 10, 64)
	if err != nil ***REMOVED***
		return 0, 0, 0, fmt.Errorf("unexpected statusline %s: %s", statusline, err)
	***REMOVED***

	active, err = strconv.ParseInt(matches[3], 10, 64)
	if err != nil ***REMOVED***
		return 0, 0, 0, fmt.Errorf("unexpected statusline %s: %s", statusline, err)
	***REMOVED***

	return active, total, size, nil
***REMOVED***

func evalBuildline(buildline string) (syncedBlocks int64, err error) ***REMOVED***
	matches := buildlineRE.FindStringSubmatch(buildline)
	if len(matches) != 2 ***REMOVED***
		return 0, fmt.Errorf("unexpected buildline: %s", buildline)
	***REMOVED***

	syncedBlocks, err = strconv.ParseInt(matches[1], 10, 64)
	if err != nil ***REMOVED***
		return 0, fmt.Errorf("%s in buildline: %s", err, buildline)
	***REMOVED***

	return syncedBlocks, nil
***REMOVED***
