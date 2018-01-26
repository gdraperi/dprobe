package cgroups

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	units "github.com/docker/go-units"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

var isUserNS = runningInUserNS()

// runningInUserNS detects whether we are currently running in a user namespace.
// Copied from github.com/lxc/lxd/shared/util.go
func runningInUserNS() bool ***REMOVED***
	file, err := os.Open("/proc/self/uid_map")
	if err != nil ***REMOVED***
		// This kernel-provided file only exists if user namespaces are supported
		return false
	***REMOVED***
	defer file.Close()

	buf := bufio.NewReader(file)
	l, _, err := buf.ReadLine()
	if err != nil ***REMOVED***
		return false
	***REMOVED***

	line := string(l)
	var a, b, c int64
	fmt.Sscanf(line, "%d %d %d", &a, &b, &c)
	/*
	 * We assume we are in the initial user namespace if we have a full
	 * range - 4294967295 uids starting at uid 0.
	 */
	if a == 0 && b == 0 && c == 4294967295 ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

// defaults returns all known groups
func defaults(root string) ([]Subsystem, error) ***REMOVED***
	h, err := NewHugetlb(root)
	if err != nil && !os.IsNotExist(err) ***REMOVED***
		return nil, err
	***REMOVED***
	s := []Subsystem***REMOVED***
		NewNamed(root, "systemd"),
		NewFreezer(root),
		NewPids(root),
		NewNetCls(root),
		NewNetPrio(root),
		NewPerfEvent(root),
		NewCputset(root),
		NewCpu(root),
		NewCpuacct(root),
		NewMemory(root),
		NewBlkio(root),
	***REMOVED***
	// only add the devices cgroup if we are not in a user namespace
	// because modifications are not allowed
	if !isUserNS ***REMOVED***
		s = append(s, NewDevices(root))
	***REMOVED***
	// add the hugetlb cgroup if error wasn't due to missing hugetlb
	// cgroup support on the host
	if err == nil ***REMOVED***
		s = append(s, h)
	***REMOVED***
	return s, nil
***REMOVED***

// remove will remove a cgroup path handling EAGAIN and EBUSY errors and
// retrying the remove after a exp timeout
func remove(path string) error ***REMOVED***
	delay := 10 * time.Millisecond
	for i := 0; i < 5; i++ ***REMOVED***
		if i != 0 ***REMOVED***
			time.Sleep(delay)
			delay *= 2
		***REMOVED***
		if err := os.RemoveAll(path); err == nil ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***
	return fmt.Errorf("cgroups: unable to remove path %q", path)
***REMOVED***

// readPids will read all the pids in a cgroup by the provided path
func readPids(path string, subsystem Name) ([]Process, error) ***REMOVED***
	f, err := os.Open(filepath.Join(path, cgroupProcs))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer f.Close()
	var (
		out []Process
		s   = bufio.NewScanner(f)
	)
	for s.Scan() ***REMOVED***
		if t := s.Text(); t != "" ***REMOVED***
			pid, err := strconv.Atoi(t)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			out = append(out, Process***REMOVED***
				Pid:       pid,
				Subsystem: subsystem,
				Path:      path,
			***REMOVED***)
		***REMOVED***
	***REMOVED***
	return out, nil
***REMOVED***

func hugePageSizes() ([]string, error) ***REMOVED***
	var (
		pageSizes []string
		sizeList  = []string***REMOVED***"B", "kB", "MB", "GB", "TB", "PB"***REMOVED***
	)
	files, err := ioutil.ReadDir("/sys/kernel/mm/hugepages")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	for _, st := range files ***REMOVED***
		nameArray := strings.Split(st.Name(), "-")
		pageSize, err := units.RAMInBytes(nameArray[1])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		pageSizes = append(pageSizes, units.CustomSize("%g%s", float64(pageSize), 1024.0, sizeList))
	***REMOVED***
	return pageSizes, nil
***REMOVED***

