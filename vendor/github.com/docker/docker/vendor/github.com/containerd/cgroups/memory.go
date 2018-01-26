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
	"syscall"

	"golang.org/x/sys/unix"

	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func NewMemory(root string) *memoryController ***REMOVED***
	return &memoryController***REMOVED***
		root: filepath.Join(root, string(Memory)),
	***REMOVED***
***REMOVED***

type memoryController struct ***REMOVED***
	root string
***REMOVED***

func (m *memoryController) Name() Name ***REMOVED***
	return Memory
***REMOVED***

func (m *memoryController) Path(path string) string ***REMOVED***
	return filepath.Join(m.root, path)
***REMOVED***

func (m *memoryController) Create(path string, resources *specs.LinuxResources) error ***REMOVED***
	if err := os.MkdirAll(m.Path(path), defaultDirPerm); err != nil ***REMOVED***
		return err
	***REMOVED***
	if resources.Memory == nil ***REMOVED***
		return nil
	***REMOVED***
	if resources.Memory.Kernel != nil ***REMOVED***
		// Check if kernel memory is enabled
		// We have to limit the kernel memory here as it won't be accounted at all
		// until a limit is set on the cgroup and limit cannot be set once the
		// cgroup has children, or if there are already tasks in the cgroup.
		for _, i := range []int64***REMOVED***1, -1***REMOVED*** ***REMOVED***
			if err := ioutil.WriteFile(
				filepath.Join(m.Path(path), "memory.kmem.limit_in_bytes"),
				[]byte(strconv.FormatInt(i, 10)),
				defaultFilePerm,
			); err != nil ***REMOVED***
				return checkEBUSY(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return m.set(path, getMemorySettings(resources))
***REMOVED***

func (m *memoryController) Update(path string, resources *specs.LinuxResources) error ***REMOVED***
	if resources.Memory == nil ***REMOVED***
		return nil
	***REMOVED***
	g := func(v *int64) bool ***REMOVED***
		return v != nil && *v > 0
	***REMOVED***
	settings := getMemorySettings(resources)
	if g(resources.Memory.Limit) && g(resources.Memory.Swap) ***REMOVED***
		// if the updated swap value is larger than the current memory limit set the swap changes first
		// then set the memory limit as swap must always be larger than the current limit
		current, err := readUint(filepath.Join(m.Path(path), "memory.limit_in_bytes"))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if current < uint64(*resources.Memory.Swap) ***REMOVED***
			settings[0], settings[1] = settings[1], settings[0]
		***REMOVED***
	***REMOVED***
	return m.set(path, settings)
***REMOVED***

func (m *memoryController) Stat(path string, stats *Metrics) error ***REMOVED***
	f, err := os.Open(filepath.Join(m.Path(path), "memory.stat"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer f.Close()
	stats.Memory = &MemoryStat***REMOVED***
		Usage:     &MemoryEntry***REMOVED******REMOVED***,
		Swap:      &MemoryEntry***REMOVED******REMOVED***,
		Kernel:    &MemoryEntry***REMOVED******REMOVED***,
		KernelTCP: &MemoryEntry***REMOVED******REMOVED***,
	***REMOVED***
	if err := m.parseStats(f, stats.Memory); err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, t := range []struct ***REMOVED***
		module string
		entry  *MemoryEntry
	***REMOVED******REMOVED***
		***REMOVED***
			module: "",
			entry:  stats.Memory.Usage,
		***REMOVED***,
		***REMOVED***
			module: "memsw",
			entry:  stats.Memory.Swap,
		***REMOVED***,
		***REMOVED***
			module: "kmem",
			entry:  stats.Memory.Kernel,
		***REMOVED***,
		***REMOVED***
			module: "kmem.tcp",
			entry:  stats.Memory.KernelTCP,
		***REMOVED***,
	***REMOVED*** ***REMOVED***
		for _, tt := range []struct ***REMOVED***
			name  string
			value *uint64
		***REMOVED******REMOVED***
			***REMOVED***
				name:  "usage_in_bytes",
				value: &t.entry.Usage,
			***REMOVED***,
			***REMOVED***
				name:  "max_usage_in_bytes",
				value: &t.entry.Max,
			***REMOVED***,
			***REMOVED***
				name:  "failcnt",
				value: &t.entry.Failcnt,
			***REMOVED***,
			***REMOVED***
				name:  "limit_in_bytes",
				value: &t.entry.Limit,
			***REMOVED***,
		***REMOVED*** ***REMOVED***
			parts := []string***REMOVED***"memory"***REMOVED***
			if t.module != "" ***REMOVED***
				parts = append(parts, t.module)
			***REMOVED***
			parts = append(parts, tt.name)
			v, err := readUint(filepath.Join(m.Path(path), strings.Join(parts, ".")))
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			*tt.value = v
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (m *memoryController) OOMEventFD(path string) (uintptr, error) ***REMOVED***
	root := m.Path(path)
	f, err := os.Open(filepath.Join(root, "memory.oom_control"))
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	defer f.Close()
	fd, _, serr := unix.RawSyscall(unix.SYS_EVENTFD2, 0, unix.EFD_CLOEXEC, 0)
	if serr != 0 ***REMOVED***
		return 0, serr
	***REMOVED***
	if err := writeEventFD(root, f.Fd(), fd); err != nil ***REMOVED***
		unix.Close(int(fd))
		return 0, err
	***REMOVED***
	return fd, nil
***REMOVED***

func writeEventFD(root string, cfd, efd uintptr) error ***REMOVED***
	f, err := os.OpenFile(filepath.Join(root, "cgroup.event_control"), os.O_WRONLY, 0)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err = f.WriteString(fmt.Sprintf("%d %d", efd, cfd))
	f.Close()
	return err
***REMOVED***

func (m *memoryController) parseStats(r io.Reader, stat *MemoryStat) error ***REMOVED***
	var (
		raw  = make(map[string]uint64)
		sc   = bufio.NewScanner(r)
		line int
	)
	for sc.Scan() ***REMOVED***
		if err := sc.Err(); err != nil ***REMOVED***
			return err
		***REMOVED***
		key, v, err := parseKV(sc.Text())
		if err != nil ***REMOVED***
			return fmt.Errorf("%d: %v", line, err)
		***REMOVED***
		raw[key] = v
		line++
	***REMOVED***
	stat.Cache = raw["cache"]
	stat.RSS = raw["rss"]
	stat.RSSHuge = raw["rss_huge"]
	stat.MappedFile = raw["mapped_file"]
	stat.Dirty = raw["dirty"]
	stat.Writeback = raw["writeback"]
	stat.PgPgIn = raw["pgpgin"]
	stat.PgPgOut = raw["pgpgout"]
	stat.PgFault = raw["pgfault"]
	stat.PgMajFault = raw["pgmajfault"]
	stat.InactiveAnon = raw["inactive_anon"]
	stat.ActiveAnon = raw["active_anon"]
	stat.InactiveFile = raw["inactive_file"]
	stat.ActiveFile = raw["active_file"]
	stat.Unevictable = raw["unevictable"]
	stat.HierarchicalMemoryLimit = raw["hierarchical_memory_limit"]
	stat.HierarchicalSwapLimit = raw["hierarchical_memsw_limit"]
	stat.TotalCache = raw["total_cache"]
	stat.TotalRSS = raw["total_rss"]
	stat.TotalRSSHuge = raw["total_rss_huge"]
	stat.TotalMappedFile = raw["total_mapped_file"]
	stat.TotalDirty = raw["total_dirty"]
	stat.TotalWriteback = raw["total_writeback"]
	stat.TotalPgPgIn = raw["total_pgpgin"]
	stat.TotalPgPgOut = raw["total_pgpgout"]
	stat.TotalPgFault = raw["total_pgfault"]
	stat.TotalPgMajFault = raw["total_pgmajfault"]
	stat.TotalInactiveAnon = raw["total_inactive_anon"]
	stat.TotalActiveAnon = raw["total_active_anon"]
	stat.TotalInactiveFile = raw["total_inactive_file"]
	stat.TotalActiveFile = raw["total_active_file"]
	stat.TotalUnevictable = raw["total_unevictable"]
	return nil
***REMOVED***

func (m *memoryController) set(path string, settings []memorySettings) error ***REMOVED***
	for _, t := range settings ***REMOVED***
		if t.value != nil ***REMOVED***
			if err := ioutil.WriteFile(
				filepath.Join(m.Path(path), fmt.Sprintf("memory.%s", t.name)),
				[]byte(strconv.FormatInt(*t.value, 10)),
				defaultFilePerm,
			); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

type memorySettings struct ***REMOVED***
	name  string
	value *int64
***REMOVED***

func getMemorySettings(resources *specs.LinuxResources) []memorySettings ***REMOVED***
	mem := resources.Memory
	var swappiness *int64
	if mem.Swappiness != nil ***REMOVED***
		v := int64(*mem.Swappiness)
		swappiness = &v
	***REMOVED***
	return []memorySettings***REMOVED***
		***REMOVED***
			name:  "limit_in_bytes",
			value: mem.Limit,
		***REMOVED***,
		***REMOVED***
			name:  "memsw.limit_in_bytes",
			value: mem.Swap,
		***REMOVED***,
		***REMOVED***
			name:  "kmem.limit_in_bytes",
			value: mem.Kernel,
		***REMOVED***,
		***REMOVED***
			name:  "kmem.tcp.limit_in_bytes",
			value: mem.KernelTCP,
		***REMOVED***,
		***REMOVED***
			name:  "oom_control",
			value: getOomControlValue(mem),
		***REMOVED***,
		***REMOVED***
			name:  "swappiness",
			value: swappiness,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func checkEBUSY(err error) error ***REMOVED***
	if pathErr, ok := err.(*os.PathError); ok ***REMOVED***
		if errNo, ok := pathErr.Err.(syscall.Errno); ok ***REMOVED***
			if errNo == unix.EBUSY ***REMOVED***
				return fmt.Errorf(
					"failed to set memory.kmem.limit_in_bytes, because either tasks have already joined this cgroup or it has children")
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return err
***REMOVED***

func getOomControlValue(mem *specs.LinuxMemory) *int64 ***REMOVED***
	if mem.DisableOOMKiller != nil && *mem.DisableOOMKiller ***REMOVED***
		i := int64(1)
		return &i
	***REMOVED***
	return nil
***REMOVED***
