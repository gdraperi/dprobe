package cgroups

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func NewBlkio(root string) *blkioController ***REMOVED***
	return &blkioController***REMOVED***
		root: filepath.Join(root, string(Blkio)),
	***REMOVED***
***REMOVED***

type blkioController struct ***REMOVED***
	root string
***REMOVED***

func (b *blkioController) Name() Name ***REMOVED***
	return Blkio
***REMOVED***

func (b *blkioController) Path(path string) string ***REMOVED***
	return filepath.Join(b.root, path)
***REMOVED***

func (b *blkioController) Create(path string, resources *specs.LinuxResources) error ***REMOVED***
	if err := os.MkdirAll(b.Path(path), defaultDirPerm); err != nil ***REMOVED***
		return err
	***REMOVED***
	if resources.BlockIO == nil ***REMOVED***
		return nil
	***REMOVED***
	for _, t := range createBlkioSettings(resources.BlockIO) ***REMOVED***
		if t.value != nil ***REMOVED***
			if err := ioutil.WriteFile(
				filepath.Join(b.Path(path), fmt.Sprintf("blkio.%s", t.name)),
				t.format(t.value),
				defaultFilePerm,
			); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (b *blkioController) Update(path string, resources *specs.LinuxResources) error ***REMOVED***
	return b.Create(path, resources)
***REMOVED***

func (b *blkioController) Stat(path string, stats *Metrics) error ***REMOVED***
	stats.Blkio = &BlkIOStat***REMOVED******REMOVED***
	settings := []blkioStatSettings***REMOVED***
		***REMOVED***
			name:  "throttle.io_serviced",
			entry: &stats.Blkio.IoServicedRecursive,
		***REMOVED***,
		***REMOVED***
			name:  "throttle.io_service_bytes",
			entry: &stats.Blkio.IoServiceBytesRecursive,
		***REMOVED***,
	***REMOVED***
	// Try to read CFQ stats available on all CFQ enabled kernels first
	if _, err := os.Lstat(filepath.Join(b.Path(path), fmt.Sprintf("blkio.io_serviced_recursive"))); err == nil ***REMOVED***
		settings = append(settings,
			blkioStatSettings***REMOVED***
				name:  "sectors_recursive",
				entry: &stats.Blkio.SectorsRecursive,
			***REMOVED***,
			blkioStatSettings***REMOVED***
				name:  "io_service_bytes_recursive",
				entry: &stats.Blkio.IoServiceBytesRecursive,
			***REMOVED***,
			blkioStatSettings***REMOVED***
				name:  "io_serviced_recursive",
				entry: &stats.Blkio.IoServicedRecursive,
			***REMOVED***,
			blkioStatSettings***REMOVED***
				name:  "io_queued_recursive",
				entry: &stats.Blkio.IoQueuedRecursive,
			***REMOVED***,
			blkioStatSettings***REMOVED***
				name:  "io_service_time_recursive",
				entry: &stats.Blkio.IoServiceTimeRecursive,
			***REMOVED***,
			blkioStatSettings***REMOVED***
				name:  "io_wait_time_recursive",
				entry: &stats.Blkio.IoWaitTimeRecursive,
			***REMOVED***,
			blkioStatSettings***REMOVED***
				name:  "io_merged_recursive",
				entry: &stats.Blkio.IoMergedRecursive,
			***REMOVED***,
			blkioStatSettings***REMOVED***
				name:  "time_recursive",
				entry: &stats.Blkio.IoTimeRecursive,
			***REMOVED***,
		)
	***REMOVED***

	devices, err := getDevices("/dev")
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, t := range settings ***REMOVED***
		if err := b.readEntry(devices, path, t.name, t.entry); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (b *blkioController) readEntry(devices map[deviceKey]string, path, name string, entry *[]*BlkIOEntry) error ***REMOVED***
	f, err := os.Open(filepath.Join(b.Path(path), fmt.Sprintf("blkio.%s", name)))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() ***REMOVED***
		if err := sc.Err(); err != nil ***REMOVED***
			return err
		***REMOVED***
		// format: dev type amount
		fields := strings.FieldsFunc(sc.Text(), splitBlkIOStatLine)
		if len(fields) < 3 ***REMOVED***
			if len(fields) == 2 && fields[0] == "Total" ***REMOVED***
				// skip total line
				continue
			***REMOVED*** else ***REMOVED***
				return fmt.Errorf("Invalid line found while parsing %s: %s", path, sc.Text())
			***REMOVED***
		***REMOVED***
		major, err := strconv.ParseUint(fields[0], 10, 64)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		minor, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		op := ""
		valueField := 2
		if len(fields) == 4 ***REMOVED***
			op = fields[2]
			valueField = 3
		***REMOVED***
		v, err := strconv.ParseUint(fields[valueField], 10, 64)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*entry = append(*entry, &BlkIOEntry***REMOVED***
			Device: devices[deviceKey***REMOVED***major, minor***REMOVED***],
			Major:  major,
			Minor:  minor,
			Op:     op,
			Value:  v,
		***REMOVED***)
	***REMOVED***
	return nil
***REMOVED***

