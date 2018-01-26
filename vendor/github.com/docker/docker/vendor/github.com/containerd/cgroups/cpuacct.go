package cgroups

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

const nanosecondsInSecond = 1000000000

var clockTicks = getClockTicks()

func NewCpuacct(root string) *cpuacctController ***REMOVED***
	return &cpuacctController***REMOVED***
		root: filepath.Join(root, string(Cpuacct)),
	***REMOVED***
***REMOVED***

type cpuacctController struct ***REMOVED***
	root string
***REMOVED***

func (c *cpuacctController) Name() Name ***REMOVED***
	return Cpuacct
***REMOVED***

func (c *cpuacctController) Path(path string) string ***REMOVED***
	return filepath.Join(c.root, path)
***REMOVED***

func (c *cpuacctController) Stat(path string, stats *Metrics) error ***REMOVED***
	user, kernel, err := c.getUsage(path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	total, err := readUint(filepath.Join(c.Path(path), "cpuacct.usage"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	percpu, err := c.percpuUsage(path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	stats.CPU.Usage.Total = total
	stats.CPU.Usage.User = user
	stats.CPU.Usage.Kernel = kernel
	stats.CPU.Usage.PerCPU = percpu
	return nil
***REMOVED***

func (c *cpuacctController) percpuUsage(path string) ([]uint64, error) ***REMOVED***
	var usage []uint64
	data, err := ioutil.ReadFile(filepath.Join(c.Path(path), "cpuacct.usage_percpu"))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	for _, v := range strings.Fields(string(data)) ***REMOVED***
		u, err := strconv.ParseUint(v, 10, 64)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		usage = append(usage, u)
	***REMOVED***
	return usage, nil
***REMOVED***

func (c *cpuacctController) getUsage(path string) (user uint64, kernel uint64, err error) ***REMOVED***
	statPath := filepath.Join(c.Path(path), "cpuacct.stat")
	data, err := ioutil.ReadFile(statPath)
	if err != nil ***REMOVED***
		return 0, 0, err
	***REMOVED***
	fields := strings.Fields(string(data))
	if len(fields) != 4 ***REMOVED***
		return 0, 0, fmt.Errorf("%q is expected to have 4 fields", statPath)
	***REMOVED***
	for _, t := range []struct ***REMOVED***
		index int
		name  string
		value *uint64
	***REMOVED******REMOVED***
		***REMOVED***
			index: 0,
			name:  "user",
			value: &user,
		***REMOVED***,
		***REMOVED***
			index: 2,
			name:  "system",
			value: &kernel,
		***REMOVED***,
	***REMOVED*** ***REMOVED***
		if fields[t.index] != t.name ***REMOVED***
			return 0, 0, fmt.Errorf("expected field %q but found %q in %q", t.name, fields[t.index], statPath)
		***REMOVED***
		v, err := strconv.ParseUint(fields[t.index+1], 10, 64)
		if err != nil ***REMOVED***
			return 0, 0, err
		***REMOVED***
		*t.value = v
	***REMOVED***
	return (user * nanosecondsInSecond) / clockTicks, (kernel * nanosecondsInSecond) / clockTicks, nil
***REMOVED***
