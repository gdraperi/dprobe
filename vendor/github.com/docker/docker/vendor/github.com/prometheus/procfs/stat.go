package procfs

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Stat represents kernel/system statistics.
type Stat struct ***REMOVED***
	// Boot time in seconds since the Epoch.
	BootTime int64
***REMOVED***

// NewStat returns kernel/system statistics read from /proc/stat.
func NewStat() (Stat, error) ***REMOVED***
	fs, err := NewFS(DefaultMountPoint)
	if err != nil ***REMOVED***
		return Stat***REMOVED******REMOVED***, err
	***REMOVED***

	return fs.NewStat()
***REMOVED***

// NewStat returns an information about current kernel/system statistics.
func (fs FS) NewStat() (Stat, error) ***REMOVED***
	f, err := os.Open(fs.Path("stat"))
	if err != nil ***REMOVED***
		return Stat***REMOVED******REMOVED***, err
	***REMOVED***
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() ***REMOVED***
		line := s.Text()
		if !strings.HasPrefix(line, "btime") ***REMOVED***
			continue
		***REMOVED***
		fields := strings.Fields(line)
		if len(fields) != 2 ***REMOVED***
			return Stat***REMOVED******REMOVED***, fmt.Errorf("couldn't parse %s line %s", f.Name(), line)
		***REMOVED***
		i, err := strconv.ParseInt(fields[1], 10, 32)
		if err != nil ***REMOVED***
			return Stat***REMOVED******REMOVED***, fmt.Errorf("couldn't parse %s: %s", fields[1], err)
		***REMOVED***
		return Stat***REMOVED***BootTime: i***REMOVED***, nil
	***REMOVED***
	if err := s.Err(); err != nil ***REMOVED***
		return Stat***REMOVED******REMOVED***, fmt.Errorf("couldn't parse %s: %s", f.Name(), err)
	***REMOVED***

	return Stat***REMOVED******REMOVED***, fmt.Errorf("couldn't parse %s, missing btime", f.Name())
***REMOVED***