func createBlkioSettings(blkio *specs.LinuxBlockIO) []blkioSettings ***REMOVED***
	settings := []blkioSettings***REMOVED***
		***REMOVED***
			name:   "weight",
			value:  blkio.Weight,
			format: uintf,
		***REMOVED***,
		***REMOVED***
			name:   "leaf_weight",
			value:  blkio.LeafWeight,
			format: uintf,
		***REMOVED***,
	***REMOVED***
	for _, wd := range blkio.WeightDevice ***REMOVED***
		settings = append(settings,
			blkioSettings***REMOVED***
				name:   "weight_device",
				value:  wd,
				format: weightdev,
			***REMOVED***,
			blkioSettings***REMOVED***
				name:   "leaf_weight_device",
				value:  wd,
				format: weightleafdev,
			***REMOVED***)
	***REMOVED***
	for _, t := range []struct ***REMOVED***
		name string
		list []specs.LinuxThrottleDevice
	***REMOVED******REMOVED***
		***REMOVED***
			name: "throttle.read_bps_device",
			list: blkio.ThrottleReadBpsDevice,
		***REMOVED***,
		***REMOVED***
			name: "throttle.read_iops_device",
			list: blkio.ThrottleReadIOPSDevice,
		***REMOVED***,
		***REMOVED***
			name: "throttle.write_bps_device",
			list: blkio.ThrottleWriteBpsDevice,
		***REMOVED***,
		***REMOVED***
			name: "throttle.write_iops_device",
			list: blkio.ThrottleWriteIOPSDevice,
		***REMOVED***,
	***REMOVED*** ***REMOVED***
		for _, td := range t.list ***REMOVED***
			settings = append(settings, blkioSettings***REMOVED***
				name:   t.name,
				value:  td,
				format: throttleddev,
			***REMOVED***)
		***REMOVED***
	***REMOVED***
	return settings
***REMOVED***

type blkioSettings struct ***REMOVED***
	name   string
	value  interface***REMOVED******REMOVED***
	format func(v interface***REMOVED******REMOVED***) []byte
***REMOVED***

type blkioStatSettings struct ***REMOVED***
	name  string
	entry *[]*BlkIOEntry
***REMOVED***

func uintf(v interface***REMOVED******REMOVED***) []byte ***REMOVED***
	return []byte(strconv.FormatUint(uint64(*v.(*uint16)), 10))
***REMOVED***

func weightdev(v interface***REMOVED******REMOVED***) []byte ***REMOVED***
	wd := v.(specs.LinuxWeightDevice)
	return []byte(fmt.Sprintf("%d:%d %d", wd.Major, wd.Minor, wd.Weight))
***REMOVED***

func weightleafdev(v interface***REMOVED******REMOVED***) []byte ***REMOVED***
	wd := v.(specs.LinuxWeightDevice)
	return []byte(fmt.Sprintf("%d:%d %d", wd.Major, wd.Minor, wd.LeafWeight))
***REMOVED***

func throttleddev(v interface***REMOVED******REMOVED***) []byte ***REMOVED***
	td := v.(specs.LinuxThrottleDevice)
	return []byte(fmt.Sprintf("%d:%d %d", td.Major, td.Minor, td.Rate))
***REMOVED***

func splitBlkIOStatLine(r rune) bool ***REMOVED***
	return r == ' ' || r == ':'
***REMOVED***

type deviceKey struct ***REMOVED***
	major, minor uint64
***REMOVED***

// getDevices makes a best effort attempt to read all the devices into a map
// keyed by major and minor number. Since devices may be mapped multiple times,
// we err on taking the first occurrence.
func getDevices(path string) (map[deviceKey]string, error) ***REMOVED***
	// TODO(stevvooe): We are ignoring lots of errors. It might be kind of
	// challenging to debug this if we aren't mapping devices correctly.
	// Consider logging these errors.
	devices := map[deviceKey]string***REMOVED******REMOVED***
	if err := filepath.Walk(path, func(p string, fi os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		switch ***REMOVED***
		case fi.IsDir():
			switch fi.Name() ***REMOVED***
			case "pts", "shm", "fd", "mqueue", ".lxc", ".lxd-mounts":
				return filepath.SkipDir
			default:
				return nil
			***REMOVED***
		case fi.Name() == "console":
			return nil
		default:
			if fi.Mode()&os.ModeDevice == 0 ***REMOVED***
				// skip non-devices
				return nil
			***REMOVED***

			st, ok := fi.Sys().(*syscall.Stat_t)
			if !ok ***REMOVED***
				return fmt.Errorf("%s: unable to convert to system stat", p)
			***REMOVED***

			key := deviceKey***REMOVED***major(st.Rdev), minor(st.Rdev)***REMOVED***
			if _, ok := devices[key]; ok ***REMOVED***
				return nil // skip it if we have already populated the path.
			***REMOVED***

			devices[key] = p
		***REMOVED***

		return nil
	***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return devices, nil
***REMOVED***

func major(devNumber uint64) uint64 ***REMOVED***
	return (devNumber >> 8) & 0xfff
***REMOVED***

func minor(devNumber uint64) uint64 ***REMOVED***
	return (devNumber & 0xff) | ((devNumber >> 12) & 0xfff00)
***REMOVED***
