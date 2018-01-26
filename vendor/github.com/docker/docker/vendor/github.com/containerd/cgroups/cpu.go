package cgroups

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func NewCpu(root string) *cpuController ***REMOVED***
	return &cpuController***REMOVED***
		root: filepath.Join(root, string(Cpu)),
	***REMOVED***
***REMOVED***

type cpuController struct ***REMOVED***
	root string
***REMOVED***

func (c *cpuController) Name() Name ***REMOVED***
	return Cpu
***REMOVED***

func (c *cpuController) Path(path string) string ***REMOVED***
	return filepath.Join(c.root, path)
***REMOVED***

func (c *cpuController) Create(path string, resources *specs.LinuxResources) error ***REMOVED***
	if err := os.MkdirAll(c.Path(path), defaultDirPerm); err != nil ***REMOVED***
		return err
	***REMOVED***
	if cpu := resources.CPU; cpu != nil ***REMOVED***
		for _, t := range []struct ***REMOVED***
			name   string
			ivalue *int64
			uvalue *uint64
		***REMOVED******REMOVED***
			***REMOVED***
				name:   "rt_period_us",
				uvalue: cpu.RealtimePeriod,
			***REMOVED***,
			***REMOVED***
				name:   "rt_runtime_us",
				ivalue: cpu.RealtimeRuntime,
			***REMOVED***,
			***REMOVED***
				name:   "shares",
				uvalue: cpu.Shares,
			***REMOVED***,
			***REMOVED***
				name:   "cfs_period_us",
				uvalue: cpu.Period,
			***REMOVED***,
			***REMOVED***
				name:   "cfs_quota_us",
				ivalue: cpu.Quota,
			***REMOVED***,
		***REMOVED*** ***REMOVED***
			var value []byte
			if t.uvalue != nil ***REMOVED***
				value = []byte(strconv.FormatUint(*t.uvalue, 10))
			***REMOVED*** else if t.ivalue != nil ***REMOVED***
				value = []byte(strconv.FormatInt(*t.ivalue, 10))
			***REMOVED***
			if value != nil ***REMOVED***
				if err := ioutil.WriteFile(
					filepath.Join(c.Path(path), fmt.Sprintf("cpu.%s", t.name)),
					value,
					defaultFilePerm,
				); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (c *cpuController) Update(path string, resources *specs.LinuxResources) error ***REMOVED***
	return c.Create(path, resources)
***REMOVED***

func (c *cpuController) Stat(path string, stats *Metrics) error ***REMOVED***
	f, err := os.Open(filepath.Join(c.Path(path), "cpu.stat"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer f.Close()
	// get or create the cpu field because cpuacct can also set values on this struct
	sc := bufio.NewScanner(f)
	for sc.Scan() ***REMOVED***
		if err := sc.Err(); err != nil ***REMOVED***
			return err
		***REMOVED***
		key, v, err := parseKV(sc.Text())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		switch key ***REMOVED***
		case "nr_periods":
			stats.CPU.Throttling.Periods = v
		case "nr_throttled":
			stats.CPU.Throttling.ThrottledPeriods = v
		case "throttled_time":
			stats.CPU.Throttling.ThrottledTime = v
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