func readUint(path string) (uint64, error) ***REMOVED***
	v, err := ioutil.ReadFile(path)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return parseUint(strings.TrimSpace(string(v)), 10, 64)
***REMOVED***

func parseUint(s string, base, bitSize int) (uint64, error) ***REMOVED***
	v, err := strconv.ParseUint(s, base, bitSize)
	if err != nil ***REMOVED***
		intValue, intErr := strconv.ParseInt(s, base, bitSize)
		// 1. Handle negative values greater than MinInt64 (and)
		// 2. Handle negative values lesser than MinInt64
		if intErr == nil && intValue < 0 ***REMOVED***
			return 0, nil
		***REMOVED*** else if intErr != nil &&
			intErr.(*strconv.NumError).Err == strconv.ErrRange &&
			intValue < 0 ***REMOVED***
			return 0, nil
		***REMOVED***
		return 0, err
	***REMOVED***
	return v, nil
***REMOVED***

func parseKV(raw string) (string, uint64, error) ***REMOVED***
	parts := strings.Fields(raw)
	switch len(parts) ***REMOVED***
	case 2:
		v, err := parseUint(parts[1], 10, 64)
		if err != nil ***REMOVED***
			return "", 0, err
		***REMOVED***
		return parts[0], v, nil
	default:
		return "", 0, ErrInvalidFormat
	***REMOVED***
***REMOVED***

func parseCgroupFile(path string) (map[string]string, error) ***REMOVED***
	f, err := os.Open(path)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer f.Close()
	return parseCgroupFromReader(f)
***REMOVED***

func parseCgroupFromReader(r io.Reader) (map[string]string, error) ***REMOVED***
	var (
		cgroups = make(map[string]string)
		s       = bufio.NewScanner(r)
	)
	for s.Scan() ***REMOVED***
		if err := s.Err(); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		var (
			text  = s.Text()
			parts = strings.SplitN(text, ":", 3)
		)
		if len(parts) < 3 ***REMOVED***
			return nil, fmt.Errorf("invalid cgroup entry: %q", text)
		***REMOVED***
		for _, subs := range strings.Split(parts[1], ",") ***REMOVED***
			if subs != "" ***REMOVED***
				cgroups[subs] = parts[2]
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return cgroups, nil
***REMOVED***

func getCgroupDestination(subsystem string) (string, error) ***REMOVED***
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() ***REMOVED***
		if err := s.Err(); err != nil ***REMOVED***
			return "", err
		***REMOVED***
		fields := strings.Fields(s.Text())
		for _, opt := range strings.Split(fields[len(fields)-1], ",") ***REMOVED***
			if opt == subsystem ***REMOVED***
				return fields[3], nil
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return "", ErrNoCgroupMountDestination
***REMOVED***

func pathers(subystems []Subsystem) []pather ***REMOVED***
	var out []pather
	for _, s := range subystems ***REMOVED***
		if p, ok := s.(pather); ok ***REMOVED***
			out = append(out, p)
		***REMOVED***
	***REMOVED***
	return out
***REMOVED***

func initializeSubsystem(s Subsystem, path Path, resources *specs.LinuxResources) error ***REMOVED***
	if c, ok := s.(creator); ok ***REMOVED***
		p, err := path(s.Name())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := c.Create(p, resources); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else if c, ok := s.(pather); ok ***REMOVED***
		p, err := path(s.Name())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		// do the default create if the group does not have a custom one
		if err := os.MkdirAll(c.Path(p), defaultDirPerm); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func cleanPath(path string) string ***REMOVED***
	if path == "" ***REMOVED***
		return ""
	***REMOVED***
	path = filepath.Clean(path)
	if !filepath.IsAbs(path) ***REMOVED***
		path, _ = filepath.Rel(string(os.PathSeparator), filepath.Clean(string(os.PathSeparator)+path))
	***REMOVED***
	return filepath.Clean(path)
***REMOVED***
