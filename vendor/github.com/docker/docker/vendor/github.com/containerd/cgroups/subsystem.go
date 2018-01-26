package cgroups

import (
	"fmt"

	specs "github.com/opencontainers/runtime-spec/specs-go"
)

// Name is a typed name for a cgroup subsystem
type Name string

const (
	Devices   Name = "devices"
	Hugetlb   Name = "hugetlb"
	Freezer   Name = "freezer"
	Pids      Name = "pids"
	NetCLS    Name = "net_cls"
	NetPrio   Name = "net_prio"
	PerfEvent Name = "perf_event"
	Cpuset    Name = "cpuset"
	Cpu       Name = "cpu"
	Cpuacct   Name = "cpuacct"
	Memory    Name = "memory"
	Blkio     Name = "blkio"
)

// Subsystems returns a complete list of the default cgroups
// avaliable on most linux systems
func Subsystems() []Name ***REMOVED***
	n := []Name***REMOVED***
		Hugetlb,
		Freezer,
		Pids,
		NetCLS,
		NetPrio,
		PerfEvent,
		Cpuset,
		Cpu,
		Cpuacct,
		Memory,
		Blkio,
	***REMOVED***
	if !isUserNS ***REMOVED***
		n = append(n, Devices)
	***REMOVED***
	return n
***REMOVED***

type Subsystem interface ***REMOVED***
	Name() Name
***REMOVED***

type pather interface ***REMOVED***
	Subsystem
	Path(path string) string
***REMOVED***

type creator interface ***REMOVED***
	Subsystem
	Create(path string, resources *specs.LinuxResources) error
***REMOVED***

type deleter interface ***REMOVED***
	Subsystem
	Delete(path string) error
***REMOVED***

type stater interface ***REMOVED***
	Subsystem
	Stat(path string, stats *Metrics) error
***REMOVED***

type updater interface ***REMOVED***
	Subsystem
	Update(path string, resources *specs.LinuxResources) error
***REMOVED***

// SingleSubsystem returns a single cgroup subsystem within the base Hierarchy
func SingleSubsystem(baseHierarchy Hierarchy, subsystem Name) Hierarchy ***REMOVED***
	return func() ([]Subsystem, error) ***REMOVED***
		subsystems, err := baseHierarchy()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		for _, s := range subsystems ***REMOVED***
			if s.Name() == subsystem ***REMOVED***
				return []Subsystem***REMOVED***
					s,
				***REMOVED***, nil
			***REMOVED***
		***REMOVED***
		return nil, fmt.Errorf("unable to find subsystem %s", subsystem)
	***REMOVED***
***REMOVED***
